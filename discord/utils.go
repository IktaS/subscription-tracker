package discord

import (
	"fmt"

	"github.com/IktaS/subscription-tracker/entity"
	"github.com/bwmarrin/discordgo"
)

func (b *DiscordBot) generateEmbedTableFromSubs(title string, subs []entity.Subscription) *discordgo.MessageSend {
	embed := &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeRich,
		Title: title,
		//Headers
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Title",
				Inline: true,
			},
			{
				Name:   "\tPayment Method",
				Inline: true,
			},
			{
				Name:   "\tAmount",
				Inline: true,
			},
			{
				Name:   "\tCurrency",
				Inline: true,
			},
		},
	}
	for _, v := range subs {
		embed.Fields = append(embed.Fields, []*discordgo.MessageEmbedField{
			{},
			{
				Value:  v.Title,
				Inline: true,
			},
			{
				Value:  v.PaymentMethod,
				Inline: true,
			},
			{
				Value:  fmt.Sprintf("%.2f", v.Amount.Value),
				Inline: true,
			},
			{
				Value:  v.Amount.Currency,
				Inline: true,
			},
		}...)
	}
	return &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{embed},
	}
}
