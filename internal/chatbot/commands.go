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
		{
			Name:        "set-gpt-model",
			Description: "Set the GPT model for the user",
			DescriptionLocalizations: &map[discord.Locale]string{
				discord.ChineseTW: "設定用戶使用的 GPT 模型",
			},
			Options: []*discord.ApplicationCommandOption{
				{
					Name:        "model",
					Description: "The GPT model to use",
					Type:        discord.ApplicationCommandOptionString,
					Required:    true,
					Choices: []*discord.ApplicationCommandOptionChoice{
						{
							Name:  "gpt-3.5-turbo",
							Value: "gpt-3.5-turbo",
						},
						{
							Name:  "gpt-4-turbo-preview",
							Value: "gpt-4-turbo-preview",
						},
						{
							Name:  "gpt-4-vision-preview",
							Value: "gpt-4-vision-preview",
						},
					},
				},
			},
		},
		{
			Name:        "clear-gpt-history",
			Description: "Clear GPT chat history for this channel",
			DescriptionLocalizations: &map[discord.Locale]string{
				discord.ChineseTW: "清除此頻道的 GPT 聊天歷史",
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
		"set-gpt-model":        setGptModel,
		"clear-gpt-history":    clearGptHistory,
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
		interactionRespond(s, i, notActiveGptChannelMessage)
		return
	}

	interactionRespond(s, i, fmt.Sprintf("System prompt:\n```%s```", activeGptChannels[i.ChannelID].GPT.SysPrompt))
}

func setGptSysPrompt(s *discord.Session, i *discord.InteractionCreate) {
	if !isActiveGptChannel(i.ChannelID) {
		interactionRespond(s, i, notActiveGptChannelMessage)
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
		interactionRespond(s, i, notActiveGptChannelMessage)
		return
	}

	activeGptChannels[i.ChannelID].GPT.SysPrompt = config.GPT.DefaultSysPrompt
	saveGptChannels()
	interactionRespond(s, i, fmt.Sprintf("System prompt update:\n```%s```", activeGptChannels[i.ChannelID].GPT.SysPrompt))
}

func setGptModel(s *discord.Session, i *discord.InteractionCreate) {
	input := i.ApplicationCommandData().Options[0].Value

	if value, ok := input.(string); ok {
		user, userExist := userdata.GetUser(i.Member.User.ID)
		if !userExist {
			user = userdata.SetUser(i.Member.User.ID, userdata.NewUserInfo())
		}

		if user.HasPrivilege(value) {
			user.Model = value
			userdata.SetUser(i.Member.User.ID, user)
			userdata.SaveUserData()
			interactionRespondEphemeral(s, i, fmt.Sprintf("GPT model set:\n```%s```", value))
		} else {
			interactionRespondEphemeral(s, i, modelPermissionDeniedMessage)
		}

	} else {
		interactionRespondEphemeral(s, i, "Invaild input.")
	}
}

func clearGptHistory(s *discord.Session, i *discord.InteractionCreate) {
	if !isActiveGptChannel(i.ChannelID) {
		interactionRespond(s, i, notActiveGptChannelMessage)
		return
	}

	activeGptChannels[i.ChannelID].GPT.ClearHistory()
	interactionRespond(s, i, "Chat history cleared.")
}
