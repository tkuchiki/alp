package flags

import (
	"github.com/tkuchiki/alp/options"
	"gopkg.in/alecthomas/kingpin.v2"
)

type RegexpFlags struct {
	Pattern            string
	UriSubexp          string
	MethodSubexp       string
	TimeSubexp         string
	ResponseTimeSubexp string
	BodyBytesSubexp    string
	StatusSubexp       string
}

func NewRegexpFlags() *RegexpFlags {
	return &RegexpFlags{}
}

func (f *RegexpFlags) InitFlags(app *kingpin.CmdClause) {
	app.Flag("pattern", "Regular expressions pattern matching the log").
		PlaceHolder(options.DefaultPatternOption).PlaceHolder("PATTERN").StringVar(&f.Pattern)
	app.Flag("uri-subexp", "Change the uri sub expression").
		PlaceHolder(options.DefaultUriSubexpOption).StringVar(&f.UriSubexp)
	app.Flag("method-subexp", "Change the method sub expression").
		PlaceHolder(options.DefaultMethodSubexpOption).StringVar(&f.MethodSubexp)
	app.Flag("time-subexp", "Change the time sub expression").
		PlaceHolder(options.DefaultTimeSubexpOption).StringVar(&f.TimeSubexp)
	app.Flag("restime-subexp", "Change the response_time sub expression").
		PlaceHolder(options.DefaultResponseTimeSubexpOption).StringVar(&f.ResponseTimeSubexp)
	app.Flag("body-bytes-subexp", "Change the body_bytes sub expression").
		PlaceHolder(options.DefaultBodyBytesSubexpOption).StringVar(&f.BodyBytesSubexp)
	app.Flag("status-subexp", "Change the status sub expression").
		PlaceHolder(options.DefaultStatusSubexpOption).StringVar(&f.StatusSubexp)
}
