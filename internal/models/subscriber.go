package models

import (
	"fmt"
	"strings"
)

type Subscriber struct {
	ID                  int    `json:"id" db:"id"`
	ChatId              int64  `json:"chat_id" db:"chat_id" binding:"required"`
	DialogStatus        string `json:"dialog_status" db:"dialog_status" binding:"required"`
	ServiceInProgressID int    `json:"service_in_progress_id" db:"service_in_progress_id"`
}

const (
	DialogStatusNone               = "none"
	DialogStatusSetServiceName     = "set service_name"
	DialogStatusSetServiceLogin    = "set service_login"
	DialogStatusSetServicePassword = "set service_password"

	DialogStatusGetServiceName = "get service_name"
)

type SubscriberService struct {
	ID              int    `json:"id" db:"id"`
	UserID          int    `json:"user_id" db:"user_id"`
	ServiceName     string `json:"service_name" db:"service_name"`
	ServiceLogin    string `json:"service_login" db:"service_login"`
	ServicePassword string `json:"service_password" db:"service_password"`
}

type SubscriberServiceOutput struct {
	ServiceName     string `json:"service_name" db:"service_name"`
	ServiceLogin    string `json:"service_login" db:"service_login"`
	ServicePassword string `json:"service_password" db:"service_password"`
}

func (s SubscriberServiceOutput) String() string {
	return fmt.Sprintf("Название сервиса: `%s`\nЛогин: `%s`\nПароль: `%s`\n",
		s.ServiceName,
		s.ServiceLogin,
		s.ServicePassword)
}

func FormatAllSubscriberServiceOutput(slice []SubscriberServiceOutput) string {
	var sb = strings.Builder{}
	for i := 0; i < len(slice); i++ {
		sb.WriteString(slice[i].String())
	}
	return sb.String()
}
