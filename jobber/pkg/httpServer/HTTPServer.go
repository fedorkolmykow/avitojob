package httpServer

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	m "github.com/fedorkolmykow/avitojob/pkg/models"
)

type service interface {
	ChangeBalance(Req *m.ChangeBalanceReq) (Resp *m.ChangeBalanceResp, err error)
	Transfer(Req *m.TransferReq) (Resp *m.TransferResp, err error)
	GetBalance(Req *m.GetBalanceReq) (Resp *m.GetBalanceResp, err error)
	GetTransactions(Req *m.GetTransactionsReq) (Resp *m.GetTransactionsResp, err error)
}

type server struct {
	svc service
}


func (s *server) HandleChangeBalance(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	UserID, err := strconv.Atoi(vars["user_id"])
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req := &m.ChangeBalanceReq{}
	err = req.UnmarshalJSON(body)
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.UserId = UserID
	log.Trace("Received data: " + fmt.Sprintf("%+v", req))
	resp, err := s.svc.ChangeBalance(req)
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err = resp.MarshalJSON()
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) HandleTransfer(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	UserID, err := strconv.Atoi(vars["user_id"])
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req := &m.TransferReq{}
	err = req.UnmarshalJSON(body)
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.UserId = UserID
	log.Trace("Received data: " + fmt.Sprintf("%+v", req))
	resp, err := s.svc.Transfer(req)
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err = resp.MarshalJSON()
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) HandleBalanceGet(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	UserID, err := strconv.Atoi(vars["user_id"])
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req := &m.GetBalanceReq{}
	req.UserId = UserID
	req.Currency = r.FormValue("currency")
	log.Trace("Received data: " + fmt.Sprintf("%+v", req))
	resp, err := s.svc.GetBalance(req)
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err := resp.MarshalJSON()
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *server) HandleTransactionsGet(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	UserID, err := strconv.Atoi(vars["user_id"])
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req := &m.GetTransactionsReq{}
	err = req.UnmarshalJSON(body)
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.UserId = UserID
	log.Trace("Received data: " + fmt.Sprintf("%+v", req))
	resp, err := s.svc.GetTransactions(req)
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	body, err = resp.MarshalJSON()
	if err != nil{
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(body)
	if err != nil {
		log.Warn(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NewHTTPServer(svc service) (httpServer *mux.Router) {
	router := mux.NewRouter()
    s := &server{svc: svc}
	router.HandleFunc("/users/{user_id:[0-9]+}/balance", s.HandleChangeBalance).
		Methods("PATCH")
	router.HandleFunc("/users/{user_id:[0-9]+}/balance/transfer", s.HandleTransfer).
		Methods("PATCH")
	router.HandleFunc("/users/{user_id:[0-9]+}/balance", s.HandleBalanceGet).
		Methods("GET")
	router.HandleFunc("/users/{user_id:[0-9]+}/transactions", s.HandleTransactionsGet).
		Methods("POST")
	return router
}