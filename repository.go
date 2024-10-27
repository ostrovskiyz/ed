package messagesService

import "gorm.io/gorm"

type MessageRepository interface {
	CreateMessage(message Message) (Message, error)
	GetAllMessages() ([]Message, error)
	UpdateMessageByID(id int, message Message) (Message, error)
	DeleteMessageByID(id int) error
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *messageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) CreateMessage(message Message) (Message, error) {
	result := r.db.Create(&message)
	if result.Error != nil {
		return Message{}, result.Error
	}
	return message, nil
}

func (r *messageRepository) GetAllMessages() ([]Message, error) {
	var messages []Message
	err := r.db.Find(&messages).Error
	return messages, err
}

func (r *messageRepository) UpdateMessageByID(id int, message Message) (Message, error) {
	var oldMessage Message
	err := r.db.First(&oldMessage, id).Error
	if err != nil {
		return Message{}, err
	}

	oldMessage.Text = message.Text

	if err := r.db.Save(&oldMessage).Error; err != nil {
		return Message{}, err
	}

	return oldMessage, nil
}

func (r *messageRepository) DeleteMessageByID(id int) error {
	var message Message
	err := r.db.First(&message, id).Error
	if err != nil {
		return err
	}

	err = r.db.Delete(&message).Error
	return err
}
