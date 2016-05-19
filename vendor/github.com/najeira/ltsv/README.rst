ltsv
====

LTSV (Labeled Tab-separated Values) reader/writer for Go language.

About LTSV: http://ltsv.org/

	Labeled Tab-separated Values (LTSV) format is a variant of 
	Tab-separated Values (TSV).Each record in a LTSV file is represented 
	as a single line. Each field is separated by TAB and has a label and
	 a value. The label and the value have been separated by ':'. With 
	the LTSV format, you can parse each line by spliting with TAB (like 
	original TSV format) easily, and extend any fields with unique labels 
	in no particular order.


Example
=======

Reader
------

::

  package main
  
  import (
  	"bytes"
  	"fmt"
  	"github.com/najeira/ltsv"
  )
  
  func main() {
  	data := `
  time:05/Feb/2013:15:34:47 +0000 host:192.168.50.1	req:GET / HTTP/1.1	status:200
  time:05/Feb/2013:15:35:15 +0000 host:192.168.50.1	req:GET /foo HTTP/1.1	status:200
  time:05/Feb/2013:15:35:54 +0000 host:192.168.50.1	req:GET /bar HTTP/1.1	status:404
  `
  	b := bytes.NewBufferString(data)
  	
  	// Read LTSV file into map[string]string
  	reader := ltsv.NewReader(b)
  	records, err := reader.ReadAll()
  	if err != nil {
  		panic(err)
  	}
  	
  	// dump
  	for i, record := range records {
  		fmt.Printf("===== Data %d\n", i)
  		for k, v := range record {
  			fmt.Printf("\t%s --> %s\n", k, v)
  		}
  	}
  }


Writer
------

::

  package main
  	
  import (
  	"fmt"
  	"bytes"
  	"github.com/najeira/ltsv"
  )
  	
  func main() {
  	data := []map[string]string {
  		{"time": "05/Feb/2013:15:34:47 +0000", "host": "192.168.50.1", "req": "GET / HTTP/1.1", "status": "200"},
  		{"time": "05/Feb/2013:15:35:15 +0000", "host": "192.168.50.1", "req": "GET /foo HTTP/1.1", "status": "200"},
  		{"time": "05/Feb/2013:15:35:54 +0000", "host": "192.168.50.1", "req": "GET /bar HTTP/1.1", "status": "404"},
  	}
  	
  	b := &bytes.Buffer{}
  	writer := ltsv.NewWriter(b)
  	err := writer.WriteAll(data)
  	if err != nil {
  		panic(err)
  	}
  	fmt.Printf("%v", b.String())
  }


License
=======

New BSD License.


Links
=====

- http://ltsv.org/
- https://github.com/ymotongpoo/goltsv  LTSV package by ymotongpoo
