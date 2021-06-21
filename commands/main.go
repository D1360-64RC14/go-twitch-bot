package commands

import (
	"github.com/d1360-64rc14/twitch-bot/chat"
	"github.com/d1360-64rc14/twitch-bot/commands/objects"
	"github.com/gempir/go-twitch-irc/v2"
	"gorm.io/gorm"
)

func Handler(client *twitch.Client, db *gorm.DB) {
	var commandList = []objects.Command{
		Bot,
		Working,
		LastUse,
	}

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		go chat.Log(message)

		for _, cmmnd := range commandList {
			if cmmnd.Validate(message) {
				cmmnd.Exec(message, client, db)
			}
		}
	})
}