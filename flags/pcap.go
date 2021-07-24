package flags

import (
	"fmt"
	"strconv"

	"github.com/tkuchiki/alp/options"
	"gopkg.in/alecthomas/kingpin.v2"
)

type PcapFlags struct {
	ServerIPs  []string
	ServerPort uint16
}

func NewPcapFlags() *PcapFlags {
	return &PcapFlags{}
}

func (f *PcapFlags) InitFlags(app *kingpin.CmdClause) {
	app.Flag("pcap-server-ip", "HTTP server IP address of the captured packets").
		PlaceHolder(options.DefaultPcapServerIPsOption[0]).StringsVar(&f.ServerIPs)
	app.Flag("pcap-server-port", "HTTP server TCP port of the captured packets").Default(strconv.FormatUint(options.DefaultPcapServerPortOption, 10)).
		PlaceHolder(fmt.Sprint(options.DefaultPcapServerPortOption)).Uint16Var(&f.ServerPort)
}
