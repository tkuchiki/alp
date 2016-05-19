package main

import (
	"sort"
)

type Percentail struct {
	RequestTime float64
}

type Percentails []Percentail

type Profile struct {
	Uri         string       `yaml:"uri"`
	Cnt         int          `yaml:"cnt"`
	Max         float64      `yaml:"max"`
	Min         float64      `yaml:"min"`
	Sum         float64      `yaml:"sum"`
	Avg         float64      `yaml:"anv"`
	Method      string       `yaml:"method"`
	MaxBody     float64      `yaml:"max_body"`
	MinBody     float64      `yaml:"min_body"`
	SumBody     float64      `yaml:"sum_body"`
	AvgBody     float64      `yaml:"avg_body"`
	Percentails []Percentail `yaml:"percentails"`
	P1          float64      `yaml:"p1"`
	P50         float64      `yaml:"p50"`
	P99         float64      `yaml:"p99"`
	Stddev      float64      `yaml:"stddev"`
}

type Profiles []Profile
type ByMax struct{ Profiles }
type ByMin struct{ Profiles }
type BySum struct{ Profiles }
type ByAvg struct{ Profiles }
type ByUri struct{ Profiles }
type ByCnt struct{ Profiles }
type ByMethod struct{ Profiles }
type ByMaxBody struct{ Profiles }
type ByMinBody struct{ Profiles }
type BySumBody struct{ Profiles }
type ByAvgBody struct{ Profiles }
type ByP1 struct{ Profiles }
type ByP50 struct{ Profiles }
type ByP99 struct{ Profiles }
type ByStddev struct{ Profiles }

func (s Profiles) Len() int      { return len(s) }
func (s Profiles) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByMax) Less(i, j int) bool     { return s.Profiles[i].Max < s.Profiles[j].Max }
func (s ByMin) Less(i, j int) bool     { return s.Profiles[i].Min < s.Profiles[j].Min }
func (s BySum) Less(i, j int) bool     { return s.Profiles[i].Sum < s.Profiles[j].Sum }
func (s ByAvg) Less(i, j int) bool     { return s.Profiles[i].Avg < s.Profiles[j].Avg }
func (s ByUri) Less(i, j int) bool     { return s.Profiles[i].Uri < s.Profiles[j].Uri }
func (s ByCnt) Less(i, j int) bool     { return s.Profiles[i].Cnt < s.Profiles[j].Cnt }
func (s ByMethod) Less(i, j int) bool  { return s.Profiles[i].Method < s.Profiles[j].Method }
func (s ByMaxBody) Less(i, j int) bool { return s.Profiles[i].MaxBody < s.Profiles[j].MaxBody }
func (s ByMinBody) Less(i, j int) bool { return s.Profiles[i].MinBody < s.Profiles[j].MinBody }
func (s BySumBody) Less(i, j int) bool { return s.Profiles[i].SumBody < s.Profiles[j].SumBody }
func (s ByAvgBody) Less(i, j int) bool { return s.Profiles[i].AvgBody < s.Profiles[j].AvgBody }
func (s ByP1) Less(i, j int) bool      { return s.Profiles[i].P1 < s.Profiles[j].P1 }
func (s ByP50) Less(i, j int) bool     { return s.Profiles[i].P50 < s.Profiles[j].P50 }
func (s ByP99) Less(i, j int) bool     { return s.Profiles[i].P99 < s.Profiles[j].P99 }
func (s ByStddev) Less(i, j int) bool  { return s.Profiles[i].Stddev < s.Profiles[j].Stddev }

type ByRequestTime struct{ Percentails }

func (s Percentails) Len() int      { return len(s) }
func (s Percentails) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s ByRequestTime) Less(i, j int) bool {
	return s.Percentails[i].RequestTime < s.Percentails[j].RequestTime
}

func SortByMax(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByMax{ps}))
	} else {
		sort.Sort(ByMax{ps})
	}
	Output(ps, c)
}

func SortByMin(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByMin{ps}))
	} else {
		sort.Sort(ByMin{ps})
	}
	Output(ps, c)
}

func SortByAvg(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByAvg{ps}))
	} else {
		sort.Sort(ByAvg{ps})
	}
	Output(ps, c)
}

func SortBySum(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(BySum{ps}))
	} else {
		sort.Sort(BySum{ps})
	}
	Output(ps, c)
}

func SortByCnt(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByCnt{ps}))
	} else {
		sort.Sort(ByCnt{ps})
	}
	Output(ps, c)
}

func SortByUri(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByUri{ps}))
	} else {
		sort.Sort(ByUri{ps})
	}
	Output(ps, c)
}

func SortByMethod(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByMethod{ps}))
	} else {
		sort.Sort(ByMethod{ps})
	}
	Output(ps, c)
}

func SortByMaxBody(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByMaxBody{ps}))
	} else {
		sort.Sort(ByMaxBody{ps})
	}
	Output(ps, c)
}

func SortByMinBody(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByMinBody{ps}))
	} else {
		sort.Sort(ByMinBody{ps})
	}
	Output(ps, c)
}

func SortByAvgBody(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByAvgBody{ps}))
	} else {
		sort.Sort(ByAvgBody{ps})
	}
	Output(ps, c)
}

func SortBySumBody(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(BySumBody{ps}))
	} else {
		sort.Sort(BySumBody{ps})
	}
	Output(ps, c)
}

func SortByP1(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByP1{ps}))
	} else {
		sort.Sort(ByP1{ps})
	}
	Output(ps, c)
}

func SortByP50(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByP50{ps}))
	} else {
		sort.Sort(ByP50{ps})
	}
	Output(ps, c)
}

func SortByP99(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByP99{ps}))
	} else {
		sort.Sort(ByP99{ps})
	}
	Output(ps, c)
}

func SortByStddev(ps Profiles, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByStddev{ps}))
	} else {
		sort.Sort(ByStddev{ps})
	}
	Output(ps, c)
}

func SortProfiles(accessLog Profiles, c Config) {
	switch c.Sort {
	case "max":
		SortByMax(accessLog, c)
	case "min":
		SortByMin(accessLog, c)
	case "avg":
		SortByAvg(accessLog, c)
	case "sum":
		SortBySum(accessLog, c)
	case "cnt":
		SortByCnt(accessLog, c)
	case "uri":
		SortByUri(accessLog, c)
	case "method":
		SortByMethod(accessLog, c)
	case "max-body":
		SortByMaxBody(accessLog, c)
	case "min-body":
		SortByMinBody(accessLog, c)
	case "avg-body":
		SortByAvgBody(accessLog, c)
	case "sum-body":
		SortBySumBody(accessLog, c)
	case "p1":
		SortByP1(accessLog, c)
	case "p50":
		SortByP50(accessLog, c)
	case "p99":
		SortByP99(accessLog, c)
	case "stddev":
		SortByStddev(accessLog, c)
	}
}
