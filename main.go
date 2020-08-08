package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"os"
)

func main() {
	checkEnv()

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhookWithCert("https://"+os.Getenv("SERVER_IP")+":8443/"+bot.Token,
		os.Getenv("SSL_CERT")))
	if err != nil {
		log.Fatal(err)
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	go http.ListenAndServeTLS("0.0.0.0:8443", os.Getenv("SSL_CERT"),
		os.Getenv("SSL_CERT_KEY"), nil)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "echo: "+update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func checkEnv() {
	if _, ok := os.LookupEnv("BOT_TOKEN"); !ok {
		log.Panic("BOT_TOKEN is not set!")
	}

	if _, ok := os.LookupEnv("SSL_CERT"); !ok {
		log.Panic("SSL_CERT is not set!")
	}

	if _, ok := os.LookupEnv("SSL_CERT_KEY"); !ok {
		log.Panic("SSL_CERT_KEY is not set!")
	}

	if _, ok := os.LookupEnv("SERVER_IP"); !ok {
		log.Panic("SERVER_IP is not set!")
	}

	if !fileExists(os.Getenv("SSL_CERT")) || !fileExists(os.Getenv("SSL_CERT_KEY")) {
		log.Panic("cannot find ssl-certificate key pair!")
	}
}
