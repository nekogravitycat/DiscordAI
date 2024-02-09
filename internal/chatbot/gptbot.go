package chatbot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	discord "github.com/bwmarrin/discordgo"
	"github.com/nekogravitycat/DiscordAI/internal/config"
	"github.com/nekogravitycat/DiscordAI/internal/gpt"
	"github.com/nekogravitycat/DiscordAI/internal/pricing"
	"github.com/nekogravitycat/DiscordAI/internal/userdata"
)

type gptChannel struct {
	GPT   gpt.GPT
	queue []*discord.MessageCreate
}

func newGptChannel() *gptChannel {
	chat := &gptChannel{
		GPT:   gpt.NewGPT(),
		queue: []*discord.MessageCreate{},
	}
	return chat
}

func gptReply(s *discord.Session, m *discord.MessageCreate) {
	user, ok := userdata.GetUser(m.Author.ID)
	if !ok {
		userdata.SetUser(m.Author.ID, userdata.NewUserInfo())
	}

	if gpt.CountToken(m.Content, user.Model) > config.GPT.Limits.PromptTokens {
		messageReply(s, m, "Prompt too long.")
		return
	}

	activeGptChannels[m.ChannelID].queueMessage(m)
}

func (c *gptChannel) queueMessage(m *discord.MessageCreate) {
	c.queue = append(c.queue, m)
	if len(c.queue) == 1 {
		c.replyNext()
	}
}

func (c *gptChannel) replyNext() {
	if len(c.queue) <= 0 {
		return
	}

	m := c.queue[0]

	// Remove prefixs
	trim1 := strings.TrimPrefix(m.Content, "!")
	trim2 := strings.TrimPrefix(trim1, "！")

	c.GPT.AddMessage(trim2)

	if m.Content == trim2 {
		user, _ := userdata.GetUser(m.Author.ID)

		if user.Credit <= 0 {
			messageReply(bot, m, "Not enough credits.")

		} else if !user.HasPrivilege(user.Model) {
			messageReply(bot, m, "Permission denied. Please switch to other models.")

		} else {
			bot.ChannelTyping(m.ChannelID)
			reply, usage, _ := c.GPT.Generate(user.Model, m.Author.ID)

			// Segment the reply if its longer than 2000 characters
			for len(reply) > 2000 {
				messageReply(bot, m, reply[:2000])
				reply = reply[2000:]
			}

			messageReply(bot, m, reply)

			// Update userdata to follow up possible simultaneous operations
			user, _ := userdata.GetUser(m.Author.ID)
			user.Credit -= pricing.GetGPTCost(user.Model, usage)
			userdata.SetUser(m.Author.ID, user)
			fmt.Printf("User credits: %f\n", user.Credit)
			userdata.SaveUserData()
		}
	}
	// else stack prompts if they start with "!"

	c.queue = c.queue[1:]
	c.replyNext()
}

var gptChannelData = map[string]gpt.GPT{}

const GPTCHANNELFILE = "./data/gptchannels.json"

func LoadGptChannels() {
	if _, err := os.Stat(GPTCHANNELFILE); errors.Is(err, os.ErrNotExist) {
		fmt.Println("No gptchannels.json found, creating one.")
		gptChannelData["0"] = gpt.NewGPT()
		saveGptChannels()
	}

	jsonFile, err := os.Open(GPTCHANNELFILE)
	if err != nil {
		fmt.Println("Error reading gptchannels.json")
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading bytes of gptchannels.json")
	}

	err = json.Unmarshal(byteValue, &gptChannelData)
	if err != nil {
		fmt.Println("Error parsing gptchannels.json into GPT struct.")
	}
}

func saveGptChannels() {
	jsonFile, err := os.Create(GPTCHANNELFILE)
	if err != nil {
		fmt.Println("Error writing gptchannels.json")
		return
	}
	defer jsonFile.Close()

	jsonData, err := json.MarshalIndent(gptChannelData, "", "  ")
	if err != nil {
		fmt.Println("Error parsing Users struct into json data.")
		return
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing gptchannels.json file.")
	}
}
