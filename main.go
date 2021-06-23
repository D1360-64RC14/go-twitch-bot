package main

import (
	"log"
	"os"

	"github.com/d1360-64rc14/twitch-bot/chat"
	"github.com/d1360-64rc14/twitch-bot/commands"
	"github.com/d1360-64rc14/twitch-bot/database"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	log.Println("--- INICIALIZAÇÃO ---")

	var err = godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(".env carregado.")

	// ----- CONEXÃO COM BASE DE DADOS ----- //
	var db *gorm.DB
	db, err = database.NewDatabase(database.Options{
		Host:     os.Getenv("DATABASE_HOST")    ,
		Username: os.Getenv("DATABASE_USERNAME"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		DBName:   os.Getenv("DATABASE_DBNAME")  ,
		Port:     os.Getenv("DATABASE_PORT")    ,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Conexão com base de dados realizada.")

	database.Migrations(db)
	log.Println("Migrations verificadas.")

	// ----- CONEXÃO COM CHAT DA TWITCH ----- //
	var client *twitch.Client
	client = chat.NewChat(chat.Options{
		Username: os.Getenv("TWITCH_USERNAME"),
		OAuth:    os.Getenv("TWITCH_OAUTH")   ,
		Channel:  os.Getenv("TWITCH_CHANNEL") ,
	})
	log.Println("Inicialização do chat da twitch realizada.")

	commands.Handler(client, db)

	log.Print("--- INICIADO COM SUCESSO! ---\n\n")

	err = client.Connect()
	if err != nil {
		log.Fatalf("Não foi possível se conectar ao chat da Twitch:\n%s\n", err)
	}
}
