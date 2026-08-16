package stats

import (
	"io"

	"gopkg.in/yaml.v2"
)

func (hs *HTTPStats) LoadStats(r io.Reader) error {
	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	var stats []*HTTPStat
	if err := yaml.Unmarshal(buf, &stats); err != nil {
		return err
	}
	for _, stat := range stats {
		migrateLegacyBodyBytes(stat)
		restoreLegacySampleCounts(stat)
	}

	hs.stats = stats

	return nil
}

func restoreLegacySampleCounts(stat *HTTPStat) {
	restoreLegacySampleCount(stat.ResponseTime, stat.Cnt)
	restoreLegacySampleCount(stat.RequestBodyBytes, stat.Cnt)
	restoreLegacySampleCount(stat.ResponseBodyBytes, stat.Cnt)
}

func restoreLegacySampleCount(stats *floatStats, fallback int) {
	if stats == nil || stats.SampleCount != nil {
		return
	}

	stats.SampleCount = &fallback
}

func migrateLegacyBodyBytes(stat *HTTPStat) {
	if hasFloatStatsData(stat.ResponseBodyBytes) || !hasFloatStatsData(stat.RequestBodyBytes) {
		return
	}

	// Historical ALP dumps stored parsed response sizes in request_body_bytes.
	stat.RequestBodyBytes, stat.ResponseBodyBytes = stat.ResponseBodyBytes, stat.RequestBodyBytes
}

func hasFloatStatsData(stats *floatStats) bool {
	if stats == nil {
		return false
	}

	return stats.Max != 0 || stats.Min != 0 || stats.Sum != 0 || len(stats.Percentiles) != 0
}
