package flags

import (
	"github.com/tkuchiki/alp/options"
	"gopkg.in/alecthomas/kingpin.v2"
)

type LTSVFlags struct {
	UriLabel     string
	MethodLabel  string
	TimeLabel    string
	ApptimeLabel string
	SizeLabel    string
	StatusLabel  string
}

func NewLTSVFlags() *LTSVFlags {
	return &LTSVFlags{}
}

func (f *LTSVFlags) InitFlags(app *kingpin.CmdClause) {
	app.Flag("uri-label", "uri label").
		Default(options.DefaultUriLabelOption).StringVar(&f.UriLabel)
	app.Flag("method-label", "method label").
		Default(options.DefaultMethodLabelOption).StringVar(&f.MethodLabel)
	app.Flag("time-label", "time label").
		Default(options.DefaultTimeLabelOption).StringVar(&f.TimeLabel)
	app.Flag("apptime-label", "apptime label").
		Default(options.DefaultApptimeLabelOption).StringVar(&f.ApptimeLabel)
	app.Flag("size-label", "size label").
		Default(options.DefaultSizeLabelOption).StringVar(&f.SizeLabel)
	app.Flag("status-label", "status label").
		Default(options.DefaultStatusLabelOption).StringVar(&f.StatusLabel)
}
