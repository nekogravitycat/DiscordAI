package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

type User struct {
	Model          string  `json:"model"`
	Credit         float64 `json:"credit"`
	PrivilegeLevel int     `json:"privilege-level"`
}

var Users map[string]User

func AddUser(discordID string, user User) {
	Users[discordID] = user
}

func ReadUserData() {
	if _, err := os.Stat("./data/users.json"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("No users.json found, creating one.")
		WriteUserData()
	}

	jsonFile, err := os.Open("./data/users.json")
	if err != nil {
		fmt.Println("Error reading user.json")
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading bytes of user.json")
	}

	err = json.Unmarshal(byteValue, &Users)
	if err != nil {
		fmt.Println("Error parsing user.json into Users struct.")
	}
}

func WriteUserData() {
	jsonFile, err := os.Create("./data/users.json")
	if err != nil {
		fmt.Println("Error writing user.json")
		return
	}
	defer jsonFile.Close()

	jsonData, err := json.MarshalIndent(Users, "", "  ")
	if err != nil {
		fmt.Println("Error parsing Users struct into json data.")
		return
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing user.json file.")
	}
}
