package storage

import (
	"encoding/json"
	"log"
	"os"
)

var DNSRegistry = map[string]string{}
const dnsRegistryFile = "data/dnsRegistry.json"

func LoadDNSRegistry() {
	data, err := os.ReadFile(dnsRegistryFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("No dnsRegistry.json found, using empty registry.")
			return
		}
		log.Fatalf("Failed to read DNS registry file: %v", err)
	}

	err = json.Unmarshal(data, &DNSRegistry)
	if err != nil {
		log.Fatalf("Failed to parse DNS registry: %v", err)
	}
}

func SaveDNSRegistry() {
	data, err := json.MarshalIndent(DNSRegistry, "", "  ")
	if err != nil {
		log.Printf("Failed to serialize DNS registry: %v", err)
		return
	}

	err = os.WriteFile(dnsRegistryFile, data, 0644)
	if err != nil {
		log.Printf("Failed to write DNS registry file: %v", err)
	}
}
