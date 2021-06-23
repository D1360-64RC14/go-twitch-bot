package commands

import (
	"regexp"

	"github.com/d1360-64rc14/twitch-bot/commands/objects"
	"github.com/gempir/go-twitch-irc/v2"
	"gorm.io/gorm"
)

var Working = objects.Command{
	Name: "Working?",
	Description: "Comando utilizado por moderadores para verificar a atividade do bot.",
	Pattern: regexp.MustCompile(`^!working\?$`),
	CaseSensitive: false,
	Cooldown: &objects.Cooldown{Global: 0, User: 0},
	Behavior: WorkingBehavior,
}

func WorkingBehavior(message twitch.PrivateMessage, client *twitch.Client, database *gorm.DB, command *objects.Command) string {
	return "Estou Funcionando VoHiYo"
}
