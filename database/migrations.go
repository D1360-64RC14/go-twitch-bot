package database

import (
	"log"

	"github.com/d1360-64rc14/twitch-bot/database/tables"
	"gorm.io/gorm"
)

func Migrations(db *gorm.DB) {
	// Cria a tabela 'LastUse' caso não exista.
	if !db.Migrator().HasTable(&tables.LastUse{}) {
		var err = db.Migrator().CreateTable(&tables.LastUse{})
		if err != nil {
			log.Fatal(err)
		}
	}

	// Adiciona um elemento à tabela caso não tenha.
	// As colunas 'count' e 'last_time' já estão definidas
	// como auto create '0' e 'time.Now()', respectivamente.
	if db.Take(&tables.LastUse{}).Error == gorm.ErrRecordNotFound {
		db.Create(&tables.LastUse{})
	}
}