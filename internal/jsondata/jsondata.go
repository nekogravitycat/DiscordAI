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
		fmt.Printf("'%s' does not exist, creating one with default values\n", file)
		Save(file, defaultData)
	}
}

func Load(file string, dataPointer any) {
	jsonFile, err := os.Open(file)
	if err != nil {
		fmt.Printf("Error reading '%s'\n", file)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Printf("Error reading bytes of '%s'\n", file)
	}

	err = json.Unmarshal(byteValue, dataPointer)
	if err != nil {
		fmt.Printf("Error parsing '%s'\n", file)
	}
}

func Save(file string, data any) {
	jsonFile, err := os.Create(file)
	if err != nil {
		fmt.Printf("Error writing '%s'\n", file)
		return
	}
	defer jsonFile.Close()

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error parsing data into json.")
		return
	}

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		fmt.Printf("Error writing '%s'\n", file)
	}
}
