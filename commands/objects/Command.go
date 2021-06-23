package objects

import (
	"regexp"
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

func (c *Command) Validate(message twitch.PrivateMessage) bool {
	return c.Pattern.MatchString(message.Message)
}

func (c *Command) Exec(message twitch.PrivateMessage, client *twitch.Client, database *gorm.DB) {
	if c.OnCooldownAll(message.User) {
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

// TO;DO : CaseSensitive
func (c *Command) OnCooldownGlobal() bool {
	var cooldown = c.Cooldown.globalSave

	return checkCooldown(&cooldown)
}

func (c *Command) OnCooldownUser(user twitch.User) bool {
	var cooldown = c.Cooldown.userSave[user.ID]

	return checkCooldown(&cooldown)
}

func (c *Command) OnCooldownAll(user twitch.User) bool {
	return c.OnCooldownGlobal() || c.OnCooldownUser(user)
}

func (c *Command) UpdateGlobalCooldown(secs uint32) {
	var cooldown = time.Duration(secs)

	c.Cooldown.globalSave = time.Now().Add(time.Second * cooldown)
}

func (c *Command) UpdateUserCooldown(user twitch.User, secs uint32) {
	var cooldown = time.Duration(secs)

	if len(c.Cooldown.userSave) == 0 {
		c.Cooldown.userSave = make(map[string]time.Time)
	}

	c.Cooldown.userSave[user.ID] = time.Now().Add(time.Second * cooldown)
}

func (c *Command) UpdateCooldownAll(user twitch.User) {
	c.UpdateGlobalCooldown(c.Cooldown.Global)
	c.UpdateUserCooldown(user, c.Cooldown.User)
}

func checkCooldown(cooldown *time.Time) bool {
	if !(cooldown.Unix() < time.Now().Unix()) {
		// Ainda em cooldown
		return true
	}

	// Cooldown terminou ou não está na fila.
	return false
}

// Talvez seja utilizado.
func NoBehavior(_ twitch.PrivateMessage, _ *twitch.Client, _ *gorm.DB, _ *Command) string {
	return ""
}