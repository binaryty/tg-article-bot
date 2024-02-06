package mw

import (
	"context"
	"github.com/binaryty/tg-bot/internal/models"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

// AdmOnly middleware to restrict access only for group admins.
func AdmOnly(chatId int64, next models.CmdFunc) models.CmdFunc {
	return func(ctx context.Context, bot *tgbotapi.BotAPI, update *tgbotapi.Update) error {
		admins, err := bot.GetChatAdministrators(
			tgbotapi.ChatConfig{
				ChatID: chatId,
			},
		)

		if err != nil {
			return err
		}

		for _, adm := range admins {
			if adm.User.ID == update.Message.From.ID {
				return next(ctx, bot, update)
			}
		}

		if _, err := bot.Send(tgbotapi.NewMessage(
			update.Message.Chat.ID,
			"У вас нет прав на выполнение этой команды.",
		)); err != nil {
			return err
		}

		return nil
	}
}
