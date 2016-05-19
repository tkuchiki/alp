package parsetime

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var loc = createLocation("US/Arizona")

var iso8601Times = []TestTime{
	{
		Value: "2006-01-02 15:04",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:00", time.Local),
	},
	{
		Value: "2006-01-02 15:04-07:00",
		Time:  createTime("2006-01-02T15:04-07:00", "2006-01-02T15:04-07:00"),
	},
	{
		Value: "2006-01-02 15:04 -07:00",
		Time:  createTime("2006-01-02T15:04-07:00", "2006-01-02T15:04-07:00"),
	},
	{
		Value: "2006-01-02 15:04:05",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", time.Local),
	},
	{
		Value: "2006-01-02 15:04:05-07:00",
		Time:  createTime("2006-01-02T15:04:05-07:00", "2006-01-02T15:04:05-07:00"),
	},
	{
		Value: "2006-01-02 15:04:05 -07:00",
		Time:  createTime("2006-01-02T15:04:05-07:00", "2006-01-02T15:04:05-07:00"),
	},
	{
		Value: "2006-01-02 15:04:05-07:00 MST",
		Time:  createTime("2006-01-02T15:04:05-07:00 MST", "2006-01-02T15:04:05-07:00 MST"),
	},
	{
		Value: "2006-01-02 15:04:05 -07:00 MST",
		Time:  createTime("2006-01-02T15:04:05-07:00 MST", "2006-01-02T15:04:05-07:00 MST"),
	},
	{
		Value: "2006-01-02 15:04:05.999999999",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999999", "2006-01-02T15:04:05.999999999", time.Local),
	},
	{
		Value: "2006-01-02 15:04:05.999999-07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.999999-07:00 MST", "2006-01-02T15:04:05.999999-07:00 MST"),
	},
	{
		Value: "2006-01-02 15:04:05.9-07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.9-07:00 MST", "2006-01-02T15:04:05.9-07:00 MST"),
	},
	{
		Value: "2006-01-02 15:04:05.9 -07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.9-07:00 MST", "2006-01-02T15:04:05.9-07:00 MST"),
	},
	{
		Value: "2006-01-02 15:04:05.999-07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.999-07:00 MST", "2006-01-02T15:04:05.999-07:00 MST"),
	},
	{
		Value: "2006-01-02 15:04:05.999 -07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.999-07:00 MST", "2006-01-02T15:04:05.999-07:00 MST"),
	},
	{
		Value: "2006-01-02 15:04:05.999999-07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.999999-07:00 MST", "2006-01-02T15:04:05.999999-07:00 MST"),
	},
	{
		Value: "2006-01-02 15:04:05.999999 -07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.999999-07:00 MST", "2006-01-02T15:04:05.999999-07:00 MST"),
	},
	{
		Value: "2006-01-02 15:04:05.999999999-07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.999999999-07:00 MST", "2006-01-02T15:04:05.999999999-07:00 MST"),
	},
	{
		Value: "2006-01-02 15:04:05.999999999 -07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.999999999-07:00 MST", "2006-01-02T15:04:05.999999999-07:00 MST"),
	},
	{
		Value: "2006-01-02T15:04",
		Time:  createTimeInLocation("2006-01-02T15:04:00", "2006-01-02T15:04:00", time.Local),
	},
	{
		Value: "2006-01-02T15:04-07:00",
		Time:  createTime("2006-01-02T15:04-07:00", "2006-01-02T15:04-07:00"),
	},
	{
		Value: "2006-01-02T15:04 -07:00",
		Time:  createTime("2006-01-02T15:04-07:00", "2006-01-02T15:04-07:00"),
	},
	{
		Value: "2006-01-02T15:04:05",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", time.Local),
	},
	{
		Value: "2006-01-02T15:04:05-07:00",
		Time:  createTime("2006-01-02T15:04:05-07:00", "2006-01-02T15:04:05-07:00"),
	},
	{
		Value: "2006-01-02T15:04:05 -07:00",
		Time:  createTime("2006-01-02T15:04:05-07:00", "2006-01-02T15:04:05-07:00"),
	},
	{
		Value: "2006-01-02T15:04:05-07:00 MST",
		Time:  createTime("2006-01-02T15:04:05-07:00 MST", "2006-01-02T15:04:05-07:00 MST"),
	},
	{
		Value: "2006-01-02T15:04:05 -07:00 MST",
		Time:  createTime("2006-01-02T15:04:05-07:00 MST", "2006-01-02T15:04:05-07:00 MST"),
	},
	{
		Value: "2006-01-02T15:04:05.999999999",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999999", "2006-01-02T15:04:05.999999999", time.Local),
	},
	{
		Value: "2006-01-02T15:04:05.999999999-07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.999999999-07:00 MST", "2006-01-02T15:04:05.999999999-07:00 MST"),
	},
	{
		Value: "2006-01-02T15:04:05.999999999 -07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.999999999-07:00 MST", "2006-01-02T15:04:05.999999999-07:00 MST"),
	},
	{
		Value: "2006-01-02T15:04:05.999999-07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.999999-07:00 MST", "2006-01-02T15:04:05.999999-07:00 MST"),
	},
	{
		Value: "2006-01-02T15:04:05.999999 -07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.999999-07:00 MST", "2006-01-02T15:04:05.999999-07:00 MST"),
	},
	{
		Value: "2006-01-02T15:04:05.9-07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.9-07:00 MST", "2006-01-02T15:04:05.9-07:00 MST"),
	},
	{
		Value: "2006-01-02T15:04:05.9 -07:00 MST",
		Time:  createTime("2006-01-02T15:04:05.9-07:00 MST", "2006-01-02T15:04:05.9-07:00 MST"),
	},
	{
		Value: "2006-01-02",
		Time:  createTimeInLocation("2006-01-02", "2006-01-02", time.Local),
	},
	{
		Value: "20060102",
		Time:  createTimeInLocation("20060102", "20060102", time.Local),
	},
	{
		Value: "20060102150405",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", time.Local),
	},
	{
		Value: "20060102 150405",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", time.Local),
	},
	{
		Value: "20060102T150405",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", time.Local),
	},
	{
		Value: "15:04:05",
		Time:  createCurrentDateInLocation("15:04:05", "15:04:05", time.Local),
	},
	{
		Value: "15:04:05-07:00 MST",
		Time:  createCurrentDateInLocation("15:04:05-07:00 MST", "15:04:05-07:00 MST", loc),
	},
	{
		Value: "15:04:05 -07:00 MST",
		Time:  createCurrentDateInLocation("15:04:05-07:00 MST", "15:04:05-07:00 MST", loc),
	},
	{
		Value: "15:04:05.9-07:00 MST",
		Time:  createCurrentDateInLocation("15:04:05.9-07:00 MST", "15:04:05.9-07:00 MST", loc),
	},
	{
		Value: "15:04:05.9 -07:00 MST",
		Time:  createCurrentDateInLocation("15:04:05.9-07:00 MST", "15:04:05.9-07:00 MST", loc),
	},
	{
		Value: "15:04:05.999-07:00 MST",
		Time:  createCurrentDateInLocation("15:04:05.999-07:00 MST", "15:04:05.999-07:00 MST", loc),
	},
	{
		Value: "15:04:05.999 -07:00 MST",
		Time:  createCurrentDateInLocation("15:04:05.999-07:00 MST", "15:04:05.999-07:00 MST", loc),
	},
	{
		Value: "15:04:05.999999-07:00 MST",
		Time:  createCurrentDateInLocation("15:04:05.999999-07:00 MST", "15:04:05.999999-07:00 MST", loc),
	},
	{
		Value: "15:04:05.999999 -07:00 MST",
		Time:  createCurrentDateInLocation("15:04:05.999999-07:00 MST", "15:04:05.999999-07:00 MST", loc),
	},
	{
		Value: "15:04:05.999999999-07:00 MST",
		Time:  createCurrentDateInLocation("15:04:05.999999999-07:00 MST", "15:04:05.999999999-07:00 MST", loc),
	},
	{
		Value: "15:04:05.999999999 -07:00 MST",
		Time:  createCurrentDateInLocation("15:04:05.999999999-07:00 MST", "15:04:05.999999999-07:00 MST", loc),
	},
	{
		Value: "150405-07:00 MST",
		Time:  createCurrentDateInLocation("150405-07:00 MST", "150405-07:00 MST", loc),
	},
	{
		Value: "150405 -07:00 MST",
		Time:  createCurrentDateInLocation("150405-07:00 MST", "150405-07:00 MST", loc),
	},
	{
		Value: "150405.9-07:00 MST",
		Time:  createCurrentDateInLocation("150405.9-07:00 MST", "150405.9-07:00 MST", loc),
	},
	{
		Value: "150405.9 -07:00 MST",
		Time:  createCurrentDateInLocation("150405.9-07:00 MST", "150405.9-07:00 MST", loc),
	},
	{
		Value: "150405.999-07:00 MST",
		Time:  createCurrentDateInLocation("150405.999-07:00 MST", "150405.999-07:00 MST", loc),
	},
	{
		Value: "150405.999 -07:00 MST",
		Time:  createCurrentDateInLocation("150405.999-07:00 MST", "150405.999-07:00 MST", loc),
	},
	{
		Value: "150405.999999-07:00 MST",
		Time:  createCurrentDateInLocation("150405.999999-07:00 MST", "150405.999999-07:00 MST", loc),
	},
	{
		Value: "150405.999999 -07:00 MST",
		Time:  createCurrentDateInLocation("150405.999999-07:00 MST", "150405.999999-07:00 MST", loc),
	},
	{
		Value: "150405.999999999-07:00 MST",
		Time:  createCurrentDateInLocation("150405.999999999-07:00 MST", "150405.999999999-07:00 MST", loc),
	},
	{
		Value: "150405.999999999 -07:00 MST",
		Time:  createCurrentDateInLocation("150405.999999999-07:00 MST", "150405.999999999-07:00 MST", loc),
	},
	{
		Value: "2006-01-02 15:04:05Z",
		Time:  createTime("2006-01-02T15:04:05Z0700", "2006-01-02T15:04:05Z"),
	},
	{
		Value: "2006-01-02T15:04:05Z",
		Time:  createTime("2006-01-02T15:04:05Z0700", "2006-01-02T15:04:05Z"),
	},
	{
		Value: "2006-01-02 15:04:05.9Z",
		Time:  createTime("2006-01-02T15:04:05.9Z0700", "2006-01-02T15:04:05.9Z"),
	},
	{
		Value: "2006-01-02T15:04:05.9Z",
		Time:  createTime("2006-01-02T15:04:05.9Z0700", "2006-01-02T15:04:05.9Z"),
	},
	{
		Value: "2006-01-02 15:04:05.999Z",
		Time:  createTime("2006-01-02T15:04:05.999Z0700", "2006-01-02T15:04:05.999Z"),
	},
	{
		Value: "2006-01-02T15:04:05.999Z",
		Time:  createTime("2006-01-02T15:04:05.999Z0700", "2006-01-02T15:04:05.999Z"),
	},
	{
		Value: "2006-01-02 15:04:05.999999Z",
		Time:  createTime("2006-01-02T15:04:05.999999Z0700", "2006-01-02T15:04:05.999999Z"),
	},
	{
		Value: "2006-01-02T15:04:05.999999Z",
		Time:  createTime("2006-01-02T15:04:05.999999Z0700", "2006-01-02T15:04:05.999999Z"),
	},
	{
		Value: "2006-01-02 15:04:05.999999999Z",
		Time:  createTime("2006-01-02T15:04:05.999999999Z0700", "2006-01-02T15:04:05.999999999Z"),
	},
	{
		Value: "2006-01-02T15:04:05.999999999Z",
		Time:  createTime("2006-01-02T15:04:05.999999999Z0700", "2006-01-02T15:04:05.999999999Z"),
	},
}

var rfc8xx1123Times = []TestTime{
	{
		Value: "02-Jan-06 1504 MST",
		Time:  createTimeInLocation("02-Jan-06 15:04:05 MST", "02-Jan-06 15:04:00 MST", loc),
	},
	{
		Value: "02-Jan-06 15:04 MST",
		Time:  createTimeInLocation("02-Jan-06 15:04:05 MST", "02-Jan-06 15:04:00 MST", loc),
	},
	{
		Value: "02-Jan-06 150405 MST",
		Time:  createTimeInLocation("02-Jan-06 15:04:05 MST", "02-Jan-06 15:04:05 MST", loc),
	},
	{
		Value: "02-Jan-06 15:04:05 MST",
		Time:  createTimeInLocation("02-Jan-06 15:04:05 MST", "02-Jan-06 15:04:05 MST", loc),
	},
	{
		Value: "02-Jan-06 1504-0700",
		Time:  createTimeInLocation("02-Jan-06 15:04:05 MST", "02-Jan-06 15:04:00 MST", loc),
	},
	{
		Value: "02-Jan-06 15:04-0700",
		Time:  createTimeInLocation("02-Jan-06 15:04:05 MST", "02-Jan-06 15:04:00 MST", loc),
	},
	{
		Value: "02-Jan-06 150405-0700",
		Time:  createTimeInLocation("02-Jan-06 15:04:05 MST", "02-Jan-06 15:04:05 MST", loc),
	},
	{
		Value: "02-Jan-06 15:04:05-0700",
		Time:  createTimeInLocation("02-Jan-06 15:04:05 MST", "02-Jan-06 15:04:05 MST", loc),
	},
	{
		Value: "02-Jan-06 15:04 -0700",
		Time:  createTimeInLocation("02-Jan-06 15:04:05 MST", "02-Jan-06 15:04:00 MST", loc),
	},
	{
		Value: "02-Jan-06 15:04:05 -0700",
		Time:  createTimeInLocation("02-Jan-06 15:04:05 MST", "02-Jan-06 15:04:05 MST", loc),
	},
	{
		Value: "Monday, 02-Jan-06 15:04 MST",
		Time:  createTimeInLocation("Monday, 02-Jan-06 15:04:05 MST", "Monday, 02-Jan-06 15:04:00 MST", loc),
	},
	{
		Value: "Monday, 02-Jan-06 15:04:05 MST",
		Time:  createTimeInLocation("Monday, 02-Jan-06 15:04:05 MST", "Monday, 02-Jan-06 15:04:05 MST", loc),
	},
	{
		Value: "Mon, 02-Jan-06 15:04 MST",
		Time:  createTimeInLocation("Mon, 02-Jan-06 15:04:05 MST", "Mon, 02-Jan-06 15:04:00 MST", loc),
	},
	{
		Value: "Mon, 02-Jan-06 15:04:05 MST",
		Time:  createTimeInLocation("Mon, 02-Jan-06 15:04:05 MST", "Mon, 02-Jan-06 15:04:05 MST", loc),
	},
	{
		Value: "Mon, 02-Jan-06 15:04-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-06 15:04:05 -07:00", "Mon, 02-Jan-06 15:04:00 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-06 15:04:05-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-06 15:04:05 -07:00", "Mon, 02-Jan-06 15:04:05 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-06 15:04 -07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-06 15:04:05 -07:00", "Mon, 02-Jan-06 15:04:00 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-06 15:04:05 -07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-06 15:04:05 -07:00", "Mon, 02-Jan-06 15:04:05 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-2006 15:04-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-2006 15:04:00 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-2006 15:04:05-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-2006 15:04:05 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-2006 15:04 -07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-2006 15:04:00 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-2006 15:04:05 -07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-2006 15:04:05 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-70 15:04-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-1970 15:04:00 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-70 15:04:05-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-1970 15:04:05 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-70 15:04 -07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-1970 15:04:00 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-70 15:04:05 -07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-1970 15:04:05 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-99 15:04-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-1999 15:04:00 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-99 15:04:05-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-1999 15:04:05 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-99 15:04:05 -07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-1999 15:04:05 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-00 15:04-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-2000 15:04:00 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-00 15:04:05-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-2000 15:04:05 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-00 15:04:05 -07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05 -07:00", "Mon, 02-Jan-2000 15:04:05 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-00 15:04:05.9-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05.9 -07:00", "Mon, 02-Jan-2000 15:04:05.9 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-00 15:04:05.9 -07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05.9 -07:00", "Mon, 02-Jan-2000 15:04:05.9 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-00 15:04:05.999-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05.999 -07:00", "Mon, 02-Jan-2000 15:04:05.999 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-00 15:04:05.999 -07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05.999 -07:00", "Mon, 02-Jan-2000 15:04:05.999 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-00 15:04:05.999999-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05.999999 -07:00", "Mon, 02-Jan-2000 15:04:05.999999 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-00 15:04:05.999999 -07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05.999999 -07:00", "Mon, 02-Jan-2000 15:04:05.999999 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-00 15:04:05.999999999-07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05.999999999 -07:00", "Mon, 02-Jan-2000 15:04:05.999999999 -07:00", loc),
	},
	{
		Value: "Mon, 02-Jan-00 15:04:05.999999999 -07:00",
		Time:  createTimeInLocation("Mon, 02-Jan-2006 15:04:05.999999999 -07:00", "Mon, 02-Jan-2000 15:04:05.999999999 -07:00", loc),
	},
}

var ansicTimes = []TestTime{
	{
		Value: "Mon Jan 02 150405 2006",
		Time:  createTimeInLocation("Mon Jan 02 15:04:05 2006", "Mon Jan 02 15:04:05 2006", time.Local),
	},
	{
		Value: "Mon Jan 02 15:04:05 2006",
		Time:  createTimeInLocation("Mon Jan 02 15:04:05 2006", "Mon Jan 02 15:04:05 2006", time.Local),
	},
	{
		Value: "Mon Jan 02 150405 MST 2006",
		Time:  createTimeInLocation("Mon Jan 02 15:04:05 MST 2006", "Mon Jan 02 15:04:05 MST 2006", loc),
	},
	{
		Value: "Mon Jan 02 15:04:05 MST 2006",
		Time:  createTimeInLocation("Mon Jan 02 15:04:05 MST 2006", "Mon Jan 02 15:04:05 MST 2006", loc),
	},
	{
		Value: "Mon Jan 02 1504-07:00 2006",
		Time:  createTimeInLocation("Mon Jan 02 15:04:05 -07:00 2006", "Mon Jan 02 15:04:00 -07:00 2006", loc),
	},
	{
		Value: "Mon Jan 02 15:04-07:00 2006",
		Time:  createTimeInLocation("Mon Jan 02 15:04:05 -07:00 2006", "Mon Jan 02 15:04:00 -07:00 2006", loc),
	},
	{
		Value: "Mon Jan 02 1504 -07:00 2006",
		Time:  createTimeInLocation("Mon Jan 02 15:04:05 -07:00 2006", "Mon Jan 02 15:04:00 -07:00 2006", loc),
	},
	{
		Value: "Mon Jan 02 15:04 -07:00 2006",
		Time:  createTimeInLocation("Mon Jan 02 15:04:05 -07:00 2006", "Mon Jan 02 15:04:00 -07:00 2006", loc),
	},
	{
		Value: "Mon Jan 02 150405-07:00 2006",
		Time:  createTimeInLocation("Mon Jan 02 15:04:05 -07:00 2006", "Mon Jan 02 15:04:05 -07:00 2006", loc),
	},
	{
		Value: "Mon Jan 02 15:04:05-07:00 2006",
		Time:  createTimeInLocation("Mon Jan 02 15:04:05 -07:00 2006", "Mon Jan 02 15:04:05 -07:00 2006", loc),
	},
	{
		Value: "Mon Jan 02 150405 -07:00 2006",
		Time:  createTimeInLocation("Mon Jan 02 15:04:05 -07:00 2006", "Mon Jan 02 15:04:05 -07:00 2006", loc),
	},
	{
		Value: "Mon Jan 02 15:04:05 -07:00 2006",
		Time:  createTimeInLocation("Mon Jan 02 15:04:05 -07:00 2006", "Mon Jan 02 15:04:05 -07:00 2006", loc),
	},
	{
		Value: "Jan 02 150405",
		Time:  createCurrentYearInLocation("Jan 02 15:04:05", "Jan 02 15:04:05", time.Local),
	},
	{
		Value: "Jan 02 15:04:05",
		Time:  createCurrentYearInLocation("Jan 02 15:04:05", "Jan 02 15:04:05", time.Local),
	},
	{
		Value: "Jan 02 150405.9",
		Time:  createCurrentYearInLocation("Jan 02 15:04:05.9", "Jan 02 15:04:05.9", time.Local),
	},
	{
		Value: "Jan 02 15:04:05.9",
		Time:  createCurrentYearInLocation("Jan 02 15:04:05.9", "Jan 02 15:04:05.9", time.Local),
	},
	{
		Value: "Jan 02 150405.999",
		Time:  createCurrentYearInLocation("Jan 02 15:04:05.999", "Jan 02 15:04:05.999", time.Local),
	},
	{
		Value: "Jan 02 15:04:05.999",
		Time:  createCurrentYearInLocation("Jan 02 15:04:05.999", "Jan 02 15:04:05.999", time.Local),
	},
	{
		Value: "Jan 02 150405.999999",
		Time:  createCurrentYearInLocation("Jan 02 15:04:05.999999", "Jan 02 15:04:05.999999", time.Local),
	},
	{
		Value: "Jan 02 15:04:05.999999",
		Time:  createCurrentYearInLocation("Jan 02 15:04:05.999999", "Jan 02 15:04:05.999999", time.Local),
	},
	{
		Value: "Jan 02 150405.999999999",
		Time:  createCurrentYearInLocation("Jan 02 15:04:05.999999999", "Jan 02 15:04:05.999999999", time.Local),
	},
	{
		Value: "Jan 02 15:04:05.999999999",
		Time:  createCurrentYearInLocation("Jan 02 15:04:05.999999999", "Jan 02 15:04:05.999999999", time.Local),
	},
}

var usTimes = []TestTime{
	{
		Value: "11:04AM",
		Time:  createCurrentDateInLocation("15:04:05", "11:04:00", time.Local),
	},
	{
		Value: "11:04PM",
		Time:  createCurrentDateInLocation("15:04:05", "23:04:00", time.Local),
	},
	{
		Value: "11:04 AM",
		Time:  createCurrentDateInLocation("15:04:05", "11:04:00", time.Local),
	},
	{
		Value: "11:04 PM",
		Time:  createCurrentDateInLocation("15:04:05", "23:04:00", time.Local),
	},
	{
		Value: "11:04:05 AM",
		Time:  createCurrentDateInLocation("15:04:05", "11:04:05", time.Local),
	},
	{
		Value: "11:04:05 PM",
		Time:  createCurrentDateInLocation("15:04:05", "23:04:05", time.Local),
	},
	{
		Value: "11:04:05.9AM",
		Time:  createCurrentDateInLocation("15:04:05.9", "11:04:05.9", time.Local),
	},
	{
		Value: "11:04:05.9 AM",
		Time:  createCurrentDateInLocation("15:04:05.9", "11:04:05.9", time.Local),
	},
	{
		Value: "11:04:05.9PM",
		Time:  createCurrentDateInLocation("15:04:05.9", "23:04:05.9", time.Local),
	},
	{
		Value: "11:04:05.9 PM",
		Time:  createCurrentDateInLocation("15:04:05.9", "23:04:05.9", time.Local),
	},
	{
		Value: "11:04:05.999AM",
		Time:  createCurrentDateInLocation("15:04:05.999", "11:04:05.999", time.Local),
	},
	{
		Value: "11:04:05.999 AM",
		Time:  createCurrentDateInLocation("15:04:05.999", "11:04:05.999", time.Local),
	},
	{
		Value: "11:04:05.999PM",
		Time:  createCurrentDateInLocation("15:04:05.999", "23:04:05.999", time.Local),
	},
	{
		Value: "11:04:05.999 PM",
		Time:  createCurrentDateInLocation("15:04:05.999", "23:04:05.999", time.Local),
	},
	{
		Value: "11:04:05.999999AM",
		Time:  createCurrentDateInLocation("15:04:05.999999", "11:04:05.999999", time.Local),
	},
	{
		Value: "11:04:05.999999 AM",
		Time:  createCurrentDateInLocation("15:04:05.999999", "11:04:05.999999", time.Local),
	},
	{
		Value: "11:04:05.999999PM",
		Time:  createCurrentDateInLocation("15:04:05.999999", "23:04:05.999999", time.Local),
	},
	{
		Value: "11:04:05.999999 PM",
		Time:  createCurrentDateInLocation("15:04:05.999999", "23:04:05.999999", time.Local),
	},
	{
		Value: "11:04:05.999999999AM",
		Time:  createCurrentDateInLocation("15:04:05.999999999", "11:04:05.999999999", time.Local),
	},
	{
		Value: "11:04:05.999999999 AM",
		Time:  createCurrentDateInLocation("15:04:05.999999999", "11:04:05.999999999", time.Local),
	},
	{
		Value: "11:04:05.999999999PM",
		Time:  createCurrentDateInLocation("15:04:05.999999999", "23:04:05.999999999", time.Local),
	},
	{
		Value: "11:04:05.999999999 PM",
		Time:  createCurrentDateInLocation("15:04:05.999999999", "23:04:05.999999999", time.Local),
	},
	{
		Value: "01-02-06 3:04AM",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:00", time.Local),
	},
	{
		Value: "01-02-06 3:04 AM",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:00", time.Local),
	},
	{
		Value: "01-02-06 3:04PM",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:00", time.Local),
	},
	{
		Value: "01-02-06 3:04 PM",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:00", time.Local),
	},
	{
		Value: "01-02-06 03:04:05AM",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:05", time.Local),
	},
	{
		Value: "01-02-06 03:04:05 AM",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:05", time.Local),
	},
	{
		Value: "01-02-06 03:04:05PM",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", time.Local),
	},
	{
		Value: "01-02-06 03:04:05 PM",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.9AM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.9", "2006-01-02T03:04:05.9", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.9 AM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.9", "2006-01-02T03:04:05.9", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.9PM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.9", "2006-01-02T15:04:05.9", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.9 PM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.9", "2006-01-02T15:04:05.9", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.999AM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999", "2006-01-02T03:04:05.999", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.999 AM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999", "2006-01-02T03:04:05.999", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.999PM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999", "2006-01-02T15:04:05.999", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.999 PM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999", "2006-01-02T15:04:05.999", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.999999AM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999", "2006-01-02T03:04:05.999999", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.999999 AM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999", "2006-01-02T03:04:05.999999", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.999999PM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999", "2006-01-02T15:04:05.999999", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.999999 PM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999", "2006-01-02T15:04:05.999999", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.999999999AM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999999", "2006-01-02T03:04:05.999999999", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.999999999 AM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999999", "2006-01-02T03:04:05.999999999", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.999999999PM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999999", "2006-01-02T15:04:05.999999999", time.Local),
	},
	{
		Value: "01-02-06 03:04:05.999999999 PM",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999999", "2006-01-02T15:04:05.999999999", time.Local),
	},
	{
		Value: "Jan 2, 2006",
		Time:  createTimeInLocation("2006-01-02", "2006-01-02", time.Local),
	},
	{
		Value: "Jan 2, 2006 at 3:04am (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:00", loc),
	},
	{
		Value: "Jan 2, 2006 at 03:04am (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:00", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04pm (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:00", loc),
	},
	{
		Value: "Jan 2, 2006 at 03:04pm (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:00", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04 am (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:00", loc),
	},
	{
		Value: "Jan 2, 2006 at 03:04 am (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:00", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04 pm (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:00", loc),
	},
	{
		Value: "Jan 2, 2006 at 03:04 pm (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:00", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05am (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:05", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05pm (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05 am (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:05", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05 pm (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05.9am (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05.9", "2006-01-02T03:04:05.9", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05.9pm (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05.9", "2006-01-02T15:04:05.9", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05.999am (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999", "2006-01-02T03:04:05.999", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05.999pm (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999", "2006-01-02T15:04:05.999", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05.999999am (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999", "2006-01-02T03:04:05.999999", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05.999999pm (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999", "2006-01-02T15:04:05.999999", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05.999999999am (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999999", "2006-01-02T03:04:05.999999999", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05.999999999pm (MST)",
		Time:  createTimeInLocation("2006-01-02T15:04:05.999999999", "2006-01-02T15:04:05.999999999", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04am MST",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:00", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04pm MST",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:00", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04am -07:00",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:00", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04pm -07:00",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:00", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05am MST",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:05", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05pm MST",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05am -07:00",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T03:04:05", loc),
	},
	{
		Value: "Jan 2, 2006 at 3:04:05pm -07:00",
		Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", loc),
	},
}

type TestTime struct {
	Value string
	Time  time.Time
}

func createLocation(name string) *time.Location {
	loc, _ := time.LoadLocation(name)
	return loc
}

func getOffset(t time.Time) int {
	_, offset := t.Zone()
	return offset
}

func localTime() time.Time {
	t, _ := time.ParseInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", time.Local)
	return t
}

func createTime(layout, value string) time.Time {
	t, _ := time.Parse(layout, value)
	return t
}

func createTimeInLocation(layout, value string, loc *time.Location) time.Time {
	t, _ := time.ParseInLocation(layout, value, loc)
	return t
}

func createCurrentDateInLocation(layout, value string, loc *time.Location) time.Time {
	tmpT, _ := time.ParseInLocation(layout, value, loc)

	now := time.Now().In(loc)
	year := now.Year()
	month := time.Month(now.Month())
	day := now.Day()

	return time.Date(year, month, day, tmpT.Hour(), tmpT.Minute(), tmpT.Second(), tmpT.Nanosecond(), loc)
}

func createCurrentYearInLocation(layout, value string, loc *time.Location) time.Time {
	tmpT, _ := time.ParseInLocation(layout, value, loc)

	now := time.Now().In(loc)
	year := now.Year()

	return time.Date(year, time.Month(tmpT.Month()), tmpT.Day(), tmpT.Hour(), tmpT.Minute(), tmpT.Second(), tmpT.Nanosecond(), loc)
}

func testTimes(times []TestTime, parseType string, test *testing.T) {
	var t time.Time
	var err error
	assert := assert.New(test)

	for _, tt := range times {
		p, _ := NewParseTime()

		switch parseType {
		case "ISO8601":
			t, err = p.ISO8601(tt.Value)
		case "RFC8xx1123":
			t, err = p.RFC8xx1123(tt.Value)
		case "ANSIC":
			t, err = p.ANSIC(tt.Value)
		case "US":
			t, err = p.US(tt.Value)
		case "Parse":
			t, err = p.Parse(tt.Value)
		}

		assert.Equal(nil, err, "Invalid date/time")

		t2 := tt.Time

		assert.Equal(getOffset(t), getOffset(t2), "Incorrect offset")
		assert.Equal(t.Unix(), t2.Unix(), "Parse error")
	}
}

func TestNewParseTime(test *testing.T) {
	assert := assert.New(test)

	p, _ := NewParseTime()
	loc := p.GetLocation()
	zone, offset := time.Now().In(time.Local).Zone()
	loc2 := time.FixedZone(zone, offset)

	assert.Equal(loc.String(), loc2.String(), "Incorrect location")
}

func TestNewParseTimeEmptyString(test *testing.T) {
	assert := assert.New(test)

	p, _ := NewParseTime("")

	loc := p.GetLocation()
	zone, offset := time.Now().In(time.Local).Zone()
	loc2 := time.FixedZone(zone, offset)

	assert.Equal(loc.String(), loc2.String(), "Incorrect location")
}

func TestNewParseTimeLocation(test *testing.T) {
	assert := assert.New(test)

	loc, _ := time.LoadLocation("Etc/GMT+12")
	p, _ := NewParseTime(loc)
	loc2 := p.GetLocation()

	assert.Equal(loc.String(), loc2.String(), "Incorrect location")
}

func TestNewParseTimeTimezone(test *testing.T) {
	assert := assert.New(test)

	p, _ := NewParseTime("GMT+12")
	loc2, _ := time.LoadLocation("Etc/GMT+12")
	t, _ := p.Parse("2006-01-02T15:04:05")
	t2, _ := time.ParseInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", loc2)

	assert.Equal(t.Unix(), t2.Unix(), "Incorrect location")
}

func TestNewParseTimeLoadLocation(test *testing.T) {
	assert := assert.New(test)

	p, _ := NewParseTime("Etc/GMT+12")
	loc := p.GetLocation()
	loc2, _ := time.LoadLocation("Etc/GMT+12")

	assert.Equal(loc.String(), loc2.String(), "Incorrect location")
}

func TestNewParseTimeFixedZone(test *testing.T) {
	assert := assert.New(test)

	p, _ := NewParseTime("Etc/GMT+12", -43200)
	loc := p.GetLocation()
	loc2, _ := time.LoadLocation("Etc/GMT+12")

	assert.Equal(loc.String(), loc2.String(), "Incorrect location")
}

func TestISO8601(test *testing.T) {
	testTimes(iso8601Times, "ISO8601", test)
}

func TestRFC8xx1123(test *testing.T) {
	testTimes(rfc8xx1123Times, "RFC8xx1123", test)
}

func TestANSIC(test *testing.T) {
	testTimes(ansicTimes, "ANSIC", test)
}

func TestUS(test *testing.T) {
	testTimes(usTimes, "US", test)

	times := []TestTime{
		{
			Value: "01/02/2006 15:04:05",
			Time:  createTimeInLocation("2006-01-02T15:04:05", "2006-01-02T15:04:05", time.Local),
		},
	}

	testTimes(times, "US", test)
}

func TestParse(test *testing.T) {
	testTimes(iso8601Times, "Parse", test)
	testTimes(rfc8xx1123Times, "Parse", test)
	testTimes(ansicTimes, "Parse", test)
	testTimes(usTimes, "Parse", test)
}
