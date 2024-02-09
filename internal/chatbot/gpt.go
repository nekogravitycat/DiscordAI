package chatbot

import (
	"fmt"
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
		s.ChannelMessageSendReply(m.ChannelID, "Prompt too long.", m.Reference())
		return
	}

	activeGptChannels[m.ChannelID].QueueMessage(m)
}

func (c *gptChannel) QueueMessage(m *discord.MessageCreate) {
	c.queue = append(c.queue, m)
	if len(c.queue) == 1 {
		c.replyNext()
	}
}

func handleReplyError(err error) {
	if err != nil {
		fmt.Println("Error replying: " + err.Error())
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
			bot.ChannelMessageSendReply(m.ChannelID, "Not enough credits.", m.Reference())

		} else {
			bot.ChannelTyping(m.ChannelID)
			reply, usage, _ := c.GPT.Generate("gpt-3.5-turbo", m.Author.ID)

			// Segment the reply if its longer than 2000 characters
			for len(reply) > 2000 {
				_, err := bot.ChannelMessageSendReply(m.ChannelID, reply[:2000], m.Reference())
				handleReplyError(err)
				reply = reply[2000:]
			}

			_, err := bot.ChannelMessageSendReply(m.ChannelID, reply, m.Reference())
			handleReplyError(err)

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
