package discord

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/IktaS/subscription-tracker/entity"
	"github.com/IktaS/subscription-tracker/service"
	"github.com/bwmarrin/discordgo"
)

func (b *DiscordBot) RegisterService(srv service.Service) error {
	b.service = srv
	return nil
}

var (
	cmdHelp = &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title: "Available commands",
				Fields: []*discordgo.MessageEmbedField{
					{
						Value: "set-sub",
					},
					{
						Value: "set-log",
					},
					{
						Value: "set-payday",
					},
					{
						Value: "get-sub-this-payday-cycle",
					},
					{
						Value: "get-all-sub",
					},
				},
			},
		},
	}
)

func (b *DiscordBot) processCommands(cmd string, args []string, info *discordgo.MessageCreate) {
	switch cmd {
	case "help":
		b.session.ChannelMessageSendComplex(info.ChannelID, cmdHelp)
	case "set-sub":
		b.setSubscription(info, args)
	case "set-log":
		b.setLogChannel(info)
	case "set-payday":
		b.setPayday(info, args)
	}
}

func (b *DiscordBot) setLogChannel(info *discordgo.MessageCreate) {
	action := "set this channel as log channel"
	b.logMutex.Lock()
	defer b.logMutex.Unlock()
	b.logChannelID = info.ChannelID
	err := b.store.SetDefaultLogChannel(context.Background(), b.logChannelID)
	if err != nil {
		b.logger.Println(err)
		b.sendFailedCommandResponse(action, info.ChannelID, info.Author.ID, err)
		return
	}
	b.sendSuccessCommandResponse(action, info.ChannelID, info.Author.ID)
}

var (
	setSubscriptionHelp = &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Set Subscription Command",
				Description: "set-sub [{field}=\"{value}\"]",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name: "Available field:",
					},
					{
						Name:  "title",
						Value: "mandatory (string)",
					},
					{
						Name:  "payment_method",
						Value: "mandatory (string)",
					},
					{
						Name:  "amount_currency",
						Value: "mandatory (ISO 4217)",
					},
					{
						Name:  "amount_value",
						Value: "mandatory (float64)",
					},
					{
						Name:  "last_paid",
						Value: "mandatory (RFC 3339)",
					},
					{
						Name:  "duration_value",
						Value: "mandatory (int)",
					},
					{
						Name:  "duration_unit",
						Value: "mandatory ('year', 'month', 'day', 'hour', 'minute')",
					},
				},
			},
		},
	}
)

func (b *DiscordBot) setSubscription(info *discordgo.MessageCreate, args []string) {
	action := "setting subscription"
	if len(args) == 1 && args[0] == "help" {
		b.session.ChannelMessageSendComplex(info.ChannelID, setSubscriptionHelp)
		return
	}
	sub, err := b.parseSetSubscriptionArgs(args)
	if err != nil {
		b.sendFailedCommandResponse(action, info.ChannelID, info.Author.ID, err)
		return
	}
	err = b.service.SetSubscription(context.Background(), sub)
	if err != nil {
		b.sendFailedCommandResponse(action, info.ChannelID, info.Author.ID, err)
		return
	}
	b.sendSuccessCommandResponse(action, info.ChannelID, info.Author.ID)
}

func (b *DiscordBot) parseSetSubscriptionArgs(args []string) (entity.Subscription, error) {
	var sub entity.Subscription
	for _, v := range args {
		strs := strings.Split(v, "=")
		if len(strs) != 2 {
			return sub, errors.New("invalid argument")
		}
		field, value := strs[0], strs[1]
		err := b.parseFieldToSub(&sub, field, value)
		if err != nil {
			return sub, err
		}
	}
	if !sub.IsValid() {
		return sub, errors.New("parsed arguments is not valid")
	}
	return sub, nil
}

func (b *DiscordBot) parseFieldToSub(sub *entity.Subscription, field, value string) error {
	switch field {
	case "title":
		sub.Title = value
	case "payment_method":
		sub.PaymentMethod = value
	case "amount_currency":
		sub.Amount.Currency = value
	case "amount_value":
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		sub.Amount.Value = f
	case "last_paid":
		tm, err := time.Parse(time.RFC3339, value)
		if err != nil {
			return err
		}
		sub.LastPaidDate = tm
	case "duration_value":
		i, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		sub.Duration.Value = i
	case "duration_unit":
		u, err := entity.StringToSubDurationUnit(value)
		if err != nil {
			return err
		}
		sub.Duration.Unit = u
	default:
		return errors.New("unknown field")
	}
	return nil
}

var (
	setPaydayHelp = &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{
			{
				Title:       "Set Payday Command",
				Description: "set-payday [{field}=\"{value}\"]",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name: "Available field:",
					},
					{
						Name:  "time",
						Value: "mandatory (1-31 or 'end')",
					},
				},
			},
		},
	}
)

func (b *DiscordBot) setPayday(info *discordgo.MessageCreate, args []string) {
	action := "setting subscription"
	if len(args) == 1 && args[0] == "help" {
		b.session.ChannelMessageSendComplex(info.ChannelID, setPaydayHelp)
		return
	}
	payday, err := b.parseSetPaydayArgs(args)
	if err != nil {
		b.sendFailedCommandResponse(action, info.ChannelID, info.Author.ID, err)
		return
	}
	err = b.service.SetPaydayTime(context.Background(), entity.User{
		ID:     info.Author.ID,
		Payday: payday,
	})
	if err != nil {
		b.sendFailedCommandResponse(action, info.ChannelID, info.Author.ID, err)
		return
	}
	b.sendSuccessCommandResponse(action, info.ChannelID, info.Author.ID)
}

func (b *DiscordBot) parseSetPaydayArgs(args []string) (entity.Payday, error) {
	for _, v := range args {
		strs := strings.Split(v, "=")
		if len(strs) != 2 {
			return "", errors.New("invalid argument")
		}
		field, value := strs[0], strs[1]
		if field != "time" {
			return "", errors.New("unknown argument field")
		}
		return entity.StringToPayday(value)
	}
	return "", errors.New("failed to find 'time' argument")
}
