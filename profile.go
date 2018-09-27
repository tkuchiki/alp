package alp

import (
	"github.com/tkuchiki/gohttpstats"
	"github.com/tkuchiki/http-profile-helper"
)

type Profiler struct {
	stats  *httpstats.HTTPStats
	helper *helper.HTTPProfileHelper
}

func NewProfiler() *Profiler {
	po := httpstats.NewPrintOption()
	return &Profiler{
		stats: httpstats.NewHTTPStats(true, false, false, po),
		helper: helper.NewHTTPProfileHelper(),
	}
}


