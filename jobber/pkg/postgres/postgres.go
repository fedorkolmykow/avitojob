package postgres

import (
	"errors"
	"fmt"
	m "github.com/fedorkolmykow/avitojob/pkg/models"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

const(
	CheckExistence = `SELECT EXISTS(SELECT user_id FROM Users WHERE user_id=$1) ;`
	SelectUserBalance = `SELECT balance FROM Users WHERE user_id=$1;`
	InsertUser = `INSERT INTO Users (user_id, balance) VALUES ($1, $2) RETURNING user_id;`
	UpdateUserBalance = `UPDATE Users SET balance = balance + $1 WHERE user_id = $2 RETURNING balance;`
	SetIsolationSerializable = `SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;`
	InsertTrans = `INSERT INTO Transactions (user_id, init_balance, change, time, comment)  
                     VALUES (:user_id, :init_balance, :change, :time, :comment);`
)

type DbClient interface{
	UpdateBalance(Req *m.ChangeBalanceReq) (Resp *m.ChangeBalanceResp, err error)
	UpdateBalances(Req *m.TransferReq) (Resp *m.TransferResp, err error)
	SelectBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error)
	SelectTransactions(Req *m.GetTransactionsReq) (Resp *m.Transactions, err error)
}

type dbClient struct{
    db *sqlx.DB
}

func insertTransaction(tx *sqlx.Tx, trans *m.Transaction) error {
	trans.ChangeTime = time.Now().Format(time.RFC822)
	_, err := tx.NamedExec(InsertTrans, &trans)
	log.Trace("inserted transaction with data: " + fmt.Sprintf("%#v", trans))
	return err
}

func rollAndErr(tx *sqlx.Tx, err error) error{
	log.Trace("Rollback")
	log.Warn(err)
	errRoll := tx.Rollback()
	if errRoll != nil{
		log.Warn(errRoll)
		return errRoll
	}
	return err
}

func createNewUser(tx *sqlx.Tx, tr *m.Transaction) (exists bool, err error){
	err = tx.QueryRow(CheckExistence, tr.UserId).Scan(&exists)
	if err != nil{
		return
	}
	if !exists{
		if tr.Change < 0{
			_ = rollAndErr(tx, err)
			err = errors.New("negative initial balance")
			return
		}
		_, err = tx.Exec(InsertUser, tr.UserId ,tr.Change)
		if err != nil{
			return
		}
		log.Trace("Created new user")
		err = insertTransaction(tx, tr)
		return
	}
	return
}

func changeBalance(tx *sqlx.Tx, tr *m.Transaction) (balance float64 ,err error){
	if tr.InitialBalance + tr.Change < 0{
		_ = rollAndErr(tx, err)
		err = errors.New("negative balance")
		return
	}
	err = tx.QueryRow(UpdateUserBalance, tr.Change, tr.UserId).Scan(&balance)
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	err = insertTransaction(tx, tr)
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	return
}

func (d *dbClient) UpdateBalance(Req *m.ChangeBalanceReq) (Resp *m.ChangeBalanceResp, err error){
	trans := &m.Transaction{
		Change: Req.Change,
		UserId: Req.UserId,
		Comment: Req.Comment,
		Source: Req.Source,
	}
	Resp = &m.ChangeBalanceResp{
		UserId: Req.UserId,
		Balance: Req.Change,
	}
	tx, err := d.db.Beginx()
	if err != nil{
		log.Warn(err)
		return
	}
	_, err = tx.Exec(SetIsolationSerializable)
	if err != nil{
		log.Warn(err)
		return
	}
	log.Trace("start transaction with data: " + fmt.Sprintf("%#v", Req))
	existed, err := createNewUser(tx, trans)
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	if existed{
		err = tx.QueryRow(SelectUserBalance, Req.UserId).Scan(&trans.InitialBalance)
		if err != nil{
			err = rollAndErr(tx, err)
			return
		}
		Resp.Balance, err = changeBalance(tx, trans)
		if err != nil{
			err = rollAndErr(tx, err)
			return
		}
	}
	err = tx.Commit()
	if err != nil{
		log.Warn(err)
		return
	}
	log.Trace("changed balance, result: " + fmt.Sprintf("%#v", Resp))
	return
}

func (d *dbClient) UpdateBalances(Req *m.TransferReq) (Resp *m.TransferResp, err error){
	var exists bool
	Resp = &m.TransferResp{
		Source: m.ChangeBalanceResp{UserId: Req.UserId},
		Target: m.ChangeBalanceResp{UserId: Req.TargetId},
	}
	sourceTrans := &m.Transaction{
		Change: -Req.Change,
		UserId: Req.UserId,
		Comment: Req.Comment,
		Source: strconv.Itoa(Req.UserId),
	}
	targetTrans := &m.Transaction{
		Change: Req.Change,
		UserId: Req.TargetId,
		Comment: Req.Comment,
		Source: strconv.Itoa(Req.UserId),
	}
	tx, err := d.db.Beginx()
	if err != nil{
		log.Warn(err)
		return
	}
	_, err = tx.Exec(SetIsolationSerializable)
	if err != nil{
		log.Warn(err)
		return
	}
	log.Trace("start transaction with data: " + fmt.Sprintf("%#v", Req))
	err = tx.QueryRow(CheckExistence, sourceTrans.UserId).Scan(&exists)
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	if !exists{
		_ = rollAndErr(tx, err)
		err = errors.New("user don't have balance")
		return
	}
	err = tx.QueryRow(SelectUserBalance, Req.UserId).Scan(&sourceTrans.InitialBalance)
	if err != nil{
			err = rollAndErr(tx, err)
			return
		}
	Resp.Source.Balance, err = changeBalance(tx, sourceTrans)
	if err != nil{
			err = rollAndErr(tx, err)
			return
		}

	existed, err := createNewUser(tx, targetTrans)
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	if existed{
		err = tx.QueryRow(SelectUserBalance, Req.TargetId).Scan(&targetTrans.InitialBalance)
		if err != nil{
			err = rollAndErr(tx, err)
			return
		}
		Resp.Target.Balance, err = changeBalance(tx, targetTrans)
		if err != nil{
			err = rollAndErr(tx, err)
			return
		}
	}
	err = tx.Commit()
	if err != nil{
		log.Warn(err)
		return
	}
	log.Trace("changed balances, result: " + fmt.Sprintf("%#v", Resp))
	return
}
func (d *dbClient) SelectBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error){
	var exists bool
	Resp = &m.GetBalanceResp{UserId: Req.UserId}
	tx, err := d.db.Beginx()
	if err != nil{
		log.Warn(err)
		return
	}
	err = tx.QueryRow(CheckExistence, Req.UserId).Scan(&exists)
	if !exists{
		_ = rollAndErr(tx, err)
		err = errors.New("user don't have balance")
		return
	}
	err = tx.QueryRow(SelectUserBalance, Req.UserId).Scan(&Resp.Balance)
	if err != nil{
		err = rollAndErr(tx, err)
		return
	}
	err = tx.Commit()
	if err != nil{
		log.Warn(err)
		return
	}
	return
}
func (d *dbClient) SelectTransactions(Req *m.GetTransactionsReq) (Resp *m.GetTransactionsResp, err error){
	return
}

func NewDbClient() DbClient{
	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	//db.SetMaxIdleConns(n int)
	//db.SetMaxOpenConns(n int)


	return &dbClient{db: db}
}