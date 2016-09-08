package main

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"math"
	"os"
	"strings"
	"time"
)

func Round(f float64) string {
	return fmt.Sprintf("%.3f", f)
}

func Output(ps Profiles, c Config) {
	if c.Tsv {
		if !c.NoHeaders {
			fmt.Printf("Count\tMin\tMax\tSum\tAvg\tP1\tP50\tP99\tStddev\tMin(Body)\tMax(Body)\tSum(Body)\tAvg(Body)\tMethod\tUri")
			fmt.Println("")
		}

		for _, p := range ps {
			fmt.Printf("%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v\t%v",
				p.Cnt, Round(p.Min), Round(p.Max), Round(p.Sum), Round(p.Avg),
				Round(p.P1), Round(p.P50), Round(p.P99), Round(p.Stddev),
				Round(p.MinBody), Round(p.MaxBody), Round(p.SumBody), Round(p.AvgBody),
				p.Method, p.Uri)
			fmt.Println("")
		}
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Count", "Min", "Max", "Sum", "Avg",
			"P1", "P50", "P99", "Stddev",
			"Min(Body)", "Max(Body)", "Sum(Body)", "Avg(Body)",
			"Method", "Uri"})
		for _, p := range ps {
			data := []string{
				fmt.Sprint(p.Cnt), Round(p.Min), Round(p.Max), Round(p.Sum), Round(p.Avg),
				Round(p.P1), Round(p.P50), Round(p.P99), Round(p.Stddev),
				Round(p.MinBody), Round(p.MaxBody), Round(p.SumBody), Round(p.AvgBody),
				p.Method, p.Uri}
			table.Append(data)
		}
		table.Render()
	}
}

func SetCursor(index string, uri string) {
	if _, ok := uriHints[index]; ok {
		cursor = uriHints[index]
	} else {
		uriHints[index] = length
		cursor = length
		length++
		accessLog = append(accessLog, Profile{Uri: uri})
	}
}

func TimeCmp(startTimeNano int64, endTimeNano int64, timeNano int64) bool {
	if startTimeNano > 0 && endTimeNano == 0 {
		return startTimeNano <= timeNano
	} else if endTimeNano > 0 && startTimeNano == 0 {
		return endTimeNano >= timeNano
	} else if startTimeNano > 0 && endTimeNano > 0 {
		return startTimeNano <= timeNano && endTimeNano >= timeNano
	}

	return false
}

func TimeDurationSub(duration string) (t time.Time, err error) {
	var d time.Duration
	d, err = time.ParseDuration(duration)
	if err != nil {
		return t, err
	}

	t = time.Now().Add(-1 * d)

	return t, err
}

func LenPercentail(l int, n int) (pLen int) {
	pLen = (l * n / 100) - 1
	if pLen < 0 {
		pLen = 0
	}

	return pLen
}

func RequestTimeStddev(requests Percentails, sum, avg float64) (stddev float64) {
	n := float64(len(requests))
	for _, r := range requests {
		stddev += (r.RequestTime - avg) * (r.RequestTime - avg)
	}

	return math.Sqrt(stddev / n)
}

func Split(val, sep string) []string {
	strs := strings.Split(val, sep)
	trimedStrs := make([]string, 0, len(strs))

	for _, s := range strs {
		trimedStrs = append(trimedStrs, strings.Trim(s, " "))
	}

	return trimedStrs
}
