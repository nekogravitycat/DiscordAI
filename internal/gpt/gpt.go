package gpt

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/liuzl/gocc"
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

func (g *GPT) ClearHistory() {
	g.history = []openai.ChatCompletionMessage{}
}

func (g *GPT) removeHistoryIndex(index int) {
	g.history = append(g.history[:index], g.history[index+1:]...)
}

func (g *GPT) removeHistoryImages() {
	for i, m := range g.history {
		if m.MultiContent != nil {
			g.removeHistoryIndex(i)
		}
	}
}

func (g *GPT) trimOldHistory() {
	excess := len(g.history) - config.GPT.Limits.HistoryLength
	if excess > 0 {
		g.history = g.history[excess:]
		fmt.Printf("Trim %d history messages\n", excess)
	}
}

func (g *GPT) downgradeHistoryImages() {
	for i, msg := range g.history {
		if msg.MultiContent != nil {
			fmt.Printf("Image detail downgraded: %s", msg.MultiContent[0].ImageURL.URL)
			g.history[i].MultiContent[0].ImageURL.Detail = openai.ImageURLDetailLow
		}
	}
}

var s2t, _ = gocc.New("s2tw")

func (g *GPT) Generate(model string, user string) (reply string, usage openai.Usage, err error) {

	if len(g.history) <= 0 {
		return "", openai.Usage{}, errors.New("empty history")
	}

	g.trimOldHistory()

	var history = []openai.ChatCompletionMessage{}
	if model != "gpt-4-vision-preview" {
		history = append(history, historyWithoutImages(g.history)...)
		fmt.Println("Images ignored due to model limitation.")
	} else {
		history = append(history, g.history...)
	}

	response, err := g.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:     model,
			Messages:  append(g.sysPromptMsg(), history...),
			User:      user,
			MaxTokens: config.GPT.Limits.ReplyTokens,
		},
	)

	if err != nil {
		g.removeHistoryImages()
		return fmt.Sprintf("Something went wrong, please try again later. (Any image is removed from history.)\n```%s```", err.Error()), openai.Usage{}, err
	}

	reply = response.Choices[0].Message.Content

	if config.GPT.ConvertSCtoTC {
		fmt.Println("Converting SC to TC.")
		if tc, err := s2t.Convert(reply); err == nil {
			reply = tc
		} else {
			fmt.Println("Error converting SC to TC:")
			fmt.Println(err)
		}
	}

	usage = response.Usage

	g.addReply(reply)
	g.downgradeHistoryImages()

	return reply, usage, err
}

func historyWithoutImages(history []openai.ChatCompletionMessage) []openai.ChatCompletionMessage {
	h := []openai.ChatCompletionMessage{}
	for _, m := range history {
		if m.MultiContent == nil {
			h = append(h, m)
		} else {
			replacement := openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: "[Image removed]",
			}
			h = append(h, replacement)
		}
	}
	return h
}

func CountToken(prompt string, model string) int {
	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return 0
	}

	token := tkm.Encode(prompt, nil, nil)
	return len(token)
}
