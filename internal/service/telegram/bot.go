package telegram

import (
	"github.com/foxleren/password-manager-bot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/siruspen/logrus"
	"time"
)

type Bot struct {
	bot       *tgbotapi.BotAPI
	botConfig *BotConfig
	repo      repository.Repository
}

type BotConfig struct {
	MessageTTLInMinutes               int
	MessageTTLInHours                 int
	OutdatedMessagesSleepLimitInHours int
}

func NewBot(bot *tgbotapi.BotAPI, repo repository.Repository, botConfig *BotConfig) *Bot {
	return &Bot{
		bot:       bot,
		repo:      repo,
		botConfig: botConfig,
	}
}

func (b *Bot) Start() error {
	logrus.Printf("Level: telegram; func Start(): authorized on account %s", b.bot.Self.UserName)

	updates := b.initUpdatesChannel()
	go b.deleteOutdatedMessages()
	err := b.handleUpdates(updates)
	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) error {
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			if err := b.handleCommand(&update); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
			continue
		} else {
			if err := b.handleMessage(&update); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}

		}
	}

	return nil
}

func (b *Bot) initUpdatesChannel() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func (b *Bot) deleteOutdatedMessages() {
	for {
		lastFreshDate := time.Now()
		lastFreshDate = lastFreshDate.Add(-time.Duration(b.botConfig.MessageTTLInMinutes) * time.Minute)
		lastFreshDate = lastFreshDate.Add(-time.Duration(b.botConfig.MessageTTLInHours) * time.Hour)

		logrus.Printf("Level: telegram; func deleteOutdatedMessages(): set fresh date: %v", lastFreshDate)

		outdatedMessages, err := b.repo.Message.GetAllOutdatedMessages(lastFreshDate)
		if err != nil {
			return
		}

		for _, msg := range outdatedMessages {
			deleteMsg := tgbotapi.NewDeleteMessage(msg.ChatId, msg.MessageId)
			_, err = b.bot.Send(deleteMsg)
		}

		logrus.Printf("Level: telegram; func deleteOutdatedMessages(): taking timeout for %d hours...", b.botConfig.OutdatedMessagesSleepLimitInHours)

		time.Sleep(time.Duration(b.botConfig.OutdatedMessagesSleepLimitInHours) * time.Hour)
	}
}
