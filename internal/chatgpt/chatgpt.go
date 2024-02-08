package chatgpt

import (
	"context"
	"errors"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type GPT struct {
	client    openai.Client
	SysPrompt string
	history   []openai.ChatCompletionMessage
}

func NewGPT() GPT {
	gpt := GPT{
		client:    *openai.NewClient(os.Getenv("OPENAI_TOKEN")),
		SysPrompt: "You are a helpful assistant",
		history:   []openai.ChatCompletionMessage{},
	}
	return gpt
}

func (g *GPT) AddMessage(prompt string) {
	msg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompt,
	}
	g.history = append(g.history, msg)
}

func (g *GPT) AddImage(imageURL string, imageDetail string) {
	msg := openai.ChatCompletionMessage{
		Role: openai.ChatMessageRoleUser,
		MultiContent: []openai.ChatMessagePart{
			{
				Type: openai.ChatMessagePartTypeImageURL,
				ImageURL: &openai.ChatMessageImageURL{
					URL:    imageURL,
					Detail: openai.ImageURLDetail(imageDetail),
				},
			},
		},
	}
	g.history = append(g.history, msg)
}

func (g *GPT) addReply(reply string) {
	msg := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: reply,
	}
	g.history = append(g.history, msg)
}

func (g *GPT) trimHistory() {
	for len(g.history) > 10 {
		g.history = g.history[1:]
	}
}

func (g *GPT) Generate(model string, user string) (reply string, usage int, err error) {
	if len(g.history) <= 0 {
		return "", 0, errors.New("empty history")
	}

	g.trimHistory()

	fmt.Println(g.history)

	response, err := g.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     model,
			Messages:  g.history,
			User:      user,
			MaxTokens: 1500,
		},
	)

	if err != nil {
		return "```Something went wrong, please try again later.```", 0, err
	}

	reply = response.Choices[0].Message.Content
	usage = response.Usage.TotalTokens

	g.addReply(reply)

	return reply, usage, err
}
