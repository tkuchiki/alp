// based on encoding/csv
package ltsv

import (
	"bufio"
	"errors"
	"io"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

// A Writer writes records to a LTSV encoded file.
//
// As returned by NewWriter, a Writer writes records terminated by a
// newline and uses '\t' as the field delimiter.  The exported fields can be
// changed to customize the details before the first call to Write or WriteAll.
//
// Delimiter is the field delimiter.
//
// If UseCRLF is true, the Writer ends each record with \r\n instead of \n.
type Writer struct {
	Delimiter rune // Label delimiter (set to to '\t' by NewWriter)
	UseCRLF   bool // True to use \r\n as the line terminator
	w         *bufio.Writer
}

type structIndexLabelMap map[int]string

type formatter func(reflect.Value) (string, error)

var (
	ErrUnsupportedType        = errors.New("unsupported type")
	ErrLabelInvalid           = errors.New("label is invalid")
	ErrFieldInvalid           = errors.New("field is invalid")
	structIndexLabelCache     = make(map[reflect.Type]structIndexLabelMap)
	structIndexLabelCacheLock sync.RWMutex
)

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Delimiter: '\t',
		UseCRLF:   false,
		w:         bufio.NewWriter(w),
	}
}

// Writer writes a single CSV record to w along with any necessary quoting.
// A record is a slice of strings with each string being one field.
func (w *Writer) Write(record interface{}) error {
	m, ok := record.(map[string]string)
	if ok {
		return w.writeMapString(m)
	}
	return w.writeAny(reflect.ValueOf(record))
}

func (w *Writer) writeMapString(record map[string]string) error {
	var err error
	cnt := 0
	for label, field := range record {
		if cnt >= 1 {
			if _, err = w.w.WriteRune(w.Delimiter); err != nil {
				return err
			}
		}
		if err = w.writeLabelAndField(label, field); err != nil {
			return err
		}
		cnt++
	}
	return w.writeLineEnd()
}

func (w *Writer) writeAny(v reflect.Value) error {
	var err error
	switch v.Kind() {
	case reflect.Struct:
		err = w.writeStruct(v)
	case reflect.Map:
		err = w.writeMap(v)
	case reflect.Interface, reflect.Ptr:
		if !v.IsNil() {
			return w.writeAny(v.Elem())
		}
	default:
		return ErrUnsupportedType
	}
	if err != nil {
		return err
	}
	return w.writeLineEnd()
}

// Flush writes any buffered data to the underlying io.Writer.
func (w *Writer) Flush() {
	w.w.Flush()
}

// WriteAll writes multiple LTSV records to w using Write and then calls Flush.
func (w *Writer) WriteAll(records interface{}) error {
	m, ok := records.([]map[string]string)
	if ok {
		return w.writeMapStringAll(m)
	}
	return w.writeAnyAll(reflect.ValueOf(records))
}

func (w *Writer) writeMapStringAll(records []map[string]string) error {
	var err error
	for _, record := range records {
		if err = w.writeMapString(record); err != nil {
			break
		}
	}
	w.Flush()
	return err
}

func (w *Writer) writeAnyAll(v reflect.Value) error {
	k := v.Kind()
	if k == reflect.Slice {
		if v.IsNil() {
			return nil
		}
	} else if k != reflect.Array {
		return ErrUnsupportedType
	}
	var err error
	n := v.Len()
	for i := 0; i < n; i++ {
		if err = w.writeAny(v.Index(i)); err != nil {
			break
		}
	}
	w.Flush()
	return err
}

func (w *Writer) writeLabelAndField(label, field string) error {
	var err error
	if err = w.writeLabel(label); err != nil {
		return err
	}
	if _, err = w.w.WriteRune(':'); err != nil {
		return err
	}
	if err = w.writeField(field); err != nil {
		return err
	}
	return nil
}

func (w *Writer) writeLineEnd() error {
	var line_end string
	if w.UseCRLF {
		line_end = "\r\n"
	} else {
		line_end = "\n"
	}
	_, err := w.w.WriteString(line_end)
	return err
}

func (w *Writer) writeStruct(v reflect.Value) error {
	labels := structIndexLabel(v)
	n := v.NumField()
	cnt := 0
	for i := 0; i < n; i++ {
		label, ok := labels[i]
		if !ok {
			continue
		}
		field, err := reflectValue(v.Field(i))
		if err != nil {
			return err
		}
		if cnt >= 1 {
			if _, err = w.w.WriteRune(w.Delimiter); err != nil {
				return err
			}
		}
		if err = w.writeLabelAndField(label, field); err != nil {
			return err
		}
		cnt++
	}
	return nil
}

func (w *Writer) writeMap(v reflect.Value) error {
	fm := labelFormatter(v.Type().Key())
	if fm == nil {
		return ErrUnsupportedType
	}
	if !v.IsNil() {
		var err error
		for i, key := range v.MapKeys() {
			if i >= 1 {
				if _, err = w.w.WriteRune(w.Delimiter); err != nil {
					return err
				}
			}
			label, err := fm(key)
			if err != nil {
				return err
			}
			field, err := reflectValue(v.MapIndex(key))
			if err != nil {
				return err
			}
			if err = w.writeLabelAndField(label, field); err != nil {
				return err
			}
		}
	}
	return nil
}

func formatInt(v reflect.Value) (string, error) {
	return strconv.FormatInt(v.Int(), 10), nil
}

func formatUint(v reflect.Value) (string, error) {
	return strconv.FormatUint(v.Uint(), 10), nil
}

func formatFloat(v reflect.Value) (string, error) {
	return strconv.FormatFloat(v.Float(), 'g', -1, v.Type().Bits()), nil
}

func formatString(v reflect.Value) (string, error) {
	return v.String(), nil
}

func labelFormatter(t reflect.Type) formatter {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return formatInt
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return formatUint
	case reflect.Float32, reflect.Float64:
		return formatFloat
	case reflect.String:
		return formatString
	}
	return nil
}

func reflectValue(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.Bool:
		if v.Bool() {
			return "true", nil
		} else {
			return "false", nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return formatInt(v)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return formatUint(v)
	case reflect.Float32, reflect.Float64:
		return formatFloat(v)
	case reflect.String:
		return formatString(v)
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return "", nil
		}
		return reflectValue(v.Elem())
	}
	return "", ErrUnsupportedType
}

func structIndexLabel(v reflect.Value) structIndexLabelMap {
	t := v.Type()
	structIndexLabelCacheLock.RLock()
	fs, ok := structIndexLabelCache[t]
	structIndexLabelCacheLock.RUnlock()
	if ok {
		return fs
	}
	structIndexLabelCacheLock.Lock()
	defer structIndexLabelCacheLock.Unlock()
	fs, ok = structIndexLabelCache[t]
	if ok {
		return fs
	}
	labels := make(structIndexLabelMap)
	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if f.Anonymous {
			continue
		}
		label := f.Name
		tv := f.Tag.Get("ltsv")
		if tv != "" {
			if tv == "-" {
				continue
			}
			label, _ = parseTag(tv)
		}
		labels[i] = label
	}
	structIndexLabelCache[t] = labels
	return labels
}

func parseTag(tag string) (string, []string) {
	ss := strings.Split(tag, ",")
	for i, s := range ss {
		ss[i] = strings.TrimSpace(s)
	}
	if len(ss) >= 2 {
		return ss[0], ss[1:]
	}
	return ss[0], []string{}
}

func (w *Writer) writeLabel(s string) error {
	if s == "" {
		return ErrLabelInvalid
	}
	var err error
	for _, c := range s {
		if c == ':' || c == w.Delimiter || c == '\n' || unicode.IsControl(c) || !unicode.IsPrint(c) {
			return ErrLabelInvalid
		}
		if _, err = w.w.WriteRune(c); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeField(s string) error {
	var err error
	for _, c := range s {
		if c == w.Delimiter || c == '\n' {
			return ErrFieldInvalid
		}
		if _, err = w.w.WriteRune(c); err != nil {
			return err
		}
	}
	return nil
}
