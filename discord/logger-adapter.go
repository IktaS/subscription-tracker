package discord

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Necroforger/dgwidgets"
	"github.com/bwmarrin/discordgo"
)

func (b *DiscordBot) Write(p []byte) (int, error) {
	if b.logChannelID == "" {
		return -1, ErrNoLogChannel
	}
	_, err := b.session.ChannelMessageSend(b.logChannelID, string(p))
	if err != nil {
		return -1, err
	}
	return len(p), nil
}

func (b *DiscordBot) Info(ctx context.Context, msg string, args ...interface{}) {
	if b.logChannelID == "" {
		return
	}
	format := "INFO: %s"
	color := 8558253
	var fields []*discordgo.MessageEmbedField
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprint(args[i]),
			Value: fmt.Sprint(args[i+1]),
		})
	}
	p := dgwidgets.NewWidget(b.session, b.logChannelID, &discordgo.MessageEmbed{
		Type:      discordgo.EmbedTypeRich,
		Color:     color,
		Fields:    fields,
		Title:     fmt.Sprintf(format, msg),
		Timestamp: time.Now().Format(time.RFC3339),
	})
	p.Spawn()
}

func (b *DiscordBot) Warning(ctx context.Context, msg string, args ...interface{}) {
	if b.logChannelID == "" {
		return
	}
	format := "WARNING: %s"
	color := 16776960
	var fields []*discordgo.MessageEmbedField
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprint(args[i]),
			Value: fmt.Sprint(args[i+1]),
		})
	}
	p := dgwidgets.NewWidget(b.session, b.logChannelID, &discordgo.MessageEmbed{
		Type:      discordgo.EmbedTypeRich,
		Color:     color,
		Fields:    fields,
		Title:     fmt.Sprintf(format, msg),
		Timestamp: time.Now().Format(time.RFC3339),
	})
	p.Spawn()
}

func (b *DiscordBot) Error(ctx context.Context, msg string, args ...interface{}) {
	if b.logChannelID == "" {
		return
	}
	format := "ERROR: %s"
	color := 16711680
	var fields []*discordgo.MessageEmbedField
	for i := 0; i < len(args); i += 2 {
		if i+1 >= len(args) {
			break
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprint(args[i]),
			Value: fmt.Sprint(args[i+1]),
		})
	}
	p := dgwidgets.NewWidget(b.session, b.logChannelID, &discordgo.MessageEmbed{
		Type:      discordgo.EmbedTypeRich,
		Color:     color,
		Fields:    fields,
		Title:     fmt.Sprintf(format, msg),
		Timestamp: time.Now().Format(time.RFC3339),
	})
	p.Spawn()
}

func (b *DiscordBot) NewDiscordBotLogger() *log.Logger {
	return log.New(b, "", log.Ldate|log.Ltime|log.Llongfile)
}
