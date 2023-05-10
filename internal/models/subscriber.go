package models

type Subscriber struct {
	ID                  int    `json:"id" db:"id"`
	ChatId              int64  `json:"chat_id" db:"chat_id" binding:"required"`
	DialogStatus        string `json:"dialog_status" db:"dialog_status" binding:"required"`
	ServiceInProgressID int    `json:"service_in_progress_id" db:"service_in_progress_id"`
}

const (
	DialogStatusNone                      = "none"
	DialogStatusWaitingForServiceName     = "wait service_name"
	DialogStatusWaitingForServiceLogin    = "wait service_login"
	DialogStatusWaitingForServicePassword = "wait service_password"
)

type SubscriberService struct {
	ID              int    `json:"id" db:"id"`
	UserID          int    `json:"user_id" db:"user_id"`
	ServiceName     string `json:"service_name" db:"service_name"`
	ServiceLogin    string `json:"service_login" db:"service_login"`
	ServicePassword string `json:"service_password" db:"service_password"`
}
