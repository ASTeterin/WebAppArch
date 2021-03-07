package transport

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type order struct {
	Id    string `json:"id"`
	Items []menuItem
}

type orderResponse struct {
	order
	OrderedAtTimeStamp string `json:"orderedAtTimeStamp"`
	Cost               int    `json:"cost"`
}

type menuItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type orders struct {
	Orders []order
}

var menuitem = menuItem{
	Id:       "f290d1ce-6c234-4b31-90e6-d701748fo851",
	Quantity: 1,
}

var newOrder = order{
	Id:  "d290f1ee-6c54-4b01-90E6-d701748fo851",
	Items: []menuItem{
		menuitem,
	},
}

var testOrderResponse = orderResponse{
	order: newOrder,
	OrderedAtTimeStamp: "1613758423",
	Cost:               999,
}

var driver = "mysql"
var dataSourceName = "root:Qwerty123@/task"

func getOrders(w http.ResponseWriter, r *http.Request) {
	orders := orders{
		Orders: []order {newOrder,
		},
	}
	b, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendResponse(w, b, err)
}

func sendResponse(w http.ResponseWriter, b []byte, err interface{}) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, string(b)); err != nil {
		log.WithField("err", err).Error("write responce error")
	}
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["ID"]
	if id == newOrder.Id {
		b, err := json.Marshal(newOrder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendResponse(w, b, err)
	}
}


func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/orders", getOrders)
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID}", getOrder)
	//s.HandleFunc("/order", createOrder).Methods(http.MethodPost)
	return logMiddleware(r)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		endTime := time.Now()
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"time":       int(endTime.Sub(startTime)),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
		endTime = time.Now()
	})
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "hello world")
}
