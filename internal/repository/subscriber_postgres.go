package repository

import (
	"fmt"
	"github.com/foxleren/password-manager-bot/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/siruspen/logrus"
)

type SubscriberPostgres struct {
	db *sqlx.DB
}

func NewSubscriberPostgres(db *sqlx.DB) *SubscriberPostgres {
	return &SubscriberPostgres{db: db}
}

func (p *SubscriberPostgres) CreateSubscriber(subscriber models.Subscriber) (int, error) {
	var id int
	createSubscriberQuery := fmt.Sprintf("INSERT INTO %s (chat_id, dialog_status, service_in_progress_id) VALUES ($1, $2, $3) RETURNING id", subscribersTable)
	row := p.db.QueryRow(createSubscriberQuery,
		subscriber.ChatId,
		models.DialogStatusNone,
		0)

	if err := row.Scan(&id); err != nil {
		logrus.Printf("repo: CreateSubscriber(): %v", err.Error())
		return 0, err
	}

	return id, nil
}

func (p *SubscriberPostgres) GetAllSubscribers() ([]models.Subscriber, error) {
	logrus.Printf("%v", p.db != nil)

	var subscribers []models.Subscriber
	getAllQuery := fmt.Sprintf("SELECT id, chat_id FROM %s", subscribersTable)
	err := p.db.Select(&subscribers, getAllQuery)

	if err != nil {
		logrus.Printf("repo: GetAllSubscribers(): %v", err.Error())
	}

	return subscribers, err
}

func (p *SubscriberPostgres) GetSubscriber(chatId int64) (models.Subscriber, error) {
	var subscriber models.Subscriber
	query := fmt.Sprintf("SELECT id, chat_id, dialog_status, service_in_progress_id FROM %s WHERE chat_id = $1", subscribersTable)
	err := p.db.Get(&subscriber, query, chatId)

	if err != nil {
		logrus.Printf("repo: GetSubscriber(): subscriber with chat_id: %v does not exist", chatId)
	}

	return subscriber, err
}

//func (p *SubscriberPostgres) GetSubscriberServiceInProgress(chatId int64) (models.Subscriber, error) {
//	var subscriber models.Subscriber
//	query := fmt.Sprintf("SELECT id, chat_id, dialog_status FROM %s WHERE chat_id = $1", subscribersTable)
//	err := p.db.Get(&subscriber, query, chatId)
//
//	if err != nil {
//		logrus.Printf("repo: GetSubscriber(): subscriber with chat_id: %v does not exist", chatId)
//	}
//
//	return subscriber, err
//}

func (p *SubscriberPostgres) DeleteSubscriber(chatId int64) error {
	deleteCartItemByIDQuery := fmt.Sprintf("DELETE FROM %s WHERE chat_id = %d", subscribersTable, chatId)
	_, err := p.db.Exec(deleteCartItemByIDQuery)

	if err != nil {
		logrus.Printf("repo: DeleteSubscriber(): %v", err.Error())
	}

	return err
}

func (p *SubscriberPostgres) UpdateSubscriberDialogStatus(chatId int64, dialogStatus string) error {
	updateSubscriberStatusQuery := fmt.Sprintf("UPDATE %s SET dialog_status = $1 WHERE chat_id = %d", subscribersTable, chatId)
	_, err := p.db.Exec(updateSubscriberStatusQuery, dialogStatus)

	if err != nil {
		logrus.Printf("repo: UpdateSubscriberDialogStatus(): %v", err.Error())
	}

	return err
}

func (p *SubscriberPostgres) UpdateSubscriberServiceInProgressID(chatId int64, serviceID string) error {
	updateSubscriberStatusQuery := fmt.Sprintf("UPDATE %s SET service_in_progress_id = $1 WHERE chat_id = %d", subscribersTable, chatId)
	_, err := p.db.Exec(updateSubscriberStatusQuery, serviceID)

	if err != nil {
		logrus.Printf("repo: UpdateSubscriberServiceInProgressID(): %v", err.Error())
	}

	return err
}
