package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	discord "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/nekogravitycat/DiscordAI/internal/config"
	"github.com/nekogravitycat/DiscordAI/internal/pricing"
	"github.com/nekogravitycat/DiscordAI/internal/userdata"
)

var bot *discord.Session
var chats = map[string]*Chat{}

func init() {
	// Load enviroment variables from .env file if exist
	if _, err := os.Stat(".env"); err == nil {
		fmt.Println(".env file found.")
		err = godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file.")
		}
	} else {
		fmt.Println("No .env file found, using system env.")
	}

	// Create ./configs folder if not exist
	if err := os.MkdirAll("configs", os.ModePerm); err != nil {
		log.Fatal("Error creating ./configs folder.")
	}

	// Create ./data folder if not exist
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		log.Fatal("Error creating ./data folder.")
	}

	// Load config
	config.LoadConfig()

	// Load user data
	userdata.LoadUserData()

	// Load pricing table
	pricing.LoadPricingTable()
}

func main() {
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
	registeredCommands := make([]*discord.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := bot.ApplicationCommandCreate(bot.State.User.ID, "", v)
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

	userdata.SaveUserData()

	// Remove commands before shut down
	fmt.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := bot.ApplicationCommandDelete(bot.State.User.ID, "", v.ID)
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

	if _, ok := chats[m.ChannelID]; ok {
		gptReply(s, m)
	}
}
