package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"projectgo/internal/database"
	"projectgo/internal/handlers"
	"projectgo/internal/messagesService"
	"projectgo/internal/web/messages"
)

func main() {
	database.InitDB()

	repo := messagesService.NewMessageRepository(database.DB)
	service := messagesService.NewService(repo)

	handler := handlers.NewHandler(service)

	// Инициализируем echo
	e := echo.New()

	// используем Logger и Recover
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	messages.RegisterHandlers(e, handler)
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start with err: %v", err)
	}
}
