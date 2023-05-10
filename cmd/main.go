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

	//bot.Debug = true

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})

	repos := repository.NewRepository(db)
	//services := service.NewService(repos)

	messageTTLInMinutes, err := strconv.Atoi(viper.GetString("bot.messageTTLInMinutes"))
	if err != nil {
		logrus.Fatalf("Caught error while parsing bot.messageTTLInMinutes ", err.Error())
	}

	tgBot := telegram.NewBot(bot, *repos, messageTTLInMinutes)
	if err = tgBot.Start(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
