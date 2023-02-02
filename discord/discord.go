package discord

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IktaS/subscription-tracker/service"
	"github.com/bwmarrin/discordgo"
)

const (
	defaultPrefix = "~"
)

type DiscordBot struct {
	authToken    string
	logChannelID string
	prefix       string
	logMutex     *sync.Mutex
	service      service.Service
	session      *discordgo.Session
	logger       *log.Logger
	store        Store
}

func NewDiscordBot(ctx context.Context, authToken string, store Store) (*DiscordBot, error) {
	bot := &DiscordBot{
		authToken: authToken,
		prefix:    defaultPrefix,
		logger:    log.Default(),
		store:     store,
		logMutex:  &sync.Mutex{},
	}
	err := bot.initDiscordBot()
	if err != nil {
		return nil, err
	}
	logChannel, err := store.GetDefaultLogChannel(ctx)
	if err != nil {
		return nil, err
	}
	bot.logMutex.Lock()
	defer bot.logMutex.Unlock()
	bot.logChannelID = logChannel

	return bot, nil
}

func (b *DiscordBot) initDiscordBot() error {
	var err error
	b.session, err = discordgo.New("Bot " + b.authToken)
	if err != nil {
		return err
	}
	b.session.AddHandler(b.handleMessage)
	b.session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionMessageComponent:
			if h, ok := b.getInteractionHandlers()[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})
	b.session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	return nil
}

func (b *DiscordBot) StartBot() error {
	// Open a websocket connection to Discord and begin listening.
	err := b.session.Open()
	if err != nil {
		b.logger.Println("error opening connection,", err)
		return err
	}
	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	b.session.Close()
	return nil
}

func (b *DiscordBot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !strings.HasPrefix(m.Content, b.prefix) {
		return
	}

	msgContent := m.Content[len(b.prefix):]
	splits := strings.Split(msgContent, " ")
	cmd := splits[0]
	args := splits[1:]
	b.processCommands(cmd, args, m)
}
