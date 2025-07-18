package storage

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"posServer/internal/models"
)

var Bindings = map[string]models.Binding{}

const bindingsFile = "bindings.json"

func LoadBindings() {
	data, err := ioutil.ReadFile(bindingsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return // No file = first run
		}
		log.Fatalf("Failed to read bindings file: %v", err)
	}

	err = json.Unmarshal(data, &Bindings)
	if err != nil {
		log.Fatalf("Failed to parse bindings JSON: %v", err)
	}
}

func SaveBindings() {
	data, err := json.MarshalIndent(Bindings, "", "  ")
	if err != nil {
		log.Printf("Failed to serialize bindings: %v", err)
		return
	}

	err = ioutil.WriteFile(bindingsFile, data, 0644)
	if err != nil {
		log.Printf("Failed to write bindings file: %v", err)
	}
}
