package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	cwd, _ := os.Getwd()
	logFile := filepath.Join(cwd, os.Getenv("LOG_FILE"))
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		logger.Fatal(err, "Failed opening or creating logs file")
	}
	defer file.Close()
	logger.SetOutput(file)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		logger.Fatal(err, "Failed conection to RabbitMQ")
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal(err, "Failed to open a channel")
	}
	defer ch.Close()
	q, err := ch.QueueDeclare(
		os.Getenv("QUEUE_NAME"), // name
		false,                   // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		logger.Fatal(err, "Failed to declare a queue")
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logger.Fatal(err, "Failed to register a consumer")
	}

	var forever chan struct{}

	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		logger.Fatal(err, "Failed to connect a bot")
	}

	chatId, err := strconv.Atoi(os.Getenv("TELEGRAM_CHAT_ID"))
	if err != nil {
		logger.Fatal(err, "Failed to connect a chat")
	}
	go func() {
		for d := range msgs {
			msg := tgbotapi.NewMessage(int64(chatId), "")
			msg.Text = string(d.Body)
			logger.Info(msg.Text)
			if _, err := bot.Send(msg); err != nil {
				logger.Error(err, "Failed to send a message")
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
