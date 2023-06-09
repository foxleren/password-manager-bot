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
		return b.handleCommandDel(update)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleMessage(update *tgbotapi.Update) error {
	subscriber, err := b.repo.GetSubscriber(update.Message.Chat.ID)
	if err != nil {
		logrus.Println("Level: service.telegram; handleMessage(): error while getting subscriber with subscriber_id %v", subscriber.ID)
		return errDisabledCommand
	}

	//logrus.Println("Level: service.telegram; handleMessage(): got subscriber %v", subscriber)

	switch subscriber.DialogStatus {
	case models.DialogStatusNone:
		{

			err := b.sendReply(update.Message.Chat.ID, ReplySendServiceNameToSet)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusGetServiceName)
			if err != nil {
				logrus.Println("Level: service.telegram; handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusSetServiceName:
		{
			_, err = b.repo.SubscriberService.CreateSubscriberServiceByName(subscriber.ID, update.Message.Text)
			if err != nil {
				logrus.Println("Level: service.telegram; handleMessage(): error while creating service by name for subscriber: %v", subscriber.ID)
				return errUnableToSetService
			}
			err := b.sendReply(update.Message.Chat.ID, ReplyServiceNameIsSet)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusSetServiceLogin)
			if err != nil {
				logrus.Println("Level: service.telegram; handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusSetServiceLogin:
		{
			err = b.repo.SubscriberService.UpdateSubscriberServiceLogin(subscriber.ID, subscriber.ServiceInProgressID, update.Message.Text)
			if err != nil {
				logrus.Println("Level: service.telegram; handleMessage(): error while updating service login for subscriber: %v", subscriber.ID)
				return err
			}
			err := b.sendReply(update.Message.Chat.ID, ReplyServiceLoginIsSet)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusSetServicePassword)
			if err != nil {
				logrus.Println("Level: service.telegram; handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusSetServicePassword:
		{
			err = b.repo.SubscriberService.UpdateSubscriberServicePassword(subscriber.ID, subscriber.ServiceInProgressID, update.Message.Text)
			if err != nil {
				logrus.Println("Level: service.telegram; handleMessage(): error while updating service password for subscriber: %v", subscriber.ID)
				return err
			}
			err := b.sendReply(update.Message.Chat.ID, ReplyServicePasswordIsSet)
			if err != nil {
				return err
			}
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusNone)
			if err != nil {
				logrus.Println("Level: service.telegram; handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusGetServiceName:
		{
			services, err := b.repo.SubscriberService.GetSubscriberServiceByName(update.Message.Chat.ID, update.Message.Text)
			if err != nil {
				logrus.Println("Level: service.telegram; handleMessage(): error while getting service by name for subscriber: %v", subscriber.ID)
				return errUnableToGetService
			}

			reply := fmt.Sprintf("Результат запроса:\n\n%v", services.String())

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			msg.ParseMode = tgbotapi.ModeMarkdown
			sentMsg, err := b.bot.Send(msg)

			_, err = b.repo.CreateMessage(update.Message.Chat.ID, sentMsg.MessageID)
			if err != nil {
				logrus.Println("Level: service.telegram; handleMessage(): error while saving message subscriber: %v", subscriber.ID)
				return err
			}

			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusNone)
			if err != nil {
				logrus.Println("Level: service.telegram; handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
		}
	case models.DialogStatusDelService:
		{
			err := b.repo.SubscriberService.DeleteSubscriberService(update.Message.Chat.ID, update.Message.Text)
			if err != nil {
				logrus.Println("Level: service.telegram; handleMessage(): error while deleting service for subscriber: %v", subscriber.ID)
				return errUnableToDeleteService
			}

			err = b.sendReply(update.Message.Chat.ID, ReplyServiceIsDeleted)
			if err != nil {
				return err
			}

			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusNone)
			if err != nil {
				logrus.Println("Level: service.telegram; handleMessage(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
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
		logrus.Printf("Level: service.telegram; func handleCommandSubscribe(): error while creating subscriber: err=%v", err.Error())
		return errUnableToSubscribe
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, successfulSubscription)
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}

	logrus.Println("Subscribed successfully. ID: %d", id)

	return nil
}

func (b *Bot) handleCommandUnsubscribe(message *tgbotapi.Message) error {
	err := b.repo.DeleteSubscriber(message.Chat.ID)
	if err != nil {
		logrus.Printf("Level: service.telegram; func handleCommandUnsubscribe(): error while deleting subscriber: err=%v", err.Error())
		return errUnableToUnsubscribe
	}

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
			logrus.Printf("Level: service.telegram; func handleCommandCheckSubscribe(): error while checking subscriber status: err=%v", err.Error())
			return err
		}
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, subscriptionStatusGood)
		_, err = b.bot.Send(msg)
		if err != nil {
			logrus.Printf("Level: service.telegram; func handleCommandCheckSubscribe(): error while checking subscriber status: err=%v", err.Error())
			return err
		}
	}

	return nil
}

func (b *Bot) handleCommandSet(update *tgbotapi.Update) error {
	subscriber, err := b.repo.GetSubscriber(update.Message.Chat.ID)
	if err != nil {
		logrus.Printf("Level: service.telegram; func handleCommandSet(): error while getting subscriber: err=%v", err.Error())
		return errDisabledCommand
	}

	if subscriber.ServiceInProgressID != 0 {
		return b.sendReply(update.Message.Chat.ID, ReplyFinishSettingService)
	}

	switch subscriber.DialogStatus {
	case models.DialogStatusNone:
		{
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusSetServiceName)
			if err != nil {
				logrus.Println("Level: service.telegram; func handleCommandSet(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
			err := b.sendReply(update.Message.Chat.ID, ReplySendServiceNameToSet)
			if err != nil {
				return err
			}
		}
	case models.DialogStatusSetServiceName:
		{
			return b.sendReply(update.Message.Chat.ID, ReplySendServiceNameToSet)
		}
	case models.DialogStatusSetServiceLogin:
		{
			return b.sendReply(update.Message.Chat.ID, ReplySendServiceLogin)
		}
	case models.DialogStatusSetServicePassword:
		{
			return b.sendReply(update.Message.Chat.ID, ReplySendServicePassword)
		}
	default:
		{
			err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusSetServiceName)
			if err != nil {
				logrus.Println("Level: service.telegram; func handleCommandSet(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
				return err
			}
			err := b.sendReply(update.Message.Chat.ID, ReplySendServiceNameToSet)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Bot) handleCommandGet(update *tgbotapi.Update) error {
	subscriber, err := b.repo.GetSubscriber(update.Message.Chat.ID)
	if err != nil {
		logrus.Printf("Level: service.telegram; func handleCommandGet(): error while getting subscriber: err=%v", err.Error())
		return errDisabledCommand
	}

	if subscriber.ServiceInProgressID != 0 {
		return b.sendReply(update.Message.Chat.ID, ReplyFinishSettingService)
	}

	err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusGetServiceName)
	if err != nil {
		return err
	}

	err = b.sendReply(update.Message.Chat.ID, ReplySendServiceNameToGet)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleCommandDel(update *tgbotapi.Update) error {
	subscriber, err := b.repo.GetSubscriber(update.Message.Chat.ID)
	if err != nil {
		return errDisabledCommand
	}
	logrus.Println("handleCommandSet(): got subscriber %v", subscriber)

	if subscriber.ServiceInProgressID != 0 {
		return b.sendReply(update.Message.Chat.ID, ReplyFinishSettingService)
	}

	err = b.repo.UpdateSubscriberDialogStatus(update.Message.Chat.ID, models.DialogStatusDelService)
	if err != nil {
		logrus.Println("Level: service.telegram; func handleCommandDel(): error while updating dialog_status for subscriber_id: %v", subscriber.ID)
		return err
	}

	err = b.sendReply(update.Message.Chat.ID, ReplySendServiceNameToDelete)
	if err != nil {
		return err
	}

	return nil
}
