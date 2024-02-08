package config

type gptLimit struct {
	PromptTokens    int `json:"prompt-tokens"`
	SysPromptTokens int `json:"sys-prompt-tokens"`
	ReplyTokens     int `json:"reply-tokens"`
	HistoryLength   int `json:"history-length"`
}

type gptConfig struct {
	DefaultSysPrompt string   `json:"default-sys-prompt"`
	Limits           gptLimit `json:"limits"`
}

type config struct {
	InitCredits   float64   `json:"init-credits"`
	InitPrivilege int       `json:"init-privilege"`
	GPT           gptConfig `json:"gpt"`
}

func newConfig() config {
	c := config{
		InitCredits:   0.05,
		InitPrivilege: 1,
		GPT: gptConfig{
			DefaultSysPrompt: "You have a great sense of humor and are an independent thinker who likes to chat.",
			Limits: gptLimit{
				PromptTokens:    650,
				SysPromptTokens: 500,
				ReplyTokens:     1500,
				HistoryLength:   12,
			},
		},
	}
	return c
}

var (
	InitCredits   float64
	InitPrivilege int
	GPT           gptConfig
)

const CONFIGFILE string = "./config/config.json"
