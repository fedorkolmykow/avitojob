package service

import (
	"strconv"

	m "github.com/fedorkolmykow/avitojob/pkg/models"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	ChangeBalance(Req *m.ChangeBalanceReq) (Resp *m.ChangeBalanceResp, err error)
	Transfer(Req *m.TransferReq) (Resp *m.GetTransactionsResp, err error)
	GetBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error)
	GetTransactions(Req *m.GetTransactionsReq) (Resp *m.GetTransactionsResp, err error)
}

type dbClient interface{
	UpdateBalance(Req *m.ChangeBalanceReq) (Resp *m.ChangeBalanceResp, err error)
	UpdateBalances(Req *m.TransferReq) (Resp *m.GetTransactionsResp, err error)
	SelectBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error)
	SelectTransactions(Req *m.GetTransactionsReq) (Resp *m.GetTransactionsResp, err error)
}

type cashClient interface{
	Get(key string) (value string, err error)
	Set(key string, value string) (err error)
	Delete(key string) (err error)
}

type service struct{
	db dbClient
	cash cashClient
}

func (s *service) ChangeBalance(Req *m.ChangeBalanceReq) (Resp *m.ChangeBalanceResp, err error) {
	id := strconv.Itoa(Req.UserId)
	err = s.cash.Delete("user:" + id + ":balance")
	if err!= nil{
		log.Trace(err)
	}
	Resp, err = s.db.UpdateBalance(Req)
	return
}

func (s *service) Transfer(Req *m.TransferReq) (Resp *m.GetTransactionsResp, err error) {
	idSource := strconv.Itoa(Req.UserId)
	idTarget := strconv.Itoa(Req.TargetId)
	err = s.cash.Delete("user:" + idSource + ":balance")
	if err!= nil{
		log.Trace(err)
	}
	err = s.cash.Delete("user:" + idTarget + ":balance")
	if err!= nil{
		log.Trace(err)
	}
	Resp, err = s.db.UpdateBalances(Req)
	return
}

func (s *service) GetBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error) {
	Resp, err = s.getCashedBalance(Req)
	if err != nil{
		log.Trace(err)
	} else{
		log.Trace("Get balance from cash")
		return
	}
	Resp, err = s.db.SelectBalance(Req)
	return
}

func (s *service) getCashedBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error){
	id := strconv.Itoa(Req.UserId)
	balance, err := s.cash.Get("user:" + id + ":balance")
	if err != nil{
		return
	}
	Resp = &m.GetBalanceResp{}
	Resp.Balance, err = strconv.Atoi(balance)
	if err != nil {
		return
	}
	Resp.UserId, err = strconv.Atoi(balance)
	return
}

func (s *service) GetTransactions(Req *m.GetTransactionsReq) (Resp *m.GetTransactionsResp, err error) {
	return
}

func NewService(db dbClient, cash cashClient) Service{
    svc := &service{
    	db: db,
    	cash: cash,
	}
    return svc
}