package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

type privilegeConfig struct {
	Models []string `json:"models"`
}

func newPrivilegeConfig() privilegeConfig {
	pc := privilegeConfig{
		Models: []string{},
	}
	return pc
}

var privileges = map[string]privilegeConfig{}

func GetPrivilegeConfig(level int) (c privilegeConfig, ok bool) {
	c, ok = privileges[strconv.Itoa(level)]
	return c, ok
}

const PRIVILEGEFILE = "./configs/privilege.json"

func loadPrivilegeConfig() {
	if _, err := os.Stat(PRIVILEGEFILE); errors.Is(err, os.ErrNotExist) {
		fmt.Println("No privilege.json found, creating one with default values.")
		privileges["0"] = newPrivilegeConfig()
		savePrivilegeConfig()
	}

	jsonFile, err := os.Open(PRIVILEGEFILE)
	if err != nil {
		fmt.Println("Error reading privilege.json")
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading bytes of privilege.json")
	}

	err = json.Unmarshal(byteValue, &privileges)
	if err != nil {
		fmt.Println("Error parsing privilege.json into privilegeConfig struct.")
	}
}

func savePrivilegeConfig() {
	jsonFile, err := os.Create(PRIVILEGEFILE)
	if err != nil {
		fmt.Println("Error writing privilege.json")
		return
	}
	defer jsonFile.Close()

	jsonData, err := json.MarshalIndent(privileges, "", "  ")
	if err != nil {
		fmt.Println("Error parsing privilegeConfig struct into json data.")
		return
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing privilege.json file.")
	}
}
