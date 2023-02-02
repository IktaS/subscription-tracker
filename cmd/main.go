package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/IktaS/subscription-tracker/discord"
	"github.com/IktaS/subscription-tracker/forex"
	"github.com/IktaS/subscription-tracker/service"
	"github.com/IktaS/subscription-tracker/store/sqlite"
	"github.com/joho/godotenv"
)

// Variables used for command line parameters
var (
	Env string
)

func init() {
	flag.StringVar(&Env, "e", "", "Env file")
	flag.Parse()
}

func main() {
	if Env == "" {
		Env = ".env"
	}
	err := godotenv.Load(Env)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()

	store, err := sqlite.NewSQLiteStore(os.Getenv("SQLITEDB"))
	if err != nil {
		log.Fatal(err)
	}

	bot, err := discord.NewDiscordBot(ctx, os.Getenv("DISCORDKEY"), store)
	if err != nil {
		log.Fatal(err)
	}

	expr, err := strconv.Atoi(os.Getenv("FOREXEXPR"))
	if err != nil {
		log.Fatal(err)
	}
	forex := forex.NewForexService(os.Getenv("FOREXKEY"), expr, bot.NewDiscordBotLogger())
	srv := service.NewService(bot, bot, forex)
	bot.RegisterService(srv)
	bot.StartBot()
	store.Shutdown()
}
