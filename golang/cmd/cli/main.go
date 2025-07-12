package main

import (
	"log"
	"os"

	"github.com/damianko135/curseforge-autoupdate/golang/helper"
	"github.com/damianko135/curseforge-autoupdate/golang/internal/api"
	"github.com/pelletier/go-toml"
)

func main() {
	apiKey, modID := initVars()

	client := api.NewClient(apiKey)
	exists, err := client.CheckIfExists(modID)
	if err != nil {
		log.Fatalf("Error checking if mod exists: %v", err)
	}
	if exists {
		log.Printf("Mod with ID %d found!", modID)
	} else {
		log.Fatalf("Mod with ID %d not found!", modID)
	}

}

func createTomlFile() {
	// Example of creating a TOML file using the pelletier/go-toml library
	tomlData := map[string]interface{}{
		"key": "value",
	}
	data, err := toml.Marshal(tomlData)
	if err != nil {
		log.Fatalf("Error marshalling TOML: %v", err)
	}
	err = os.WriteFile("config.toml", data, 0644)
	if err != nil {
		log.Fatalf("Error writing TOML file: %v", err)
	}
}

func initVars() (apiKey string, modID int) {
	apiKey = helper.GetEnvVar("API_KEY", "")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is not set")
	}

	modID = helper.GetIntVar("MOD_ID", 0)
	if modID == 0 {
		log.Fatal("MOD_ID environment variable is not set or invalid")
	}
	return apiKey, modID
}
