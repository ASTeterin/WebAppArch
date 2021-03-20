package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	model2 "orderservice/pkg/orderservice/model"
	"time"
	)


type order struct {
	Id    string `json:"id"`
	menuItems []menuItem `json:"menuItems"`
}

type orderResponse struct {
	order
	orderedAtTimeStamp string `json:"orderedAtTimeStamp"`
	Cost               int    `json:"cost"`
}

type menuItem struct {
	Id       string `json:"id"`
	Quantity int    `json:"quantity"`
}



type orders struct {
	orders []order
}

var menuitem = menuItem{
	Id:       "f290d1ce-6c234-4b31-90e6-d701748fo851",
	Quantity: 1,
}

var newOrder = order{
	Id:  "d290f1ee-6c54-4b01-90E6-d701748fo851",
	menuItems: []menuItem{
		menuitem,
	},
}

var driver = "mysql"
var dataSourceName = "root:Qwerty123@/order"


func getOrders(w http.ResponseWriter, r *http.Request) {
	orders := orders{
		orders: []order {newOrder,
			},
	}
	b, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendResponse(w, b, http.StatusOK)
}

/*func getOrders(w http.ResponseWriter, r *http.Request) {
	s := createDBConnection()
	orders := s.GetOrders()
	b, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendResponse(w, b, http.StatusOK)
	defer s.Database.Close()
}
*/


func sendResponse(w http.ResponseWriter, b []byte, status int) {
	w.Header().Set("Content-type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	if _, err := io.WriteString(w, string(b)); err != nil {
		log.WithField("err", err).Error("write responce error")
	}
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["ID"] == "" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	id := vars["ID"]

	if id == newOrder.Id {
		b, err := json.Marshal(newOrder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sendResponse(w, b, http.StatusOK)
	} else {
		sendResponse(w, []byte("Not found"), http.StatusNotFound)
	}
}

func createDBConnection() model2.Server {
	db, err := sql.Open(driver, dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return model2.Server{Database: db}
}

func  createOrder(w http.ResponseWriter, r *http.Request) {
	s := createDBConnection()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var msg model2.MenuItems
	err = json.Unmarshal(b, &msg)
	fmt.Println(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	guid := uuid.New().String()
	timestamp := int(time.Now().Unix())
	cost := 1000
	s.CreateOrder(guid, timestamp, cost, msg.MenuItems)
	defer s.Database.Close()
}

func deleteOrder(w http.ResponseWriter, r *http.Request) {
	s := createDBConnection()
	vars := mux.Vars(r)
	id := vars["ID"]
	s.DeleteOrder(id)
	defer s.Database.Close()
}



func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/orders", getOrders).Methods(http.MethodGet)
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/order/{ID}", getOrder).Methods(http.MethodGet)
	s.HandleFunc("/order", createOrder).Methods(http.MethodPost)
	s.HandleFunc("/order/{ID}", deleteOrder).Methods(http.MethodDelete)
	return logMiddleware(r)
}

func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"duration":    time.Since(startTime).String(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "hello world")
}
