package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nekogravitycat/DiscordAI/internal/chatbot"
	"github.com/nekogravitycat/DiscordAI/internal/config"
	"github.com/nekogravitycat/DiscordAI/internal/pricing"
	"github.com/nekogravitycat/DiscordAI/internal/userdata"
)

func init() {
	// Load enviroment variables from .env file if exist
	if _, err := os.Stat("./.env"); err == nil {
		fmt.Println(".env file found.")
		err = godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file.")
		}
	} else {
		fmt.Println("No .env file found, using system env.")
	}

	// Create ./configs folder if not exist
	if err := os.MkdirAll("./configs", os.ModePerm); err != nil {
		log.Fatal("Error creating ./configs folder.")
	}

	// Create ./data folder if not exist
	if err := os.MkdirAll("./data", os.ModePerm); err != nil {
		log.Fatal("Error creating ./data folder.")
	}

	config.LoadConfig()

	userdata.LoadUserData()

	pricing.LoadPricingTable()

	chatbot.LoadGptChannels()
}

func main() {
	chatbot.Run()
	// Run() will not exit unitl receive CTRL-C signal
	chatbot.Stop()
}
