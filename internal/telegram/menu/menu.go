package menu

import (
	"github.com/binaryty/tg-bot/internal/models"
)

type Menu struct {
	Commands map[string]models.CmdFunc
}

func New() *Menu {
	cmds := make(map[string]models.CmdFunc)

	return &Menu{
		Commands: cmds,
	}
}

// RegisterCmd register a command in the bot menu.
func (m Menu) RegisterCmd(cmd string, cmdFunc models.CmdFunc) {
	if m.Commands == nil {
		m.Commands = make(map[string]models.CmdFunc)
	}

	m.Commands[cmd] = cmdFunc
}

// Command get a command and ok from menu.
func (m Menu) Command(cmd string) (models.CmdFunc, bool) {
	cf, ok := m.Commands[cmd]

	return cf, ok
}
