package chatbot

import (
	discord "github.com/bwmarrin/discordgo"
)

func messageReply(s *discord.Session, m *discord.MessageCreate, content string) {
	s.ChannelMessageSendReply(m.ChannelID, content, m.Reference())
}

func interactionRespond(s *discord.Session, i *discord.InteractionCreate, content string) error {
	return s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Content: content,
		},
	})
}

func interactionRespondEphemeral(s *discord.Session, i *discord.InteractionCreate, content string) {
	s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Flags:   discord.MessageFlagsEphemeral,
			Content: content,
		},
	})
}
