package flags

import (
	"github.com/tkuchiki/alp/options"
	"gopkg.in/alecthomas/kingpin.v2"
)

type JSONFlags struct {
	UriKey          string
	MethodKey       string
	TimeKey         string
	ResponseTimeKey string
	BodyBytesKey    string
	StatusKey       string
}

func NewJSONFlags() *JSONFlags {
	return &JSONFlags{}
}

func (f *JSONFlags) InitFlags(app *kingpin.CmdClause) {
	app.Flag("uri-key", "Change the uri key").
		Default(options.DefaultUriKeyOption).StringVar(&f.UriKey)
	app.Flag("method-key", "Change the method key").
		Default(options.DefaultMethodKeyOption).StringVar(&f.MethodKey)
	app.Flag("time-key", "Change the time key").
		Default(options.DefaultTimeKeyOption).StringVar(&f.TimeKey)
	app.Flag("restime-key", "Change the response_time key").
		Default(options.DefaultResponseTimeKeyOption).StringVar(&f.ResponseTimeKey)
	app.Flag("body-bytes-key", "Change the body_bytes key").
		Default(options.DefaultBodyBytesKeyOption).StringVar(&f.BodyBytesKey)
	app.Flag("status-key", "Change the status key").
		Default(options.DefaultStatusKeyOption).StringVar(&f.StatusKey)
}
