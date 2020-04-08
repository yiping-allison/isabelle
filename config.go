package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// BotConfig represents bot configurations
type BotConfig struct {
	BotKey    string         `json:"botKey"`
	BotPrefix string         `json:"botPrefix"`
	Database  PostgresConfig `json:"database"`
	AdminRole string         `json:"adminRole"`
}

// PostgresConfig represents metadata required to start and maintain postgres
// database connection
type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
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

// ConnectionInfo prints out the connection line used for PostgreSQL
func (pc PostgresConfig) ConnectionInfo() string {
	if pc.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s "+
			"sslmode=disable", pc.Host, pc.Port, pc.User, pc.Name)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s "+
		"sslmode=disable", pc.Host, pc.Port, pc.User, pc.Password, pc.Name)
}

// Dialect represent the database type
func (pc PostgresConfig) Dialect() string {
	return "postgres"
}
