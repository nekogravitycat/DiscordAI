package chatbot

import (
	"fmt"

	discord "github.com/bwmarrin/discordgo"
	"github.com/nekogravitycat/DiscordAI/internal/config"
	"github.com/nekogravitycat/DiscordAI/internal/gpt"
	"github.com/nekogravitycat/DiscordAI/internal/userdata"
	"github.com/sashabaranov/go-openai"
)

var regularCommands = []*discord.ApplicationCommand{
	{
		Name:        "gpt",
		Description: "Commands for GPT",
		Options: []*discord.ApplicationCommandOption{
			{
				Name:        "activate",
				Description: "Start ChatGPT on this channel",
				DescriptionLocalizations: map[discord.Locale]string{
					discord.ChineseTW: "開始在此頻道使用 ChatGPT",
				},
				Type: discord.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "deactivate",
				Description: "Stop ChatGPT on this channel",
				DescriptionLocalizations: map[discord.Locale]string{
					discord.ChineseTW: "停止在此頻道使用 ChatGPT",
				},
				Type: discord.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "set-model",
				Description: "Set the GPT model for the user",
				DescriptionLocalizations: map[discord.Locale]string{
					discord.ChineseTW: "設定用戶使用的 GPT 模型",
				},
				Type: discord.ApplicationCommandOptionSubCommand,
				Options: []*discord.ApplicationCommandOption{
					{
						Name:        "model",
						Description: "The GPT model to use",
						Type:        discord.ApplicationCommandOptionString,
						Required:    true,
						Choices: []*discord.ApplicationCommandOptionChoice{
							{
								Name:  "GPT-3.5 Turbo",
								Value: openai.GPT3Dot5Turbo,
							},
							{
								Name:  "GPT-4 Turbo Preview",
								Value: openai.GPT4TurboPreview,
							},
							{
								Name:  "GPT-4 Vision Preview",
								Value: openai.GPT4VisionPreview,
							},
						},
					},
				},
			},
			{
				Name:        "clear-history",
				Description: "Clear GPT chat history for this channel",
				DescriptionLocalizations: map[discord.Locale]string{
					discord.ChineseTW: "清除此頻道的 GPT 聊天歷史",
				},
				Type: discord.ApplicationCommandOptionSubCommand,
			},
			{
				Name:        "sys-prompt",
				Description: "Commands for the system prompt of GPT",
				Type:        discord.ApplicationCommandOptionSubCommandGroup,
				Options: []*discord.ApplicationCommandOption{
					{
						Name:        "show",
						Description: "Show GPT system prompt for this channel",
						DescriptionLocalizations: map[discord.Locale]string{
							discord.ChineseTW: "查看此頻道的 GPT 系統設定",
						},
						Type: discord.ApplicationCommandOptionSubCommand,
					},
					{
						Name:        "set",
						Description: "Set GPT system prompt for this channel",
						DescriptionLocalizations: map[discord.Locale]string{
							discord.ChineseTW: "更改此頻道的 GPT 系統設定",
						},
						Type: discord.ApplicationCommandOptionSubCommand,
						Options: []*discord.ApplicationCommandOption{
							{
								Name:        "sys-prompt",
								Description: "The system prompt of GPT",
								Type:        discord.ApplicationCommandOptionString,
								Required:    true,
							},
						},
					},
					{
						Name:        "reset",
						Description: "Reset GPT system prompt for this channel to default",
						DescriptionLocalizations: map[discord.Locale]string{
							discord.ChineseTW: "還原此頻道的 GPT 系統設定至預設值",
						},
						Type: discord.ApplicationCommandOptionSubCommand,
					},
				},
			},
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
		Name:        "dall-e-2-generate",
		Description: "Generate an image using DALL·E 2",
		DescriptionLocalizations: &map[discord.Locale]string{
			discord.ChineseTW: "使用 DALL·E 2 生成圖片",
		},
		Options: []*discord.ApplicationCommandOption{
			{
				Name:        "prompt",
				Description: "Prompt for image generation",
				Type:        discord.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "size",
				Description: "Image size to generate",
				Type:        discord.ApplicationCommandOptionString,
				Required:    true,
				Choices: []*discord.ApplicationCommandOptionChoice{
					{
						Name:  "256x256",
						Value: openai.CreateImageSize256x256,
					},
					{
						Name:  "512x512",
						Value: openai.CreateImageSize512x512,
					},
					{
						Name:  "1024x1024",
						Value: openai.CreateImageSize1024x1024,
					},
				},
			},
		},
	},
	{
		Name:        "dall-e-3-generate",
		Description: "Generate an image using DALL·E 3",
		DescriptionLocalizations: &map[discord.Locale]string{
			discord.ChineseTW: "使用 DALL·E 3 生成圖片",
		},
		Options: []*discord.ApplicationCommandOption{
			{
				Name:        "prompt",
				Description: "Prompt for image generation",
				Type:        discord.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "size",
				Description: "Image size to generate",
				Type:        discord.ApplicationCommandOptionString,
				Required:    true,
				Choices: []*discord.ApplicationCommandOptionChoice{
					{
						Name:  "1024x1024",
						Value: openai.CreateImageSize1024x1024,
					},
					{
						Name:  "1024x1792",
						Value: openai.CreateImageSize1024x1792,
					},
					{
						Name:  "1792x1024",
						Value: openai.CreateImageSize1024x1792,
					},
				},
			},
			{
				Name:        "quality",
				Description: "Image quality to generate",
				Type:        discord.ApplicationCommandOptionString,
				Required:    true,
				Choices: []*discord.ApplicationCommandOptionChoice{
					{
						Name:  "SD",
						Value: openai.CreateImageQualityStandard,
					},
					{
						Name:  "HD",
						Value: openai.CreateImageQualityHD,
					},
				},
			},
			{
				Name:        "style",
				Description: "Image style to generate",
				Type:        discord.ApplicationCommandOptionString,
				Required:    true,
				Choices: []*discord.ApplicationCommandOptionChoice{
					{
						Name:  "vivid",
						Value: openai.CreateImageStyleVivid,
					},
					{
						Name:  "natural",
						Value: openai.CreateImageStyleNatural,
					},
				},
			},
		},
	},
}

var adminCommands = []*discord.ApplicationCommand{
	{
		Name:        "add-credits",
		Description: "Add credits for a user",
		Options: []*discord.ApplicationCommandOption{
			{
				Name:        "user-id",
				Description: "User ID",
				Type:        discord.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "amount",
				Description: "Amount to add (in USD)",
				Type:        discord.ApplicationCommandOptionNumber,
				Required:    true,
			},
		},
	},
	{
		Name:        "set-user-privilege",
		Description: "Set privilege level for the user",
		Options: []*discord.ApplicationCommandOption{
			{
				Name:        "user-id",
				Description: "User ID",
				Type:        discord.ApplicationCommandOptionString,
				Required:    true,
			},
			{
				Name:        "privilege-level",
				Description: "Privilege level",
				Type:        discord.ApplicationCommandOptionString,
				Required:    true,
			},
		},
	},
}

var commandHandlers = map[string]func(s *discord.Session, i *discord.InteractionCreate){
	"gpt":               gptGroup,
	"credits":           credits,
	"dall-e-2-generate": dalle2Generate,
	"dall-e-3-generate": dalle3Generate,
}

var adminCommandHandlers = map[string]func(s *discord.Session, i *discord.InteractionCreate){
	"add-credits":        addCredit,
	"set-user-privilege": setPrivilege,
}

func mapInteractionOptions(options []*discord.ApplicationCommandInteractionDataOption) map[string]*discord.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discord.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

func gptGroup(s *discord.Session, i *discord.InteractionCreate) {
	gptOptions := i.ApplicationCommandData().Options

	switch gptOptions[0].Name {
	case "activate":
		activateGPT(s, i)
	case "deactivate":
		deactivateGPT(s, i)
	case "set-model":
		setGptModel(s, i)
	case "clear-history":
		clearGptHistory(s, i)
	case "sys-prompt":
		sysPromptOptions := gptOptions[0].Options
		switch sysPromptOptions[0].Name {
		case "show":
			showGptSysPrompt(s, i)
		case "set":
			setGptSysPrompt(s, i)
		case "reset":
			resetGptSysPrompt(s, i)
		}
	}
	interactionRespondEphemeral(s, i, "Unknown command")
}

func activateGPT(s *discord.Session, i *discord.InteractionCreate) {
	if isActiveGptChannel(i.ChannelID) {
		interactionRespondEphemeral(s, i, "Channel already activated")
		return
	}

	interactionRespond(s, i, "Hello!")
	activeGptChannels[i.ChannelID] = newGptChannel()
	saveGptChannels()
}

func deactivateGPT(s *discord.Session, i *discord.InteractionCreate) {
	if !isActiveGptChannel(i.ChannelID) {
		interactionRespondEphemeral(s, i, "Channel is not active")
		return
	}

	interactionRespond(s, i, "Bye!")
	delete(activeGptChannels, i.ChannelID)
	saveGptChannels()
}

func showGptSysPrompt(s *discord.Session, i *discord.InteractionCreate) {
	if !isActiveGptChannel(i.ChannelID) {
		interactionRespondEphemeral(s, i, notActiveGptChannelMessage)
		return
	}

	interactionRespond(s, i, fmt.Sprintf("System prompt:\n```%s```", activeGptChannels[i.ChannelID].GPT.SysPrompt))
}

func setGptSysPrompt(s *discord.Session, i *discord.InteractionCreate) {
	if !isActiveGptChannel(i.ChannelID) {
		interactionRespondEphemeral(s, i, notActiveGptChannelMessage)
		return
	}

	options := i.ApplicationCommandData().Options[0].Options[0].Options
	fmt.Println(options)
	optionMap := mapInteractionOptions(options)
	fmt.Println(optionMap)

	inputPrompt, ok := optionMap["sys-prompt"]
	if !ok {
		interactionRespondEphemeral(s, i, "Invaild input for system prompt.")
		return
	}

	prompt := inputPrompt.StringValue()
	if gpt.CountToken(prompt, openai.GPT3Dot5Turbo) > config.GPT.Limits.SysPromptTokens {
		interactionRespondEphemeral(s, i, "System prompt is too long.")
		return
	}

	activeGptChannels[i.ChannelID].GPT.SysPrompt = prompt

	saveGptChannels()
	interactionRespond(s, i, fmt.Sprintf("System prompt updated:\n```%s```", prompt))
}

func resetGptSysPrompt(s *discord.Session, i *discord.InteractionCreate) {
	if !isActiveGptChannel(i.ChannelID) {
		interactionRespond(s, i, notActiveGptChannelMessage)
		return
	}

	activeGptChannels[i.ChannelID].GPT.SysPrompt = config.GPT.DefaultSysPrompt
	saveGptChannels()
	interactionRespond(s, i, fmt.Sprintf("System prompt updated:\n```%s```", activeGptChannels[i.ChannelID].GPT.SysPrompt))
}

func setGptModel(s *discord.Session, i *discord.InteractionCreate) {
	options := i.ApplicationCommandData().Options[0].Options
	optionMap := mapInteractionOptions(options)

	inputModel, ok := optionMap["model"]
	if !ok {
		interactionRespondEphemeral(s, i, "Invaild input for model.")
		return
	}

	model := inputModel.StringValue()
	user, userExist := userdata.GetUser(i.Member.User.ID)
	if !userExist {
		user = userdata.SetUser(i.Member.User.ID, userdata.NewUserInfo())
	}

	if !user.HasModelPrivilege(model) {
		interactionRespondEphemeral(s, i, modelPermissionDeniedMessage)
		return
	}

	user.Model = model
	userdata.SetUser(i.Member.User.ID, user)
	userdata.SaveUserData()
	interactionRespondEphemeral(s, i, fmt.Sprintf("GPT model set:\n```%s```", model))
}

func clearGptHistory(s *discord.Session, i *discord.InteractionCreate) {
	if !isActiveGptChannel(i.ChannelID) {
		interactionRespondEphemeral(s, i, notActiveGptChannelMessage)
		return
	}

	activeGptChannels[i.ChannelID].GPT.ClearHistory()
	interactionRespond(s, i, "Chat history cleared.")
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

func dalle2Generate(s *discord.Session, i *discord.InteractionCreate) {
	dalleReply(s, i, openai.CreateImageModelDallE2)
}

func dalle3Generate(s *discord.Session, i *discord.InteractionCreate) {
	dalleReply(s, i, openai.CreateImageModelDallE3)
}

func addCredit(s *discord.Session, i *discord.InteractionCreate) {
	operator, ok := userdata.GetUser(i.Member.User.ID)
	if !ok {
		interactionRespondEphemeral(s, i, "Unknown opeartor.")
		return
	}

	if !operator.IsAdmin() {
		interactionRespondEphemeral(s, i, "Permission denied: not admin.")
		return
	}

	options := i.ApplicationCommandData().Options
	optionMap := mapInteractionOptions(options)

	inputUserId, ok := optionMap["user-id"]
	if !ok {
		interactionRespondEphemeral(s, i, "Invaild input for user-id.")
		return
	}
	userId := inputUserId.StringValue()

	user, ok := userdata.GetUser(userId)
	if !ok {
		interactionRespondEphemeral(s, i, fmt.Sprintf("User id does not exist: `%s`", userId))
		return
	}

	inputAmount, ok := optionMap["amount"]
	if !ok {
		interactionRespondEphemeral(s, i, "Invaild input for amount.")
		return
	}
	amount := inputAmount.FloatValue()

	user.Credit += float32(amount)
	userdata.SetUser(userId, user)
	userdata.SaveUserData()
	interactionRespond(s, i, fmt.Sprintf("User (`%s`) credit updated: `$%f USD`", userId, user.Credit))
}

func setPrivilege(s *discord.Session, i *discord.InteractionCreate) {
	operator, ok := userdata.GetUser(i.Member.User.ID)
	if !ok {
		interactionRespondEphemeral(s, i, "Unknown opeartor.")
		return
	}

	if !operator.IsAdmin() {
		interactionRespondEphemeral(s, i, "Permission denied: not admin.")
		return
	}

	options := i.ApplicationCommandData().Options
	optionMap := mapInteractionOptions(options)

	inputUserId, ok := optionMap["user-id"]
	if !ok {
		interactionRespondEphemeral(s, i, "Invaild input for user-id.")
		return
	}
	userId := inputUserId.StringValue()

	user, ok := userdata.GetUser(userId)
	if !ok {
		interactionRespondEphemeral(s, i, fmt.Sprintf("User id does not exist: `%s`", userId))
		return
	}

	inputLevel, ok := optionMap["privilege-level"]
	if !ok {
		interactionRespondEphemeral(s, i, "Invaild input for privilege level.")
		return
	}
	level := inputLevel.StringValue()

	if !config.VaildPrivilegeLevel(level) {
		interactionRespondEphemeral(s, i, fmt.Sprintf("Unrecognized privilege level: `%s`", level))
		return
	}

	user.PrivilegeLevel = level
	userdata.SetUser(userId, user)
	userdata.SaveUserData()
	interactionRespond(s, i, fmt.Sprintf("User (`%s`) privilege level updated: `%s`", userId, user.PrivilegeLevel))
}
