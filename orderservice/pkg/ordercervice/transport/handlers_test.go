package transport

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestList (t *testing.T) {
	w := httptest.NewRecorder()
	getOrders(w, nil)
	response := w.Result()
	if response.StatusCode != http.StatusOK {
		t.Errorf("Status code is wrong. Hame %d, want %d.", response.StatusCode, http.StatusOK)
	}

	jsonString, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	//items := make([]OrderListItem, 10)
	var orders orders
	if err = json.Unmarshal(jsonString, &orders); err != nil {
		t.Errorf("Can't parse json response with error %v", err)
	}
}
