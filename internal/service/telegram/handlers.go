package telegram

import (
	"fmt"
	"github.com/foxleren/password-manager-bot/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/siruspen/logrus"
)

const (
	commandStart          = "start"
	commandSet            = "set"
	commandGet            = "get"
	commandDel            = "del"
	commandSubscribe      = "subscribe"
	commandCheckSubscribe = "check_subscription"
	commandUnsubscribe    = "unsubscribe"
	//commandGetData        = "get_data"

	//subscriptionErrorReply = "Подпишитесь для использования этой функции!"

	replyStart          = "Добро пожаловать!\n"
	replyUnknownCommand = "Неизвестная команда."

	successfulSet = "Вы успешно сохранили логин и пароль для сервиса: %s"

	successfulSubscription   = "Вы успешно подписались на рассылку!"
	successfulUnsubscription = "Вы успешно отписались от рассылки."
	subscriptionStatusGood   = "Статус подписки: активирована."
	subscriptionStatusBad    = "Статус подписки: деактивирована."
)

func (b *Bot) handleCommand(update *tgbotapi.Update) error {
	message := update.Message
	switch update.Message.Command() {
	case commandStart:
		return b.handleCommandStart(message)
	case commandSubscribe:
		return b.handleCommandSubscribe(message)
	case commandUnsubscribe:
		return b.handleCommandUnsubscribe(message)
	case commandCheckSubscribe:
		return b.handleCommandCheckSubscribe(message)
	case commandSet:
		return b.handleCommandSet(update)
	case commandGet:
		return b.handleCommandGet(update)
	case commandDel:
		return b.handleCommandDel(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(update *tgbotapi.Update) error {
	subscriber, err := b.repo.GetSubscriber(update.Message.Chat.ID)
	if err != nil {
		logrus.Println("handleMessage(): error with subscriber_id %v", subscriber.ID)
		return errDisabledCommand
	}

	logrus.Println("handleMessage(): got subscriber %v", subscriber)

	switch subscriber.DialogStatus {
	case models.DialogStatusNone:
		{

			err := b.sendReply(update.Message.Chat.ID, ReplySendServiceName)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusGetServiceName)
			if err != nil {
				logrus.Println("handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusSetServiceName:
		{
			_, err = b.repo.SubscriberService.CreateSubscriberServiceByName(subscriber.ID, update.Message.Text)
			if err != nil {
				logrus.Printf("error in handler")
				return err
			}
			err := b.sendReply(update.Message.Chat.ID, ReplyServiceNameIsSet)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusSetServiceLogin)
			if err != nil {
				logrus.Println("handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusSetServiceLogin:
		{
			err = b.repo.SubscriberService.UpdateSubscriberServiceLogin(subscriber.ID, subscriber.ServiceInProgressID, update.Message.Text)
			if err != nil {
				logrus.Printf("error in handler")
				return err
			}
			err := b.sendReply(update.Message.Chat.ID, ReplyServiceLoginIsSet)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusSetServicePassword)
			if err != nil {
				logrus.Println("handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusSetServicePassword:
		{
			err = b.repo.SubscriberService.UpdateSubscriberServicePassword(subscriber.ID, subscriber.ServiceInProgressID, update.Message.Text)
			if err != nil {
				logrus.Printf("error in handler")
				return err
			}
			err := b.sendReply(update.Message.Chat.ID, ReplyServicePasswordIsSet)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusNone)
			if err != nil {
				logrus.Println("handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusGetServiceName:
		{
			services, err := b.repo.GetAllSubscriberServicesByName(update.Message.Chat.ID, update.Message.Text)
			if err != nil {
				return err
			}

			reply := fmt.Sprintf("Результат запроса:\n\n %v", models.FormatAllSubscriberServiceOutput(services))
			if len(services) == 0 {
				reply = fmt.Sprintf("Поиск не дал результатов.\n")
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

			b.bot.Send(msg)

			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusNone)
			if err != nil {
				logrus.Println("handleCommandSet(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	default:
		logrus.Printf("Unknown status(")
	}

	return nil
}

func (b *Bot) handleCommandStart(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, replyStart)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, replyUnknownCommand)
	_, err := b.bot.Send(msg)

	return err
}

func (b *Bot) handleCommandSubscribe(message *tgbotapi.Message) error {
	subscriber := models.Subscriber{ChatId: message.Chat.ID}

	var id int
	id, err := b.repo.CreateSubscriber(subscriber)
	if err != nil {
		logrus.Printf("Error in  handleCommandSubscribe(): %v", err.Error())
		return errUnableToSubscribe
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, successfulSubscription)
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}

	logrus.Println("Subscribed successfully. ID: %d", id)

	b.sendData(message.Chat.ID)

	return nil
}

func (b *Bot) handleCommandUnsubscribe(message *tgbotapi.Message) error {
	err := b.repo.DeleteSubscriber(message.Chat.ID)
	if err != nil {
		logrus.Printf("Error in  handleCommandUnsubscribe(): %v", err.Error())
		return errUnableToUnsubscribe
	}

	logrus.Println("Unsubscribed successfully.")

	msg := tgbotapi.NewMessage(message.Chat.ID, successfulUnsubscription)
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleCommandCheckSubscribe(message *tgbotapi.Message) error {
	_, err := b.repo.GetSubscriber(message.Chat.ID)
	if err != nil {
		msg := tgbotapi.NewMessage(message.Chat.ID, subscriptionStatusBad)
		_, err = b.bot.Send(msg)
		if err != nil {
			logrus.Printf("Error in  handleCommandCheckSubscribe(): %v", err.Error())
			return err
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, subscriptionStatusGood)
		_, err = b.bot.Send(msg)
		if err != nil {
			logrus.Printf("Error in  handleCommandCheckSubscribe(): %v", err.Error())
			return err
		}
	}

	return nil
}

func (b *Bot) handleCommandSet(update *tgbotapi.Update) error {
	subscriber, err := b.repo.GetSubscriber(update.Message.Chat.ID)
	if err != nil {
		return errDisabledCommand
	}
	logrus.Println("handleCommandSet(): got subscriber %v", subscriber)

	switch subscriber.DialogStatus {
	case models.DialogStatusNone:
		{
			err := b.sendReply(update.Message.Chat.ID, ReplySendServiceName)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusSetServiceName)
			if err != nil {
				logrus.Println("handleCommandSet(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusSetServiceName:
		{
			return b.sendReply(update.Message.Chat.ID, ReplySendServiceName)
		}
	case models.DialogStatusSetServiceLogin:
		{
			return b.sendReply(update.Message.Chat.ID, ReplySendServiceLogin)
		}
	case models.DialogStatusSetServicePassword:
		{
			return b.sendReply(update.Message.Chat.ID, ReplySendServicePassword)
		}
	}

	return nil
}

func (b *Bot) handleCommandGet(update *tgbotapi.Update) error {
	subscriber, err := b.repo.GetSubscriber(update.Message.Chat.ID)
	if err != nil {
		return errDisabledCommand
	}
	logrus.Println("handleCommandSet(): got subscriber %v", subscriber)

	if subscriber.ServiceInProgressID != 0 {
		return b.sendReply(update.Message.Chat.ID, ReplyFinishSettingService)
	}

	err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusGetServiceName)
	if err != nil {
		return err
	}

	err = b.sendReply(update.Message.Chat.ID, ReplySendServiceName)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleCommandDel(message *tgbotapi.Message) error {
	//subscriber := models.Subscriber{ChatId: message.Chat.ID}

	//var id int
	//id, err := b.repo.CreateSubscriber(subscriber)
	//if err != nil {
	//	logrus.Printf("Error in  handleCommandSubscribe(): %v", err.Error())
	//	return errUnableToSubscribe
	//}
	//
	//msg := tgbotapi.NewMessage(message.Chat.ID, successfulSubscription)
	//_, err = b.bot.Send(msg)
	//if err != nil {
	//	return err
	//}
	//
	//logrus.Println("Subscribed successfully. ID: %d", id)
	//
	//b.sendData(message.Chat.ID)

	return nil
}
