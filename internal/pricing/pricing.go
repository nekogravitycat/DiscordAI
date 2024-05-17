package pricing

import (
	"github.com/nekogravitycat/DiscordAI/internal/jsondata"
	"github.com/sashabaranov/go-openai"
)

type gptRate struct {
	Input  float32 `json:"input"`
	Output float32 `json:"output"`
}

type dallE3Resolutions struct {
	R1024X1024 float32 `json:"r1024x1024"`
	R1024X1792 float32 `json:"r1024x1792"`
	R1792X1024 float32 `json:"r1792x1024"`
}

type dallE3Pricing struct {
	Standard dallE3Resolutions `json:"standard"`
	Hd       dallE3Resolutions `json:"hd"`
}
type dallE2Pricing struct {
	R1024X1024 float32 `json:"r1024x1024"`
	R512X512   float32 `json:"r512x512"`
	R256X256   float32 `json:"r256x256"`
}

type pricingTable struct {
	Rate          float32       `json:"rate"`
	Gpt4o         gptRate       `json:"gpt-4o"`
	Gpt3Dot5Turbo gptRate       `json:"gpt-3.5-turbo"`
	DallE3        dallE3Pricing `json:"dall-e-3"`
	DallE2        dallE2Pricing `json:"dall-e-2"`
}

func newPricingTable() pricingTable {
	pt := pricingTable{
		Rate: 1.2,
		Gpt4o: gptRate{
			Input:  0.005,
			Output: 0.015,
		},
		Gpt3Dot5Turbo: gptRate{
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

var pricingData pricingTable = newPricingTable()

func GetGPTCost(model string, usage openai.Usage) float32 {
	var cost float32
	switch model {
	case openai.GPT4o:
		cost = float32(usage.PromptTokens)/1000*pricingData.Gpt4o.Input +
			float32(usage.CompletionTokens)/1000*pricingData.Gpt4o.Output
	case openai.GPT3Dot5Turbo:
		cost = float32(usage.PromptTokens)/1000*pricingData.Gpt3Dot5Turbo.Input +
			float32(usage.CompletionTokens)/1000*pricingData.Gpt3Dot5Turbo.Output
	default:
		cost = float32(usage.PromptTokens)/1000*pricingData.Gpt4o.Input +
			float32(usage.CompletionTokens)/1000*pricingData.Gpt4o.Output
	}
	return cost * pricingData.Rate
}

func GetDalleCost(model string, size string, quality string) float32 {
	if model == openai.CreateImageModelDallE3 {
		if quality == "standard" {
			if size == "1024x1024" {
				return pricingData.DallE3.Standard.R1024X1024 * pricingData.Rate
			} else if size == "1024x1792" {
				return pricingData.DallE3.Standard.R1024X1792 * pricingData.Rate
			} else if size == "1792x1024" {
				return pricingData.DallE3.Standard.R1792X1024 * pricingData.Rate
			}
		} else if quality == "hd" {
			if size == "1024x1024" {
				return pricingData.DallE3.Hd.R1024X1024 * pricingData.Rate
			} else if size == "1024x1792" {
				return pricingData.DallE3.Hd.R1024X1792 * pricingData.Rate
			} else if size == "1792x1024" {
				return pricingData.DallE3.Hd.R1792X1024 * pricingData.Rate
			}
		}
	} else if model == openai.CreateImageModelDallE2 {
		if size == "1024x1024" {
			return pricingData.DallE2.R1024X1024 * pricingData.Rate
		} else if size == "512x512" {
			return pricingData.DallE2.R512X512 * pricingData.Rate
		} else if size == "256x256" {
			return pricingData.DallE2.R256X256 * pricingData.Rate
		}
	}
	return pricingData.DallE3.Hd.R1024X1792 * pricingData.Rate
}

const PRICINGFILE = "./configs/pricing.json"

func LoadPricingTable() {
	jsondata.Check(PRICINGFILE, pricingData)
	jsondata.Load(PRICINGFILE, &pricingData)
}
