package main

import (
	gokitlog "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"github.com/prometheus/prometheus/tsdb/chunks"
	"log"
	"math"
	"os"
)

func ImportTSDB(blockPath string) (Block, error) {

	logger := gokitlog.NewLogfmtLogger(os.Stderr)

	var newBlock Block

	block, err := tsdb.OpenBlock(logger, blockPath, chunkenc.NewPool())
	if err != nil {
		log.Println("Failed to open Block", err)
		return Block{}, errors.Wrap(err, "tsdb.OpenBlock")
	}

	metaInfo := block.Meta()
	newBlock.NumSeries = metaInfo.Stats.NumSeries
	newBlock.NumChunks = metaInfo.Stats.NumChunks
	newBlock.NumSamples = metaInfo.Stats.NumSamples
	newBlock.NumTombstones = metaInfo.Stats.NumTombstones

	newBlock.MaxTime = metaInfo.MaxTime
	newBlock.MinTime = metaInfo.MinTime

	newBlock.ULID = metaInfo.ULID.String()

	indexr, err := block.Index()
	if err != nil {
		return Block{}, errors.Wrap(err, "block.Index")
	}
	defer indexr.Close()

	newBlock.LabelNames, err = indexr.LabelNames()
	newBlock.Symbols = indexr.Symbols()

	newBlock.Postings, err = indexr.Postings("", "")

	chunkr, err := block.Chunks()

	if err != nil {
		return Block{}, errors.Wrap(err, "block.Chunks")
	}

	var it chunkenc.Iterator

	for newBlock.Postings.Next() {
		var customTimeSeries TimeSeries

		ref := newBlock.Postings.At()
		lset := labels.Labels{}
		chks := []chunks.Meta{}

		if err := indexr.Series(ref, &lset, &chks); err != nil {
			return Block{}, errors.Wrap(err, "index.Series")
		}

		customTimeSeries.Ref = ref
		customTimeSeries.Labels = lset
		for _, meta := range chks {
			chunk, err := chunkr.Chunk(meta.Ref)
			if err != nil {
				return Block{}, errors.Wrap(err, "chunkr.Chunk")
			}
			var customChunk Chunk
			customChunk.Ref = meta.Ref
			customChunk.NumSamples = chunk.NumSamples()

			it := chunk.Iterator(it)
			for it.Next() {
				t, v := it.At()
				if math.IsNaN(v) {
					continue
				}
				if math.IsInf(v, -1) || math.IsInf(v, 1) {
					continue
				}
				customChunk.TimeStamps = append(customChunk.TimeStamps, t)
				customChunk.Values = append(customChunk.Values, v)
			}

			if it.Err() != nil {
				return Block{}, errors.Wrap(err, "iterator.Err")
			}
			if len(customChunk.TimeStamps) == 0 {
				continue
			}
			customTimeSeries.Chunks = append(customTimeSeries.Chunks, customChunk)
		}

		newBlock.CustomTimeSeries = append(newBlock.CustomTimeSeries, customTimeSeries)

	}

	return newBlock, nil
}
