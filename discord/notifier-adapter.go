package discord

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IktaS/subscription-tracker/entity"
	"github.com/bwmarrin/discordgo"
)

func (b *DiscordBot) NotifySubsription(ctx context.Context, sub entity.Subscription) error {
	// We create the private channel with the user who sent the message.
	channel, err := b.session.UserChannelCreate(sub.User.ID)
	if err != nil {
		log.Println("error creating channel:", err)
		return err
	}
	msg, err := b.generateDiscordMessageForSubscription(ctx, sub)
	if err != nil {
		return err
	}
	_, err = b.session.ChannelMessageSendComplex(channel.ID, msg)
	return err
}

func (b *DiscordBot) generateDiscordMessageForSubscription(ctx context.Context, sub entity.Subscription) (*discordgo.MessageSend, error) {
	msg := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Type:        discordgo.EmbedTypeRich,
				Description: "Subscription is due",
				Color:       15795975,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Subscription",
						Value:  sub.Title,
						Inline: true,
					},
					{
						Name:   "Payment Method",
						Value:  sub.PaymentMethod,
						Inline: true,
					},
					{
						Name:   "Payment Amount (IDR)",
						Value:  fmt.Sprintf("%.2f", sub.Amount.Value),
						Inline: true,
					},
				},
				Title:     "Subscription Alert",
				Timestamp: time.Now().Format(time.RFC3339),
			},
		},
		Components: []discordgo.MessageComponent{
			SubscriptionNotificationActionButtons,
		},
	}
	return msg, nil
}

var (
	PaidID                                = "sub_notif_yes"
	SubscriptionNotificationActionButtons = discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				// Label is what the user will see on the button.
				Label: "Paid",
				// Style provides coloring of the button. There are not so many styles tho.
				Style: discordgo.SuccessButton,
				// Disabled allows bot to disable some buttons for users.
				Disabled: false,
				// CustomID is a thing telling Discord which data to send when this button will be pressed.
				CustomID: PaidID,
			},
		},
	}
)

func (b *DiscordBot) getInteractionHandlers() map[string]func(*discordgo.Session, *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"sub_notif_yes": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if i.Message != nil && len(i.Message.Embeds) > 0 && len(i.Message.Embeds[0].Fields) > 0 {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseUpdateMessage,
					Data: &discordgo.InteractionResponseData{
						Embeds: []*discordgo.MessageEmbed{
							b.getSubscriptionEmbedAndCreatePaidEmbed(i.Message.Embeds[0].Fields[0]),
						},
						Components: []discordgo.MessageComponent{},
					},
				})
			}
		},
	}
}

func (b *DiscordBot) getSubscriptionEmbedAndCreatePaidEmbed(field *discordgo.MessageEmbedField) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Type:      discordgo.EmbedTypeRich,
		Title:     fmt.Sprintf("Subscription to %s has been paid", field.Value),
		Timestamp: time.Now().Format(time.RFC3339),
		Color:     3553599,
	}
}
