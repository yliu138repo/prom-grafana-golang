package httpserver

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func getOrderHandler(w http.ResponseWriter, r *http.Request) {

	resp := make(map[string]string)
	resp["message"] = "The new Order at time " + time.Now().Format("2006-01-02T15:04:05Z07:00")

	json, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("error handling json marshal req")
	}

	w.Write(json)

}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {

	resp := make(map[string]string)
	resp["message"] = "Its not here what you looking for"

	json, err := json.Marshal(resp)
	if err != nil {
		log.Fatal("error handling json marshal req")
	}

	w.Write(json)

}
