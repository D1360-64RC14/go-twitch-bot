package commands

import (
	"fmt"
	"regexp"

	"github.com/d1360-64rc14/twitch-bot/commands/objects"
	"github.com/gempir/go-twitch-irc/v2"
	"gorm.io/gorm"
)

func BotBehavior(message twitch.PrivateMessage, client *twitch.Client, database *gorm.DB, command *objects.Command) string {
	var msg = fmt.Sprintf(
		"Oi @%s, fui feito utilizando GoLang, estou rodando em um Docker num Raspberry Pi 4B e " +
		"uso PostgreSQL como base de dados VoHiYo",
		message.User.DisplayName,
	)
	
	return msg
}

var Bot = objects.Command{
	Name: "Bot",
	Pattern: regexp.MustCompile(`^!bot.*`),
	CaseSensitive: false,
	Cooldown: &objects.Cooldown{Global: 60, User: 0},
	Behavior: BotBehavior,
}