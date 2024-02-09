package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

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

type mainConfig struct {
	InitCredits   float32   `json:"init-credits"`
	InitPrivilege int       `json:"init-privilege"`
	GPT           gptConfig `json:"gpt"`
}

func newMainConfig() mainConfig {
	c := mainConfig{
		InitCredits:   0.05,
		InitPrivilege: 1,
		GPT: gptConfig{
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

var config mainConfig

var (
	InitCredits   float32
	InitPrivilege int
	GPT           gptConfig
)

const CONFIGFILE string = "./configs/config.json"

func LoadConfig() {
	if _, err := os.Stat(CONFIGFILE); errors.Is(err, os.ErrNotExist) {
		fmt.Println("No config.json found, creating one with default values.")
		config = newMainConfig()
		saveConfig()
	}

	jsonFile, err := os.Open(CONFIGFILE)
	if err != nil {
		fmt.Println("Error reading config.json")
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading bytes of config.json")
	}

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println("Error parsing config.json into mainConfig struct.")
	}

	InitCredits = config.InitCredits
	InitPrivilege = config.InitPrivilege
	GPT = config.GPT
}

func saveConfig() {
	jsonFile, err := os.Create(CONFIGFILE)
	if err != nil {
		fmt.Println("Error writing config.json")
		return
	}
	defer jsonFile.Close()

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error parsing mainConfig struct into json data.")
		return
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing config.json file.")
	}
}
