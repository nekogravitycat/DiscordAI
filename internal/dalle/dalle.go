package dalle

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

func Generate(client *openai.Client, model string, prompt string, size string, quality string, style string, user string) (url string, err error) {
	request := openai.ImageRequest{
		Model:          model,
		Prompt:         prompt,
		Size:           size,
		Quality:        quality,
		Style:          style,
		N:              1,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		User:           user,
	}

	response, err := client.CreateImage(context.Background(), request)
	if err != nil {
		return "https://t.gravitycat.tw/errorimg", err
	}

	return response.Data[0].URL, err
}
