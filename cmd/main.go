package main

import (
	"github.com/foxleren/password-manager-bot/internal/repository"
	"github.com/foxleren/password-manager-bot/internal/service/telegram"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/siruspen/logrus"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if err := initConfig(); err != nil {
		logrus.Fatalf("Caught error while initializing config: ", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Caught error while loading .env file: ", err.Error())
	}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	if err != nil {
		logrus.Fatalf("Caught error while creating database: ", err.Error())
	}

	repos := repository.NewRepository(db)

	msgTTLInMinutes, err := strconv.Atoi(viper.GetString("bot.messageTTLInMinutes"))
	if err != nil {
		logrus.Fatalf("Caught error while parsing bot.messageTTLInMinutes ", err.Error())
	}

	msgTTLInHours, err := strconv.Atoi(viper.GetString("bot.messageTTLInHours"))
	if err != nil {
		logrus.Fatalf("Caught error while parsing bot.messageTTLInHours ", err.Error())
	}

	outdatedMsgSleepInHours, err := strconv.Atoi(viper.GetString("bot.outdatedMessagesSleepIntervalInHours"))
	if err != nil {
		logrus.Fatalf("Caught error while parsing bot.outdatedMessagesSleepIntervalInHours ", err.Error())
	}

	tgBot := telegram.NewBot(bot, *repos, &telegram.BotConfig{
		MessageTTLInMinutes:               msgTTLInMinutes,
		MessageTTLInHours:                 msgTTLInHours,
		OutdatedMessagesSleepLimitInHours: outdatedMsgSleepInHours,
	})
	if err = tgBot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
