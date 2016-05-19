package parsetime

import (
	"strings"
)

const (
	year         = `(2[0-9]{3}|19[7-9][0-9])`
	month        = `(1[012]|0?[1-9])`
	day          = `([12][0-9]|3[01]|0?[1-9])`
	hour         = `(2[0-3]|[01]?[0-9])`
	min          = `([0-5]?[0-9])`
	sec          = min
	nsec         = `(?:[.])?([0-9]{1,9})?`
	weekday      = `(?:Mon|Monday|Tue|Tuesday|Wed|Wednesday|Thu|Thursday|Fri|Friday|Sat|Saturday|Sun|Sunday)`
	monthAbbr    = `(Jan|January|Feb|Februray|Mar|March|Apr|April|May|Jun|June|Jul|July|Aug|August|Sep|September|Oct|October|Nov|November|Dec|December|1[012]|0?[1-9])`
	offset       = `(Z|[+-][01][1-9]:[0-9]{2})?`
	zone         = `(?:[a-zA-Z0-9+-]{3,6})?`
	ymdSep       = `[ /.-]?`
	hmsSep       = `[ :.]?`
	t            = `(?:t|T|\s*)?`
	s            = `(?:\s*)?`
	ampm         = `([aApP][mM])`
	ampmHour     = `(1[01]|[0]?[0-9])`
	shortYear    = `(2[0-9]{3}|19[7-9][0-9]|[0-9]{2})`
	offsetZone   = `([+-][01][1-9]:[0-9]{2}|[a-zA-Z0-9+-]{3,6})?`
	usOffsetZone = `(?:[(])?([+-][01][1-9]:[0-9]{2}|[a-zA-Z0-9+-]{3,6})?(?:[)])?`
)

// Regular expressions
var (
	// ISO8601, RFC3339
	ISO8601 = strings.Join([]string{
		`(?:`, year, ymdSep, month, ymdSep, day, `)?`, t,
		`(?:`, hour, hmsSep, min, hmsSep, sec, `?`, nsec, `)?`,
		s, offset, s, zone,
	}, "")

	// RFC822, RFC850, RFC1123
	RFC8xx1123 = strings.Join([]string{
		`(?:`, weekday, `,?`, s, `)?`, day, ymdSep, monthAbbr, ymdSep, shortYear,
		hmsSep, `(?:`, hour, hmsSep, min, hmsSep, sec, `?`, nsec, `)?`,
		s, offsetZone,
	}, "")

	ANSIC = strings.Join([]string{
		`(?:`, weekday, s, `)?`, monthAbbr, ymdSep, day, ymdSep,
		`(?:`, hour, hmsSep, min, hmsSep, sec, `?`, nsec, `)?`,
		s, `(?:`, offsetZone, s, year, `)?`,
	}, "")

	US = strings.Join([]string{
		`(?:`, monthAbbr, ymdSep, day, `(?:,)?`, ymdSep, shortYear, `)?`, s, `(?:at)?`, s,
		`(?:`, hour, hmsSep, min, hmsSep, sec, `?`, nsec, `)?`,
		s, ampm, `?`, s, usOffsetZone,
	}, "")

	Months = map[string]int{
		"Jan":       1,
		"January":   1,
		"Feb":       2,
		"Februray":  2,
		"Mar":       3,
		"March":     3,
		"Apr":       4,
		"April":     4,
		"May":       5,
		"Jun":       6,
		"June":      6,
		"Jul":       7,
		"July":      7,
		"Aug":       8,
		"August":    8,
		"Sep":       9,
		"September": 9,
		"Oct":       10,
		"October":   10,
		"Nov":       11,
		"November":  11,
		"Dec":       12,
		"December":  12,
	}
)
