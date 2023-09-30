package log_reader

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/tkuchiki/alp/errors"
	"github.com/tkuchiki/alp/helpers"
	"github.com/tkuchiki/alp/options"
	"github.com/tkuchiki/alp/parsers"
	"github.com/tkuchiki/alp/stats"
)

type AccessLog struct {
	Uri          string
	Method       string
	ResponseTime float64
	BodyBytes    float64
	Status       int
	TimeStr      string
}

type AccessLogReader struct {
	logs      []*AccessLog
	options   *options.Options
	outWriter io.Writer
	errWriter io.Writer
	inReader  *os.File
	printer   *Printer
	numOfTopN int
}

func NewAccessLogReader(outw, errw io.Writer, opts *options.Options, numOfTopN int) *AccessLogReader {
	printOptions := NewPrintOptions(opts.NoHeaders, opts.DecodeUri, opts.PaginationLimit)
	printer := NewPrinter(outw, opts.Format, printOptions)

	opts = options.SetOptions(opts,
		options.QueryString(true),
	)

	return &AccessLogReader{
		options:   opts,
		outWriter: outw,
		errWriter: errw,
		inReader:  os.Stdin,
		printer:   printer,
		numOfTopN: numOfTopN,
	}
}

func (a *AccessLogReader) SetInReader(f *os.File) {
	a.inReader = f
}

func (a *AccessLogReader) Open(filename string) (*os.File, error) {
	var f *os.File
	var err error

	if filename != "" {
		f, err = os.Open(filename)
	} else {
		f = a.inReader
	}

	return f, err
}

func (a *AccessLogReader) OpenPosFile(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
}

func (a *AccessLogReader) ReadPosFile(f *os.File) (int, error) {
	reader := bufio.NewReader(f)
	pos, _, err := reader.ReadLine()
	if err != nil {
		return 0, err
	}

	return helpers.StringToInt(string(pos))
}

func (a *AccessLog) UriWithOptions(decode bool) string {
	if !decode {
		return a.Uri
	}

	u, err := url.Parse(a.Uri)
	if err != nil {
		return a.Uri
	}

	if u.RawQuery == "" {
		unescaped, _ := url.PathUnescape(u.EscapedPath())
		return unescaped
	}

	unescaped, _ := url.PathUnescape(u.EscapedPath())
	decoded, _ := url.QueryUnescape(u.Query().Encode())

	return fmt.Sprintf("%s?%s", unescaped, decoded)
}

func (a *AccessLogReader) Append(uri, method, time string, responseTime, bodyBytes float64, status int) {
	a.logs = append(a.logs, &AccessLog{
		Uri:          uri,
		Method:       method,
		TimeStr:      time,
		ResponseTime: responseTime,
		BodyBytes:    bodyBytes,
		Status:       status,
	})
}

func (a *AccessLogReader) ReadAll(parser parsers.Parser) error {
	var err error
	var posfile *os.File
	if a.options.PosFile != "" {
		posfile, err = a.OpenPosFile(a.options.PosFile)
		if err != nil {
			return err
		}
		defer posfile.Close()

		pos, err := a.ReadPosFile(posfile)
		if err != nil && err != io.EOF {
			return err
		}

		err = parser.Seek(pos)
		if err != nil {
			return err
		}

		parser.SetReadBytes(pos)
	}

	sts := stats.NewHTTPStats(true, false, false)
	err = sts.InitFilter(a.options)
	if err != nil {
		return err
	}
	sts.SetOptions(a.options)

Loop:
	for {
		s, err := parser.Parse()
		if err != nil {
			if err == io.EOF {
				break
			} else if err == errors.SkipReadLineErr {
				continue Loop
			}

			return err
		}

		var b bool
		b, err = sts.DoFilter(s)
		if err != nil {
			return err
		}

		if !b {
			continue Loop
		}

		a.Append(s.Uri, s.Method, s.Time, s.ResponseTime, s.BodyBytes, s.Status)
	}

	if !a.options.NoSavePos && a.options.PosFile != "" {
		posfile.Seek(0, 0)
		_, err = posfile.Write([]byte(fmt.Sprint(parser.ReadBytes())))
		if err != nil {
			return err
		}
	}

	err = a.Sort(a.options.TopN.Sort, a.options.TopN.Reverse)

	return err
}

func (a *AccessLogReader) Print() {
	var n int
	numOfLogs := len(a.logs)

	if a.numOfTopN > numOfLogs {
		n = numOfLogs
	} else {
		n = a.numOfTopN
	}

	a.printer.Print(a.logs[0:n])
}
