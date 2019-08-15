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
		PlaceHolder(options.DefaultUriLabelOption).StringVar(&f.UriLabel)
	app.Flag("method-label", "Change the method label").
		PlaceHolder(options.DefaultMethodLabelOption).StringVar(&f.MethodLabel)
	app.Flag("time-label", "Change the time label").
		PlaceHolder(options.DefaultTimeLabelOption).StringVar(&f.TimeLabel)
	app.Flag("apptime-label", "Change the apptime label").
		PlaceHolder(options.DefaultApptimeLabelOption).StringVar(&f.ApptimeLabel)
	app.Flag("size-label", "Change the size label").
		PlaceHolder(options.DefaultSizeLabelOption).StringVar(&f.SizeLabel)
	app.Flag("status-label", "Change the status label").
		PlaceHolder(options.DefaultStatusLabelOption).StringVar(&f.StatusLabel)
}
