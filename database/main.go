package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Options struct {
	Host     string
	Username string
	Password string
	DBName   string
	Port     string
}

// NewDatabase cria e faz a conexão com a base de dados.
func NewDatabase(dbopt Options) (*gorm.DB, error) {
	var dsn = fmt.Sprintf(
		"host=%s port=%s"            +
		" user=%s password=%s"       +
		" dbname=%s"                 +
		" sslmode=disable"           +
		" TimeZone=America/Sao_Paulo",
		dbopt.Host,
		dbopt.Port,
		dbopt.Username,
		dbopt.Password,
		dbopt.DBName,
	)

	var db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil,
			fmt.Errorf("Não foi possível se conectar à base de dados:\n%s\n", err)
	}

	return db, nil
}
