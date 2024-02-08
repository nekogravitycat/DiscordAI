package myjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// Please pass the address to toStruct
func ReadData(file string, toStruct any) {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		fmt.Println("No users.json found, creating one.")
		WriteData(file, toStruct)
	}

	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Error reading user.json")
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading bytes of user.json")
	}

	err = json.Unmarshal(byteValue, toStruct)
	if err != nil {
		fmt.Println("Error parsing user.json into Users struct.")
	}
}

func WriteData(file string, fromStruct any) {
	jsonFile, err := os.Create(file)
	if err != nil {
		fmt.Println("Error writing user.json")
		return
	}
	defer jsonFile.Close()

	jsonData, err := json.MarshalIndent(fromStruct, "", "  ")
	if err != nil {
		fmt.Println("Error parsing Users struct into json data.")
		return
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing user.json file.")
	}
}
