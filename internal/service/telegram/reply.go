package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	replyStart               = "Добро пожаловать! Для использования сервиса подпишитесь командой /subscribe.\n"
	replyUnknownCommand      = "Неизвестная команда."
	successfulSubscription   = "Вы успешно подписались!"
	successfulUnsubscription = "Вы успешно отписались."
	subscriptionStatusGood   = "Статус подписки: активирована."
	subscriptionStatusBad    = "Статус подписки: деактивирована."
)

const (
	ReplySendServiceNameToSet    = "Пришлите название сервиса, который хотите сохранить."
	ReplySendServiceNameToGet    = "Пришлите название сервиса, который хотите получить."
	ReplySendServiceNameToDelete = "Пришлите название сервиса, который хотите удалить."
	ReplySendServiceLogin        = "Пришлите логин сервиса"
	ReplySendServicePassword     = "Пришлите пароль сервиса"
	ReplyFinishSettingService    = "Функция не доступна! Необходимо завершить добавление сервиса!"
	ReplyServiceIsDeleted        = "Сервис успешно удален!"
)

const (
	ReplyServiceNameIsSet     = "Название сервиса задано! Пришлите логин сервиса."
	ReplyServiceLoginIsSet    = "Логин для сервиса задан! Пришлите пароль сервиса."
	ReplyServicePasswordIsSet = "Пароль для сервиса задан! Сервис успешно добавлен!"
)

func (b *Bot) sendReply(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)

	_, err := b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}
