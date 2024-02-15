package chatbot

import (
	discord "github.com/bwmarrin/discordgo"
)

func messageReply(session *discord.Session, message *discord.Message, content string) (*discord.Message, error) {
	return session.ChannelMessageSendReply(message.ChannelID, content, message.Reference())
}

func messageReplyEmbedImage(session *discord.Session, channelID string, title string, description string, url string) (*discord.Message, error) {
	embed := discord.MessageEmbed{
		Title:       title,
		Description: description,
		URL:         url,
		Type:        discord.EmbedTypeImage,
		Image: &discord.MessageEmbedImage{
			URL: url,
		},
	}

	return session.ChannelMessageSendEmbed(channelID, &embed)
}

func interactionRespond(session *discord.Session, interactionCreate *discord.InteractionCreate, content string) error {
	response := discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Content: content,
		},
	}

	return session.InteractionRespond(interactionCreate.Interaction, &response)
}

func interactionRespondEphemeral(session *discord.Session, InteractionCreate *discord.InteractionCreate, content string) error {
	response := discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Flags:   discord.MessageFlagsEphemeral,
			Content: content,
		},
	}

	return session.InteractionRespond(InteractionCreate.Interaction, &response)
}
