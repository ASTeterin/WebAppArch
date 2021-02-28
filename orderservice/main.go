package main

import (
	"fmt"
	"net/http"
	"orderservice/transport"
)



func HelloHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(w, "hello")
}

func main() {
	router := transport.Router()
	fmt.Println(http.ListenAndServe(":8000", router))
	/*http.HandleFunc("/hello-world", HelloHandler)
	http.HandleFunc("/api/v1/orders", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, orders)
	})
	http.HandleFunc("/api/v1/order/d290f1ee-6c54-4b01-90E6-d701748fo851", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, order)
	})
	//http.HandlerFunc("/hello-world", HelloHandler)
	http.ListenAndServe(":8000", nil)*/
}
