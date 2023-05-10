package repository

import (
	"fmt"
	"github.com/foxleren/password-manager-bot/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/siruspen/logrus"
	"time"
)

type MessagePostgres struct {
	db *sqlx.DB
}

func NewMessagePostgres(db *sqlx.DB) *MessagePostgres {
	return &MessagePostgres{db: db}
}

func (p *MessagePostgres) CreateMessage(chatId int64, messageId int) (int, error) {
	var id int
	addMessageInfoQuery := fmt.Sprintf("INSERT INTO %s (message_id, chat_id, message_date) VALUES ($1, $2, $3) RETURNING id", messagesTable)
	row := p.db.QueryRow(addMessageInfoQuery,
		messageId,
		chatId,
		time.Now(),
	)

	if err := row.Scan(&id); err != nil {
		logrus.Printf("Level: repos; func CreateMessage(): error while creating message with id: %d", messageId)
		return 0, err
	}

	return id, nil
}

func (p *MessagePostgres) GetAllOutdatedMessages(lastFreshDate time.Time) ([]models.Message, error) {
	var messages []models.Message
	getAllOutdatedMessagesQuery := fmt.Sprintf("SELECT * FROM %s WHERE message_date <= $1", messagesTable)
	err := p.db.Select(&messages, getAllOutdatedMessagesQuery, lastFreshDate)

	if err != nil {
		logrus.Printf("Level: repos; GetAllOutdatedMessages(): err=%v", err.Error())
		return nil, err
	}

	deleteOutdatedMessagesQuery := fmt.Sprintf("DELETE FROM %s WHERE message_date <= $1", messagesTable)
	_, err = p.db.Exec(deleteOutdatedMessagesQuery, lastFreshDate)

	if err != nil {
		logrus.Printf("Level: repos; GetAllOutdatedMessages(): err=%v", err.Error())
		return nil, err
	}

	logrus.Printf("Level: repos; func GetAllOutdatedMessages(): status=success")

	return messages, err
}
