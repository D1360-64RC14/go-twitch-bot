package chat

import (
	"fmt"
	"log"

	"github.com/gempir/go-twitch-irc/v2"
)

type Options struct {
	Username string
	OAuth    string
	Channel  string
}

// NewChat cria e faz a conexão com o cliente de chat.
func NewChat(chatopt Options) *twitch.Client {
	var client = twitch.NewClient(chatopt.Username, chatopt.OAuth)
	client.Join(chatopt.Channel)

	return client
}

// Say envia uma mensagem no chat do canal.
func Say(client *twitch.Client, channel, message string) {
	client.Say(channel, fmt.Sprintf("/me (bot): %s", message))
}

func Permission(message twitch.PrivateMessage, perms ...string) bool {
	var accumulator = false

	for _, permission := range perms {
		accumulator = accumulator || message.User.Badges[permission] == 1
	}

	return accumulator
}
func PermissionAND(message twitch.PrivateMessage, perms ...string) bool {
	var accumulator = false

	for _, permission := range perms {
		accumulator = accumulator && message.User.Badges[permission] == 1
	}

	return accumulator
}

func Int2Bool(num uint8) bool {
	if num == 1 {
		return true
	}
	return false
}

func Log(message twitch.PrivateMessage) {
	log.Printf("%s | %s: %s",
		message.Channel,
		message.User.DisplayName,
		message.Message,
	)
}
