package userdata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/nekogravitycat/DiscordAI/internal/config"
	openai "github.com/sashabaranov/go-openai"
)

type UserInfo struct {
	Model          string  `json:"model"`
	Credit         float32 `json:"credit"`
	PrivilegeLevel int     `json:"privilege-level"`
}

func NewUserInfo() UserInfo {
	u := UserInfo{
		Model:          openai.GPT3Dot5Turbo,
		Credit:         config.InitCredits,
		PrivilegeLevel: config.InitPrivilege,
	}
	return u
}

var users = map[string]UserInfo{}

func GetUser(discordID string) (user UserInfo, ok bool) {
	user, ok = users[discordID]
	return user, ok
}

func SetUser(discordID string, user UserInfo) {
	users[discordID] = user
}

const USERFILE = "./data/users.json"

func LoadUserData() {
	if _, err := os.Stat(USERFILE); errors.Is(err, os.ErrNotExist) {
		fmt.Println("No users.json found, creating one.")
		users["0"] = NewUserInfo()
		SaveUserData()
	}

	jsonFile, err := os.Open(USERFILE)
	if err != nil {
		fmt.Println("Error reading user.json")
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading bytes of user.json")
	}

	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		fmt.Println("Error parsing user.json into Users struct.")
	}
}

func SaveUserData() {
	jsonFile, err := os.Create(USERFILE)
	if err != nil {
		fmt.Println("Error writing user.json")
		return
	}
	defer jsonFile.Close()

	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		fmt.Println("Error parsing Users struct into json data.")
		return
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing user.json file.")
	}
}
