package cmds

import (
	"context"
	"fmt"
	"github.com/binaryty/tg-bot/internal/models"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type RndProvider interface {
	ReadRandom() (*models.Article, error)
}

// CmdRnd get a random article from storage and sent to telegram.
func CmdRnd(r RndProvider) models.CmdFunc {
	return func(ctx context.Context, api *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		article, err := r.ReadRandom()
		if err != nil {
			return err
		}

		text := fmt.Sprintf(
			"[%s](%s)\nДата публикации: _%s_",
			article.Title,
			article.Link,
			article.PublishedAt,
		)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ParseMode = tgbotapi.ModeMarkdown

		if _, err := api.Send(msg); err != nil {
			return err
		}

		return nil
	}
}
