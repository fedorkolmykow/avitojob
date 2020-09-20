package httpServer

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	m "github.com/fedorkolmykow/avitojob/pkg/models"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const(
	changeBalance = iota
	transfer
	getBalance
	getTransactions
)

type correctService struct{
}

type errorService struct{
}

type TestCase struct {
	Vars   map[string]string
	Req     []byte
	Resp    string
	Status  int
	S       server
	Handle  int
}

func TestHandles(t *testing.T){
	cases := []TestCase{
		{
			Vars:        map[string]string{"user_id":"0"},
			Req:          []byte(`{"change":200,"comment":"My First","source":"Sberbank"}`),
			Resp:         `{"user_id":0,"balance":200}`,
			Status:       http.StatusOK,
			S:            server{svc: &correctService{}},
			Handle:       changeBalance,
		},
		{
			Vars:        map[string]string{"user_id":"Here is error"},
			Req:          []byte(``),
			Resp:         ``,
			Status:       http.StatusBadRequest,
			S:            server{svc: &correctService{}},
			Handle:       changeBalance,
		},
		{
			Vars:        map[string]string{"user_id":"0"},
			Req:          []byte(`Here is error`),
			Resp:         ``,
			Status:       http.StatusBadRequest,
			S:            server{svc: &correctService{}},
			Handle:       changeBalance,
		},
		{
			Vars:        map[string]string{"user_id":"0"},
			Req:          []byte(`{"change":"200","comment":0,"source":0}`),
			Resp:         `{"user_id":0,"balance":200}`,
			Status:       http.StatusBadRequest,
			S:            server{svc: &correctService{}},
			Handle:       changeBalance,
		},
		{
			Vars:        map[string]string{"user_id":"0"},
			Req:          []byte(`{}`),
			Resp:         ``,
			Status:       http.StatusInternalServerError,
			S:            server{svc: &errorService{}},
			Handle:       changeBalance,
		},
		{
			Vars:        map[string]string{"user_id":"0"},
			Req:          []byte(`{"change":200,"comment":"My First","target_id":2}`),
			Resp:         `{"source":{"user_id":0,"balance":0},"target":{"user_id":2,"balance":200}}`,
			Status:       http.StatusOK,
			S:            server{svc: &correctService{}},
			Handle:       transfer,
		},
		{
			Vars:        map[string]string{"user_id":"Here is error"},
			Req:          []byte(``),
			Resp:         ``,
			Status:       http.StatusBadRequest,
			S:            server{svc: &correctService{}},
			Handle:       transfer,
		},
		{
			Vars:        map[string]string{"user_id":"0"},
			Req:          []byte(`Here is error`),
			Resp:         ``,
			Status:       http.StatusBadRequest,
			S:            server{svc: &correctService{}},
			Handle:       transfer,
		},
		{
			Vars:        map[string]string{"user_id":"0"},
			Req:          []byte(`{}`),
			Resp:         ``,
			Status:       http.StatusInternalServerError,
			S:            server{svc: &errorService{}},
			Handle:       transfer,
		},
		{
			Vars:        map[string]string{"user_id":"0"},
			Req:          []byte(``),
			Resp:         `{"user_id":0,"balance":0,"currency":"RUB"}`,
			Status:       http.StatusOK,
			S:            server{svc: &correctService{}},
			Handle:       getBalance,
		},
		{
			Vars:        map[string]string{"user_id":"Here is error"},
			Req:          []byte(``),
			Resp:         ``,
			Status:       http.StatusBadRequest,
			S:            server{svc: &correctService{}},
			Handle:       getBalance,
		},
		{
			Vars:        map[string]string{"user_id":"0"},
			Req:          []byte(`{}`),
			Resp:         ``,
			Status:       http.StatusInternalServerError,
			S:            server{svc: &errorService{}},
			Handle:       getBalance,
		},
		{
			Vars:        map[string]string{"user_id":"0"},
			Req:          []byte(`{"page":1,"per_page":1,"change_sort":false,"time_sort":false}`),
			Resp:         `{"user_id":0,"transactions":[{"trans_id":0,"init_balance":0,"change":100,"change_time":"Now","source":"Sberbank","comment":"Test"}]}`,
			Status:       http.StatusOK,
			S:            server{svc: &correctService{}},
			Handle:       getTransactions,
		},
		{
			Vars:        map[string]string{"user_id":"Here is error"},
			Req:          []byte(``),
			Resp:         ``,
			Status:       http.StatusBadRequest,
			S:            server{svc: &correctService{}},
			Handle:       getTransactions,
		},
		{
			Vars:        map[string]string{"user_id":"0"},
			Req:          []byte(`Here is error`),
			Resp:         ``,
			Status:       http.StatusBadRequest,
			S:            server{svc: &correctService{}},
			Handle:       getTransactions,
		},
		{
			Vars:        map[string]string{"user_id":"0"},
			Req:          []byte(`{}`),
			Resp:         ``,
			Status:       http.StatusInternalServerError,
			S:            server{svc: &errorService{}},
			Handle:       getTransactions,
		},
	}
	log.SetLevel(log.FatalLevel)
	for num, c := range cases{
		req := httptest.NewRequest(
			"NotImportant",
			"http://localhost",
			bytes.NewBuffer(c.Req),
		)
		req = mux.SetURLVars(req, c.Vars)
		w := httptest.NewRecorder()
		switch c.Handle {
		case changeBalance:     c.S.HandleChangeBalance(w, req)
		case transfer:     		c.S.HandleTransfer(w, req)
		case getBalance:  		c.S.HandleBalanceGet(w, req)
		case getTransactions:   c.S.HandleTransactionsGet(w, req)
	}

		if w.Result().StatusCode != c.Status{
			t.Errorf("[%d] unexpected status: %d, expected: %d",num, w.Result().StatusCode,  c.Status)
		}
		if c.Status == http.StatusOK{
			if c.Resp != w.Body.String(){
				t.Errorf("[%d] unexpected result:\n%s\nexpected:\n%s ", num, w.Body.String(), c.Resp)
			}
		}
	}
}

//correctService
func (s *correctService)  ChangeBalance(Req *m.ChangeBalanceReq) (Resp *m.ChangeBalanceResp, err error) {
	return &m.ChangeBalanceResp{
		UserId:  Req.UserId,
		Balance: Req.Change,
	}, nil
}


func (s *correctService) Transfer(Req *m.TransferReq) (Resp *m.TransferResp, err error){
	return &m.TransferResp{
		Source: m.ChangeBalanceResp{
			UserId:  Req.UserId,
			Balance: 0,
		},
		Target: m.ChangeBalanceResp{
			UserId:  Req.TargetId,
			Balance: Req.Change,
		},
	}, nil
}


func (s *correctService) GetBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error){
	return &m.GetBalanceResp{
		UserId:   0,
		Balance:  0,
		Currency: "RUB",
	}, nil
}


func (s *correctService) GetTransactions(Req *m.GetTransactionsReq) (Resp *m.GetTransactionsResp, err error){
	return &m.GetTransactionsResp{
		UserId:       Req.UserId,
		Transactions: []m.Transaction{
			{
				TransId:        0,
				UserId:         Req.UserId,
				InitialBalance: 0,
				Change:         100,
				ChangeTime:     "Now",
				Source:         "Sberbank",
				Comment:        "Test",
			},
		},
	}, nil
}


//errorService
func (s *errorService)  ChangeBalance(Req *m.ChangeBalanceReq) (Resp *m.ChangeBalanceResp, err error) {
	return nil, errors.New("test error")
}


func (s *errorService) Transfer(Req *m.TransferReq) (Resp *m.TransferResp, err error){
	return nil, errors.New("test error")
}


func (s *errorService) GetBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error){
	return nil, errors.New("test error")
}


func (s *errorService) GetTransactions(Req *m.GetTransactionsReq) (Resp *m.GetTransactionsResp, err error){
	return nil, errors.New("test error")
}