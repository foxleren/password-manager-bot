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
		logrus.Printf("Level: repos; func CreateSubscriberServiceByName(): error while adding service to subscriber with id: %s, err=%v", subscriberID, err.Error())
		tx.Rollback()
		return 0, err
	}

	addServiceNameToSubscriberQuery := fmt.Sprintf("INSERT INTO %s (subscriber_id, service_id, service_name) VALUES ($1, $2, $3)", subscribersServiceNamesTable)
	_, err = tx.Exec(addServiceNameToSubscriberQuery, subscriberID, sbsServiceID, subscriberServiceName)
	if err != nil {
		logrus.Printf("Level: repos; func CreateSubscriberServiceByName(): error while adding service_name to subscriber with id: %s, err=%v", subscriberID, err.Error())
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

	updateSubscriberStatusQuery := fmt.Sprintf("UPDATE %s SET service_in_progress_id = $1 WHERE id = %d", subscribersTable, subscriberId)
	_, err = tx.Exec(updateSubscriberStatusQuery, 0)

	if err != nil {
		logrus.Printf("Level: repos; func UpdateSubscriberServicePassword(): error while updating service_in_progress_id for subscriber with default id: %s", 0)
		tx.Rollback()
	}

	if err != nil {
		logrus.Printf("repo: UpdateSubscriberServiceInProgressID(): %v", err.Error())
	}

	logrus.Printf("Level: repos; func UpdateSubscriberServicePassword(): service_id=%d", subscriberServiceID)

	return tx.Commit()
}

func (p *SubscriberServicePostgres) GetSubscriberServiceByName(chatId int64, serviceName string) (*models.SubscriberServiceOutput, error) {
	var services models.SubscriberServiceOutput
	query := fmt.Sprintf(
		`SELECT service_name, service_login, service_password
					FROM (
					(%s subs JOIN %s ss on subs.id = ss.subscriber_id) as subs_and_serv_id
					JOIN %s s on subs_and_serv_id.service_id = s.id
					) WHERE service_name = $1 AND chat_id = $2
		`, subscribersTable, subscribersServicesTable, servicesTable)
	err := p.db.Get(&services, query, serviceName, chatId)

	if err != nil {
		logrus.Printf("Level: repos; func GetSubscriberServiceByName(): error while getting all services (err=%v)", err.Error())
	}

	logrus.Printf("Level: repos; func GetSubscriberServiceByName(): service=%v", services)

	return &services, err
}

func (p *SubscriberServicePostgres) DeleteSubscriberService(chatId int64, serviceName string) error {
	_, err := p.GetSubscriberServiceByName(chatId, serviceName)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(
		`DELETE FROM %s WHERE id =
				(  
    				SELECT service_id FROM 
    				(
        			(%s subs JOIN %s ss on subs.id = ss.subscriber_id) as subs_and_serv_id
        			JOIN %s s on subs_and_serv_id.service_id = s.id
    				)   WHERE s.service_name = $1 AND chat_id = $2
				)
		`, servicesTable, subscribersTable, subscribersServicesTable, servicesTable)
	_, err = p.db.Exec(query, serviceName, chatId)

	if err != nil {
		logrus.Printf("repo: DeleteSubscriber(): %v", err.Error())
	}

	return err
}
