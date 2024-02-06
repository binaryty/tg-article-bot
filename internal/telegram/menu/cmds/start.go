package cmds

import (
	"context"

	"github.com/binaryty/tg-bot/internal/models"
	"github.com/binaryty/tg-bot/internal/telegram/menu"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// CmdStart displays information on the time from the start of the bot.
func CmdStart() models.CmdFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {

		if _, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, menu.MsgStart)); err != nil {
			return err
		}
		return nil
	}
}
