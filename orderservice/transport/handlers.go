package transport

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

var orders = `[{
	"id": "d290f1ee-6c54-4b01-90E6-d701748fo851",
	"menuitems": [{
		"id": "f290d1ce-6c234-4b31-90e6-d701748fo851",
		"quantity": 1
	}]
}]`

var order = `[{	
	"id": "d290f1ee-6c54-4b01-90E6-d701748fo851",
	"menuitems": [{
		"id": "f290d1ce-6c234-4b31-90e6-d701748fo851",
		"quantity": 1
	}]
	"orderedAtTimestamp": 1613758423,
	"cost": 999
}]`

type Menu struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	Id        string `json:"id"`
	Menuitems Menu`json:"menuitems"`
	OrderedAtTimestamp int `json:"orderedAtTimestamp"`
	Cost int `json:"cost"`
}


func getOrder(w http.ResponseWriter, r *http.Request) {
	var menuitem = Menu{
		Id:       "f290d1ce-6c234-4b31-90e6-d701748fo851",
		Quantity: 1,
	}
	var order = Order {
		Id:        "d290f1ee-6c54-4b01-90E6-d701748fo851",
		Menuitems: menuitem,
		OrderedAtTimestamp: 1613758423,
		Cost: 999,
	}
	vars := mux.Vars(r)
	id := vars["ID"]
	if id == order.Id {
		b, err := json.Marshal(order)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if _, err = io.WriteString(w, string(b)); err != nil {
			log.WithField("err", err).Error("write responce error")
		}
	}
}

func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/orders/{ID}", getOrder)
	return logMiddleware(r)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		endTime := time.Now()
		log.WithFields(log.Fields{
			"method": r.Method,
			"url": r.URL,
			"remoteAddr": r.RemoteAddr,
			"time": int(endTime.Sub(startTime)),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
		endTime = time.Now()
	})
}

func getOrders(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, orders)
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "hello world")
}
