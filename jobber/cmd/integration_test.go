// +build integration

package main

import (
	"bytes"
	"context"
	"github.com/fedorkolmykow/avitojob/pkg/httpServer"
	"github.com/fedorkolmykow/avitojob/pkg/postgres"
	"github.com/fedorkolmykow/avitojob/pkg/redis"
	"github.com/fedorkolmykow/avitojob/pkg/service"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func start(t *testing.T) *http.Server{

	redCon := redis.NewDb()
	dbCon := postgres.NewDbClient()
	swc := service.NewService(dbCon, redCon)
	router := httpServer.NewHTTPServer(swc)
	srv := &http.Server{
		Addr:    os.Getenv("HTTP_PORT"),
		Handler: router,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed{
			t.Fatal(err)
		}
	}()
	return srv
}

type testCase struct{
	RespExpData string
	ReqData		[]byte
	Url         string
	Method      string
}

func Test_All(t *testing.T){
	monkey.Patch(time.Now, func() time.Time {
		return time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	})
	log.SetLevel(log.FatalLevel)
	log.SetOutput(os.Stdout)
	srv := start(t)
	client := &http.Client{}

	cases := []testCase{
		{
			RespExpData: `{"user_id":0,"balance":400}`,
			ReqData:     []byte(`{"change":400,"comment":"My First","source":"Sberbank"}`),
			Url:         "http://testserver:9001/users/0/balance",
			Method:      "PATCH",
		},
		{
			RespExpData: `{"user_id":0,"balance":200}`,
			ReqData:     []byte(`{"change":-200,"comment":"My First","source":"Sberbank"}`),
			Url:         "http://testserver:9001/users/0/balance",
			Method:      "PATCH",
		},
		{
			RespExpData: `{"source":{"user_id":0,"balance":0},"target":{"user_id":1,"balance":200}}`,
			ReqData:     []byte(`{"change":200,"comment":"My First","target_id":1}`),
			Url:         "http://testserver:9001/users/0/balance/transfer",
			Method:      "PATCH",
		},
		{
			RespExpData: `{"user_id":0,"balance":0,"currency":"RUB"}`,
			ReqData:     []byte(``),
			Url:         "http://testserver:9001/users/0/balance",
			Method:      "GET",
		},
		{
			RespExpData: `{"user_id":1,"transactions":[{"init_balance":0,"change":200,"change_time":"17 Nov 09 20:34 UTC","source":"0","comment":"My First"}]}`,
			ReqData:     []byte(`{"page":1,"per_page":1,"change_sort":false,"time_sort":false}`),
			Url:         "http://testserver:9001/users/1/transactions",
			Method:      "POST",
		},
	}

	for num, c := range cases {
		req, err := http.NewRequest(c.Method, c.Url, bytes.NewBuffer(c.ReqData))
		if err != nil {
			t.Error(err)
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			t.Error(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}
		if string(body) != c.RespExpData {
			t.Errorf("[%d] unexpected result:\n%s\nexpected:\n%s ",num ,body ,c.RespExpData)
		}
	}

	err := srv.Shutdown(context.Background())
	if err != nil{
		t.Error(err)
	}
}
