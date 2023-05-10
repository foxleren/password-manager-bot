package models

import "time"

type Message struct {
	Id          int       `json:"id" db:"id"`
	MessageId   int       `json:"message_id" db:"message_id"`
	ChatId      int64     `json:"chat_id" db:"chat_id"`
	MessageDate time.Time `json:"message_date" db:"message_date"`
}
