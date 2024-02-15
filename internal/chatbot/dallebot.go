package chatbot

import (
	"fmt"

	discord "github.com/bwmarrin/discordgo"
	"github.com/nekogravitycat/DiscordAI/internal/dalle"
	"github.com/nekogravitycat/DiscordAI/internal/pricing"
	"github.com/nekogravitycat/DiscordAI/internal/userdata"
	"github.com/sashabaranov/go-openai"
)

func dalleReply(s *discord.Session, i *discord.InteractionCreate, model string) {
	user, _ := userdata.GetUser(i.Member.User.ID)
	if user.Credit <= 0 {
		interactionRespond(s, i, "Not enough credits.")
		return
	}
	if !user.HasModelPrivilege(user.Model) {
		interactionRespond(s, i, "Permission denied. You cannot use this model.")
		return
	}

	options := i.ApplicationCommandData().Options
	optionMap := mapInteractionOptions(options)

	inputPrompt, ok := optionMap["prompt"]
	if !ok {
		interactionRespondEphemeral(s, i, "Invaild input for prompt.")
		return
	}
	prompt := inputPrompt.StringValue()

	inputSize, ok := optionMap["size"]
	if !ok {
		interactionRespondEphemeral(s, i, "Invaild input for size.")
		return
	}
	size := inputSize.StringValue()

	inputQuality, ok := optionMap["quality"]
	quality := ""
	if !ok {
		if model == openai.CreateImageModelDallE3 {
			interactionRespondEphemeral(s, i, "Invaild input for quality.")
			return
		}
	} else {
		quality = inputQuality.StringValue()
	}

	inputStyle, ok := optionMap["style"]
	style := ""
	if !ok {
		if model == openai.CreateImageModelDallE3 {
			interactionRespondEphemeral(s, i, "Invaild input for style.")
			return
		}
	} else {
		style = inputStyle.StringValue()
	}

	fmt.Printf("Model: %s, User: %s\n", model, i.Member.User.ID)

	interactionRespond(s, i, fmt.Sprintf("Your image is being generated, please wait...\n```Model: %s, Size: %s, Quality: %s, Style: %s\nPrompt: %s```", model, size, quality, style, prompt))
	result, err := dalle.Generate(openaiClient, model, prompt, size, quality, style, i.Member.User.ID)
	if err != nil {
		messageReplyEmbedImage(s, i.ChannelID, "Something went wrong", err.Error(), result)
		return
	}

	fmt.Printf("Usage: $%f USD\n", pricing.GetDalleCost(model, size, quality))
	messageReplyEmbedImage(s, i.ChannelID, fmt.Sprintf("Generated image for @%s", i.Member.User.Username), prompt, result)

	// Update userdata to follow up possible simultaneous operations
	user, _ = userdata.GetUser(i.Member.User.ID)
	user.Credit -= pricing.GetDalleCost(model, size, quality)
	userdata.SetUser(i.Member.User.ID, user)
	userdata.SaveUserData()
}
