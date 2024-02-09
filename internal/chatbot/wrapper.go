package chatbot

import (
	discord "github.com/bwmarrin/discordgo"
)

func messageReply(session *discord.Session, messageCreate *discord.MessageCreate, content string) {
	session.ChannelMessageSendReply(messageCreate.ChannelID, content, messageCreate.Reference())
}

func interactionRespond(session *discord.Session, interactionCreate *discord.InteractionCreate, content string) error {
	return session.InteractionRespond(interactionCreate.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Content: content,
		},
	})
}

func interactionRespondEphemeral(session *discord.Session, InteractionCreate *discord.InteractionCreate, content string) {
	session.InteractionRespond(InteractionCreate.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Flags:   discord.MessageFlagsEphemeral,
			Content: content,
		},
	})
}
