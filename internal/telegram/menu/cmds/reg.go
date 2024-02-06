package cmds

import (
	"context"

	"github.com/binaryty/tg-bot/internal/models"
	"github.com/binaryty/tg-bot/internal/telegram/menu/keyboard"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// CmdReg register a new event.
func CmdReg() models.CmdFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Зарегистрировать новое событие?")

		msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
			keyboard.AddKeyboardRow("Да", "BTN_YES"),
			keyboard.AddKeyboardRow("Нет", "BTN_NO"),
		)

		if _, err := bot.Send(msg); err != nil {
			return err
		}
		return nil
	}
}
