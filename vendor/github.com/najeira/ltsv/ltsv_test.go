package ltsv

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

type readerTest struct {
	value   string
	records []map[string]string
}

var readerTests = []readerTest{
	{
		`host:127.0.0.1	ident:-	user:frank	time:[10/Oct/2000:13:55:36 -0700]	req:GET /apache_pb.gif

HTTP/1.0	status:200	size:2326	referer:http://www.example.com/start.html	ua:Mozilla/4.08 [en] (Win98; I ;Nav)
`,
		[]map[string]string{
			{"host": "127.0.0.1", "ident": "-", "user": "frank", "time": "[10/Oct/2000:13:55:36 -0700]", "req": "GET /apache_pb.gif"},
			{"status": "200", "size": "2326", "referer": "http://www.example.com/start.html", "ua": "Mozilla/4.08 [en] (Win98; I ;Nav)"},
		},
	},
	{
		` trimspace :こんにちは
		 trim space :こんばんは
日本語:ラベル
nolabelnofield
ha,s.p-un_ct: おはよう `,
		[]map[string]string{
			{"trimspace": "こんにちは"},
			{"trim space": "こんばんは"},
			{"日本語": "ラベル"},
			{"ha,s.p-un_ct": " おはよう "},
		},
	},
	{
		`label:こんにちは	こんばんは
label2:こんばんは
こんにちは`,
		[]map[string]string{
			{"label": "こんにちは"},
			{"label2": "こんばんは"},
		},
	},
	{
		`
hoge:foo	bar:baz
perl:5.17.8	ruby:2.0	python:2.6
sushi:寿司	tennpura:天ぷら	ramen:ラーメン	gyoza:餃子
		`,
		[]map[string]string{
			{"hoge": "foo", "bar": "baz"},
			{"perl": "5.17.8", "ruby": "2.0", "python": "2.6"},
			{"sushi": "寿司", "tennpura": "天ぷら", "ramen": "ラーメン", "gyoza": "餃子"},
		},
	},
	{
		`
hoge:	bar:
perl:	ruby:
    `,
		[]map[string]string{
			{"hoge": "", "bar": ""},
			{"perl": "", "ruby": ""},
		},
	},
}

type structRecordA struct {
	host  string
	ident string
	user  string
	time  string
	req   string
}

type structRecordB struct {
	status  int
	size    int
	referer string
	ua      string
}

type structTest struct {
	value   string
	records []interface{}
}

var structTests = []structTest{
	{
		`host:127.0.0.1	ident:-	user:frank	time:[10/Oct/2000:13:55:36 -0700]	req:GET /apache_pb.gif
status:200	size:2326	referer:http://www.example.com/start.html	ua:Mozilla/4.08 [en] (Win98; I ;Nav)
`,
		[]interface{}{
			&structRecordA{host: "127.0.0.1", ident: "-", user: "frank", time: "[10/Oct/2000:13:55:36 -0700]", req: "GET /apache_pb.gif"},
			&structRecordB{status: 200, size: 2326, referer: "http://www.example.com/start.html", ua: "Mozilla/4.08 [en] (Win98; I ;Nav)"},
		},
	},
}

type structRecordC struct {
	Host  string `ltsv:"host"`
	Ident string `ltsv:"ident"`
	User  string `ltsv:"user"`
	Time  string `ltsv:"time"`
	Req   string `ltsv:"req"`
}

func (s *structRecordC) String() string {
	return fmt.Sprintf("host:%v	ident:%v	user:%v	time:%v	req:%v",
		s.Host, s.Ident, s.User, s.Time, s.Req)
}

type structLoadTest struct {
	value   string
	records []structRecordC
}

var structLoadTests = []structLoadTest{
	{
		`host:127.0.0.1	ident:-	user:frank	time:[10/Oct/2000:13:55:36 -0700]	req:GET /apache_pb.gif
`,
		[]structRecordC{
			structRecordC{Host: "127.0.0.1", Ident: "-", User: "frank", Time: "[10/Oct/2000:13:55:36 -0700]", Req: "GET /apache_pb.gif"},
		},
	},
}

func TestReaderRead(t *testing.T) {
	for n, test := range readerTests {
		reader := NewReader(bytes.NewBufferString(test.value))
		for i, result := range test.records {
			record, err := reader.Read()
			if err != nil {
				t.Errorf("error %v at test %d, line %d", err, n, i)
			}
			for label, field := range result {
				if record[label] != field {
					t.Errorf("wrong field %v at test %d, line %d, label %s, field %s", record[label], n, i, label, field)
				}
			}
			if len(result) != len(record) {
				t.Errorf("wrong size of record %d at test %d, line %d", len(record), n, i)
			}
		}
		_, err := reader.Read()
		if err == nil || err != io.EOF {
			t.Errorf("expected EOF got %v at test %d", err, n)
		}
	}
}

func TestReaderLoad(t *testing.T) {
	for n, test := range structLoadTests {
		reader := NewReader(bytes.NewBufferString(test.value))
		for i, result := range test.records {
			st := structRecordC{}
			err := reader.Load(&st)
			if err != nil {
				t.Errorf("error %v at test %d, line %d", err, n, i)
			}
			if st.String() != result.String() {
				t.Errorf("got: %s", st.String())
			}
		}
		st := structRecordC{}
		err := reader.Load(&st)
		if err == nil || err != io.EOF {
			t.Errorf("expected EOF got %v at test %d", err, n)
		}
	}
}

func TestReaderReadAll(t *testing.T) {
	for n, test := range readerTests {
		reader := NewReader(bytes.NewBufferString(test.value))
		records, err := reader.ReadAll()
		if err != nil {
			t.Errorf("error %v at test %d", err, n)
		}
		if len(test.records) != len(records) {
			t.Errorf("wrong size of records %d at test %d", len(records), n)
		} else {
			for i, result := range test.records {
				record := records[i]
				for label, field := range result {
					if record[label] != field {
						t.Errorf("wrong field %v at test %d, line %d, label %s, field %s", record[label], n, i, label, field)
					}
				}
				if len(result) != len(record) {
					t.Errorf("wrong size of record %d at test %d, line %d", len(record), n, i)
				}
			}
		}
	}
}

func TestWriterWrite(t *testing.T) {
	var buf bytes.Buffer
	for n, test := range readerTests {
		buf.Reset()
		writer := NewWriter(&buf)
		for i, record := range test.records {
			err := writer.Write(record)
			if err != nil {
				t.Errorf("error %v at test %d, line %d", err, n, i)
			}
		}
		writer.Flush()

		reader := NewReader(&buf)
		records, err := reader.ReadAll()
		if err != nil {
			t.Errorf("error %v at test %d", err, n)
			continue
		}
		if len(records) != len(test.records) {
			t.Errorf("wrong size of records %d at test %d", len(records), n)
		} else {
			for i := 0; i < len(test.records); i++ {
				record := records[i]
				result := test.records[i]
				for label, field := range result {
					if field != record[label] {
						t.Errorf("wrong field %s at test %d, line %d, label %s, field %s", record[label], n, i, label, field)
					}
				}
			}
		}
	}
}

func TestWriterWriteStruct(t *testing.T) {
	var buf bytes.Buffer
	for n, test := range structTests {
		buf.Reset()
		writer := NewWriter(&buf)
		for i, record := range test.records {
			err := writer.Write(record)
			if err != nil {
				t.Errorf("error %v at test %d, line %d", err, n, i)
			}
		}
		writer.Flush()
		if buf.String() != test.value {
			t.Errorf("expect:\n%s\ngot:\n%v\n", test.value, buf.String())
		}
	}
}

func TestWriterWriteAll(t *testing.T) {
	var buf bytes.Buffer
	for n, test := range readerTests {
		buf.Reset()
		writer := NewWriter(&buf)
		writer.WriteAll(test.records)

		reader := NewReader(&buf)
		records, err := reader.ReadAll()
		if err != nil {
			t.Errorf("error %v at test %d", err, n)
			continue
		}
		if len(records) != len(test.records) {
			t.Errorf("wrong size of records %d at test %d", len(records), n)
		} else {
			for i := 0; i < len(test.records); i++ {
				record := records[i]
				result := test.records[i]
				for label, field := range result {
					if field != record[label] {
						t.Errorf("wrong field %s at test %d, line %d, label %s, field %s", record[label], n, i, label, field)
					}
				}
			}
		}
	}
}

func TestWriterWriteAllStruct(t *testing.T) {
	var buf bytes.Buffer
	for _, test := range structTests {
		buf.Reset()
		writer := NewWriter(&buf)
		writer.WriteAll(test.records)
		if buf.String() != test.value {
			t.Errorf("expect:\n%s\ngot:\n%v\n", test.value, buf.String())
		}
	}
}

func BenchmarkReaderRead(b *testing.B) {
	for i := 0; i < 10000; i++ {
		reader := NewReader(bytes.NewBufferString(readerTests[3].value))
		_, e := reader.Read()
		if e != nil {
			break
		}
	}
}

func BenchmarkWriterWrite(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < 10000; i++ {
		for _, test := range readerTests {
			buf.Reset()
			writer := NewWriter(&buf)
			for _, record := range test.records {
				writer.Write(record)
			}
			writer.Flush()
		}
	}
}
