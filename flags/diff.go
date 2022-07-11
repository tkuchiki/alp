package flags

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

type DiffFlags struct {
	From string
	To   string
}

func NewDiffFlags() *DiffFlags {
	return &DiffFlags{}
}

func (f *DiffFlags) InitFlags(app *kingpin.CmdClause) {
	app.Arg("from", "The comparison source file").Required().StringVar(&f.From)
	app.Arg("to", "The comparison target file").Required().StringVar(&f.To)
}
