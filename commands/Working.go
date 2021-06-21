package commands

import (
	"fmt"
	"regexp"

	"github.com/d1360-64rc14/twitch-bot/commands/objects"
	"github.com/gempir/go-twitch-irc/v2"
	"gorm.io/gorm"
)

var Working = objects.Command{
	Name: "Working",
	Pattern: regexp.MustCompile(`^!working\?$`),
	CaseSensitive: false,
	Cooldown: &objects.Cooldown{Global: 0, User: 0, Behavior: WorkingCooldownBehavior},
	Behavior: WorkingBehavior,
}

func WorkingBehavior(message twitch.PrivateMessage, client *twitch.Client, database *gorm.DB, command *objects.Command) string {
	return "Estou Funcionando VoHiYo"
}

func WorkingCooldownBehavior(_ twitch.PrivateMessage, _ *twitch.Client, _ *gorm.DB, command *objects.Command) string {
	return fmt.Sprintf("O Comando '%s' está em cooldown!", command.Name)
}