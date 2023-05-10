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
		return b.handleCommandGet(message)
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
			reply := fmt.Sprintf("Пришлите название сервиса, для которого вы хотите сохранить пароль")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			_, err := b.bot.Send(msg)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusWaitingForServiceName)
			if err != nil {
				logrus.Println("handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusWaitingForServiceName:
		{
			_, err = b.repo.SubscriberService.CreateSubscriberServiceByName(subscriber.ID, update.Message.Text)
			if err != nil {
				logrus.Printf("error in handler")
				return err
			}
			reply := fmt.Sprintf("Пришлите логин для сервиса")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			_, err = b.bot.Send(msg)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusWaitingForServiceLogin)
			if err != nil {
				logrus.Println("handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusWaitingForServiceLogin:
		{
			err = b.repo.SubscriberService.UpdateSubscriberServiceLogin(subscriber.ID, subscriber.ServiceInProgressID, update.Message.Text)
			if err != nil {
				logrus.Printf("error in handler")
				return err
			}
			reply := fmt.Sprintf("Пришлите пароль для сервиса")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			_, err := b.bot.Send(msg)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusWaitingForServicePassword)
			if err != nil {
				logrus.Println("handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusWaitingForServicePassword:
		{
			err = b.repo.SubscriberService.UpdateSubscriberServicePassword(subscriber.ID, subscriber.ServiceInProgressID, update.Message.Text)
			if err != nil {
				logrus.Printf("error in handler")
				return err
			}
			reply := fmt.Sprintf("Сервис успешно добавлен!")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			_, err := b.bot.Send(msg)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusNone)
			if err != nil {
				logrus.Println("handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
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
	//subscriber := models.Subscriber{ChatId: message.Chat.ID}

	//var id int
	//id, err := b.repo.CreateSubscriber(subscriber)
	//if err != nil {
	//	logrus.Printf("Error in  handleCommandSubscribe(): %v", err.Error())
	//	return errUnableToSubscribe
	//}
	//

	subscriber, err := b.repo.GetSubscriber(update.Message.Chat.ID)
	if err != nil {
		return errDisabledCommand
	}

	logrus.Println("handleCommandSet(): got subscriber %v", subscriber)

	switch subscriber.DialogStatus {
	case models.DialogStatusNone:
		{
			reply := fmt.Sprintf("Пришлите название сервиса, для которого вы хотите сохранить пароль")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			_, err := b.bot.Send(msg)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusWaitingForServiceName)
			if err != nil {
				logrus.Println("handleCommandSet(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusWaitingForServiceName:
		{
			reply := fmt.Sprintf("Пришлите название сервиса")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			_, err := b.bot.Send(msg)
			if err != nil {
				return err
			}
		}
	case models.DialogStatusWaitingForServiceLogin:
		{
			reply := fmt.Sprintf("Пришлите логин для сервиса")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			_, err := b.bot.Send(msg)
			if err != nil {
				return err
			}
		}
	case models.DialogStatusWaitingForServicePassword:
		{
			reply := fmt.Sprintf("Пришлите пароль для сервиса")
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			_, err := b.bot.Send(msg)
			if err != nil {
				return err
			}
		}
	}

	//message := update.Message
	//
	//reply := fmt.Sprintf("Пришлите название сервиса, для которого вы хотите сохранить пароль")
	//msg := tgbotapi.NewMessage(message.Chat.ID, reply)
	//_, err := b.bot.Send(msg)
	//if err != nil {
	//	return err
	//}

	//userReply := make(chan string)
	//go func() {
	//	for {
	//		select {
	//		case update := <-updates:
	//			if update.Message != nil {
	//				userReply <- update.Message.Text
	//			}
	//		}
	//	}
	//}()

	// Получаем ответ от пользователя
	//password := <-reply
	//
	//answer := <-b.bot.Дшы(update.Message.Chat.ID, 60*time.Second)
	//
	//// отправляем ответ обратно пользователю
	//reply := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Вы ввели: %s", answer.Text))
	//bot.Send(reply)

	//reply := fmt.Sprintf("Вы успешно сохранили логин и пароль для сервиса: %s", "HUETA")
	//msg := tgbotapi.NewMessage(message.Chat.ID, reply)
	//_, err := b.bot.Send(msg)
	//if err != nil {
	//	return err
	//}
	//
	//logrus.Println("Subscribed successfully. ID: %d", id)
	//
	//b.sendData(message.Chat.ID)

	return nil
}

func (b *Bot) handleCommandGet(message *tgbotapi.Message) error {
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
