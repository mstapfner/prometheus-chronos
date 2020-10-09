package main

import (
	gokitlog "github.com/go-kit/kit/log"
	"github.com/prometheus/prometheus/tsdb"
	"log"
	"os"
)

func ExportCustomBlock(block Block, outputDir string) {

	blockStartTime := block.MinTime
	blockEndTime := block.MaxTime

	// Create MetricSamples
	var metricSamples []*tsdb.MetricSample
	for _, ts := range block.CustomTimeSeries {
		for _, ch := range ts.Chunks {
			for i, tstmp := range ch.TimeStamps {
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

	outputFile, err := tsdb.CreateBlock(metricSamples, outputDir, blockStartTime, blockEndTime, logger)
	if err != nil {
		log.Println(err)
	}
	log.Println(outputFile)
}
