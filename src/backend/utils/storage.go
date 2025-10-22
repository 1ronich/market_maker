package utils

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gagliardetto/solana-go"
)

type Storage struct {
	Public_key  solana.PublicKey
	Private_key solana.PrivateKey
	//Amount      uint64
}

func UnloadStorage(file_name string, storage map[string][]map[string]interface{}) error {
	fileData, err := os.ReadFile(fmt.Sprintf("../../data/%s", file_name))
	if err != nil {
		fmt.Errorf("error reading file: %v", err)
	}

	err = json.Unmarshal(fileData, &storage)
	if err != nil {
		fmt.Errorf("error unmarshalling file data: %v", err)
	}
	return nil

}
func LoadStorage(file_name string, storage map[string][]map[string]interface{}) error {
	updatedData, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		fmt.Errorf("Error marshaling JSON: %v", err)
	}

	err = os.WriteFile(fmt.Sprintf("./data/%s", file_name), updatedData, 0644)
	if err != nil {
		fmt.Errorf("Error writing to the file: %v", err)
	}
	return nil

}
func UnloadSigs(file_name string, sigs map[string][]solana.Signature) error {
	fileData, err := os.ReadFile(fmt.Sprintf("../../data/%s", file_name))
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	err = json.Unmarshal(fileData, &sigs)
	if err != nil {
		fmt.Errorf("error unmarshalling file data: %v", err)
	}
	return nil
}
func LoadSigs(file_name string, sigs map[string][]solana.Signature) error {
	updatedData, err := json.MarshalIndent(sigs, "", "  ")
	if err != nil {
		fmt.Errorf("Error marshaling JSON: %v", err)
	}

	err = os.WriteFile(fmt.Sprintf("../../data/%s", file_name), updatedData, 0644)
	if err != nil {
		fmt.Errorf("Error writing to the file: %v", err)
	}
	return nil

}
