package app

import (
	"context"
	"github.com/binaryty/tg-bot/internal/telegram"
	"github.com/binaryty/tg-bot/internal/telegram/fetcher"
	"github.com/binaryty/tg-bot/internal/telegram/menu"
	"github.com/binaryty/tg-bot/internal/telegram/menu/cmds"
	"github.com/binaryty/tg-bot/internal/telegram/menu/mw"
	"github.com/binaryty/tg-bot/internal/telegram/storage/article/postgres"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"log"
)

const (
	ChatId = -1002060428320
)

type App struct {
	habrFetcher *fetcher.Fetcher
	habrStorage *postgres.Storage
	menu        *menu.Menu
	client      *telegram.TgClient
}

// New create a new Application.
func New(botApi *tgbotapi.BotAPI) *App {
	app := &App{}

	app.habrStorage = postgres.New()

	app.habrFetcher = fetcher.New(app.habrStorage)

	app.menu = menu.New()

	app.menu.RegisterCmd("start", cmds.CmdStart())
	app.menu.RegisterCmd("help", cmds.CmdHelp())
	app.menu.RegisterCmd("load", mw.AdmOnly(ChatId, cmds.CmdLoad(app.habrFetcher)))
	app.menu.RegisterCmd("rnd", cmds.CmdRnd(app.habrStorage))
	app.menu.RegisterCmd("reg", cmds.CmdReg())

	app.client = telegram.New(botApi, app.habrStorage, app.menu)

	return app
}

// Run application.
func (a App) Run(ctx context.Context) error {
	err := a.client.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
