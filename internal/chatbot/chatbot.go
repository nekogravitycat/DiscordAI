package chatbot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	discord "github.com/bwmarrin/discordgo"
	"github.com/nekogravitycat/DiscordAI/internal/userdata"
)

var bot *discord.Session
var activeGptChannels = map[string]*gptChannel{"0": newGptChannel()}
var registeredCommands []*discord.ApplicationCommand

func Run() {
	var err error

	bot, err = discord.New("Bot " + os.Getenv("DISCORDBOT_TOKEN"))
	if err != nil {
		log.Fatal("Error creating Discord session. " + err.Error())
	}

	bot.Identify.Intents = discord.IntentsAll

	err = bot.Open()
	if err != nil {
		log.Fatal("Error opening connection. " + err.Error())
	}
	defer bot.Close()

	bot.AddHandler(func(s *discord.Session, i *discord.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	fmt.Println("Adding commands...")
	registeredCommands = make([]*discord.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := bot.ApplicationCommandCreate(bot.State.User.ID, "987988090528366602", v)
		if err != nil {
			fmt.Println("Error creating command: " + v.Name)
			fmt.Println(err.Error())
		}
		registeredCommands[i] = cmd
	}

	bot.AddHandler(messageCreate)

	// Run until receive CTRL-C signal
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
}

func Stop() {
	userdata.SaveUserData()

	// Remove commands before shut down
	fmt.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := bot.ApplicationCommandDelete(bot.State.User.ID, "987988090528366602", v.ID)
		if err != nil {
			fmt.Println("Error deleting slash command: " + v.Name)
		}
	}

	fmt.Println("Bot shut down gracefully.")
}

// Discord handlers

func messageCreate(s *discord.Session, m *discord.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	if strings.HasPrefix(m.Content, "#") || strings.HasPrefix(m.Content, "＃") {
		return
	}

	if isActiveGptChannel(m.ChannelID) {
		gptReply(s, m)
	}
}
