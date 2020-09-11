package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/fedorkolmykow/avitojob/pkg/httpServer"
	"github.com/fedorkolmykow/avitojob/pkg/redis"
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
	file, err := os.OpenFile("logs/kiddy.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
	    log.Warn("Failed to log to file, using default stderr")
	}

	redCon := redis.NewDb()

	serverHTTP := httpServer.NewHTTPServer(redCon)


	go func() {
		log.Trace("starting HTTP server at", os.Getenv("HTTP_PORT"))
		err := http.ListenAndServe(os.Getenv("HTTP_PORT"), serverHTTP)
		if err != nil{
			log.Fatal(err)
		}
	}()
}