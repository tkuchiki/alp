// based on encoding/csv
package ltsv

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"unicode"
)

// A Reader reads records from a LTSV-encoded file.
//
// As returned by NewReader, a Reader expects input LTSV-encoded file.
// The exported fields can be changed to customize the details before the
// first call to Read or ReadAll.
//
// Delimiter is the field delimiter.  It defaults to '\t'.
//
// Comment, if not 0, is the comment character. Lines beginning with the
// Comment character are ignored.
type Reader struct {
	Delimiter rune // Field delimiter (set to '\t' by NewReader)
	Comment   rune // Comment character for start of line
	line      int
	r         *bufio.Reader
	label     bytes.Buffer
	field     bytes.Buffer
}

type structLabelIndexMap map[string]int

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		Delimiter: '\t',
		r:         bufio.NewReader(r),
	}
}

type recorder interface {
	record(label, field string) error
	len() int
}

type mapRecorder struct {
	records map[string]string
}

func (r *mapRecorder) record(label, field string) error {
	r.records[label] = field
	return nil
}

func (r *mapRecorder) len() int {
	return len(r.records)
}

type structRecorder struct {
	count  int
	value  reflect.Value
	labels structLabelIndexMap
}

func (r *structRecorder) record(label, field string) error {
	idx, ok := r.labels[label]
	if !ok {
		return nil
	}
	r.value.Field(idx).SetString(field)
	r.count++
	return nil
}

func (r *structRecorder) len() int {
	return r.count
}

// Read reads one record from r.  The record is a slice of strings with each
// string representing one field.
func (r *Reader) Read() (map[string]string, error) {
	rec := mapRecorder{records: make(map[string]string)}
	err := r.readRecord(&rec)
	if err != nil {
		return nil, err
	}
	return rec.records, nil
}

func (r *Reader) readRecord(rec recorder) error {
	r.line++
	for {
		cm, err := r.readComment()
		if err != nil {
			return err
		} else if cm {
			continue
		}

		label, end, err := r.parseLabel()
		if err != nil {
			return err
		}
		if label == "" {
			if end {
				if rec.len() >= 1 {
					return nil
				}
			}
			continue // skip a empty line
		}

		field, end, err := r.parseField()
		if err != nil {
			return err
		}

		err = rec.record(label, field)
		if err != nil {
			return err
		}
		if end {
			return nil
		}
	}
	panic("unreachable")
}

func (r *Reader) readComment() (bool, error) {
	// If we are support comments and it is the comment character
	// then skip to the end of line.
	r1, err := r.readRune()
	if r.Comment != 0 && r1 == r.Comment {
		for {
			r1, err = r.readRune()
			if err != nil {
				return false, err
			} else if r1 == '\n' {
				return true, nil
			}
		}
	}
	r.r.UnreadRune()
	return false, nil
}

func (r *Reader) Load(record interface{}) error {
	v := reflect.ValueOf(record)
	k := v.Kind()
	if (k != reflect.Ptr && k != reflect.Interface) || v.IsNil() {
		return ErrUnsupportedType
	}
	e := v.Elem()
	rec := structRecorder{value: e, labels: structLabelIndex(e)}
	err := r.readRecord(&rec)
	if err != nil {
		return err
	}
	return nil
}

func (r *Reader) parseLabel() (string, bool, error) {
	r.label.Reset()
	for {
		r1, err := r.readRune()
		if err != nil {
			return "", false, err
		} else if r1 == ':' {
			return strings.TrimSpace(r.label.String()), false, nil
		} else if r1 == '\n' {
			return "", true, nil
		} else if r1 == '\t' {
			return "", false, nil // no label
		} else if unicode.IsControl(r1) || !unicode.IsPrint(r1) {
			return "", false, errors.New(fmt.Sprintf("line %d: invalid rune at label", r.line))
		}
		r.label.WriteRune(r1)
	}
	panic("unreachable")
}

func (r *Reader) parseField() (string, bool, error) {
	r.field.Reset()
	for {
		r1, err := r.readRune()
		if err != nil {
			if err == io.EOF {
				return r.field.String(), true, nil
			}
			return "", false, err
		} else if r1 == '\t' {
			return r.field.String(), false, nil
		} else if r1 == '\n' {
			return r.field.String(), true, nil
		}
		r.field.WriteRune(r1)
	}
	panic("unreachable")
}

func (r *Reader) readRune() (rune, error) {
	r1, _, err := r.r.ReadRune()
	if r1 == '\r' {
		r1, _, err = r.r.ReadRune()
		if err == nil {
			if r1 != '\n' {
				r.r.UnreadRune()
				r1 = '\r'
			}
		}
	}
	return r1, err
}

func (r *Reader) ReadAll() ([]map[string]string, error) {
	records := make([]map[string]string, 0)
	for {
		record, err := r.Read()
		if err == io.EOF {
			return records, nil
		}
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	panic("unreachable")
}

func structLabelIndex(v reflect.Value) structLabelIndexMap {
	sil := structIndexLabel(v)
	sli := make(structLabelIndexMap)
	for i, l := range sil {
		sli[l] = i
	}
	return sli
}
