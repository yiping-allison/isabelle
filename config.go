package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// BotConfig represents bot configurations
type BotConfig struct {
	BotKey string `json:"botKey"`
}

// LoadConfig loads bot configuration variables
// from the following file in main directory:
//
// .config
func LoadConfig() (BotConfig, error) {
	f, err := os.Open(".config")
	if err != nil {
		return BotConfig{}, err
	}
	var botConfig BotConfig
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&botConfig)
	if err != nil {
		return BotConfig{}, err
	}
	fmt.Println("successfully loaded .config")
	return botConfig, nil
}
