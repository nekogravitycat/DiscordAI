package chatbot

import (
	"fmt"

	discord "github.com/bwmarrin/discordgo"
	"github.com/nekogravitycat/DiscordAI/internal/config"
	"github.com/nekogravitycat/DiscordAI/internal/userdata"
)

var (
	commands = []*discord.ApplicationCommand{
		{
			Name:        "activate-gpt",
			Description: "Start ChatGPT on this channel",
			DescriptionLocalizations: &map[discord.Locale]string{
				discord.ChineseTW: "開始在此頻道使用 ChatGPT",
			},
		},
		{
			Name:        "deactivate-gpt",
			Description: "Stop ChatGPT on this channel",
			DescriptionLocalizations: &map[discord.Locale]string{
				discord.ChineseTW: "停止在此頻道使用 ChatGPT",
			},
		},
		{
			Name:        "credits",
			Description: "Check user credits",
			DescriptionLocalizations: &map[discord.Locale]string{
				discord.ChineseTW: "查看用戶使用額度",
			},
		},
	}

	commandHandlers = map[string]func(s *discord.Session, i *discord.InteractionCreate){
		"activate-gpt":   activateGPT,
		"deactivate-gpt": deactivateGPT,
		"credits":        credits,
	}
)

func activateGPT(s *discord.Session, i *discord.InteractionCreate) {
	if _, ok := activeGptChannels[i.ChannelID]; ok {
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
	activeGptChannels[i.ChannelID] = newGptChannel()
	gptChannelData[i.ChannelID] = activeGptChannels[i.ChannelID].GPT
	saveGptChannels()
}

func deactivateGPT(s *discord.Session, i *discord.InteractionCreate) {
	if _, ok := activeGptChannels[i.ChannelID]; !ok {
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
	delete(activeGptChannels, i.ChannelID)
	delete(gptChannelData, i.ChannelID)
	saveGptChannels()
}

func credits(s *discord.Session, i *discord.InteractionCreate) {
	user, ok := userdata.GetUser(i.Message.Author.ID)
	var credits float32
	if ok {
		credits = user.Credit
	} else {
		credits = config.InitCredits
	}

	s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Flags:   discord.MessageFlagsEphemeral,
			Content: fmt.Sprintf("%.5f", credits),
		},
	})
}
