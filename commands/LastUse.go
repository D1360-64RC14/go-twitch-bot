package commands

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/d1360-64rc14/twitch-bot/commands/objects"
	"github.com/d1360-64rc14/twitch-bot/database/tables"
	"github.com/gempir/go-twitch-irc/v2"
	"gorm.io/gorm"
)

var LastUse = objects.Command{
	Name: "LastUse",
	Pattern: regexp.MustCompile(`^!lastuse.*`),
	CaseSensitive: false,
	Cooldown: &objects.Cooldown{Global: 10, User: 0},
	Behavior: LastUseBehavior,
}

// WIP
func LastUseBehavior(message twitch.PrivateMessage,client *twitch.Client, db *gorm.DB, command *objects.Command) string {
	var lastuseTable = tables.LastUse{}

	db.Find(&lastuseTable)

	log.Println(
		fmt.Sprintf("Saída: %d, %s", lastuseTable.Count, lastuseTable.LastTime.Format(time.RFC3339)),
	)

	lastuseTable.Count += 1
	lastuseTable.LastTime = time.Now()

	db.Where("1 = 1").Save(&lastuseTable)

	var msg = fmt.Sprintf("Count: %d; Time: %s", lastuseTable.Count, lastuseTable.LastTime)

	return msg
}
