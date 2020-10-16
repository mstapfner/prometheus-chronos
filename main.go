package main

import (
	"flag"
	"github.com/cheggaaa/pb"
	"io/ioutil"
	"log"
	"path/filepath"
)

func main() {
	log.Println("Chronos")

	var exportTSDB bool
	var importDir string
	var outputDir string
	var jsonOutput bool
	var jsonOutputDir string
	var redateStart int
	var redateEnd int

	flag.BoolVar(&exportTSDB, "exportTSDB", false, "")
	flag.StringVar(&importDir, "importDir", "", "")
	flag.StringVar(&outputDir, "outputDir", "", "")
	flag.BoolVar(&jsonOutput, "jsonOutput", false, "")
	flag.StringVar(&jsonOutputDir, "jsonOutputDir", "", "")
	flag.IntVar(&redateStart, "redateStart", 0, "")
	flag.IntVar(&redateEnd, "redateEnd", 0, "")

	flag.Parse()

	log.Println("Redate Start")
	log.Println(redateStart)

	// Import Blocks from import dir
	var blocks []Block
	files, err := ioutil.ReadDir(importDir)
	if err != nil {
		log.Println("Failed to read the import directory: ", importDir)
		log.Println(err.Error())
	}

	amountOfFiles := len(files) - 4

	log.Println("Detected ", amountOfFiles, " blocks")
	bar := pb.StartNew(amountOfFiles)

	for _, file := range files {
		if file.IsDir() && file.Name() != "wal" && file.Name() != "chunks_head" && file.Name() != "queries.active" && file.Name() != "lock" {
			path := filepath.Join(importDir, file.Name())
			block, err := ImportTSDB(path)
			if err != nil {
				log.Println("Failed to read the block: ", path)
				log.Println(err.Error())
			}
			blocks = append(blocks, block)
			bar.Increment()
		}
	}
	bar.Finish()

	// Output as JSON files
	if jsonOutput {
		log.Println("Start JSON Export")
		ExportBlocksAsJSON(blocks, jsonOutputDir)
	}

	if exportTSDB {
		log.Println("Start TSDB Export")
		ExportBlocks(blocks, outputDir, redateStart, redateEnd)
	}
}
