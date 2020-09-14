package service

import (
	"errors"
	"math"
	"sort"
	"strconv"

	m "github.com/fedorkolmykow/avitojob/pkg/models"
	log "github.com/sirupsen/logrus"
)

type Service interface {
	ChangeBalance(Req *m.ChangeBalanceReq) (Resp *m.ChangeBalanceResp, err error)
	Transfer(Req *m.TransferReq) (Resp *m.TransferResp, err error)
	GetBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error)
	GetTransactions(Req *m.GetTransactionsReq) (Resp *m.GetTransactionsResp, err error)
}

type dbClient interface{
	UpdateBalance(Req *m.ChangeBalanceReq) (Resp *m.ChangeBalanceResp, err error)
	UpdateBalances(Req *m.TransferReq) (Resp *m.TransferResp, err error)
	SelectBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error)
	SelectTransactions(Req *m.GetTransactionsReq) (Resp *m.Transactions, err error)
}

type cashClient interface{
	Get(key interface{}) (value string, err error)
	Set(key interface{}, value string) (err error)
	Delete(key interface{}) (err error)
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

func (s *service) Transfer(Req *m.TransferReq) (Resp *m.TransferResp, err error) {
	if Req.Change < 0{
		err = errors.New("transfer cannot be negative")
		return
	}
	Resp, err = s.db.UpdateBalances(Req)
	return
}

func (s *service) GetBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error) {
	Resp, err = s.getCashedBalance(Req)
	if err != nil{
		log.Warn(err)
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
	Resp.Balance, err = strconv.ParseFloat(balance, 64)
	if err != nil {
		return
	}
	Resp.UserId, err = strconv.Atoi(balance)
	return
}

func (s *service) GetTransactions(Req *m.GetTransactionsReq) (Resp *m.GetTransactionsResp, err error) {
	if Req.Page < 0 {
		err = errors.New("negative page")
	}
	if Req.TransactionsOnPage < 0 {
		err = errors.New("negative number of transactions on page")
	}
	Resp = &m.GetTransactionsResp{}
	log.Trace(Resp)
	trs, err := s.db.SelectTransactions(Req)
	if err != nil{
		log.Warn(err)
		return
	}
	sort.Sort(trs)
	Resp.Transactions = pagination(trs.Transactions, Req.Page, Req.TransactionsOnPage)
	Resp.UserId = Req.UserId
	return
}

func pagination(trs []m.Transaction, page int, perPage int) (tr []m.Transaction){

	firsT := (page - 1) * perPage
	if firsT > len(trs){		//if page index too big then show last page
		pagesCount := math.Ceil(float64(len(trs))/float64(perPage))
		firsT = (int(pagesCount) - 1) * perPage
	}
	if len(trs) - perPage < firsT{
		perPage = len(trs)%perPage
	}
	tr = make([]m.Transaction, perPage)
	for i := 0; i < perPage; i++{
		tr[i] = trs[firsT + i]
	}
	return
}

func NewService(db dbClient, cash cashClient) Service{
    svc := &service{
    	db: db,
    	cash: cash,
	}
    return svc
}