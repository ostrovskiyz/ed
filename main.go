package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"projectgo/internal/database"
	"projectgo/internal/handlers"
	"projectgo/internal/messagesService"
)

func main() {
	database.InitDB()
	database.DB.AutoMigrate(&messagesService.Message{})

	repo := messagesService.NewMessageRepository(database.DB)
	service := messagesService.NewService(repo)

	handler := handlers.NewHandler(service)
	router := mux.NewRouter()

	router.HandleFunc("/api/get", handler.GetMessagesHandler).Methods("GET")
	router.HandleFunc("/api/post", handler.PostMessagesHandler).Methods("POST")
	router.HandleFunc("/api/delete/{id}", handler.DeleteMessagesHandler).Methods("DELETE")
	router.HandleFunc("/api/patch/{id}", handler.PatchMessagesHandler).Methods("PATCH")
	http.ListenAndServe(":8080", router)
}
