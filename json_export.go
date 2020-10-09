package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func ExportBlocksAsJSON(blocks []Block, jsonOutputDir string) {

	log.Println(jsonOutputDir)

	// Create output directory if not exists
	if _, err := os.Stat(jsonOutputDir); os.IsNotExist(err) {
		os.Mkdir(jsonOutputDir, 0644)
	}

	for _, block := range blocks {
		path := filepath.Join(jsonOutputDir, block.ULID)
		path += ".json"
		file, err := json.MarshalIndent(block, "", " ")
		err = ioutil.WriteFile(path, file, 0644)
		if err != nil {
			log.Println("Failed to write Block as JSON File: ", path)
		}
	}

}
