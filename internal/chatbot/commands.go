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
		{
			Name:        "gpt-sys-prompt",
			Description: "Show GPT system prompt for this channel",
			DescriptionLocalizations: &map[discord.Locale]string{
				discord.ChineseTW: "查看此頻道的 GPT 系統設定",
			},
		},
		{
			Name:        "set-gpt-sys-prompt",
			Description: "Set GPT system prompt for this channel",
			DescriptionLocalizations: &map[discord.Locale]string{
				discord.ChineseTW: "更改此頻道的 GPT 系統設定",
			},
			Options: []*discord.ApplicationCommandOption{
				{
					Type:        discord.ApplicationCommandOptionString,
					Name:        "sys-prompt",
					Description: "The system prompt of GPT",
					Required:    true,
				},
			},
		},
		{
			Name:        "reset-gpt-sys-prompt",
			Description: "Reset GPT system prompt for this channel to default",
			DescriptionLocalizations: &map[discord.Locale]string{
				discord.ChineseTW: "還原此頻道的 GPT 系統設定至預設值",
			},
		},
	}

	commandHandlers = map[string]func(s *discord.Session, i *discord.InteractionCreate){
		"activate-gpt":         activateGPT,
		"deactivate-gpt":       deactivateGPT,
		"credits":              credits,
		"gpt-sys-prompt":       showGptSysPrompt,
		"set-gpt-sys-prompt":   setGptSysPrompt,
		"reset-gpt-sys-prompt": resetGptSysPrompt,
	}
)

func activateGPT(s *discord.Session, i *discord.InteractionCreate) {
	if isActiveGptChannel(i.ChannelID) {
		interactionRespond(s, i, "Already in chat")
		return
	}

	interactionRespond(s, i, "Hello!")
	activeGptChannels[i.ChannelID] = newGptChannel()
	saveGptChannels()
}

func deactivateGPT(s *discord.Session, i *discord.InteractionCreate) {
	if !isActiveGptChannel(i.ChannelID) {
		interactionRespond(s, i, "Not in channel")
		return
	}

	interactionRespond(s, i, "Bye!")
	delete(activeGptChannels, i.ChannelID)
	saveGptChannels()
}

func credits(s *discord.Session, i *discord.InteractionCreate) {
	user, ok := userdata.GetUser(i.Member.User.ID)
	var credits float32 = 0

	if ok {
		credits = user.Credit
	} else {
		credits = config.InitCredits
	}

	interactionRespondEphemeral(s, i, fmt.Sprintf("Your credits: `$%.5f USD`", credits))
}

func showGptSysPrompt(s *discord.Session, i *discord.InteractionCreate) {
	if !isActiveGptChannel(i.ChannelID) {
		interactionRespond(s, i, "This is not an active GPT channel.")
		return
	}

	interactionRespond(s, i, fmt.Sprintf("System prompt:\n```%s```", activeGptChannels[i.ChannelID].GPT.SysPrompt))
}

func setGptSysPrompt(s *discord.Session, i *discord.InteractionCreate) {
	if !isActiveGptChannel(i.ChannelID) {
		interactionRespond(s, i, "This is not an active GPT channel.")
		return
	}

	input := i.ApplicationCommandData().Options[0].Value
	if value, ok := input.(string); ok {
		activeGptChannels[i.ChannelID].GPT.SysPrompt = value
		saveGptChannels()
		interactionRespond(s, i, fmt.Sprintf("System prompt update:\n```%s```", value))
	} else {
		interactionRespond(s, i, "Invaild input.")
	}
}

func resetGptSysPrompt(s *discord.Session, i *discord.InteractionCreate) {
	if !isActiveGptChannel(i.ChannelID) {
		interactionRespond(s, i, "This is not an active GPT channel.")
		return
	}

	activeGptChannels[i.ChannelID].GPT.SysPrompt = config.GPT.DefaultSysPrompt
	saveGptChannels()
	interactionRespond(s, i, fmt.Sprintf("System prompt update:\n```%s```", activeGptChannels[i.ChannelID].GPT.SysPrompt))
}
