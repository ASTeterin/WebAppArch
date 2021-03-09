package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	transport2 "orderservice/pkg/ordercervice/transport"
	"os"
	"os/signal"
	"syscall"
)

func getKillSignalChan() chan os.Signal {
	osKillSignalchan := make (chan os.Signal, 1)
	signal.Notify(osKillSignalchan, os.Interrupt, syscall.SIGTERM)
	return osKillSignalchan
}

func waitForKillSignall(killSignalChan <-chan  os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM")
	}
}

func startServer(serverURL string) *http.Server {
	router := transport2.Router()
	srv := &http.Server{Addr: serverURL, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	return srv
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}
	serverUrl := ":8000"
	log.WithFields(log.Fields{"url": serverUrl}).Info("starting server")
	killSignalChan := getKillSignalChan()
	srv := startServer(serverUrl)

	waitForKillSignall(killSignalChan)
	srv.Shutdown(context.Background())
}
