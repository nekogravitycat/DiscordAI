package chatbot

import (
	"fmt"

	discord "github.com/bwmarrin/discordgo"
)

func messageReply(session *discord.Session, messageCreate *discord.MessageCreate, content string) (*discord.Message, error) {
	return session.ChannelMessageSendReply(messageCreate.ChannelID, content, messageCreate.Reference())
}

func interactionRespond(session *discord.Session, interactionCreate *discord.InteractionCreate, content string) error {
	return session.InteractionRespond(interactionCreate.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Content: content,
		},
	})
}

func interactionRespondEphemeral(session *discord.Session, InteractionCreate *discord.InteractionCreate, content string) error {
	return session.InteractionRespond(InteractionCreate.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Flags:   discord.MessageFlagsEphemeral,
			Content: content,
		},
	})
}

func interactionRespondImage(session *discord.Session, InteractionCreate *discord.InteractionCreate, url string) error {
	return session.InteractionRespond(InteractionCreate.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Content: fmt.Sprintf("Image added in chat: %s", url),
			Embeds: []*discord.MessageEmbed{
				{
					Image: &discord.MessageEmbedImage{
						URL: url,
					},
				},
			},
		},
	})
}
