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
		logrus.Printf("Level: repos; func CreateSubscriber(): err=%v", err.Error())
		return 0, err
	}

	logrus.Printf("Level: repos; func CreateSubscriber(): created user with id=%d", id)

	return id, nil
}

func (p *SubscriberPostgres) GetAllSubscribers() ([]models.Subscriber, error) {
	logrus.Printf("%v", p.db != nil)

	var subscribers []models.Subscriber
	getAllSubscribersQuery := fmt.Sprintf("SELECT id, chat_id FROM %s", subscribersTable)
	err := p.db.Select(&subscribers, getAllSubscribersQuery)

	if err != nil {
		logrus.Printf("Level: repos; func GetAllSubscribers(): err=%v", err.Error())
	}

	logrus.Printf("Level: repos; func GetAllSubscribers(): status=success")

	return subscribers, err
}

func (p *SubscriberPostgres) GetSubscriber(chatId int64) (models.Subscriber, error) {
	var subscriber models.Subscriber
	getSubscriberQuery := fmt.Sprintf("SELECT id, chat_id, dialog_status, service_in_progress_id FROM %s WHERE chat_id = $1", subscribersTable)
	err := p.db.Get(&subscriber, getSubscriberQuery, chatId)

	if err != nil {
		logrus.Printf("Level: repos; func GetSubscriber(): subscriber with chat_id: %v does not exist", chatId)
	}

	logrus.Printf("Level: repos; func GetSubscriber(): subscriber=%v", subscriber)

	return subscriber, err
}

func (p *SubscriberPostgres) DeleteSubscriber(chatId int64) error {
	deleteSubscriberQuery := fmt.Sprintf("DELETE FROM %s WHERE chat_id = %d", subscribersTable, chatId)
	_, err := p.db.Exec(deleteSubscriberQuery)

	if err != nil {
		logrus.Printf("Level: repos; func DeleteSubscriber(): err=%v", err.Error())
	}

	logrus.Printf("Level: repos; func DeleteSubscriber(): status=success")

	return err
}

func (p *SubscriberPostgres) UpdateSubscriberDialogStatus(chatId int64, dialogStatus string) error {
	updateSubscriberStatusQuery := fmt.Sprintf("UPDATE %s SET dialog_status = $1 WHERE chat_id = %d", subscribersTable, chatId)
	_, err := p.db.Exec(updateSubscriberStatusQuery, dialogStatus)

	if err != nil {
		logrus.Printf("Level: repos; func UpdateSubscriberDialogStatus(): err=%v", err.Error())
	}

	logrus.Printf("Level: repos; func UpdateSubscriberDialogStatus(): status=success")

	return err
}

func (p *SubscriberPostgres) UpdateSubscriberServiceInProgressID(chatId int64, serviceID string) error {
	updateSubscriberStatusQuery := fmt.Sprintf("UPDATE %s SET service_in_progress_id = $1 WHERE chat_id = %d", subscribersTable, chatId)
	_, err := p.db.Exec(updateSubscriberStatusQuery, serviceID)

	if err != nil {
		logrus.Printf("Level: repos; func UpdateSubscriberServiceInProgressID(): err=%v", err.Error())
	}

	logrus.Printf("Level: repos; func UpdateSubscriberServiceInProgressID(): status=success")

	return err
}
