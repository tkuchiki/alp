package counter

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/tkuchiki/alp/options"

	"github.com/tkuchiki/alp/errors"
	"github.com/tkuchiki/alp/parsers"
)

type Counter struct {
	outWriter io.Writer
	errWriter io.Writer
	inReader  *os.File
	parser    parsers.Parser
	groups    *groups
	printer   *Printer
	options   *options.Options
}

func NewCounter(outw, errw io.Writer, opts *options.Options) *Counter {
	printOptions := NewPrintOptions(opts.NoHeaders, false, opts.PaginationLimit)
	return &Counter{
		outWriter: outw,
		errWriter: errw,
		inReader:  os.Stdin,
		printer:   NewPrinter(outw, opts.Format, printOptions),
		groups:    newGroups(),
		options:   opts,
	}
}

func (c *Counter) SetParser(parser parsers.Parser) {
	c.parser = parser
}

func (c *Counter) SetInReader(f *os.File) {
	c.inReader = f
}

func (c *Counter) Open(filename string) (*os.File, error) {
	var f *os.File
	var err error

	if filename != "" {
		f, err = os.Open(filename)
	} else {
		f = c.inReader
	}

	return f, err
}

func (c *Counter) Count(keys []string) error {
	c.groups.keys = keys

Loop:
	for {
		s, err := c.parser.Parse()
		if err != nil {
			if err == io.EOF {
				break
			} else if err == errors.SkipReadLineErr {
				continue Loop
			}

			return err
		}

		entries := make([]string, 0, len(keys))
		group := newGroup()
		for _, key := range keys {
			value, ok := s.Data.Attributes[key]
			if !ok {
				continue Loop
			}
			val := fmt.Sprint(value)
			entries = append(entries, val)
			group.values[key] = val
		}

		concatenatedKey := groupKey(entries)
		idx := c.groups.index.loadOrStore(concatenatedKey)

		if idx >= len(c.groups.groups) {
			c.groups.groups = append(c.groups.groups, group)
		}

		c.groups.groups[idx].count++
	}

	return nil
}

func groupKey(entries []string) string {
	var key strings.Builder
	for _, entry := range entries {
		key.WriteString(strconv.Itoa(len(entry)))
		key.WriteByte(':')
		key.WriteString(entry)
	}

	return key.String()
}

func (c *Counter) Sort() {
	if c.options.Reverse {
		sort.Slice(c.groups.groups, func(i, j int) bool {
			return c.groups.groups[i].count > c.groups.groups[j].count
		})
	} else {
		sort.Slice(c.groups.groups, func(i, j int) bool {
			return c.groups.groups[i].count < c.groups.groups[j].count
		})
	}
}

func (c *Counter) Print() {
	c.Sort()
	c.printer.Print(c.groups)
}

func (c *Counter) CountAndPrint(keys []string) error {
	err := c.Count(keys)
	if err != nil {
		return err
	}

	c.Print()

	return nil
}

type groupIndex struct {
	values map[string]int
	len    int
	mu     sync.Mutex
}

func newGroupIndex() *groupIndex {
	return &groupIndex{
		values: make(map[string]int),
	}
}

func (index *groupIndex) loadOrStore(key string) int {
	index.mu.Lock()
	defer index.mu.Unlock()
	_, ok := index.values[key]
	if !ok {
		index.values[key] = index.len
		index.len++
	}

	return index.values[key]
}

type groups struct {
	keys   []string
	groups []*group
	index  *groupIndex
}

func newGroups() *groups {
	return &groups{
		index: newGroupIndex(),
	}
}

type group struct {
	values map[string]string
	count  int64
}

func newGroup() *group {
	return &group{
		values: make(map[string]string),
		count:  0,
	}
}

func (g *group) incr() {
	g.count++
}

func (g *group) getCount() int64 {
	return g.count
}
