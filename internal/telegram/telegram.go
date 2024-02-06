package telegram

import (
	"context"
	"fmt"
	"github.com/binaryty/tg-bot/internal/models"
	"github.com/binaryty/tg-bot/internal/telegram/menu/keyboard"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

const bathSize = 50
const admID = 141431779

type ArticleProvider interface {
	ReadByTitle(string) ([]models.Article, error)
}

type CmdProvider interface {
	RegisterCmd(string, models.CmdFunc)
	Command(string) (models.CmdFunc, bool)
}

type TgClient struct {
	api     *tgbotapi.BotAPI
	storage ArticleProvider
	menu    CmdProvider
}

// New a constructor of TgClient.
func New(api *tgbotapi.BotAPI, storage ArticleProvider, menu CmdProvider) *TgClient {
	return &TgClient{
		api:     api,
		storage: storage,
		menu:    menu,
	}
}

// Run Start TgClient.
func (c *TgClient) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := c.api.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	log.Println("[INFO] starting handling updates")

	for {
		select {
		case update := <-updates:
			updateCtx, updateCancel := context.WithTimeout(ctx, 5*time.Second)
			c.HandleUpdate(updateCtx, &update)
			updateCancel()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// HandleUpdate handle an updates from telegram bot.
func (c *TgClient) HandleUpdate(ctx context.Context, update *tgbotapi.Update) {
	defer func() {
		if p := recover(); p != nil {
			log.Printf("[ERROR] panic recovered: %v\n%s", p, string(debug.Stack()))
		}
	}()

	switch {
	case isCallbackQuery(update):
		c.HandleCallbackQuery(update)
	case isCommand(update):
		log.Printf("[INFO] got a new update: [from]: %s [subject]: %s", update.Message.From.UserName, update.Message.Text)

		cmd := update.Message.Command()

		cmdFunc, ok := c.menu.Command(cmd)
		if !ok {
			return
		}

		if err := cmdFunc(ctx, c.api, update); err != nil {
			log.Printf("[ERROR] failed to handle update: %v", err)

			if _, err := c.api.Send(
				tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%s\n%v", "[Внутренняя ошибка]", err)),
			); err != nil {
				log.Printf("[ERROR] failed to send message: %v", err)
			}
		}
	case isInlineQuery(update):
		if err := c.HandleInlineQuery(update); err != nil {
			log.Printf("[ERROR] can't handle inline query: %v", err)
			return
		}
		log.Printf("[INFO] got a new inline query: [from]: %s [subject]: %s", update.InlineQuery.From.UserName, update.InlineQuery.Query)
	}
}

// HandleInlineQuery handle an inline queries from telegram bot.
func (c *TgClient) HandleInlineQuery(update *tgbotapi.Update) error {
	inlineQuery := update.InlineQuery
	queryOffset, _ := strconv.Atoi(inlineQuery.Offset)

	if queryOffset == 0 {
		queryOffset = 1
	}

	results := make([]interface{}, 0)

	articles, err := c.storage.ReadByTitle(strings.ToLower(inlineQuery.Query))
	if err != nil {
		return err
	}

	for _, article := range offsetResult(queryOffset, articles) {
		msg := fmt.Sprintf(
			"[%s](%s)\n"+
				"Дата публикации на habr: _%s_\n",
			article.Title,
			article.Link,
			article.PublishedAt)
		results = append(results, tgbotapi.InlineQueryResultArticle{
			Type:  "article",
			ID:    article.Id(),
			Title: article.Title,
			InputMessageContent: tgbotapi.InputTextMessageContent{
				Text:      msg,
				ParseMode: tgbotapi.ModeMarkdown,
			},
			ThumbURL: article.ThumbUrl,
		})
	}

	if len(results) < 50 {
		_, err := c.api.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID: inlineQuery.ID,
			Results:       results,
			IsPersonal:    true,
			CacheTime:     0,
		})

		if err != nil {
			return err
		}
	} else {
		_, err := c.api.AnswerInlineQuery(tgbotapi.InlineConfig{
			InlineQueryID: inlineQuery.ID,
			Results:       results,
			IsPersonal:    true,
			CacheTime:     0,
			NextOffset:    strconv.Itoa(queryOffset + bathSize),
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// showMenu ...
func (c *TgClient) showMenu(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выбрать действие")

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		keyboard.AddKeyboardRow("Добавить", "BTN_ADD"),
		keyboard.AddKeyboardRow("Удалить", "BTN_DELETE"),
		keyboard.AddKeyboardRow("Редактировать", "BTN_EDIT"),
		keyboard.AddKeyboardRow("Назад", "BTN_BACK"),
	)

	if _, err := c.api.Send(msg); err != nil {
		return
	}
}

// HandleCallbackQuery handle a callback query received from update telegram bot.
func (c *TgClient) HandleCallbackQuery(update *tgbotapi.Update) {
	switch update.CallbackQuery.Data {
	case "BTN_YES":
		if update.CallbackQuery.From.ID == admID {
			c.showMenu(update)
		} else {
			if _, err := c.api.Send(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "У вас отсутствуют права на доступ.")); err != nil {
				return
			}
		}

	case "BTN_NO":
		return
	}

	return
}

// isCommand check update for command message.
func isCommand(update *tgbotapi.Update) bool {

	return update.Message != nil && update.Message.IsCommand()
}

// isCallbackQuery check update for callback query data.
func isCallbackQuery(update *tgbotapi.Update) bool {

	return update.Message == nil && update.CallbackQuery != nil && update.CallbackQuery.Data != ""
}

// isInlineQuery check update for inline query.
func isInlineQuery(update *tgbotapi.Update) bool {

	return update.Message == nil && update.InlineQuery != nil
}

// offsetResult return pagination of slice of articles.
func offsetResult(startNum int, articles []models.Article) []models.Article {
	overallItems := len(articles)

	switch {
	case startNum >= overallItems:
		return []models.Article{}
	case startNum+bathSize >= overallItems:
		return articles[startNum:overallItems]
	default:
		return articles[startNum : startNum+bathSize]
	}
}
