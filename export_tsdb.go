package main

import (
	gokitlog "github.com/go-kit/kit/log"
	"github.com/prometheus/prometheus/tsdb"
	"log"
	"os"
	"time"
)

func ExportBlocks(blocks []Block, outputDir string, redateStart int, redateEnd int) {
	if redateStart == 0 && redateEnd == 0 {
		createBlocks(blocks, outputDir, 0)
		return
	}

	log.Println("OUTPUT DIR")
	log.Println(outputDir)

	startTime := time.Unix(int64(redateStart/1000), 0)

	endTime := time.Unix(int64(redateEnd/1000), 0)

	log.Println("Shift blocks to: ", startTime, " - ", endTime)

	interval := int64(redateEnd) - int64(redateStart)
	if interval <= 0 {
		// TODO: outsource all flag checks
		log.Println("Your specified ending time is bigger than your specified starting time")
		return
	}

	// Calculate starting and ending time from all blocks
	var blockStartingTime int64 = int64(redateEnd)
	var blockEndingTime int64 = int64(redateStart)

	// log.Println("Start block")
	// log.Println(blockStartingTime)
	// log.Println("End Block")
	// log.Println(blockEndingTime)

	for _, block := range blocks {
		if blockStartingTime > block.StartingTime {
			blockStartingTime = block.StartingTime
		}
		if blockEndingTime < block.EndingTime {
			blockEndingTime = block.EndingTime
		}
	}

	tsdbStart := time.Unix(int64(blockStartingTime/1000), 0)

	tsdbEnd := time.Unix(int64(blockEndingTime/1000), 0)

	log.Println("Current tsdb timespan: ", tsdbStart, " - ", tsdbEnd)

	allBlockInterval := blockEndingTime - blockStartingTime
	if allBlockInterval <= 0 {
		log.Println("Problem: block end time bigger than block start time")
		return
	}

	log.Println("Start block")
	log.Println(blockStartingTime)
	log.Println("End Block")
	log.Println(blockEndingTime)

	// Calculate multipliers
	counterNewBlocks := interval / allBlockInterval

	log.Println("allblock interval")
	log.Println(allBlockInterval)

	log.Println("Interval")
	log.Println(interval)

	log.Println("counter new blocks")
	log.Println(counterNewBlocks)

	log.Println("blocks")
	log.Println(len(blocks))

	if counterNewBlocks > 0 {
		var i int64
		// Create the needed blocks to fill the missing timeslot
		for i = 0; i < counterNewBlocks; i++ {
			shift := interval / counterNewBlocks
			shift = shift * (i + 1)
			log.Println("SHIFT")
			log.Println(shift)
			createBlocks(blocks, outputDir, shift)
		}
	}

}

func createBlocks(blocks []Block, outputDir string, shiftInterval int64) {
	log.Println("Outputdir")
	log.Println(outputDir)
	for _, block := range blocks {
		blockStartTime := block.StartingTime - shiftInterval
		blockEndTime := block.EndingTime - shiftInterval

		// Create MetricSamples
		var metricSamples []*tsdb.MetricSample
		for _, ts := range block.CustomTimeSeries {
			for _, ch := range ts.Chunks {
				for i, tstmp := range ch.TimeStamps {
					if shiftInterval > 0 {
						tstmp = tstmp - shiftInterval
					}
					metricSample := tsdb.MetricSample{
						TimestampMs: tstmp,
						Value:       ch.Values[i],
						Labels:      ts.Labels,
					}
					metricSamples = append(metricSamples, &metricSample)
				}
			}
		}

		logger := gokitlog.NewLogfmtLogger(os.Stderr)

		result, err := tsdb.CreateBlock(metricSamples, outputDir, blockStartTime, blockEndTime, logger)
		if err != nil {
			log.Println("Failed to create block in outputDir: ", outputDir)
			log.Println(err)
		}
		log.Println(result)

	}
}
