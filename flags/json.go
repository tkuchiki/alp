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
	app.Flag("uri-key", "uri key").
		Default(options.DefaultUriKeyOption).StringVar(&f.UriKey)
	app.Flag("method-key", "method key").
		Default(options.DefaultMethodKeyOption).StringVar(&f.MethodKey)
	app.Flag("time-key", "time key").
		Default(options.DefaultTimeKeyOption).StringVar(&f.TimeKey)
	app.Flag("restime-key", "response time key").
		Default(options.DefaultResponseTimeKeyOption).StringVar(&f.ResponseTimeKey)
	app.Flag("body-bytes-key", "body bytes key").
		Default(options.DefaultBodyBytesKeyOption).StringVar(&f.BodyBytesKey)
	app.Flag("status-key", "status key").
		Default(options.DefaultStatusKeyOption).StringVar(&f.StatusKey)
}
