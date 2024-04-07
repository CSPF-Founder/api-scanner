package views

import (
	"encoding/json"
	"log"

	"github.com/CSPF-Founder/api-scanner/code/panel/frontend"
)

type ManifestEntry struct {
	File    string `json:"file"`
	IsEntry bool   `json:"isEntry"`
	Src     string `json:"src"`
}

var assetMap map[string]ManifestEntry

func init() {
	loadManifest()
}

func loadManifest() {
	// Read the manifest file
	data, err := frontend.FileSystem.ReadFile("static/.vite/manifest.json")
	if err != nil {
		log.Fatalf("Error reading manifest file: %v", err)
	}

	// Unmarshal JSON data into the map
	err = json.Unmarshal(data, &assetMap)
	if err != nil {
		log.Fatalf("Error unmarshalling manifest data: %v", err)
	}
}

func AssetPath(originalName string) string {
	// Look up the original name in the map
	entry, exists := assetMap[originalName]
	if !exists {
		// Handle the case where the file is not found in the manifest
		log.Printf("Asset not found: %s", originalName)
		return originalName // or return an error/placeholder
	}
	return "/static/" + entry.File
}
