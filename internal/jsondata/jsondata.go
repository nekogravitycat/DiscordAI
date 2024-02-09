package jsondata

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

func Check(file string, defaultData any) {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("'%s' does not exist, creating one with default values.\n", file)
		Save(file, defaultData)
	}
}

func Load(file string, dataPointer any) {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Error reading user.json")
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading bytes of user.json")
	}

	err = json.Unmarshal(byteValue, dataPointer)
	if err != nil {
		fmt.Println("Error parsing user.json into Users struct.")
	}
}

func Save(file string, data any) {
	jsonFile, err := os.Create(file)
	if err != nil {
		fmt.Println("Error writing user.json")
		return
	}
	defer jsonFile.Close()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error parsing Users struct into json data.")
		return
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		fmt.Println("Error writing user.json file.")
	}
}
