package service

import (
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"

	log "github.com/sirupsen/logrus"

	m "github.com/fedorkolmykow/avitojob/pkg/models"

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
	Get(key string) (value string, err error)
	Set(key string, value string) (err error)
	Delete(key string) (err error)
}

type service struct{
	db dbClient
	cash cashClient
}

func (s *service) ChangeBalance(Req *m.ChangeBalanceReq) (Resp *m.ChangeBalanceResp, err error) {
	err = Req.Validate()
	if err != nil{
		return
	}
	Resp, err = s.db.UpdateBalance(Req)
	return
}

func (s *service) Transfer(Req *m.TransferReq) (Resp *m.TransferResp, err error) {
	err = Req.Validate()
	if err != nil{
		return
	}
	Resp, err = s.db.UpdateBalances(Req)
	return
}

func (s *service) GetBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error) {
	err = Req.Validate()
	if err != nil{
		return
	}
	Resp, err = s.db.SelectBalance(Req)
	if err != nil {
		return
	}
	if Req.Currency != ""{
		var rate float64
		rate, err = getRate(s.cash, Req.Currency)
		if err != nil{
			return
		}
		Resp.Currency = Req.Currency
		Resp.Balance = Resp.Balance * rate
	} else{
		Resp.Currency = "RUB"
	}
	return
}

func (s *service) GetTransactions(Req *m.GetTransactionsReq) (Resp *m.GetTransactionsResp, err error) {
	err = Req.Validate()
	if err != nil{
		return
	}
	Resp = &m.GetTransactionsResp{}
	log.Trace(Resp)
	trs, err := s.db.SelectTransactions(Req)
	if err != nil{
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

func getRate(cash cashClient, currency string) (rate float64, err error){
	var strRate string
	strRate, err = cash.Get("Rate:0:" + currency)
	if err == nil{
		rate, err = strconv.ParseFloat(strRate, 64)
		if err == nil {
			log.Trace("read rate from cash")
			return
		}
		log.Trace(err)
	} else{
		log.Trace(err)
		err = nil
	}
	var r *http.Response
	var body []byte
	r, err = http.Get(os.Getenv("CURRENCY_URL")+currency)
	if err != nil {
		return
	}
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	resp := &m.Rate{}
	err = resp.UnmarshalJSON(body)
	if err != nil {
		return
	}
	strRate = strconv.FormatFloat(resp.Rates[currency], 'e', -1, 64)
	err = cash.Set("Rate:0:" + currency, strRate)
	if err != nil {
		return
	}
	rate = resp.Rates[currency]
	return
}

func NewService(db dbClient, cash cashClient) Service{
    svc := &service{
    	db: db,
    	cash: cash,
	}
    return svc
}