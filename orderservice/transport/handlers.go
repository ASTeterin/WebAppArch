package transport

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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

func Router() *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/hello-world", helloWorld).Methods(http.MethodGet)
	s.HandleFunc("/orders", getOrders)

	return r
}

func getOrders(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, orders)
}

func helloWorld(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "hello world")
}
