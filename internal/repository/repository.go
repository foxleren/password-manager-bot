package repository

import (
	"github.com/foxleren/password-manager-bot/internal/models"
	"github.com/jmoiron/sqlx"
	"time"
)

type Subscriber interface {
	CreateSubscriber(subscriber models.Subscriber) (int, error)
	GetAllSubscribers() ([]models.Subscriber, error)
	GetSubscriber(chatId int64) (models.Subscriber, error)
	DeleteSubscriber(chatId int64) error
	UpdateSubscriberDialogStatus(chatId int64, dialogStatus string) error
	UpdateSubscriberServiceInProgressID(chatId int64, serviceID string) error
}

type SubscriberService interface {
	CreateSubscriberServiceByName(subscriberID int, subscriberServiceName string) (int, error)
	UpdateSubscriberServiceLogin(subscriberId int, subscriberServiceID int, subscriberServiceLogin string) error
	UpdateSubscriberServicePassword(subscriberId int, subscriberServiceID int, subscriberServicePassword string) error
	GetSubscriberServiceByName(chatId int64, serviceName string) (*models.SubscriberServiceOutput, error)
	DeleteSubscriberService(chatId int64, serviceName string) error
}

type Message interface {
	CreateMessage(chatId int64, messageId int) (int, error)
	GetAllOutdatedMessages(lastFreshDate time.Time) ([]models.Message, error)
}

type Repository struct {
	Subscriber
	SubscriberService
	Message
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Subscriber:        NewSubscriberPostgres(db),
		SubscriberService: NewSubscriberServicePostgres(db),
		Message:           NewMessagePostgres(db),
	}
}
