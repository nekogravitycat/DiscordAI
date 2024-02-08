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
	"github.com/nekogravitycat/DiscordAI/internal/chatgpt"
	"github.com/nekogravitycat/DiscordAI/internal/config"
	"github.com/nekogravitycat/DiscordAI/internal/user"
)

type Chat struct {
	GPT   chatgpt.GPT
	queue []*discord.MessageCreate
}

func NewChat() *Chat {
	chat := &Chat{
		GPT:   chatgpt.NewGPT(),
		queue: []*discord.MessageCreate{},
	}
	return chat
}

func (c *Chat) QueueMessage(m *discord.MessageCreate) {
	c.queue = append(c.queue, m)
	if len(c.queue) == 1 {
		c.replyNext()
	}
}

func (c *Chat) replyNext() {
	if len(c.queue) <= 0 {
		return
	}

	m := c.queue[0]

	// Remove prefixs
	trim1 := strings.TrimPrefix(m.Content, "!")
	trim2 := strings.TrimPrefix(trim1, "！")

	c.GPT.AddMessage(trim2)

	if m.Content == trim2 {
		reply, _, _ := c.GPT.Generate("gpt-3.5-turbo", m.Author.ID)
		bot.ChannelTyping(m.ChannelID)
		bot.ChannelMessageSendReply(m.ChannelID, reply, m.Reference())
	}
	// else stack prompts if they start with "!"

	c.queue = c.queue[1:]
	c.replyNext()
}

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
	user.LoadUserData()
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

	user.SaveUserData()

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

func gptReply(s *discord.Session, m *discord.MessageCreate) {
	chats[m.ChannelID].QueueMessage(m)
}

// Slash commands

func startChat(s *discord.Session, i *discord.InteractionCreate) {
	if _, ok := chats[i.ChannelID]; ok {
		s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
			Type: discord.InteractionResponseChannelMessageWithSource,
			Data: &discord.InteractionResponseData{
				Content: "Already in chat",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Content: "Hello!",
		},
	})
	chats[i.ChannelID] = NewChat()
}

func stopChat(s *discord.Session, i *discord.InteractionCreate) {
	if _, ok := chats[i.ChannelID]; !ok {
		s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
			Type: discord.InteractionResponseChannelMessageWithSource,
			Data: &discord.InteractionResponseData{
				Content: "Not in channel",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Content: "Bye!",
		},
	})
	delete(chats, i.ChannelID)
}

var (
	commands = []*discord.ApplicationCommand{
		{
			Name:        "start",
			Description: "Start ChatGPT on this channel",
		},
		{
			Name:        "stop",
			Description: "Stop ChatGPT on this channel",
		},
	}

	commandHandlers = map[string]func(s *discord.Session, i *discord.InteractionCreate){
		"start": startChat,
		"stop":  stopChat,
	}
)
