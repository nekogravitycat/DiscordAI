package main

import discord "github.com/bwmarrin/discordgo"

var (
	commands = []*discord.ApplicationCommand{
		{
			Name:        "start",
			Description: "Start ChatGPT on this channel",
		},
		{
			Name:        "stop",
			Description: "Stop ChatGPT on this channel",
		},
	}

	commandHandlers = map[string]func(s *discord.Session, i *discord.InteractionCreate){
		"start": startChat,
		"stop":  stopChat,
	}
)

func startChat(s *discord.Session, i *discord.InteractionCreate) {
	if _, ok := chats[i.ChannelID]; ok {
		s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
			Type: discord.InteractionResponseChannelMessageWithSource,
			Data: &discord.InteractionResponseData{
				Content: "Already in chat",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Content: "Hello!",
		},
	})
	chats[i.ChannelID] = NewChat()
}

func stopChat(s *discord.Session, i *discord.InteractionCreate) {
	if _, ok := chats[i.ChannelID]; !ok {
		s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
			Type: discord.InteractionResponseChannelMessageWithSource,
			Data: &discord.InteractionResponseData{
				Content: "Not in channel",
			},
		})
		return
	}

	s.InteractionRespond(i.Interaction, &discord.InteractionResponse{
		Type: discord.InteractionResponseChannelMessageWithSource,
		Data: &discord.InteractionResponseData{
			Content: "Bye!",
		},
	})
	delete(chats, i.ChannelID)
}
