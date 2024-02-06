package main

import (
	"context"
	"github.com/binaryty/tg-bot/internal/app"
	"github.com/joho/godotenv"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	token := loadToken()
	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Printf("[ERROR] token must be specifized: %v", err)
		os.Exit(1)
	}
	log.Printf("[INFO] Authorized on account %s", botApi.Self.UserName)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	App := app.New(botApi)

	log.Fatal(App.Run(ctx))
}

func loadToken() string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("[FATAL ERROR] no .env file found")
	}
	token, exists := os.LookupEnv("TOKEN")

	if !exists {
		return ""
	}
	return token
}
