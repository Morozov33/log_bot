package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

func main() {
	cwd, _ := os.Getwd()
	logFile := filepath.Join(cwd, os.Getenv("LOG_FILE"))
	logger := logrus.New()
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.SetOutput(os.Stdout)
	}
	defer file.Close()
	logger.SetOutput(file)

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		logger.Error(err)
		logger.Exit(0)
	}

	bot.Debug = true
}
