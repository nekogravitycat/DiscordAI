package config

import (
	"github.com/nekogravitycat/DiscordAI/internal/jsondata"
)

type gptLimit struct {
	PromptTokens    int `json:"prompt-tokens"`
	SysPromptTokens int `json:"sys-prompt-tokens"`
	ReplyTokens     int `json:"reply-tokens"`
	HistoryLength   int `json:"history-length"`
}

type gptConfig struct {
	ConvertSCtoTC    bool     `json:"convert-sc-to-tc"`
	DefaultSysPrompt string   `json:"default-sys-prompt"`
	Limits           gptLimit `json:"limits"`
}

type mainConfig struct {
	NoGlobalCommands bool      `json:"no-global-commands"`
	AdminServers     []string  `json:"admin-servers"`
	InitCredits      float32   `json:"init-credits"`
	InitPrivilege    string    `json:"init-privilege"`
	GPT              gptConfig `json:"gpt"`
}

func newMainConfig() mainConfig {
	c := mainConfig{
		NoGlobalCommands: false,
		AdminServers:     []string{},
		InitCredits:      0.05,
		InitPrivilege:    "1",
		GPT: gptConfig{
			ConvertSCtoTC:    true,
			DefaultSysPrompt: "You have a great sense of humor and are an independent thinker who likes to chat.",
			Limits: gptLimit{
				PromptTokens:    500,
				SysPromptTokens: 250,
				ReplyTokens:     500,
				HistoryLength:   12,
			},
		},
	}
	return c
}

var configData mainConfig = newMainConfig()

var (
	NoGlobalCommands bool
	AdminServers     []string
	InitCredits      float32
	InitPrivilege    string
	GPT              gptConfig
)

const CONFIGFILE string = "./configs/config.json"

func LoadConfig() {
	jsondata.Check(CONFIGFILE, configData)
	jsondata.Load(CONFIGFILE, &configData)

	NoGlobalCommands = configData.NoGlobalCommands
	AdminServers = configData.AdminServers
	InitCredits = configData.InitCredits
	InitPrivilege = configData.InitPrivilege
	GPT = configData.GPT

	loadPrivilegeConfig()
}
