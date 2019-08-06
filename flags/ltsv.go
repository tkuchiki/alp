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
	app.Flag("uri-label", "Change the uri label").
		Default(options.DefaultUriLabelOption).StringVar(&f.UriLabel)
	app.Flag("method-label", "Change the method label").
		Default(options.DefaultMethodLabelOption).StringVar(&f.MethodLabel)
	app.Flag("time-label", "Change the time label").
		Default(options.DefaultTimeLabelOption).StringVar(&f.TimeLabel)
	app.Flag("apptime-label", "Change the apptime label").
		Default(options.DefaultApptimeLabelOption).StringVar(&f.ApptimeLabel)
	app.Flag("size-label", "Change the size label").
		Default(options.DefaultSizeLabelOption).StringVar(&f.SizeLabel)
	app.Flag("status-label", "Change the status label").
		Default(options.DefaultStatusLabelOption).StringVar(&f.StatusLabel)
}
