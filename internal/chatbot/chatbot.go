package chatbot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	discord "github.com/bwmarrin/discordgo"
	"github.com/nekogravitycat/DiscordAI/internal/config"
	"github.com/nekogravitycat/DiscordAI/internal/userdata"
	"github.com/sashabaranov/go-openai"
)

var openaiClient *openai.Client
var bot *discord.Session
var activeGptChannels = map[string]*gptChannel{"0": newGptChannel()}
var registeredRegularCommands []*discord.ApplicationCommand
var registeredAdminCommands []*discord.ApplicationCommand

func NewOpenaiClient() {
	openaiClient = openai.NewClient(os.Getenv("OPENAI_TOKEN"))
}

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

func removeRemoteCommands(s *discord.Session) {
	servers := []string{""}
	if len(config.AdminServers) > 0 {
		servers = append(servers, config.AdminServers...)
	}

	for _, server := range servers {
		cmds, err := s.ApplicationCommands(s.State.User.ID, server)
		if err != nil {
			log.Fatalf("Could not fetch registered commands for %v: %v", server, err)
		}
		for _, v := range cmds {
			err := s.ApplicationCommandDelete(s.State.User.ID, server, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command for %v: %v", v.Name, server, err)
			} else {
				fmt.Printf("delete '%s' command for %s\n", v.Name, server)
			}
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

	// Remove all remote commands
	fmt.Println("Removing remote commands...")
	removeRemoteCommands(bot)

	// Regular commands
	fmt.Println("Adding regular commands...")
	bot.AddHandler(func(s *discord.Session, i *discord.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})
	if config.NoGlobalCommands {
		fmt.Println("'no-global-commands' enabled, adding regular commands to admin servers.")
		if len(config.AdminServers) > 0 {
			fmt.Print("Server list: ")
			for _, s := range config.AdminServers {
				fmt.Printf("%s ", s)
				addCommands(&registeredRegularCommands, regularCommands, s)
			}
			fmt.Println()
		} else {
			fmt.Println("Empty admin server list, regular commands not registered.")
		}
	} else {
		addCommands(&registeredRegularCommands, regularCommands, "")
	}

	// Admin commands
	if len(config.AdminServers) > 0 {
		fmt.Println("Adding admin commands...")
		bot.AddHandler(func(s *discord.Session, i *discord.InteractionCreate) {
			if handler, ok := adminCommandHandlers[i.ApplicationCommandData().Name]; ok {
				handler(s, i)
			}
		})
		fmt.Print("Server list: ")
		for _, s := range config.AdminServers {
			fmt.Printf("%s ", s)
			addCommands(&registeredAdminCommands, adminCommands, s)
		}
		fmt.Println()
	} else {
		fmt.Println("Empty admin server list, admin commands not registered.")
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
	fmt.Println("Removing regular commands...")
	if config.NoGlobalCommands {
		if len(config.AdminServers) > 0 {
			for _, s := range config.AdminServers {
				removeCommands(registeredRegularCommands, s)
			}
		}
	} else {
		removeCommands(registeredRegularCommands, "")
	}

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
