package cmd

import (
	"io"
	"os"
	"testing"

	"github.com/tkuchiki/alp/internal/testutil"
)

func TestProfileCommandGoldenOutputs(t *testing.T) {
	const accessLogExpected = `[["count","method","uri","min","max","avg_body"],[1,"POST","/hoge/piyo",0.234,0.234,34],[1,"GET","/diary/entry/1234",0.135,0.135,15],[1,"GET","/diary/entry/5678",0.432,0.432,30],[1,"GET","/foo/bar/5xx",60,60,15],[1,"GET","/req",0.321,0.321,15],[2,"GET","/foo/bar",0.123,0.123,56],[5,"POST","/foo/bar",0.057,0.234,25.2]]
`

	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "JSON",
			args: []string{"json", "--file", "../../../example/logs/json_access.log", "--format", "json", "--output", "count,method,uri,min,max,avg_body"},
			want: accessLogExpected,
		},
		{
			name: "LTSV",
			args: []string{"ltsv", "--file", "../../../example/logs/ltsv_access.log", "--format", "json", "--output", "count,method,uri,min,max,avg_body"},
			want: accessLogExpected,
		},
		{
			name: "pcap",
			args: []string{"pcap", "--file", "../../../example/logs/http.cap", "--pcap-server-port", "18080", "--format", "json", "--output", "count,method,uri,min,max,avg_body"},
			want: `[["count","method","uri","min","max","avg_body"],[1,"GET","/foo/bar/123",0.009,0.009,153],[1,"GET","/foo/bar/456",0.007,0.007,153]]
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := executeForOutput(t, test.args); got != test.want {
				t.Fatalf("output:\n%s\nwant:\n%s", got, test.want)
			}
		})
	}
}

func TestRegexpProfileCommandGoldenOutput(t *testing.T) {
	keys := testutil.NewRegexpLogKeys()
	logFile, err := testutil.CreateTempDirAndFile(t.TempDir(), "access.log", testutil.RegexpLog())
	if err != nil {
		t.Fatal(err)
	}

	args := []string{
		"regexp",
		"--file", logFile,
		"--pattern", testutil.RegexpPattern(keys),
		"--format", "json",
		"--output", "count,method,uri,min,max,avg_body",
	}
	want := `[["count","method","uri","min","max","avg_body"],[1,"POST","/foo/bar/123",0.057,0.057,12],[1,"POST","/foo/bar/456",0.1,0.1,34]]
`

	if got := executeForOutput(t, args); got != want {
		t.Fatalf("output:\n%s\nwant:\n%s", got, want)
	}
}

func TestTopNPreservesOriginalTimestamp(t *testing.T) {
	tests := []struct {
		name string
		time string
	}{
		{name: "explicit offset", time: "2015-09-06T05:58:05Z"},
		{name: "unusual value", time: "not-a-time"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			input := `{"time":"` + test.time + `","method":"GET","uri":"/timezone","status":200,"body_bytes":1,"response_time":0.001}`
			logFile, err := testutil.CreateTempDirAndFile(t.TempDir(), "access.json", input)
			if err != nil {
				t.Fatal(err)
			}

			got := executeForOutput(t, []string{"json", "topN", "1", "--file", logFile, "--format", "json"})
			want := `[["rank","uri","method","status","restime","bytes","time"],[1,"/timezone","GET",200,0.001,1,"` + test.time + `"]]
`
			if got != want {
				t.Fatalf("output:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

func TestCountCommandGoldenOutput(t *testing.T) {
	args := []string{"json", "count", "--file", "../../../example/logs/json_access.log", "--keys", "method", "--format", "json"}
	want := `[["sum","method"],[6,"POST"],[6,"GET"]]
`

	if got := executeForOutput(t, args); got != want {
		t.Fatalf("output:\n%s\nwant:\n%s", got, want)
	}
}

func TestLoadLegacyDumpGoldenOutput(t *testing.T) {
	args := []string{"json", "--load", "../../../example/logs/dump1.yaml", "--format", "json", "--output", "count,uri,avg_body"}
	want := `[["count","uri","avg_body"],[1,"/hoge/piyo",34],[1,"/diary/entry/1234",15],[1,"/diary/entry/5678",30],[1,"/foo/bar/5xx",15],[1,"/req",15],[2,"/foo/bar",56],[5,"/foo/bar",25.2]]
`

	if got := executeForOutput(t, args); got != want {
		t.Fatalf("output:\n%s\nwant:\n%s", got, want)
	}
}

func TestFilterAcceptsTimeOnRight(t *testing.T) {
	args := []string{
		"json",
		"--file", "../../../example/logs/json_access.log",
		"--format", "json",
		"--output", "count,uri",
		"--filters", `"2015-09-06T05:00:00+09:00" < Time`,
	}
	want := `[["count","uri"],[1,"/hoge/piyo"],[1,"/diary/entry/1234"],[1,"/diary/entry/5678"],[1,"/foo/bar/5xx"],[1,"/req"],[2,"/foo/bar"],[5,"/foo/bar"]]
`

	if got := executeForOutput(t, args); got != want {
		t.Fatalf("output:\n%s\nwant:\n%s", got, want)
	}
}

func executeForOutput(t *testing.T, args []string) string {
	t.Helper()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	originalStdout := os.Stdout
	os.Stdout = writer

	command := NewCommand("test")
	command.setArgs(args)
	executeErr := command.Execute()
	closeErr := writer.Close()
	os.Stdout = originalStdout

	output, readErr := io.ReadAll(reader)
	_ = reader.Close()
	if executeErr != nil {
		t.Fatal(executeErr)
	}
	if closeErr != nil {
		t.Fatal(closeErr)
	}
	if readErr != nil {
		t.Fatal(readErr)
	}

	return string(output)
}
