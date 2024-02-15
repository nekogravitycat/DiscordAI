package chatbot

import (
	discord "github.com/bwmarrin/discordgo"
)

func messageReply(session *discord.Session, messageCreate *discord.MessageCreate, content string) (*discord.Message, error) {
	return session.ChannelMessageSendReply(messageCreate.ChannelID, content, messageCreate.Reference())
}

func messageReplyEmbedImage(session *discord.Session, messageCreate *discord.MessageCreate, url string) (*discord.Message, error) {
	embed := discord.MessageEmbed{
		URL:  url,
		Type: discord.EmbedTypeImage,
		Image: &discord.MessageEmbedImage{
			URL: url,
		},
	}

	return session.ChannelMessageSendEmbedReply(messageCreate.ChannelID, &embed, messageCreate.Reference())
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
