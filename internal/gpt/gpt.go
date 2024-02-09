package gpt

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/nekogravitycat/DiscordAI/internal/config"
	"github.com/pkoukk/tiktoken-go"
	openai "github.com/sashabaranov/go-openai"
)

type GPT struct {
	client    openai.Client
	SysPrompt string `json:"sys-prompt"`
	history   []openai.ChatCompletionMessage
}

func NewGPT() GPT {
	gpt := GPT{
		client:    *openai.NewClient(os.Getenv("OPENAI_TOKEN")),
		SysPrompt: "You have a great sense of humor and are an independent thinker who likes to chat.",
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

func (g *GPT) sysPromptMsg() []openai.ChatCompletionMessage {
	sys := openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Content: g.SysPrompt,
	}
	return []openai.ChatCompletionMessage{sys}
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

func (g *GPT) downgradeHistoryImages() {
	for i, msg := range g.history {
		if msg.MultiContent != nil {
			g.history[i].MultiContent[0].ImageURL.Detail = openai.ImageURLDetailLow
		}
	}
}

func (g *GPT) Generate(model string, user string) (reply string, usage openai.Usage, err error) {
	fmt.Printf("Model: %s, User: %s\n", model, user)

	if len(g.history) <= 0 {
		return "", openai.Usage{}, errors.New("empty history")
	}

	g.trimHistory()
	g.downgradeHistoryImages()

	response, err := g.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     model,
			Messages:  append(g.sysPromptMsg(), g.history...),
			User:      user,
			MaxTokens: config.GPT.Limits.ReplyTokens,
		},
	)

	if err != nil {
		return "```Something went wrong, please try again later.```", openai.Usage{}, err
	}

	reply = response.Choices[0].Message.Content
	usage = response.Usage

	g.addReply(reply)

	fmt.Printf("Usage: %d\n", usage)

	return reply, usage, err
}

func CountToken(prompt string, model string) int {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return 0
	}

	token := tkm.Encode(prompt, nil, nil)
	return len(token)
}
