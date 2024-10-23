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
	ID      uint   `json:"id"`
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
	DB.Create(&Message{Text: body.Message})
}

func GetHandler(w http.ResponseWriter, r *http.Request) {

	messages := []Message{}
	DB.Find(&messages)

	response, _ := json.Marshal(messages)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Bad id", http.StatusBadRequest)
		return
	}

	var message Message
	result := DB.Where("id = ?", id).First(&message)

	if result.Error != nil {
		log.Println("Wrong id:", result.Error)
		http.Error(w, "Not found id", http.StatusNotFound)
		return
	}

	if err := DB.Delete(&message).Error; err != nil {
		log.Println("Error deleting record:", err)
		http.Error(w, "Failed to delete record", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Message deleted")
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		http.Error(w, "Bad id", http.StatusBadRequest)
		return
	}

	var message Message
	result := DB.Where("id = ?", id).First(&message)

	if result.Error != nil {
		log.Println("Wrong id:", result.Error)
		http.Error(w, "Not found id", http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newMessage requestBody
	if err := decoder.Decode(&newMessage); err != nil {
		log.Println("Error decoding request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	message.Text = newMessage.Message
	result = DB.Save(&message)

	if result.Error != nil {
		log.Println("Error update:", result.Error)
		http.Error(w, "Error update", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Message updated")

}

func main() {

	InitDB()
	DB.AutoMigrate(&Message{})

	router := mux.NewRouter()
	router.HandleFunc("/api/post", PostHandler).Methods("POST")
	router.HandleFunc("/api/hello", GetHandler).Methods("GET")
	router.HandleFunc("/api/delete/{id}", DeleteHandler).Methods("DELETE")
	router.HandleFunc("/api/patch/{id}", UpdateHandler).Methods("PATCH")
	http.ListenAndServe(":8080", router)
}
