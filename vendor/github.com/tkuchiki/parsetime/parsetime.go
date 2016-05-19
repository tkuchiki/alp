package parsetime

import (
	"errors"
	"fmt"
	"github.com/tkuchiki/go-timezone"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	errInvalidDateTime = errors.New("Invalid date/time")
	errInvalidOffset   = errors.New("Invalid offset")
	errInvalidArgs     = errors.New("Invalid arguments")
	errInvalidTimezone = errors.New("Invalid timezone")
	reISO8601          = regexp.MustCompile(ISO8601)
	reRFC8xx1123       = regexp.MustCompile(RFC8xx1123)
	reANSIC            = regexp.MustCompile(ANSIC)
	reUS               = regexp.MustCompile(US)
)

type sortedTime struct {
	time     time.Time
	priority int
}

type sortedTimes []sortedTime

func (st sortedTimes) Len() int           { return len(st) }
func (st sortedTimes) Swap(i, j int)      { st[i], st[j] = st[j], st[i] }
func (st sortedTimes) Less(i, j int) bool { return st[i].priority < st[j].priority }

// ParseTime parses the date/time string
type ParseTime struct {
	location *time.Location
}

// NewParseTime returns a new parser
func NewParseTime(location ...interface{}) (ParseTime, error) {
	var loc *time.Location
	var err error

	switch len(location) {
	case 0:
		zone, offset := time.Now().In(time.Local).Zone()
		loc = time.FixedZone(zone, offset)
	case 1:
		switch val := location[0].(type) {
		case *time.Location:
			loc = val
		case string:
			if val == "" {
				zone, offset := time.Now().In(time.Local).Zone()
				loc = time.FixedZone(zone, offset)
			} else {
				loc, err = time.LoadLocation(val)
				if err != nil {
					var offset int
					offset, err = timezone.GetOffset(val)
					if err != nil {
						return ParseTime{}, errInvalidTimezone
					}
					loc = time.FixedZone(val, offset)
				}
			}
		default:
			return ParseTime{}, fmt.Errorf("Invalid type: %T", val)
		}
	case 2:
		loc = time.FixedZone(location[0].(string), location[1].(int))
	default:
		return ParseTime{}, errInvalidArgs
	}

	return ParseTime{
		location: loc,
	}, err
}

// GetLocation returns *time.Location
func (pt *ParseTime) GetLocation() *time.Location {
	return pt.location
}

// SetLocation sets *time.Location
func (pt *ParseTime) SetLocation(loc *time.Location) {
	pt.location = loc
}

func fixedZone(t time.Time) *time.Location {
	zone, offset := t.Zone()
	return time.FixedZone(zone, offset)
}

func parseOffset(value string) (*time.Location, error) {
	var err error
	var t time.Time
	var loc *time.Location

	t, err = time.Parse("-07:00", value)
	if err == nil {
		return fixedZone(t), nil
	}

	t, err = time.Parse("-0700", value)
	if err == nil {
		return fixedZone(t), nil
	}

	_, err = time.Parse("MST", value)
	if err == nil {
		var offset int
		offset, err = timezone.GetOffset(value)
		if err != nil {
			return loc, err
		}

		return time.FixedZone(value, offset), nil
	}

	return loc, errInvalidOffset
}

func toLocation(offset string) (*time.Location, error) {
	var err error
	var loc *time.Location

	if strings.ToUpper(offset) == "Z" {
		loc = time.UTC
	} else {
		loc, err = parseOffset(offset)
	}

	return loc, err
}

func twoDigitTo4DigitYear(year string) (int, error) {
	val, err := strconv.Atoi(year)
	if err != nil {
		return 0, err
	}

	if val >= 70 && val <= 99 {
		return 1900 + val, err
	}

	return 2000 + val, err
}

func dateToInt(date string, dateType string, loc *time.Location) (int, error) {
	var err error
	var val int

	if date == "" {
		switch dateType {
		case "year":
			val = time.Now().In(loc).Year()
		case "month":
			val = int(time.Now().In(loc).Month())
		case "day":
			val = time.Now().In(loc).Day()
		case "hour":
			val = time.Now().In(loc).Hour()
		case "min":
			val = time.Now().In(loc).Minute()
		case "sec":
			if date == "" {
				val = 0
			} else {
				val = time.Now().In(loc).Second()
			}
		case "nsec":
			if date == "" {
				val = 0
			} else {
				val = time.Now().In(loc).Nanosecond()
			}
		default:
			err = errInvalidDateTime
		}
	} else {
		switch dateType {
		case "year":
			if stringLen(date) == 2 {
				return twoDigitTo4DigitYear(date)
			}
		case "month":
			if _, ok := Months[date]; ok {
				return Months[date], nil
			}
		}

		val, err = strconv.Atoi(date)
		return val, err
	}

	return val, err
}

func isOnlyDate(year, month, day, hour, min string) bool {
	return year != "" && month != "" && day != "" && hour == "" && min == ""
}

func stringLen(value string) int {
	return utf8.RuneCountInString(strings.Join(strings.Fields(value), ""))
}

func to24Hour(ampm string, value int) int {
	if strings.ToUpper(ampm) == "PM" {
		return 12 + value
	}

	return value
}

func parseISO8601(value string, loc *time.Location) (time.Time, int, error) {
	var t time.Time
	var priority int
	var err error

	group := reISO8601.FindStringSubmatch(value)

	if len(group) == 0 {
		return t, priority, errInvalidDateTime
	}

	priority = stringLen(value) - stringLen(group[0])

	var year, month, day, hour, min, sec, nsec int

	if group[8] != "" {
		loc, err = toLocation(group[8])
		if err != nil {
			return t, priority, err
		}
	}

	year, err = dateToInt(group[1], "year", loc)
	if err != nil {
		return t, priority, err
	}

	month, err = dateToInt(group[2], "month", loc)
	if err != nil {
		return t, priority, err
	}

	day, err = dateToInt(group[3], "day", loc)
	if err != nil {
		return t, priority, err
	}

	// 2006-01-02 -> 2006-01-02T00:00
	if isOnlyDate(group[1], group[2], group[3], group[4], group[5]) {
		group[4] = "0"
		group[5] = "0"
	}

	hour, err = dateToInt(group[4], "hour", loc)
	if err != nil {
		return t, priority, err
	}

	min, err = dateToInt(group[5], "min", loc)
	if err != nil {
		return t, priority, err
	}

	sec, err = dateToInt(group[6], "sec", loc)
	if err != nil {
		return t, priority, err
	}

	nsec, err = dateToInt(group[7], "nsec", loc)
	if err != nil {
		return t, priority, err
	}

	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, loc), priority, err
}

// ISO8601 parses ISO8601, RFC3339 date/time string
func (pt *ParseTime) ISO8601(value string) (time.Time, error) {
	t, _, err := parseISO8601(value, pt.location)
	return t, err
}

// RFC822, RFC850, RFC1123
func parseRFC8xx1123(value string, loc *time.Location) (time.Time, int, error) {
	var t time.Time
	var priority int
	var err error

	group := reRFC8xx1123.FindStringSubmatch(value)

	if len(group) == 0 {
		return t, priority, errInvalidDateTime
	}

	priority = stringLen(value) - stringLen(group[0])

	var year, month, day, hour, min, sec, nsec int

	if group[8] != "" {
		loc, err = toLocation(group[8])
		if err != nil {
			return t, priority, err
		}
	}

	day, err = dateToInt(group[1], "day", loc)
	if err != nil {
		return t, priority, err
	}

	month, err = dateToInt(group[2], "month", loc)
	if err != nil {
		return t, priority, err
	}

	year, err = dateToInt(group[3], "year", loc)
	if err != nil {
		return t, priority, err
	}

	// 02-Jan-06 -> 02-Jan-06 00:00
	if isOnlyDate(group[1], group[2], group[3], group[4], group[5]) {
		group[4] = "0"
		group[5] = "0"
	}

	hour, err = dateToInt(group[4], "hour", loc)
	if err != nil {
		return t, priority, err
	}

	min, err = dateToInt(group[5], "min", loc)
	if err != nil {
		return t, priority, err
	}

	sec, err = dateToInt(group[6], "sec", loc)
	if err != nil {
		return t, priority, err
	}

	nsec, err = dateToInt(group[7], "nsec", loc)
	if err != nil {
		return t, priority, err
	}

	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, loc), priority, err
}

// RFC8xx1123 parses RFC822, RFC850, RFC1123 date/time string
func (pt *ParseTime) RFC8xx1123(value string) (time.Time, error) {
	t, _, err := parseRFC8xx1123(value, pt.location)
	return t, err
}

func parseANSIC(value string, loc *time.Location) (time.Time, int, error) {
	var t time.Time
	var err error
	var priority int

	group := reANSIC.FindStringSubmatch(value)

	if len(group) == 0 {
		return t, priority, errInvalidDateTime
	}

	priority = stringLen(value) - stringLen(group[0])

	var year, month, day, hour, min, sec, nsec int

	if group[7] != "" {
		loc, err = toLocation(group[7])
		if err != nil {
			return t, priority, err
		}
	}

	month, err = dateToInt(group[1], "month", loc)
	if err != nil {
		return t, priority, err
	}

	day, err = dateToInt(group[2], "day", loc)
	if err != nil {
		return t, priority, err
	}

	hour, err = dateToInt(group[3], "hour", loc)
	if err != nil {
		return t, priority, err
	}

	min, err = dateToInt(group[4], "min", loc)
	if err != nil {
		return t, priority, err
	}

	sec, err = dateToInt(group[5], "sec", loc)
	if err != nil {
		return t, priority, err
	}

	nsec, err = dateToInt(group[6], "nsec", loc)
	if err != nil {
		return t, priority, err
	}

	year, err = dateToInt(group[8], "year", loc)
	if err != nil {
		return t, priority, err
	}

	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, loc), priority, err
}

// ANSIC parses ANSIC date/time string
func (pt *ParseTime) ANSIC(value string) (time.Time, error) {
	t, _, err := parseANSIC(value, pt.location)
	return t, err
}

func parseUS(value string, loc *time.Location) (time.Time, int, error) {
	var t time.Time
	var priority int
	var err error

	group := reUS.FindStringSubmatch(value)

	if len(group) == 0 {
		return t, priority, errInvalidDateTime
	}

	priority = stringLen(value) - stringLen(group[0])

	var year, month, day, hour, min, sec, nsec int

	if group[9] != "" {
		loc, err = toLocation(group[9])
		if err != nil {
			return t, priority, err
		}
	}

	month, err = dateToInt(group[1], "month", loc)
	if err != nil {
		return t, priority, err
	}

	day, err = dateToInt(group[2], "day", loc)
	if err != nil {
		return t, priority, err
	}

	year, err = dateToInt(group[3], "year", loc)
	if err != nil {
		return t, priority, err
	}

	if isOnlyDate(group[1], group[2], group[3], group[4], group[5]) {
		group[4] = "0"
		group[5] = "0"
	}

	hour, err = dateToInt(group[4], "hour", loc)
	if err != nil {
		return t, priority, err
	}

	min, err = dateToInt(group[5], "min", loc)
	if err != nil {
		return t, priority, err
	}

	sec, err = dateToInt(group[6], "sec", loc)
	if err != nil {
		return t, priority, err
	}

	nsec, err = dateToInt(group[7], "nsec", loc)
	if err != nil {
		return t, priority, err
	}

	ampm := group[8]

	if ampm != "" {
		hour = to24Hour(ampm, hour)
	}

	return time.Date(year, time.Month(month), day, hour, min, sec, nsec, loc), priority, err
}

// US parses MM/DD/YYYY format date/time string
func (pt *ParseTime) US(value string) (time.Time, error) {
	t, _, err := parseUS(value, pt.location)
	return t, err
}

// Parse parses date/time string
func (pt *ParseTime) Parse(value string) (time.Time, error) {
	times := make(sortedTimes, 0)
	t, priority, _ := parseISO8601(value, pt.location)
	if !t.IsZero() {
		times = append(times, sortedTime{time: t, priority: priority})
	}

	t, priority, _ = parseRFC8xx1123(value, pt.location)
	if !t.IsZero() {
		times = append(times, sortedTime{time: t, priority: priority})
	}

	t, priority, _ = parseANSIC(value, pt.location)
	if !t.IsZero() {
		times = append(times, sortedTime{time: t, priority: priority})
	}

	t, priority, _ = parseUS(value, pt.location)
	if !t.IsZero() {
		times = append(times, sortedTime{time: t, priority: priority})
	}

	if len(times) == 0 {
		var tmpT time.Time
		return tmpT, errInvalidDateTime
	}

	sort.Sort(times)

	return times[0].time, nil
}
