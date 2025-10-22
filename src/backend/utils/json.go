package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func WriteToJSONFile(content any, file_path string) {
	data, err := json.MarshalIndent(content, " ", "  ")
	if err != nil {
		log.Fatalf("error writing to JSON file: %v", err)
	}

	file, err := os.Create(file_path)
	if err != nil {
		log.Fatalf("error creating json file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("error writing to json file: %v", err)
	}
}

func ReadJSONFile(expectedResponseType interface{}, file_path string) {
	file, err := os.ReadFile(file_path)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	err = json.Unmarshal(file, expectedResponseType)
	if err != nil {
		fmt.Println("error decoding JSON:", err)
	}
}
