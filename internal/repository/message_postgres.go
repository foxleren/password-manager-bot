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
		logrus.Printf("Level: repos; func UpdateSubscriberServicePassword(): error while creating service with name: %s")
		return 0, err
	}

	return id, nil
}

func (p *MessagePostgres) GetAllOutdatedMessages(messageTTLInMinutes int) ([]models.Message, error) {
	var messages []models.Message
	getAllQuery := fmt.Sprintf("SELECT * FROM %s WHERE message_date <= NOW() - INTERVAL '%d minutes'", messagesTable, messageTTLInMinutes)
	err := p.db.Select(&messages, getAllQuery)

	if err != nil {
		logrus.Printf("repo: GetAllOutdatedMessages(): err=%v", err.Error())
		return nil, err
	}

	deleteCartItemByIDQuery := fmt.Sprintf("DELETE FROM %s WHERE message_date <= NOW() - INTERVAL '%d minutes'", messagesTable, messageTTLInMinutes)
	_, err = p.db.Exec(deleteCartItemByIDQuery)

	if err != nil {
		logrus.Printf("repo: GetAllOutdatedMessages(): err=%v", err.Error())
		return nil, err
	}

	logrus.Printf("GetAllOutdatedMessages(): []=", messages)

	return messages, err
}
