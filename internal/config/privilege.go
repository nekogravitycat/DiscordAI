package config

import (
	"github.com/nekogravitycat/DiscordAI/internal/jsondata"
	"github.com/sashabaranov/go-openai"
)

type privilegeConfig struct {
	Models []string `json:"models"`
}

var privilegeData = map[string]privilegeConfig{
	"admin": {
		Models: []string{
			openai.GPT4o,
			openai.GPT3Dot5Turbo,
			openai.CreateImageModelDallE3,
			openai.CreateImageModelDallE2,
		},
	},
	"0": {
		Models: []string{},
	},
	"1": {
		Models: []string{
			openai.GPT3Dot5Turbo,
		},
	},
	"2": {
		Models: []string{
			openai.GPT4o,
			openai.GPT3Dot5Turbo,
		},
	},
	"3": {
		Models: []string{
			openai.GPT4o,
			openai.GPT3Dot5Turbo,
			openai.CreateImageModelDallE3,
			openai.CreateImageModelDallE2,
		},
	},
}

func VaildPrivilegeLevel(level string) bool {
	_, ok := privilegeData[level]
	return ok
}

func GetPrivilegeConfig(level string) (c privilegeConfig, ok bool) {
	c, ok = privilegeData[level]
	return c, ok
}

const PRIVILEGEFILE = "./configs/privilege.json"

func loadPrivilegeConfig() {
	jsondata.Check(PRIVILEGEFILE, privilegeData)
	jsondata.Load(PRIVILEGEFILE, &privilegeData)
}
