package main

import (
	"context"
	"flag"
	"log"

	"github.com/IktaS/subscription-tracker/discord"
	"github.com/IktaS/subscription-tracker/service"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	bot, err := discord.NewDiscordBot(ctx, Token)
	if err != nil {
		log.Fatal(err)
	}
	srv := service.NewService(bot, bot)
	bot.RegisterService(srv)
	bot.StartBot()
}
