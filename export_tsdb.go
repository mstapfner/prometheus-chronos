package main

import (
	gokitlog "github.com/go-kit/kit/log"
	"github.com/prometheus/prometheus/tsdb"
	"log"
	"os"
	"strconv"
	"time"
)

func ExportBlocks(blocks []Block, outputDir string, redateStart string, redateEnd string) {
	if redateStart == "" && redateEnd == "" {
		createBlocks(blocks, outputDir, 0)
		return
	}

	redateStartInt, err := strconv.ParseInt(redateStart, 10, 64)
	startTime := time.Unix(redateStartInt, 0)
	if err != nil {
		log.Println("Failed to parse the starting timestamp")
	}
	redateEndInt, err := strconv.ParseInt(redateEnd, 10, 64)
	endTime := time.Unix(redateEndInt, 0)
	if err != nil {
		log.Println("Failed to parse the ending timestamp")
	}

	log.Println("Shift blocks to: ", startTime, " - ", endTime)

	interval := redateEndInt - redateStartInt
	if interval <= 0 {
		// TODO: outsource all flag checks
		log.Println("Your specified ending time is bigger than your specified starting time")
		return
	}

	// Calculate starting and ending time from all blocks
	var blockStartingTime int64 = 0
	var blockEndingTime int64 = 0

	for _, block := range blocks {
		if blockStartingTime > block.StartingTime {
			blockStartingTime = block.StartingTime
		}
		if blockEndingTime < block.EndingTime {
			blockEndingTime = block.EndingTime
		}
	}

	allBlockInterval := blockEndingTime - blockStartingTime
	if allBlockInterval <= 0 {
		log.Println("Problem: block end time bigger than block start time")
		return
	}

	// Calculate multipliers
	counterNewBlocks := allBlockInterval / interval
	if counterNewBlocks > 0 {
		var i int64
		// Create the needed blocks to fill the missing timeslot
		for i = 0; i < counterNewBlocks; i++ {
			// TODO: add calculation of shift interval
			createBlocks(blocks, outputDir, 0)
		}
	}

}

func createBlocks(blocks []Block, outputDir string, shiftInterval int64) {
	for _, block := range blocks {
		blockStartTime := block.StartingTime
		blockEndTime := block.EndingTime

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

		_, err := tsdb.CreateBlock(metricSamples, outputDir, blockStartTime, blockEndTime, logger)
		if err != nil {
			log.Println(err)
		}

	}
}
