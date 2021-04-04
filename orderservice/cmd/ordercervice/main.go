package main

import (
	"context"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"net/http"
	model2 "orderservice/pkg/orderservice/model"
	transport2 "orderservice/pkg/orderservice/transport"
	"os"
	"os/signal"
	"syscall"
)

var driver = "mysql"
var dataSourceName = "root:Qwerty123@/order"

func getKillSignalChan() chan os.Signal {
	osKillSignalchan := make(chan os.Signal, 1)
	signal.Notify(osKillSignalchan, os.Interrupt, syscall.SIGTERM)
	return osKillSignalchan
}

func waitForKillSignall(killSignalChan <-chan os.Signal) {
	killSignal := <-killSignalChan
	switch killSignal {
	case os.Interrupt:
		log.Info("got SIGINT...")
	case syscall.SIGTERM:
		log.Info("got SIGTERM")
	}
}

func startServer(serverURL string, dbServer *model2.DBServer) *http.Server {
	server := makeServer(dbServer.Database)
	router := transport2.Router(server)
	srv := &http.Server{Addr: serverURL, Handler: router}
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	return srv
}

func makeServer(db *sql.DB) *transport2.OrderRepository {
	return &transport2.OrderRepository{
		OrderService: model2.CreateOrderServiceInterface(db),
	}
}

func main() {

	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	}

	conf, err := parseEnv()
	if err == nil {

		s := createDBConnection()

		log.WithFields(log.Fields{"url": conf.SrvRESTAddress}).Info("starting server")
		killSignalChan := getKillSignalChan()
		srv := startServer(conf.SrvRESTAddress, &s)

		waitForKillSignall(killSignalChan)
		srv.Shutdown(context.Background())
	}
	log.Println("can't load config")
}

func createDBConnection() model2.DBServer {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return model2.DBServer{Database: db}
}
