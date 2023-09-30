package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/google/go-cmp/cmp"
	"github.com/tkuchiki/alp/internal/testutil"
	"github.com/tkuchiki/alp/options"
)

func Test_createOptionsFromConfig(t *testing.T) {
	viper.Reset()
	command := NewCommand("test")

	tempDir := t.TempDir()
	sort := "max"
	dummyOpts := testutil.DummyOptions(sort)

	var err error
	command.flags.config, err = testutil.CreateTempDirAndFile(tempDir, "test_create_options_from_config_config", testutil.DummyConfigFile(sort, dummyOpts))
	if err != nil {
		t.Fatal(err)
	}

	var opts *options.Options
	opts, err = command.flags.createOptionsFromConfig(command.rootCmd)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(dummyOpts, opts); diff != "" {
		t.Errorf("%s", diff)
	}
}

func Test_createOptionsFromConfig_overwrite(t *testing.T) {
	command := NewCommand("test")

	tempDir := t.TempDir()
	sort := "max"

	overwrittenSort := "min"
	overwrittenOpts := testutil.DummyOverwrittenOptions(overwrittenSort)

	var err error
	command.flags.config, err = testutil.CreateTempDirAndFile(tempDir, "test_create_options_from_config_overwrite_config", testutil.DummyConfigFile(sort, overwrittenOpts))
	if err != nil {
		t.Fatal(err)
	}

	viper.Set("file", overwrittenOpts.File)
	viper.Set("dump", overwrittenOpts.Dump)
	viper.Set("load", overwrittenOpts.Load)
	viper.Set("sort", overwrittenSort)
	viper.Set("reverse", overwrittenOpts.Reverse)
	viper.Set("query_string", overwrittenOpts.QueryString)
	viper.Set("query_string_ignore_values", overwrittenOpts.QueryStringIgnoreValues)
	viper.Set("decode_uri", overwrittenOpts.DecodeUri)
	viper.Set("format", overwrittenOpts.Format)
	viper.Set("noheaders", overwrittenOpts.NoHeaders)
	viper.Set("show_footers", overwrittenOpts.ShowFooters)
	viper.Set("limit", overwrittenOpts.Limit)
	viper.Set("matching_groups", strings.Join(overwrittenOpts.MatchingGroups, ","))
	viper.Set("filters", overwrittenOpts.Filters)
	viper.Set("pos_file", overwrittenOpts.PosFile)
	viper.Set("nosave_pos", overwrittenOpts.NoSavePos)
	viper.Set("location", overwrittenOpts.Location)
	viper.Set("output", overwrittenOpts.Output)
	viper.Set("percentiles", testutil.IntSliceToString(overwrittenOpts.Percentiles))
	viper.Set("pagination_limit", overwrittenOpts.PaginationLimit)

	// json
	viper.Set("json.uri_key", overwrittenOpts.JSON.UriKey)
	viper.Set("json.method_key", overwrittenOpts.JSON.MethodKey)
	viper.Set("json.time_key", overwrittenOpts.JSON.TimeKey)
	viper.Set("json.restime_key", overwrittenOpts.JSON.ResponseTimeKey)
	viper.Set("json.reqtime_key", overwrittenOpts.JSON.RequestTimeKey)
	viper.Set("json.body_bytes_key", overwrittenOpts.JSON.BodyBytesKey)
	viper.Set("json.status_key", overwrittenOpts.JSON.StatusKey)

	// ltsv
	viper.Set("ltsv.uri_label", overwrittenOpts.LTSV.UriLabel)
	viper.Set("ltsv.method_label", overwrittenOpts.LTSV.MethodLabel)
	viper.Set("ltsv.time_label", overwrittenOpts.LTSV.TimeLabel)
	viper.Set("ltsv.apptime_label", overwrittenOpts.LTSV.ApptimeLabel)
	viper.Set("ltsv.reqtime_label", overwrittenOpts.LTSV.ReqtimeLabel)
	viper.Set("ltsv.size_label", overwrittenOpts.LTSV.SizeLabel)
	viper.Set("ltsv.status_label", overwrittenOpts.LTSV.StatusLabel)

	// regexp
	viper.Set("regexp.pattern", overwrittenOpts.Regexp.Pattern)
	viper.Set("regexp.uri_subexp", overwrittenOpts.Regexp.UriSubexp)
	viper.Set("regexp.method_subexp", overwrittenOpts.Regexp.MethodSubexp)
	viper.Set("regexp.time_subexp", overwrittenOpts.Regexp.TimeSubexp)
	viper.Set("regexp.restime_subexp", overwrittenOpts.Regexp.ResponseTimeSubexp)
	viper.Set("regexp.reqtime_subexp", overwrittenOpts.Regexp.RequestTimeSubexp)
	viper.Set("regexp.body_bytes_subexp", overwrittenOpts.Regexp.BodyBytesSubexp)
	viper.Set("regexp.status_subexp", overwrittenOpts.Regexp.StatusSubexp)

	// pcap
	viper.Set("pcap.server_ips", strings.Join(overwrittenOpts.Pcap.ServerIPs, ","))
	viper.Set("pcap.server_port", overwrittenOpts.Pcap.ServerPort)

	// count
	viper.Set("count.keys", overwrittenOpts.Count.Keys)

	// count
	viper.Set("topN.sort", overwrittenOpts.TopN.Sort)
	viper.Set("topN.reverse", overwrittenOpts.TopN.Reverse)

	var opts *options.Options
	opts, err = command.flags.createOptionsFromConfig(command.rootCmd)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(overwrittenOpts, opts); diff != "" {
		t.Errorf("%s", diff)
	}
}
