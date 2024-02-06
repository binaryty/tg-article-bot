package models

import (
	"context"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type CmdFunc func(ctx context.Context, api *tgbotapi.BotAPI, update *tgbotapi.Update) error
