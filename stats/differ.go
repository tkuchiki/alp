package stats

import "fmt"

type Differ struct {
	From *HTTPStat
	To   *HTTPStat
}

func NewDiffer(from, to *HTTPStat) *Differ {
	return &Differ{
		From: from,
		To:   to,
	}
}

func (d *Differ) DiffCnt() string {
	v := d.To.Cnt - d.From.Cnt
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffStatus1xx() string {
	v := d.To.Status1xx - d.From.Status1xx
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffStatus2xx() string {
	v := d.To.Status2xx - d.From.Status2xx
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffStatus3xx() string {
	v := d.To.Status3xx - d.From.Status3xx
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffStatus4xx() string {
	v := d.To.Status4xx - d.From.Status4xx
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffStatus5xx() string {
	v := d.To.Status5xx - d.From.Status5xx
	if v >= 0 {
		return fmt.Sprintf("+%d", v)
	}

	return fmt.Sprintf("%d", v)
}

func (d *Differ) DiffMaxResponseTime() string {
	v := d.To.MaxResponseTime() - d.From.MaxResponseTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffMinResponseTime() string {
	v := d.To.MinResponseTime() - d.From.MinResponseTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffSumResponseTime() string {
	v := d.To.SumResponseTime() - d.From.SumResponseTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffAvgResponseTime() string {
	v := d.To.AvgResponseTime() - d.From.AvgResponseTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffPNResponseTime(n int) string {
	v := d.To.PNResponseTime(n) - d.From.PNResponseTime(n)
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffStddevResponseTime() string {
	v := d.To.StddevResponseTime() - d.From.StddevResponseTime()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

// request
func (d *Differ) DiffMaxRequestBodyBytes() string {
	v := d.To.MaxRequestBodyBytes() - d.From.MaxRequestBodyBytes()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffMinRequestBodyBytes() string {
	v := d.To.MinRequestBodyBytes() - d.From.MinRequestBodyBytes()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffSumRequestBodyBytes() string {
	v := d.To.SumRequestBodyBytes() - d.From.SumRequestBodyBytes()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffAvgRequestBodyBytes() string {
	v := d.To.AvgRequestBodyBytes() - d.From.AvgRequestBodyBytes()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffPNRequestBodyBytes(n int) string {
	v := d.To.PNRequestBodyBytes(n) - d.From.PNRequestBodyBytes(n)
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffStddevRequestBodyBytes() string {
	v := d.To.StddevRequestBodyBytes() - d.From.StddevRequestBodyBytes()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

// response
func (d *Differ) DiffMaxResponseBodyBytes() string {
	v := d.To.MaxResponseBodyBytes() - d.From.MaxResponseBodyBytes()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffMinResponseBodyBytes() string {
	v := d.To.MinResponseBodyBytes() - d.From.MinResponseBodyBytes()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffSumResponseBodyBytes() string {
	v := d.To.SumResponseBodyBytes() - d.From.SumResponseBodyBytes()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffAvgResponseBodyBytes() string {
	v := d.To.AvgResponseBodyBytes() - d.From.AvgResponseBodyBytes()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffPNResponseBodyBytes(n int) string {
	v := d.To.PNResponseBodyBytes(n) - d.From.PNResponseBodyBytes(n)
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func (d *Differ) DiffStddevResponseBodyBytes() string {
	v := d.To.StddevResponseBodyBytes() - d.From.StddevResponseBodyBytes()
	if v >= 0 {
		return fmt.Sprintf("+%.3f", v)
	}

	return fmt.Sprintf("%.3f", v)
}

func DiffCountAll(from, to map[string]int) map[string]string {
	counts := make(map[string]string, 6)
	keys := []string{"count", "1xx", "2xx", "3xx", "4xx", "5xx"}

	for _, key := range keys {
		v := to[key] - from[key]
		if v >= 0 {
			counts[key] = fmt.Sprintf("+%d", v)
		} else {
			counts[key] = fmt.Sprintf("-%d", v)
		}
	}

	return counts
}
