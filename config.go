package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	File               string   `yaml:"file"`
	Sort               string   `yaml:"sort"`
	Reverse            bool     `yaml:"reverse"`
	QueryString        bool     `yaml:"query_string"`
	Tsv                bool     `yaml:"tsv"`
	ApptimeLabel       string   `yaml:"apptime_label"`
	ReqtimeLabel       string   `yaml:"reqtime_label"`
	StatusLabel        string   `yaml:"status_label"`
	SizeLabel          string   `yaml:"size_label"`
	MethodLabel        string   `yaml:"method_label"`
	UriLabel           string   `yaml:"uri_label"`
	TimeLabel          string   `yaml:"time_label"`
	Limit              int      `yaml:"limit"`
	Includes           []string `yaml:"includes"`
	Excludes           []string `yaml:"excludes"`
	IncludeStatuses    []string `yaml:"include_statuses"`
	ExcludeStatuses    []string `yaml:"exclude_statuses"`
	NoHeaders          bool     `yaml:no_headers`
	Aggregates         []string `yaml:"aggregates"`
	StartTime          string   `yaml:"start_time"`
	EndTime            string   `yaml:"end_time"`
	StartTimeDuration  string   `yaml:"start_time_duration"`
	EndTimeDuration    string   `yaml:"end_time_duration"`
	IncludesStr        string
	ExcludesStr        string
	IncludeStatusesStr string
	ExcludeStatusesStr string
	AggregatesStr      string
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

	if config.ApptimeLabel == "" || (config.ApptimeLabel != "" && arg.ApptimeLabel != ApptimeLabel) {
		config.ApptimeLabel = arg.ApptimeLabel
	}

	if config.ReqtimeLabel == "" || (config.ReqtimeLabel != "" && arg.ReqtimeLabel != ReqtimeLabel) {
		config.ReqtimeLabel = arg.ReqtimeLabel
	}

	if config.StatusLabel == "" || (config.StatusLabel != "" && arg.StatusLabel != StatusLabel) {
		config.StatusLabel = arg.StatusLabel
	}

	if config.SizeLabel == "" || (config.SizeLabel != "" && arg.SizeLabel != SizeLabel) {
		config.SizeLabel = arg.SizeLabel
	}

	if config.MethodLabel == "" || (config.MethodLabel != "" && arg.MethodLabel != MethodLabel) {
		config.MethodLabel = arg.MethodLabel
	}

	if config.UriLabel == "" || (config.UriLabel != "" && arg.UriLabel != UriLabel) {
		config.UriLabel = arg.UriLabel
	}

	if config.TimeLabel == "" || (config.TimeLabel != "" && arg.TimeLabel != TimeLabel) {
		config.TimeLabel = arg.TimeLabel
	}

	if config.Limit == 0 || (config.Limit != 0 && arg.Limit != Limit) {
		config.Limit = arg.Limit
	}

	if arg.IncludesStr != "" {
		config.Includes = Split(arg.IncludesStr, ",")
	}

	if arg.ExcludesStr != "" {
		config.Excludes = Split(arg.ExcludesStr, ",")
	}

	if arg.IncludeStatusesStr != "" {
		config.IncludeStatuses = Split(arg.IncludeStatusesStr, ",")
	}

	if arg.ExcludeStatusesStr != "" {
		config.ExcludeStatuses = Split(arg.ExcludeStatusesStr, ",")
	}

	if arg.NoHeaders {
		config.NoHeaders = arg.NoHeaders
	}

	if arg.AggregatesStr != "" {
		config.Aggregates = Split(arg.AggregatesStr, ",")
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

func LoadYAML(filename string) (config Config, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(buf, &config)

	return config, err
}

func DumpProfiles(filename string, ps Profiles) (err error) {
	buf, err := yaml.Marshal(&ps)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, buf, os.ModePerm)

	return err
}

func LoadProfiles(filename string) (ps Profiles, err error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return ps, err
	}

	err = yaml.Unmarshal(buf, &ps)

	return ps, err
}
