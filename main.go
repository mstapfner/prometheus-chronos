package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
)

var (
	exportTSDB    = flag.Bool("exportTSDB", false, "Flag for the import behavior")
	importDir     = flag.String("importDir", "", "")
	outputDir     = flag.String("outputDir", "", "")
	jsonOutput    = flag.Bool("jsonOutput", false, "")
	jsonOutputDir = flag.String("jsonOutputDir", "", "")
	redateStart   = flag.String("redateStart", "", "")
	redateEnd     = flag.String("redateEnd", "", "")
)

func main() {
	log.Println("Chronos")

	flag.Parse()

	// Import Blocks from import dir
	var blocks []Block
	files, err := ioutil.ReadDir(*importDir)
	if err != nil {
		log.Println("Failed to read the import directory: ", *importDir)
		log.Println(err.Error())
	}
	for _, file := range files {
		if file.IsDir() && file.Name() != "wal" && file.Name() != "chunks_head" {
			path := filepath.Join(*importDir, file.Name())
			block, err := ImportTSDB(path)
			if err != nil {
				log.Println("Failed to read the block: ", path)
				log.Println(err.Error())
			}
			blocks = append(blocks, block)
		}
	}

	// Output as JSON files
	log.Println(*jsonOutputDir)
	if *jsonOutput {
		ExportBlocksAsJSON(blocks, *jsonOutputDir)
	}

	if !*exportTSDB {
		ExportBlocks(blocks, *outputDir, *redateStart, *redateEnd)
	}
}
