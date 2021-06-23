package objects

import (
	"regexp"
	"strings"
	"time"

	"github.com/d1360-64rc14/twitch-bot/chat"
	"github.com/gempir/go-twitch-irc/v2"
	"gorm.io/gorm"
)

// Permission Enum.
// Based on Twitch Badges.
var (
	PERM_STREAMER  = []string{"broadcaster"}
	PERM_MODERATOR = append(PERM_STREAMER, "moderator")
	PERM_VIP       = append(PERM_MODERATOR, "vip")
)

type Cooldown struct {
	Global     uint32
	User       uint32
	Behavior   CommandBehavior
	globalSave time.Time
	userSave   map[string]time.Time
}

type CommandBehavior func(
	message   twitch.PrivateMessage,
	client   *twitch.Client,
	database *gorm.DB,
	command  *Command,
) string

type Command struct {
	Name              string
	Description       string
	Pattern          *regexp.Regexp
	CaseSensitive     bool
	PermissionLevel []string
	Cooldown         *Cooldown
	Behavior          CommandBehavior
}

// NewCommand cria uma instância do objeto 'Command'.
func NewCommand(
		name                string,
		description         string,
		pattern            *regexp.Regexp,
		caseSensitive       bool,
		permissionLevel   []string,
		userCooldown        uint32,
		globalCooldown      uint32,
		behavior            CommandBehavior,
		onCooldownBehavior  CommandBehavior,
	) Command {
	return Command{
		Name:            name,
		Description:     description,
		Pattern:         pattern,
		CaseSensitive:   caseSensitive,
		PermissionLevel: permissionLevel,
		Behavior:        behavior,
		Cooldown:       &Cooldown{
			globalSave: time.Time{},
			userSave:   make(map[string]time.Time),
			Behavior:   onCooldownBehavior,
			Global:     globalCooldown,
			User:       userCooldown,
		},
	}
}

// Validate verifica se o comando é encontrado na mensagem passada.
func (c *Command) Validate(message twitch.PrivateMessage) bool {
	if c.CaseSensitive {
		return c.Pattern.MatchString(message.Message)
	}
	return c.Pattern.MatchString(strings.ToLower(message.Message))
}

// Exec executa o comportamento programado para o comando.
func (c *Command) Exec(message twitch.PrivateMessage, client *twitch.Client, database *gorm.DB) {
	if c.InCooldownAll(message.User) {
		if c.Cooldown.Behavior != nil {
			if output := c.Cooldown.Behavior(message, client, database, c); len(output) > 0 {
				chat.Say(client, message.Channel, output)
			}
		}

		return
	}

	c.UpdateCooldownAll(message.User)

	var output = c.Behavior(message, client, database, c)
	if len(output) > 0 {
		chat.Say(client, message.Channel, output)
	}
}

// InCooldownGlobal verifica se o comando
// está em cooldown no escopo 'Global'.
func (c *Command) InCooldownGlobal() bool {
	var cooldown = c.Cooldown.globalSave

	return checkCooldown(&cooldown)
}

// InCooldownUser verifica se o comando
// está em cooldown no escopo de 'User'.
func (c *Command) InCooldownUser(user twitch.User) bool {
	var cooldown = c.Cooldown.userSave[user.ID]

	return checkCooldown(&cooldown)
}

// InCooldownAll verifica se o comando
// está em cooldown em ambos escopos 'Global' e 'User'.
func (c *Command) InCooldownAll(user twitch.User) bool {
	return c.InCooldownGlobal() || c.InCooldownUser(user)
}

// UpdateGlobalCooldown atualiza o cooldown do escopo 'Global'.
func (c *Command) UpdateGlobalCooldown(secs uint32) {
	var cooldown = time.Duration(secs)

	c.Cooldown.globalSave = time.Now().Add(time.Second * cooldown)
}

// UpdateUserCooldown atualiza o cooldown do escolo de 'User'.
func (c *Command) UpdateUserCooldown(user twitch.User, secs uint32) {
	var cooldown = time.Duration(secs)

	if len(c.Cooldown.userSave) == 0 {
		c.Cooldown.userSave = make(map[string]time.Time)
	}

	c.Cooldown.userSave[user.ID] = time.Now().Add(time.Second * cooldown)
}

// UpdateCooldownAll atualiza o cooldown em ambos
// escopos 'Global' e 'User', porém utiliza do tempo padrão.
func (c *Command) UpdateCooldownAll(user twitch.User) {
	c.UpdateGlobalCooldown(c.Cooldown.Global)
	c.UpdateUserCooldown(user, c.Cooldown.User)
}

// checkCooldown verifica se o tempo dado já expirou ou não.
func checkCooldown(cooldown *time.Time) bool {
	if !(cooldown.Unix() < time.Now().Unix()) {
		// Ainda em cooldown
		return true
	}

	// Cooldown terminou ou não está na fila.
	return false
}

// CheckPermission verifica se a permissão do usuário
// está de acordo com a permissão informada ao comando.
func (c *Command) CheckPermission(user twitch.User) bool {
	var initial = false

	if len(c.PermissionLevel) == 0 {
		return true
	}

	for _, p := range c.PermissionLevel {
		initial = initial || user.Badges[p] == 1
	}
	return initial
}
