package main

import (
	"strings"
)

type Config struct {
	File              string   `yaml:"file"`
	Sort              string   `yaml:"sort"`
	Reverse           bool     `yaml:"reverse"`
	QueryString       bool     `yaml:"query_string"`
	Tsv               bool     `yaml:"tsv"`
	ApptimeLabel      string   `yaml:"apptime_label"`
	SizeLabel         string   `yaml:"size_label"`
	MethodLabel       string   `yaml:"method_label"`
	UriLabel          string   `yaml:"uri_label"`
	TimeLabel         string   `yaml:"time_label"`
	Limit             int      `yaml:"limit"`
	Includes          []string `yaml:"includes"`
	Excludes          []string `yaml:"excludes"`
	NoHeaders         bool     `yaml:no_headers`
	Aggregates        []string `yaml:"aggregates"`
	StartTime         string   `yaml:"start_time"`
	EndTime           string   `yaml:"end_time"`
	StartTimeDuration string   `yaml:"start_time_duration"`
	EndTimeDuration   string   `yaml:"end_time_duration"`
	IncludesStr       string
	ExcludesStr       string
	AggregatesStr     string
}

func SetConfig(config Config, arg Config) Config {
	if arg.File != "" {
		config.File = arg.File
	}

	if arg.Reverse {
		config.Reverse = arg.Reverse
	}

	if arg.QueryString {
		config.QueryString = arg.QueryString
	}

	if arg.Tsv {
		config.Tsv = arg.Tsv
	}

	if config.ApptimeLabel == "" {
		config.ApptimeLabel = arg.ApptimeLabel
	}

	if config.SizeLabel == "" {
		config.SizeLabel = arg.SizeLabel
	}

	if config.MethodLabel == "" {
		config.MethodLabel = arg.MethodLabel
	}

	if config.UriLabel == "" {
		config.UriLabel = arg.UriLabel
	}

	if config.TimeLabel == "" {
		config.TimeLabel = arg.TimeLabel
	}

	if config.Limit == 0 {
		config.Limit = arg.Limit
	}

	if arg.IncludesStr != "" {
		config.Includes = strings.Split(arg.IncludesStr, ",")
	}

	if arg.ExcludesStr != "" {
		config.Excludes = strings.Split(arg.ExcludesStr, ",")
	}

	if arg.NoHeaders {
		config.NoHeaders = arg.NoHeaders
	}

	if arg.AggregatesStr != "" {
		config.Aggregates = strings.Split(arg.AggregatesStr, ",")
	}

	if arg.StartTime != "" {
		config.StartTime = arg.StartTime
	}

	if arg.EndTime != "" {
		config.EndTime = arg.EndTime
	}

	if arg.StartTimeDuration != "" {
		config.StartTimeDuration = arg.StartTimeDuration
	}

	if arg.EndTimeDuration != "" {
		config.EndTimeDuration = arg.EndTimeDuration
	}

	return config
}
