package handlers

import (
	"projectgo/internal/messagesService"
	"projectgo/internal/web/messages"
	"time"

	"github.com/labstack/echo/v4"
)

// Handler структура обработчика, который будет использовать сервис для выполнения операций с сообщениями
type Handler struct {
	Service *messagesService.MessageService
}

// GetMessages реализует обработчик для получения всех сообщений
func (h *Handler) GetMessages(ctx echo.Context) error {
	allMessages, err := h.Service.GetAllMessages()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	var response []messages.Message
	for _, msg := range allMessages {
		// Проверяем, если ID пустое, то передаем nil
		id := int(msg.ID)
		message := messages.Message{
			Id:        &id, // Указываем указатель на id
			Text:      &msg.Text,
			CreatedAt: &msg.CreatedAt, // Указываем указатель на время
			UpdatedAt: &msg.UpdatedAt, // Указываем указатель на время
		}

		// Для DeletedAt проверяем, что оно не нулевое
		if msg.DeletedAt.Time != (time.Time{}) {
			message.DeletedAt = &msg.DeletedAt.Time // Указываем указатель на время
		}

		response = append(response, message)
	}

	return ctx.JSON(200, response)
}

// PostMessage реализует обработчик для создания нового сообщения
func (h *Handler) PostMessage(ctx echo.Context) error {
	var body messages.MessageInput
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}

	// Создаем сообщение с полученным текстом
	messageToCreate := messagesService.Message{
		Text: body.Text,
	}

	// Создаем сообщение через сервис
	createdMessage, err := h.Service.CreateMessage(messageToCreate)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	// Формируем ответ
	id := int(createdMessage.ID)
	response := messages.Message{
		Id:        &id, // Указываем указатель на id
		Text:      &createdMessage.Text,
		CreatedAt: &createdMessage.CreatedAt, // Указываем указатель на время
		UpdatedAt: &createdMessage.UpdatedAt, // Указываем указатель на время
	}

	// Для DeletedAt проверяем, что оно не нулевое
	if createdMessage.DeletedAt.Time != (time.Time{}) {
		response.DeletedAt = &createdMessage.DeletedAt.Time // Указываем указатель на время
	}

	return ctx.JSON(201, response)
}

// PatchMessage реализует обработчик для обновления сообщения
func (h *Handler) PatchMessage(ctx echo.Context, id int) error {
	var body messages.MessageInput
	if err := ctx.Bind(&body); err != nil {
		return echo.NewHTTPError(400, "Invalid request body")
	}

	// Обновляем сообщение по ID
	updatedMessage, err := h.Service.UpdateMessageByID(id, messagesService.Message{
		Text: body.Text,
	})
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	// Формируем ответ
	id = int(updatedMessage.ID)
	response := messages.Message{
		Id:        &id, // Указываем указатель на id
		Text:      &updatedMessage.Text,
		CreatedAt: &updatedMessage.CreatedAt, // Указываем указатель на время
		UpdatedAt: &updatedMessage.UpdatedAt, // Указываем указатель на время
	}

	// Для DeletedAt проверяем, что оно не нулевое
	if updatedMessage.DeletedAt.Time != (time.Time{}) {
		response.DeletedAt = &updatedMessage.DeletedAt.Time // Указываем указатель на время
	}

	return ctx.JSON(200, response)
}

// DeleteMessage реализует обработчик для удаления сообщения
func (h *Handler) DeleteMessage(ctx echo.Context, id int) error {
	// Удаляем сообщение по ID
	err := h.Service.DeleteMessageByID(id)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return ctx.NoContent(204) // Возвращаем статус 204 No Content
}

func NewHandler(service *messagesService.MessageService) *Handler {
	return &Handler{
		Service: service,
	}
}
