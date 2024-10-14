package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var message string

type requestBody struct {
	Message string `json:"message"`
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&body); err != nil {
		log.Println("Error decoding request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message = body.Message
	fmt.Fprintf(w, "Successfully received message: %s\n", message)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!\n", message)
}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/api/post", PostHandler).Methods("POST")
	router.HandleFunc("/api/hello", HelloHandler).Methods("GET")
	http.ListenAndServe(":8080", router)
}
