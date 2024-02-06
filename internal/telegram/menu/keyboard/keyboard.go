package keyboard

import tgbotapi "gopkg.in/telegram-bot-api.v4"

// AddKeyboardRow add an inline keyboard row.
func AddKeyboardRow(btnText string, btnData string) []tgbotapi.InlineKeyboardButton {

	return tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(btnText, btnData),
	)
}
