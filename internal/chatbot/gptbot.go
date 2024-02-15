package chatbot

import (
	"fmt"
	"strings"

	discord "github.com/bwmarrin/discordgo"
	"github.com/nekogravitycat/DiscordAI/internal/config"
	"github.com/nekogravitycat/DiscordAI/internal/gpt"
	"github.com/nekogravitycat/DiscordAI/internal/jsondata"
	"github.com/nekogravitycat/DiscordAI/internal/pricing"
	"github.com/nekogravitycat/DiscordAI/internal/userdata"
)

type gptChannel struct {
	GPT   gpt.GPT `json:"gpt"`
	queue []*discord.MessageCreate
}

func newGptChannel() *gptChannel {
	chat := &gptChannel{
		GPT:   gpt.NewGPT(openaiClient),
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
		messageReply(s, m.Message, "Prompt too long")
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

	// Add image if exists
	if len(m.Attachments) > 0 {
		url := m.Attachments[0].URL
		if gpt.IsImageUrl(url) {
			c.GPT.AddImage(url, "auto")
		}
	}

	if m.Content == trim2 {
		user, _ := userdata.GetUser(m.Author.ID)

		if user.Credit <= 0 {
			messageReply(bot, m.Message, "Not enough credits.")

		} else if !user.HasModelPrivilege(user.Model) {
			messageReply(bot, m.Message, modelPermissionDeniedMessage)

		} else {
			bot.ChannelTyping(m.ChannelID)
			fmt.Printf("Model: %s, User: %s\n", user.Model, m.Author.ID)
			reply, usage, _ := c.GPT.Generate(user.Model, m.Author.ID)

			// Segment the reply if its longer than 2000 characters
			for len(reply) > 2000 {
				messageReply(bot, m.Message, reply[:2000])
				reply = reply[2000:]
			}

			messageReply(bot, m.Message, reply)
			fmt.Printf("Usage: %d (`$%f USD`)\n", usage.TotalTokens, pricing.GetGPTCost(user.Model, usage))

			// Update userdata to follow up possible simultaneous operations
			user, _ := userdata.GetUser(m.Author.ID)
			user.Credit -= pricing.GetGPTCost(user.Model, usage)
			userdata.SetUser(m.Author.ID, user)
			userdata.SaveUserData()
		}
	}
	// else stack prompts if they start with "!"

	c.queue = c.queue[1:]
	c.replyNext()
}

func isActiveGptChannel(channelID string) bool {
	_, isActive := activeGptChannels[channelID]
	return isActive
}

const GPTCHANNELFILE = "./data/gptchannels.json"

func LoadGptChannels() {
	fmt.Println("Loading GPT channels...")
	jsondata.Check(GPTCHANNELFILE, activeGptChannels)
	jsondata.Load(GPTCHANNELFILE, &activeGptChannels)

	if len(activeGptChannels) == 0 {
		return
	}

	fmt.Print("Channel list: ")
	for chId, ch := range activeGptChannels {
		sysPromptData := ch.GPT.SysPrompt
		activeGptChannels[chId] = newGptChannel()
		activeGptChannels[chId].GPT.SysPrompt = sysPromptData
		fmt.Printf("%s, ", chId)
	}
	fmt.Println()
}

func saveGptChannels() {
	jsondata.Save(GPTCHANNELFILE, activeGptChannels)
}
