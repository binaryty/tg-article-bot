package cmds

import (
	"context"
	"github.com/binaryty/tg-bot/internal/models"
	"github.com/binaryty/tg-bot/internal/telegram/menu"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
)

type ArticleLoader interface {
	Start(ctx context.Context) error
	ProcessArticles(ctx context.Context) error
}

// CmdLoad uploads articles to the database.
func CmdLoad(l ArticleLoader) models.CmdFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {

		if err := l.Start(ctx); err != nil {
			return err
		}

		log.Printf("[INFO] fetcher stopped")

		if err := l.ProcessArticles(ctx); err != nil {
			log.Printf("[ERROR] can't fetch data: %v", err)
			return err
		}

		log.Printf("[INFO] articles loaded")

		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, menu.MsgSaved)); err != nil {
			return err
		}

		return nil
	}
}
