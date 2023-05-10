package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	errUnableToSubscribe     = errors.New("error: failed while subscribing")
	errUnableToUnsubscribe   = errors.New("error: failed while unsubscribing")
	errDisabledCommand       = errors.New("error: command is disabled")
	errUnableToSetService    = errors.New("error: failed while setting service")
	errUnableToDeleteService = errors.New("error: failed while deleting service")
	errUnableToGetService    = errors.New("error: failed while getting service")
)

const (
	unableToSubscribeText        = "Не удалось подписаться на рассылку либо вы уже подписаны. Для проверки подписки используйте команду \n/check_subscription."
	unableToUnsubscribeText      = "Не удалось отписаться от рассылки либо вы уже отписаны. Для проверки подписки используйте команду \n/check_subscription."
	unknownErrorText             = "Произошла неизвестная ошибка."
	disabledCommandReply         = "Подпишитесь для использования этой функции!"
	errUnableToSetServiceText    = "Не удалось добавить сервис. Возможно сервис с таким именем уже существует."
	errUnableToGetServiceText    = "Не удалось найти сервис с таким именем. Возможно он уже удален или не существует."
	errUnableToDeleteServiceText = "Не удалось удалить сервис с таким именнем. Возможно он уже удален или не существует."
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
	case errUnableToSetService:
		msg.Text = errUnableToSetServiceText
	case errUnableToGetService:
		msg.Text = errUnableToGetServiceText
	case errUnableToDeleteService:
		msg.Text = errUnableToDeleteServiceText
	default:
		msg.Text = unknownErrorText
	}
	b.bot.Send(msg)
}
