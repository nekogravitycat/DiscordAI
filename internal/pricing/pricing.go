package pricing

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
)

type gptPricing struct {
	Input  float64 `json:"input"`
	Output float64 `json:"output"`
}

type dallE3Resolutions struct {
	R1024X1024 float64 `json:"r1024x1024"`
	R1024X1792 float64 `json:"r1024x1792"`
	R1792X1024 float64 `json:"r1792x1024"`
}

type dallE3Pricing struct {
	Standard dallE3Resolutions `json:"standard"`
	Hd       dallE3Resolutions `json:"hd"`
}
type dallE2Pricing struct {
	R1024X1024 float64 `json:"r1024x1024"`
	R512X512   float64 `json:"r512x512"`
	R256X256   float64 `json:"r256x256"`
}

type pricingTable struct {
	Gpt4TurboPreview  gptPricing    `json:"gpt-4-turbo-preview"`
	Gpt4VisionPreview gptPricing    `json:"gpt-4-vision-preview"`
	Gpt3Dot5Turbo     gptPricing    `json:"gpt-3.5-turbo"`
	DallE3            dallE3Pricing `json:"dall-e-3"`
	DallE2            dallE2Pricing `json:"dall-e-2"`
}

func newPricingTable() pricingTable {
	pt := pricingTable{
		Gpt4TurboPreview: gptPricing{
			Input:  0.01,
			Output: 0.03,
		},
		Gpt4VisionPreview: gptPricing{
			Input:  0.01,
			Output: 0.03,
		},
		Gpt3Dot5Turbo: gptPricing{
			Input:  0.0005,
			Output: 0.0015,
		},
		DallE3: dallE3Pricing{
			Standard: dallE3Resolutions{
				R1024X1024: 0.04,
				R1024X1792: 0.08,
				R1792X1024: 0.08,
			},
			Hd: dallE3Resolutions{
				R1024X1024: 0.08,
				R1024X1792: 0.12,
				R1792X1024: 0.12,
			},
		},
		DallE2: dallE2Pricing{
			R256X256:   0.016,
			R512X512:   0.018,
			R1024X1024: 0.02,
		},
	}
	return pt
}

var Table pricingTable

func GetGPTCost(model string, usage openai.Usage) float64 {
	switch model {
	case openai.GPT4TurboPreview:
		return float64(usage.PromptTokens)/1000*Table.Gpt4TurboPreview.Input +
			float64(usage.CompletionTokens)/1000*Table.Gpt4TurboPreview.Output
	case openai.GPT4VisionPreview:
		return float64(usage.PromptTokens)/1000*Table.Gpt4VisionPreview.Input +
			float64(usage.CompletionTokens)/1000*Table.Gpt4VisionPreview.Output
	case openai.GPT3Dot5Turbo:
		return float64(usage.PromptTokens)/1000*Table.Gpt3Dot5Turbo.Input +
			float64(usage.CompletionTokens)/1000*Table.Gpt3Dot5Turbo.Output
	default:
		return float64(usage.PromptTokens)/1000*Table.Gpt4TurboPreview.Input +
			float64(usage.CompletionTokens)/1000*Table.Gpt4TurboPreview.Output
	}
}

const PRICINGFILE = "./configs/pricing.json"

func LoadPricingTable() {
	if _, err := os.Stat(PRICINGFILE); errors.Is(err, os.ErrNotExist) {
		fmt.Println("No users.json found, creating one.")
		Table = newPricingTable()
		savePricingTable()
	}

	jsonFile, err := os.Open(PRICINGFILE)
	if err != nil {
		fmt.Println("Error reading user.json")
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading bytes of user.json")
	}

	err = json.Unmarshal(byteValue, &Table)
	if err != nil {
		fmt.Println("Error parsing user.json into Users struct.")
	}
}

func savePricingTable() {
	jsonFile, err := os.Create(PRICINGFILE)
	if err != nil {
		fmt.Println("Error writing user.json")
		return
	}
	defer jsonFile.Close()

	jsonData, err := json.MarshalIndent(Table, "", "  ")
	if err != nil {
		fmt.Println("Error parsing Users struct into json data.")
		return
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing user.json file.")
	}
}
