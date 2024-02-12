package chatbot

import (
	"fmt"
	"log"
	"os"
	"strings"

	discord "github.com/bwmarrin/discordgo"
	"github.com/nekogravitycat/DiscordAI/internal/config"
	"github.com/nekogravitycat/DiscordAI/internal/userdata"
)

var bot *discord.Session
var activeGptChannels = map[string]*gptChannel{"0": newGptChannel()}
var registeredRegularCommands []*discord.ApplicationCommand
var registeredAdminCommands []*discord.ApplicationCommand

func addCommands(registerList *[]*discord.ApplicationCommand, cmdToAdd []*discord.ApplicationCommand, targetServerId string) {
	*registerList = make([]*discord.ApplicationCommand, len(cmdToAdd))

	for i, c := range cmdToAdd {
		cmd, err := bot.ApplicationCommandCreate(bot.State.User.ID, targetServerId, c)
		if err != nil {
			fmt.Println("Error creating command: " + c.Name)
			fmt.Println(err.Error())
		}
		(*registerList)[i] = cmd
	}
}

func removeCommands(registerList []*discord.ApplicationCommand, targetServerId string) {
	for _, cmd := range registerList {
		err := bot.ApplicationCommandDelete(bot.State.User.ID, targetServerId, cmd.ID)
		if err != nil {
			fmt.Println("Error deleting slash command: " + cmd.Name)
		}
	}
}

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

	fmt.Println("Adding commands...")
	bot.AddHandler(func(s *discord.Session, i *discord.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})
	addCommands(&registeredRegularCommands, regularCommands, "")

	if len(config.AdminServers) > 0 {
		fmt.Println("Adding admin commands...")
		bot.AddHandler(func(s *discord.Session, i *discord.InteractionCreate) {
			if handler, ok := adminCommandHandlers[i.ApplicationCommandData().Name]; ok {
				handler(s, i)
			}
		})
		for _, s := range config.AdminServers {
			addCommands(&registeredAdminCommands, adminCommands, s)
		}
	} else {
		fmt.Println("Empty admin server list, admin commands not registered.")
	}

	bot.AddHandler(messageCreate)

	fmt.Println("Bot is now running. Type 'stop' and hit enter to exit.")
	var input string = ""
	for input != "stop" {
		fmt.Scanln(&input)
	}
}

func Stop() {
	userdata.SaveUserData()

	// Remove commands before shut down
	fmt.Println("Removing commands...")
	removeCommands(registeredRegularCommands, "")

	if len(config.AdminServers) > 0 {
		fmt.Println("Removing admin commands...")
		for _, s := range config.AdminServers {
			removeCommands(registeredAdminCommands, s)
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

const notActiveGptChannelMessage = "This is not an active GPT channel. Use `/activate-gpt` to activate GPT for this channel."
const modelPermissionDeniedMessage = "Permission denied. Please switch to other models by `/set-gpt-model [model name]`"
