package discord

import (
	"context"
	"fmt"

	"github.com/IktaS/subscription-tracker/service"
	"github.com/bwmarrin/discordgo"
)

func (b *DiscordBot) RegisterService(srv service.Service) error {
	b.service = srv
	return nil
}

func (b *DiscordBot) processCommands(cmd string, args []string, info *discordgo.MessageCreate) {
	switch cmd {
	case "ping":
		b.doPong(info)
	case "set-log":
		b.setLogChannel(info)
	}
}

func (b *DiscordBot) doPong(info *discordgo.MessageCreate) {
	if b.service == nil {
		b.logger.Print("no service available to doPong")
		return
	}
	b.service.PingUser(context.Background(), info.Author.ID)
}

func (b *DiscordBot) setLogChannel(info *discordgo.MessageCreate) {
	b.logMutex.Lock()
	defer b.logMutex.Unlock()
	b.logChannelID = info.ChannelID
	b.session.ChannelMessageSend(b.logChannelID, fmt.Sprintf("Successfully set this channel as log channel <@!%s>", info.Author.ID))
}