package main

import (
	"context"
	"github.com/fedorkolmykow/avitojob/pkg/httpServer"
	"github.com/fedorkolmykow/avitojob/pkg/postgres"
	"github.com/fedorkolmykow/avitojob/pkg/redis"
	"github.com/fedorkolmykow/avitojob/pkg/service"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)



func main() {
	log.SetFormatter(&log.JSONFormatter{})
	switch os.Getenv("LOG_LEVEL"){
		case "TRACE": log.SetLevel(log.TraceLevel)
		case "WARN": log.SetLevel(log.WarnLevel)
		case "FATAL": log.SetLevel(log.FatalLevel)
		default: log.SetLevel(log.FatalLevel)
	}
	err := os.Mkdir("logs", 0777)
	if err != nil {
		log.Warn(err)
	}
	file, err := os.OpenFile("logs/jobber.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
	    log.Warn("Failed to log to file, using default stderr")
	}

	redCon := redis.NewDb()
    dbCon := postgres.NewDbClient()
    swc := service.NewService(dbCon, redCon)
	router := httpServer.NewHTTPServer(swc)
	srv := &http.Server{
		Addr:    os.Getenv("HTTP_PORT"),
		Handler: router,
	}

	go func() {

		log.Trace("starting HTTP server at", os.Getenv("HTTP_PORT"))
		err = srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed{
			log.Fatal(err)
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	wait, err := strconv.Atoi(os.Getenv("TIME_TO_SHUTDOWN"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(wait)*time.Second)
	defer func(){
		e := redCon.Shutdown()
		if e != nil{
			log.Warn(e)
		}
		e = dbCon.Shutdown()
		if e != nil{
			log.Warn(e)
		}
		cancel()
	}()
	err = srv.Shutdown(ctx)
	if err != nil{
		log.Fatalf("Graceful Server Shutdown Failed:%+v", err)
	}
	log.Trace("Server Was Gracefully Stopped")
}