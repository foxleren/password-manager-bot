package telegram

import (
	"github.com/foxleren/password-manager-bot/internal/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/siruspen/logrus"
	"log"
	"time"
)

type Bot struct {
	bot  *tgbotapi.BotAPI
	repo repository.Repository
}

func NewBot(bot *tgbotapi.BotAPI, repo repository.Repository) *Bot {
	return &Bot{
		bot:  bot,
		repo: repo,
	}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	setParsingTime()

	updates := b.initUpdatesChannel()
	go b.sendDataToSubscribers()
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

func (b *Bot) sendDataToSubscribers() {
	for {
		err := b.compileParser()
		if err != nil {
			logrus.Println("Error while compiling python script: %s", err.Error())
		}

		subscribers, err := b.repo.GetAllSubscribers()
		if err != nil {
			log.Printf("Error in GetAllSubscribers(): %v", err.Error())
			continue
		}

		logrus.Println("Starting sending data...")
		for _, sbs := range subscribers {
			go b.sendData(sbs.ChatId)
		}
		logrus.Println("Finishing sending data...")

		logrus.Println("Taking timeout...")
		time.Sleep(24 * time.Hour)
	}
}

func (b *Bot) sendData(chatId int64) {
	//filePath := b.parserData.ExcelFile
	//
	//file, err := os.Open(filePath)
	//if !errors.Is(err, os.ErrNotExist) {
	//	defer file.Close()
	//
	//	fileInfo, err := file.Stat()
	//	if err != nil {
	//		logrus.Println(err)
	//	}
	//
	//	fileBytes := make([]byte, fileInfo.Size())
	//
	//	_, err = file.Read(fileBytes)
	//	if err != nil {
	//		logrus.Println(err)
	//	}
	//
	//	fileBytesConfig := tgbotapi.FileBytes{Name: fileInfo.Name(), Bytes: fileBytes}
	//
	//	msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("Версия от %s.\nКол-во обновлений: %s", parsingTime, parsingUpdateCounter))
	//	_, err = b.bot.Send(msg)
	//	if err != nil {
	//		return
	//	}
	//
	//	doc := tgbotapi.NewDocument(chatId, fileBytesConfig)
	//	_, err = b.bot.Send(doc)
	//	if err != nil {
	//		return
	//	}
	//}
}
