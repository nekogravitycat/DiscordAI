package chatbot

import (
	"fmt"

	discord "github.com/bwmarrin/discordgo"
	"github.com/nekogravitycat/DiscordAI/internal/config"
	"github.com/nekogravitycat/DiscordAI/internal/userdata"
)

var (
	regularCommands = []*discord.ApplicationCommand{
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
			Name:        "add-gpt-image",
			Description: "Add a image URL for GPT Vision to see",
			DescriptionLocalizations: &map[discord.Locale]string{
				discord.ChineseTW: "提供圖片 URL 給 GPT Vision 參考",
			},
			Options: []*discord.ApplicationCommandOption{
				{
					Name:        "image-url",
					Description: "The URL of the image",
					Type:        discord.ApplicationCommandOptionString,
					Required:    true,
				},
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

	adminCommands = []*discord.ApplicationCommand{
		{
			Name:        "add-credits",
			Description: "Add credits for a user",
			Options: []*discord.ApplicationCommandOption{
				{
					Name:        "user-id",
					Type:        discord.ApplicationCommandOptionString,
					Description: "User ID",
					Required:    true,
				},
				{
					Name:        "amount",
					Type:        discord.ApplicationCommandOptionNumber,
					Description: "Amount to add (in USD)",
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
					Type:        discord.ApplicationCommandOptionString,
					Description: "Privilege level",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discord.Session, i *discord.InteractionCreate){
		"activate-gpt":         activateGPT,
		"deactivate-gpt":       deactivateGPT,
		"credits":              credits,
		"add-gpt-image":        addGptImage,
		"gpt-sys-prompt":       showGptSysPrompt,
		"set-gpt-sys-prompt":   setGptSysPrompt,
		"reset-gpt-sys-prompt": resetGptSysPrompt,
		"set-gpt-model":        setGptModel,
		"clear-gpt-history":    clearGptHistory,
	}

	adminCommandHandlers = map[string]func(s *discord.Session, i *discord.InteractionCreate){
		"add-credits":        addCredit,
		"set-user-privilege": setPrivilege,
	}
)

func mapInteractionOptions(options []*discord.ApplicationCommandInteractionDataOption) map[string]*discord.ApplicationCommandInteractionDataOption {
	optionMap := make(map[string]*discord.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}

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

func addGptImage(s *discord.Session, i *discord.InteractionCreate) {
	if !isActiveGptChannel(i.ChannelID) {
		interactionRespond(s, i, notActiveGptChannelMessage)
		return
	}

	options := i.ApplicationCommandData().Options
	optionMap := mapInteractionOptions(options)

	inputImageUrl, ok := optionMap["image-url"]
	if !ok {
		interactionRespond(s, i, "Invaild input for image URL")
		return
	}
	imageUrl := inputImageUrl.StringValue()

	if err := interactionRespondImage(s, i, imageUrl); err != nil {
		interactionRespondEphemeral(s, i, "Error adding the image.")
		return
	}

	activeGptChannels[i.ChannelID].GPT.AddImage(imageUrl, "auto")
	fmt.Println("Add image: " + imageUrl)
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

	options := i.ApplicationCommandData().Options
	optionMap := mapInteractionOptions(options)

	inputPrompt, ok := optionMap["sys-prompt"]
	if !ok {
		interactionRespond(s, i, "Invaild input for system prompt.")
		return
	}

	prompt := inputPrompt.StringValue()
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
	options := i.ApplicationCommandData().Options
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
		interactionRespond(s, i, notActiveGptChannelMessage)
		return
	}

	activeGptChannels[i.ChannelID].GPT.ClearHistory()
	interactionRespond(s, i, "Chat history cleared.")
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
	interactionRespondEphemeral(s, i, fmt.Sprintf("User (`%s`) credit updated: `$%f USD`", userId, user.Credit))
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
	interactionRespondEphemeral(s, i, fmt.Sprintf("User (`%s`) privilege level updated: `%s`", userId, user.PrivilegeLevel))
}
