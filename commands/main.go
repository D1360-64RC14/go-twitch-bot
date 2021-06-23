package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/d1360-64rc14/twitch-bot/chat"
	"github.com/d1360-64rc14/twitch-bot/commands/objects"
	"github.com/gempir/go-twitch-irc/v2"
	"gorm.io/gorm"
)

// Os comandos nessa lista serão habilitados.
var AvailableCommands = []objects.Command{
	Bot,
	Working,
}

func Handler(client *twitch.Client, db *gorm.DB) {
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		go chat.Log(message)

		// Execução do comando 'AllCommands'
		if AllCommands.Validate(message) {
			AllCommands.Exec(message, client, db) // Test: use goroutine
			return
		}

		// Percorre a lista de comandos procurando
		// algum que passe pela verificação.
		for _, cmmnd := range AvailableCommands {
			if cmmnd.Validate(message) {
				cmmnd.Exec(message, client, db) // Test: use goroutine
			}
		}
	})
}

var AllCommands = objects.Command{
	Name: "AllCommands",
	Description: "Utilizado para visualizar todos os comandos disponíveis.",
	Pattern: regexp.MustCompile(`(^!commands$)|(^!comandos$)`),
	CaseSensitive: false,
	Cooldown: &objects.Cooldown{Global: 10, User: 0},
	Behavior: AllCommandsBehavior,
}

func AllCommandsBehavior(message twitch.PrivateMessage, client *twitch.Client, database *gorm.DB, command *objects.Command) string {
	var commandNames []string

	for _, cmmnd := range AvailableCommands {
		commandNames = append(commandNames, cmmnd.Name)
	}

	var msg = fmt.Sprintf("Os comandos disponíveis são: %s; " +
		"Use !info <comando> para visualizar a descrição do mesmo.",
		strings.Join(commandNames, ", "),
	)

	return msg
}
