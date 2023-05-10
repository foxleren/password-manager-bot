package repository

import (
	"fmt"
	"github.com/foxleren/password-manager-bot/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/siruspen/logrus"
)

type SubscriberServicePostgres struct {
	db *sqlx.DB
}

func NewSubscriberServicePostgres(db *sqlx.DB) *SubscriberServicePostgres {
	return &SubscriberServicePostgres{db: db}
}

// TODO: сообщать о наличии сервиса с таким именем!!!
func (p *SubscriberServicePostgres) CreateSubscriberService(subscriberService *models.SubscriberService) (int, error) {
	var id int
	createSubscriberServiceQuery := fmt.Sprintf("INSERT INTO %s (service_name, service_login, service_password) VALUES ($1, $2, $3) RETURNING id", servicesTable)
	row := p.db.QueryRow(createSubscriberServiceQuery,
		subscriberService.ServiceName,
		subscriberService.ServiceLogin,
		subscriberService.ServicePassword)

	if err := row.Scan(&id); err != nil {
		logrus.Printf("Level: repos; func CreateSubscriberService(): error while creating service with name: %s", subscriberService.ServiceName)
		return 0, err
	}

	logrus.Printf("Level: repos; func CreateSubscriberService(): service_id=%d", id)

	return id, nil
}

func (p *SubscriberServicePostgres) CreateSubscriberServiceByName(subscriberID int, subscriberServiceName string) (int, error) {
	tx, err := p.db.Begin()

	var sbsServiceID int
	createSubscriberServiceQuery := fmt.Sprintf("INSERT INTO %s (service_name, service_login, service_password) VALUES ($1, $2, $3) RETURNING id", servicesTable)
	row := tx.QueryRow(createSubscriberServiceQuery,
		subscriberServiceName,
		"",
		"")

	if err = row.Scan(&sbsServiceID); err != nil {
		logrus.Printf("Level: repos; func CreateSubscriberServiceByName(): error while creating service with name: %s", subscriberServiceName)
		return 0, err
	}

	addServiceToSubscriberQuery := fmt.Sprintf("INSERT INTO %s (subscriber_id, service_id) VALUES ($1, $2)", subscribersServicesTable)
	_, err = tx.Exec(addServiceToSubscriberQuery, subscriberID, sbsServiceID)
	if err != nil {
		logrus.Printf("Level: repos; func CreateSubscriberServiceByName(): error while adding service to subscriber with id: %s", subscriberID)
		tx.Rollback()
		return 0, err
	}

	updateSubscriberServiceInProgressQuery := fmt.Sprintf("UPDATE %s SET service_in_progress_id = $1 WHERE id = %d", subscribersTable, subscriberID)
	_, err = tx.Exec(updateSubscriberServiceInProgressQuery, sbsServiceID)

	if err != nil {
		logrus.Printf("Level: repos; func CreateSubscriberServiceByName(): error while updating service_in_progress_id for subscriber with id: %s", subscriberID)
		tx.Rollback()
	}

	logrus.Printf("Level: repos; func CreateSubscriberServiceByName(): service_id=%d", sbsServiceID)

	return sbsServiceID, tx.Commit()
}

func (p *SubscriberServicePostgres) UpdateSubscriberServiceLogin(subscriberID int, subscriberServiceID int, subscriberServiceLogin string) error {
	tx, err := p.db.Begin()
	updateSubscriberServiceLoginQuery := fmt.Sprintf("UPDATE %s SET service_login = $1 WHERE id = %d", servicesTable, subscriberServiceID)
	_, err = tx.Exec(updateSubscriberServiceLoginQuery, subscriberServiceLogin)

	if err != nil {
		logrus.Printf("Level: repos; func UpdateSubscriberServiceLogin(): error while updating service_login for service with id: %s", subscriberServiceID)
		tx.Rollback()
	}

	updateSubscriberServiceInProgressQuery := fmt.Sprintf("UPDATE %s SET service_in_progress_id = $1 WHERE id = %d", subscribersTable, subscriberID)
	_, err = tx.Exec(updateSubscriberServiceInProgressQuery, subscriberServiceID)

	if err != nil {
		logrus.Printf("Level: repos; func UpdateSubscriberServiceLogin(): error while updating service_in_progress_id for subscriber with id: %s", subscriberID)
		tx.Rollback()
	}

	logrus.Printf("Level: repos; func UpdateSubscriberServiceLogin(): service_id=%d", subscriberServiceID)

	return tx.Commit()
}

func (p *SubscriberServicePostgres) UpdateSubscriberServicePassword(subscriberId int, subscriberServiceID int, subscriberServicePassword string) error {
	tx, err := p.db.Begin()
	updateSubscriberServiceLoginQuery := fmt.Sprintf("UPDATE %s SET service_password= $1 WHERE id = %d", servicesTable, subscriberServiceID)
	_, err = tx.Exec(updateSubscriberServiceLoginQuery, subscriberServicePassword)

	if err != nil {
		logrus.Printf("Level: repos; func UpdateSubscriberServicePassword(): error while updating service_login for service with id: %s", subscriberServiceID)
		tx.Rollback()
	}

	updateSubscriberServiceInProgressQuery := fmt.Sprintf("UPDATE %s SET service_in_progress_id = $1 WHERE id = %d", subscribersTable, subscriberId)
	_, err = tx.Exec(updateSubscriberServiceInProgressQuery, subscriberServiceID)

	if err != nil {
		logrus.Printf("Level: repos; func UpdateSubscriberServicePassword(): error while updating service_in_progress_id for subscriber with id: %s", subscriberId)
		tx.Rollback()
	}

	logrus.Printf("Level: repos; func UpdateSubscriberServicePassword(): service_id=%d", subscriberServiceID)

	return tx.Commit()
}

func (p *SubscriberServicePostgres) GetSubscriberServiceByName(userId int, serviceName string) (*models.SubscriberService, error) {
	var service models.SubscriberService
	query := fmt.Sprintf(
		`SELECT service_id, service_name, service_login, service_password
					FROM (
					(%s subs JOIN %s ss on subs.id = ss.subscriber_id) as subs_and_serv_id
					JOIN %s s on subs_and_serv_id.service_id = s.id
					) WHERE service_name = $1
		`, subscribersTable, subscribersServicesTable, servicesTable)
	err := p.db.Get(&service, query, serviceName)

	if err != nil {
		logrus.Printf("repo: GetSubscriber(): subscriber with chat_id: %v does not exist", userId)
	}

	return &service, err
}