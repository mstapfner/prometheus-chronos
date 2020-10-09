package main

import (
	"github.com/prometheus/prometheus/pkg/labels"
	"github.com/prometheus/prometheus/tsdb/index"
)

type Block struct {
	MinTime int64 `json:"min_time"`
	MaxTime int64 `json:"max_time"`

	ULID string `json:"ulid"`

	LabelValues       []string         `json:"label_values,omitempty"`
	SortedLabelValues []string         `json:"sorted_label_values,omitempty"`
	Symbols           index.StringIter `json:"symbols,omitempty"`

	LabelNames []string `json:"label_names,omitempty"`

	Postings       index.Postings `json:"postings,omitempty"`
	SortedPostings index.Postings `json:"sorted_postings,omitempty"`

	NumSamples    uint64 `json:"num_samples,omitempty"`
	NumSeries     uint64 `json:"num_series,omitempty"`
	NumChunks     uint64 `json:"num_chunks,omitempty"`
	NumTombstones uint64 `json:"num_tombstones,omitempty"`

	CustomTimeSeries []TimeSeries `json:"custom_time_series"`
}

type TimeSeries struct {
	Ref    uint64        `json:"ref"`
	Labels labels.Labels `json:"labels"`
	Chunks []Chunk       `json:"chunks"`
}

type Chunk struct {
	Ref        uint64 `json:"ref"`
	NumSamples int    `json:"num_samples"`

	TimeStamps []int64   `json:"time_stamps"`
	Values     []float64 `json:"values"`
}
