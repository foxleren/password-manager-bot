package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errUnableToSubscribe   = errors.New("error: failed while subscribing")
	errUnableToUnsubscribe = errors.New("error: failed while unsubscribing")
	errDisabledCommand     = errors.New("error: command is disabled")
)

const (
	unableToSubscribeText   = "Не удалось подписаться на рассылку либо вы уже подписаны. Для проверки подписки используйте команду \n/check_subscription."
	unableToUnsubscribeText = "Не удалось отписаться от рассылки либо вы уже отписаны. Для проверки подписки используйте команду \n/check_subscription."
	unknownError            = "Произошла неизвестная ошибка."
	disabledCommandReply    = "Подпишитесь для использования этой функции!"
)

const (
	ReplySendServiceName      = "Пришлите название сервиса"
	ReplySendServiceLogin     = "Пришлите логин сервиса"
	ReplySendServicePassword  = "Пришлите пароль сервиса"
	ReplyFinishSettingService = "Функция не доступна! Необходимо завершить добавление сервиса!"
)

const (
	ReplyServiceNameIsSet     = "Название сервиса задано! Пришлите логин сервиса."
	ReplyServiceLoginIsSet    = "Логин для сервиса задан! Пришлите пароль сервиса."
	ReplyServicePasswordIsSet = "Пароль для сервиса задан! Сервис успешно добавлен!"
)

func (b *Bot) handleError(chatID int64, err error) {
	msg := tgbotapi.NewMessage(chatID, "")
	switch err {
	case errUnableToSubscribe:
		msg.Text = unableToSubscribeText
	case errUnableToUnsubscribe:
		msg.Text = unableToUnsubscribeText
	case errDisabledCommand:
		msg.Text = disabledCommandReply
	default:
		msg.Text = unknownError
	}
	b.bot.Send(msg)
}

func (b *Bot) sendReply(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)

	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
