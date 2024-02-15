package dalle

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

type RequestFormat struct {
	Client  *openai.Client
	Model   string
	Prompt  string
	Size    string
	Quality string
	Style   string
	User    string
}

func Generate(input RequestFormat) (url string, err error) {
	request := openai.ImageRequest{
		Model:          input.Model,
		Prompt:         input.Prompt,
		Size:           input.Size,
		Quality:        input.Quality,
		Style:          input.Style,
		N:              1,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		User:           input.User,
	}

	response, err := input.Client.CreateImage(context.Background(), request)
	if err != nil {
		return "https://t.gravitycat.tw/errorimg", err
	}

	return response.Data[0].URL, err
}
