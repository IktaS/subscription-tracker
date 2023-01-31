package discord

import (
	"context"
	"log"
	"time"

	"github.com/Necroforger/dgwidgets"
	"github.com/bwmarrin/discordgo"
)

// const lib = require('lib')({token: process.env.STDLIB_SECRET_TOKEN});

// await lib.discord.channels['@0.3.2'].messages.create({
//   "channel_id": `${context.params.event.channel_id}`,
//   "content": `Subscription alert /@ subscription owner`,
//   "tts": false,
//   "components": [
//     {
//       "type": 1,
//       "components": [
//         {
//           "style": 1,
//           "label": `Noted`,
//           "custom_id": `ok`,
//           "disabled": false,
//           "emoji": {
//             "id": null,
//             "name": `✅`
//           },
//           "type": 2
//         }
//       ]
//     }
//   ],
//   "embeds": [
//     {
//       "type": "rich",
//       "title": `[Subscription Name] is due in [x] days`,
//       "description": `Additional subscription info`,
//       "color": 0x00FFFF
//     }
//   ]
// });

func (b *DiscordBot) PingUser(ctx context.Context, userID string) error {
	// We create the private channel with the user who sent the message.
	channel, err := b.session.UserChannelCreate(userID)
	if err != nil {
		log.Println("error creating channel:", err)
		return err
	}
	p := dgwidgets.NewWidget(b.session, channel.ID, &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Description: "Subscription is due",
		Color:       15795975,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Subscription",
				Value:  "Northernlion Twitch",
				Inline: true,
			},
			{
				Name:   "Payment Method",
				Value:  "Jago Sub",
				Inline: true,
			},
			{
				Name:   "Payment Amount",
				Value:  "20000",
				Inline: true,
			},
		},
		Title:     "Subscription Alert",
		Timestamp: time.Now().Format(time.RFC3339),
	})
	p.Handle("✅", func(w *dgwidgets.Widget, mr *discordgo.MessageReaction) {
		w.UpdateEmbed(&discordgo.MessageEmbed{
			Type:      discordgo.EmbedTypeRich,
			Title:     "Subscription to Northernlion Twitch has been paid",
			Timestamp: time.Now().Format(time.RFC3339),
			Color:     3553599,
		})
	})
	p.Spawn()
	// Then we send the message through the channel we created.
	_, err = b.session.ChannelMessageSend(channel.ID, "Pong!")
	if err != nil {
		log.Println("error creating channel:", err)
		return err
	}
	return nil
}
