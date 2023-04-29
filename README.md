# alp

![Test](https://github.com/tkuchiki/alp/workflows/Test/badge.svg?branch=master)

alp is Access Log Profiler

[日本語](./README.ja.md)

# Installation

### Binary distribution

You can pick your download [here](https://github.com/tkuchiki/alp/releases), and install it as follows:

```bash
sudo install <downloaded file> /usr/local/bin/alp
```

### Using your distribution's package system

#### macOS (Homebrew)

Install alp with Homebrew

- `brew install alp`

### asdf

Install alp with [asdf](https://github.com/asdf-vm/asdf) and [asdf-alp](https://github.com/asdf-community/asdf-alp)

```bash
asdf plugin-add alp https://github.com/asdf-community/asdf-alp.git
asdf install alp <version>
asdf global alp <version>
```

# The difference between v0.4.0 and v1.0.0

See: [The difference between v0.4.0 and v1.0.0](./docs/the_difference_between_v0_4_0_and_v1_0_0.md)

# Usage

```console
$ alp --help
alp is the access log profiler for LTSV, JSON, Pcap, and others.

Usage:
  alp [flags]
  alp [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  diff        Show the difference between the two profile results
  help        Help about any command
  json        Profile the logs for JSON
  ltsv        Profile the logs for LTSV
  pcap        Profile the HTTP requests for captured packets
  regexp      Profile the logs that match a regular expression

Flags:
  -h, --help      help for alp
  -v, --version   version for alp

Use "alp [command] --help" for more information about a command.

$ alp ltsv --help
Profile the logs for LTSV

Usage:
  alp ltsv [flags]

Flags:
      --apptime-label string     Change the apptime label (default "apptime")
      --config string            The configuration file
      --decode-uri               Decode the URI
      --dump string              Dump profiled data as YAML
      --file string              The slowlog file
  -f, --filters string           Only the logs are profiled that match the conditions
      --format string            The output format (table, markdown, tsv, csv and html) (default "table")
  -h, --help                     help for ltsv
      --limit int                The maximum number of results to display (default 5000)
      --load string              Load the profiled YAML data
      --location string          Location name for the timezone (default "Local")
  -m, --matching-groups string   Specifies Query matching groups separated by commas
      --method-label string      Change the method label (default "method")
      --noheaders                Output no header line at all (only --format=tsv, csv)
      --nosave-pos               Do not save position file
  -o, --output string            Specifies the results to display, separated by commas (default "all")
      --page int                 Number of pages of pagination (default 100)
      --percentiles string       Specifies the percentiles separated by commas
      --pos string               The position file
      --qs-ignore-values         Ignore the value of the query string. Replace all values with xxx (only use with -q)
  -q, --query-string             Include the URI query string
      --reqtime-label string     Change the reqtime label (default "reqtime")
  -r, --reverse                  Sort results in reverse order
      --show-footers             Output footer line at all (only --format=table, markdown)
      --size-label string        Change the size label (default "size")
      --sort string              Output the results in sorted order (default "count")
      --status-label string      Change the status label (default "status")
      --time-label string        Change the time label (default "time")
      --uri-label string         Change the uri label (default "uri")
      
$ alp json --help
Profile the logs for JSON

Usage:
  alp json [flags]

Flags:
      --body-bytes-key string    Change the body_bytes key (default "body_bytes")
      --config string            The configuration file
      --decode-uri               Decode the URI
      --dump string              Dump profiled data as YAML
      --file string              The slowlog file
  -f, --filters string           Only the logs are profiled that match the conditions
      --format string            The output format (table, markdown, tsv, csv and html) (default "table")
  -h, --help                     help for json
      --limit int                The maximum number of results to display (default 5000)
      --load string              Load the profiled YAML data
      --location string          Location name for the timezone (default "Local")
  -m, --matching-groups string   Specifies Query matching groups separated by commas
      --method-key string        Change the method key (default "method")
      --noheaders                Output no header line at all (only --format=tsv, csv)
      --nosave-pos               Do not save position file
  -o, --output string            Specifies the results to display, separated by commas (default "all")
      --page int                 Number of pages of pagination (default 100)
      --percentiles string       Specifies the percentiles separated by commas
      --pos string               The position file
      --qs-ignore-values         Ignore the value of the query string. Replace all values with xxx (only use with -q)
  -q, --query-string             Include the URI query string
      --reqtime-key string       Change the request_time key (default "request_time")
      --restime-key string       Change the response_time key (default "response_time")
  -r, --reverse                  Sort results in reverse order
      --show-footers             Output footer line at all (only --format=table, markdown)
      --sort string              Output the results in sorted order (default "count")
      --status-key string        Change the status key (default "status")
      --time-key string          Change the time key (default "time")
      --uri-key string           Change the uri key (default "uri")

$ alp regexp --help
Profile the logs that match a regular expression

Usage:
  alp regexp [flags]

Flags:
      --body-bytes-subexp string   Change the body_bytes sub expression (default "body_bytes")
      --config string              The configuration file
      --decode-uri                 Decode the URI
      --dump string                Dump profiled data as YAML
      --file string                The slowlog file
  -f, --filters string             Only the logs are profiled that match the conditions
      --format string              The output format (table, markdown, tsv, csv and html) (default "table")
  -h, --help                       help for regexp
      --limit int                  The maximum number of results to display (default 5000)
      --load string                Load the profiled YAML data
      --location string            Location name for the timezone (default "Local")
  -m, --matching-groups string     Specifies Query matching groups separated by commas
      --method-subexp string       Change the method sub expression (default "method")
      --noheaders                  Output no header line at all (only --format=tsv, csv)
      --nosave-pos                 Do not save position file
  -o, --output string              Specifies the results to display, separated by commas (default "all")
      --page int                   Number of pages of pagination (default 100)
      --pattern string             Regular expressions pattern matching the log (default "^(\\S+)\\s\\S+\\s+(\\S+\\s+)+\\[(?P<time>[^]]+)\\]\\s\"(?P<method>\\S*)\\s?(?P<uri>(?:[^\"]*(?:\\\\\")?)*)\\s([^\"]*)\"\\s(?P<status>\\S+)\\s(?P<body_bytes>\\S+)\\s\"((?:[^\"]*(?:\\\\\")?)*)\"\\s\"(?:.+)\"\\s(?P<response_time>\\S+)(?:\\s(?P<request_time>\\S+))?$")
      --percentiles string         Specifies the percentiles separated by commas
      --pos string                 The position file
      --qs-ignore-values           Ignore the value of the query string. Replace all values with xxx (only use with -q)
  -q, --query-string               Include the URI query string
      --reqtime-subexp string      Change the request_time sub expression (default "request_time")
      --restime-subexp string      Change the response_time sub expression (default "response_time")
  -r, --reverse                    Sort results in reverse order
      --show-footers               Output footer line at all (only --format=table, markdown)
      --sort string                Output the results in sorted order (default "count")
      --status-subexp string       Change the status sub expression (default "status")
      --time-subexp string         Change the time sub expression (default "time")
      --uri-subexp string          Change the uri sub expression (default "uri")
      
$ alp pcap --help
Profile the HTTP requests for captured packets

Usage:
  alp pcap [flags]

Flags:
      --config string             The configuration file
      --decode-uri                Decode the URI
      --dump string               Dump profiled data as YAML
      --file string               The slowlog file
  -f, --filters string            Only the logs are profiled that match the conditions
      --format string             The output format (table, markdown, tsv, csv and html) (default "table")
  -h, --help                      help for pcap
      --limit int                 The maximum number of results to display (default 5000)
      --load string               Load the profiled YAML data
      --location string           Location name for the timezone (default "Local")
  -m, --matching-groups string    Specifies Query matching groups separated by commas
      --noheaders                 Output no header line at all (only --format=tsv, csv)
      --nosave-pos                Do not save position file
  -o, --output string             Specifies the results to display, separated by commas (default "all")
      --page int                  Number of pages of pagination (default 100)
      --pcap-server-ip strings    HTTP server IP address of the captured packets (default [127.0.0.1])
      --pcap-server-port uint16   HTTP server TCP port of the captured packets (default 80)
      --percentiles string        Specifies the percentiles separated by commas
      --pos string                The position file
      --qs-ignore-values          Ignore the value of the query string. Replace all values with xxx (only use with -q)
  -q, --query-string              Include the URI query string
  -r, --reverse                   Sort results in reverse order
      --show-footers              Output footer line at all (only --format=table, markdown)
      --sort string               Output the results in sorted order (default "count")

$ alp diff --help
Show the difference between the two profile results

Usage:
  alp diff <from> <to> [flags]

Flags:
      --from string   The comparison source file
  -h, --help          help for diff
      --to string     The comparison target file
```

## ltsv

- Parses a log in [LTSV](http://ltsv.org/) format
- By default, the following labels are parsed:
    - `time`
        - datetime
    - `method`
        - HTTP Method
    - `uri`
        - URI
    - `status`
        - HTTP Status Code
    - `apptime`
        - Response time from the upstream server
    - `reqtime`
        - Request Processing Time (Response time after receiving a request)
- The `--xxx-label` option can you change the name to any label

```console
$ cat example/logs/ltsv_access.log
time:2015-09-06T05:58:05+09:00	method:POST	uri:/foo/bar?token=xxx&uuid=1234	status:200	size:12	apptime:0.057
time:2015-09-06T05:58:41+09:00	method:POST	uri:/foo/bar?token=yyy	status:200	size:34	apptime:0.100
time:2015-09-06T06:00:42+09:00	method:GET	uri:/foo/bar?token=zzz	status:200	size:56	apptime:0.123
time:2015-09-06T06:00:43+09:00	method:GET	uri:/foo/bar	status:400	size:15	apptime:-
time:2015-09-06T05:58:44+09:00	method:POST	uri:/foo/bar?token=yyy	status:200	size:34	apptime:0.234
time:2015-09-06T05:58:44+09:00	method:POST	uri:/hoge/piyo?id=yyy	status:200	size:34	apptime:0.234
time:2015-09-06T05:58:05+09:00	method:POST	uri:/foo/bar?token=xxx&uuid=1234	status:200	size:12	apptime:0.057
time:2015-09-06T05:58:41+09:00	method:POST	uri:/foo/bar?token=yyy	status:200	size:34	apptime:0.100
time:2015-09-06T06:00:42+09:00	method:GET	uri:/foo/bar?token=zzz	status:200	size:56	apptime:0.123
time:2015-09-06T06:00:43+09:00	method:GET	uri:/foo/bar	status:400	size:15	apptime:-
time:2015-09-06T06:00:43+09:00	method:GET	uri:/diary/entry/1234	status:200	size:15	apptime:0.135
time:2015-09-06T06:00:43+09:00	method:GET	uri:/diary/entry/5678	status:200	size:30	apptime:0.432
time:2015-09-06T06:00:43+09:00	method:GET	uri:/foo/bar/5xx	status:504	size:15	apptime:60.000

$ cat example/logs/ltsv_access.log | alp ltsv
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN   |  MAX   |  SUM   |  AVG   |  P90   |  P95   |  P99   | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /req              |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | POST   | /hoge/piyo        |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.000 |    34.000 |    34.000 |    34.000 |    34.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
|     1 |   0 |   0 |   0 |   0 |   1 | GET    | /foo/bar/5xx      | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /foo/bar          |  0.123 |  0.123 |  0.246 |  0.123 |  0.123 |  0.123 |  0.123 |  0.000 |    56.000 |    56.000 |   112.000 |    56.000 |
|     5 |   0 |   5 |   0 |   0 |   0 | POST   | /foo/bar          |  0.057 |  0.234 |  0.548 |  0.110 |  0.234 |  0.234 |  0.234 |  0.065 |    12.000 |    34.000 |   126.000 |    25.200 |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
```

### Log format 

#### Apache

```
LogFormat "time:%t\tforwardedfor:%{X-Forwarded-For}i\thost:%h\treq:%r\tstatus:%>s\tmethod:%m\turi:%U%q\tsize:%B\treferer:%{Referer}i\tua:%{User-Agent}i\treqtime_microsec:%D\tapptime:%D\tcache:%{X-Cache}o\truntime:%{X-Runtime}o\tvhost:%{Host}i" ltsv
```

#### Nginx

```
log_format ltsv "time:$time_local"
                "\thost:$remote_addr"
                "\tforwardedfor:$http_x_forwarded_for"
                "\treq:$request"
                "\tstatus:$status"
                "\tmethod:$request_method"
                "\turi:$request_uri"
                "\tsize:$body_bytes_sent"
                "\treferer:$http_referer"
                "\tua:$http_user_agent"
                "\treqtime:$request_time"
                "\tcache:$upstream_http_x_cache"
                "\truntime:$upstream_http_x_runtime"
                "\tapptime:$upstream_response_time"
                "\tvhost:$host";
```

#### H2O

```
access-log:
  format: "time:%t\thost:%h\tua:\"%{User-agent}i\"\tstatus:%s\treq:%r\turi:%U\tapptime:%{duration}x\tsize:%b\tmethod:%m"
```

## json

- Parse a log with one JSON per line
- By default, the following keys are parsed:
    - `time`
        - datetime
    - `method`
        - HTTP Method
    - `uri`
        - URI
    - `status`
        - HTTP Status Code
    - `response_time`
        - Response time from the upstream server
    - `request_time`
        - Request Processing Time (Response time after receiving a request)
- The `--xxx-key` option can you change the name to any key

```console
$ cat example/logs/json_access.log
{"time":"2015-09-06T05:58:05+09:00","method":"POST","uri":"/foo/bar?token=xxx&uuid=1234","status":200,"body_bytes":12,"response_time":0.057}
{"time":"2015-09-06T05:58:41+09:00","method":"POST","uri":"/foo/bar?token=yyy","status":200,"body_bytes":34,"response_time":0.100}
{"time":"2015-09-06T06:00:42+09:00","method":"GET","uri":"/foo/bar?token=zzz","status":200,"body_bytes":56,"response_time":0.123}
{"time":"2015-09-06T06:00:43+09:00","method":"GET","uri":"/foo/bar","status":400,"body_bytes":15,"response_time":"-"}
{"time":"2015-09-06T05:58:44+09:00","method":"POST","uri":"/foo/bar?token=yyy","status":200,"body_bytes":34,"response_time":0.234}
{"time":"2015-09-06T05:58:44+09:00","method":"POST","uri":"/hoge/piyo?id=yyy","status":200,"body_bytes":34,"response_time":0.234}
{"time":"2015-09-06T05:58:05+09:00","method":"POST","uri":"/foo/bar?token=xxx&uuid=1234","status":200,"body_bytes":12,"response_time":0.057}
{"time":"2015-09-06T05:58:41+09:00","method":"POST","uri":"/foo/bar?token=yyy","status":200,"body_bytes":34,"response_time":0.100}
{"time":"2015-09-06T06:00:42+09:00","method":"GET","uri":"/foo/bar?token=zzz","status":200,"body_bytes":56,"response_time":0.123}
{"time":"2015-09-06T06:00:43+09:00","method":"GET","uri":"/foo/bar","status":400,"body_bytes":15,"response_time":"-"}
{"time":"2015-09-06T06:00:43+09:00","method":"GET","uri":"/diary/entry/1234","status":200,"body_bytes":15,"response_time":0.135}
{"time":"2015-09-06T06:00:43+09:00","method":"GET","uri":"/diary/entry/5678","status":200,"body_bytes":30,"response_time":0.432}
{"time":"2015-09-06T06:00:43+09:00","method":"GET","uri":"/foo/bar/5xx","status":504,"body_bytes":15,"response_time":60.000}

$ cat example/logs/json_access.log | alp json
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN   |  MAX   |  SUM   |  AVG   |  P90   |  P95   |  P99   | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /req              |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | POST   | /hoge/piyo        |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.000 |    34.000 |    34.000 |    34.000 |    34.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
|     1 |   0 |   0 |   0 |   0 |   1 | GET    | /foo/bar/5xx      | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /foo/bar          |  0.123 |  0.123 |  0.246 |  0.123 |  0.123 |  0.123 |  0.123 |  0.000 |    56.000 |    56.000 |   112.000 |    56.000 |
|     5 |   0 |   5 |   0 |   0 |   0 | POST   | /foo/bar          |  0.057 |  0.234 |  0.548 |  0.110 |  0.234 |  0.234 |  0.234 |  0.065 |    12.000 |    34.000 |   126.000 |    25.200 |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
```

### Log format 

#### Apache

```
LogFormat "{\"time\":\"%t\",\"forwardedfor\":\"%{X-Forwarded-For}i\",\"host\":\"%h\",\"req\":\"%r\",\"status\":%>s,\"method\":\"%m\",\"uri\":\"%U%q\",\"body_bytes\":%B,\"referer\":\"%{Referer}i\",\"ua\":\"%{User-Agent}i\",\"reqtime_microsec\":%D,\"response_time\":%D,\"cache\":\"%{X-Cache}o\",\"runtime\":\"%{X-Runtime}o\",\"vhost\":\"%{Host}i\"}" json
```

#### Nginx

```
    log_format json escape=json '{"time":"$time_local",'
                                '"host":"$remote_addr",'
                                '"forwardedfor":"$http_x_forwarded_for",'
                                '"req":"$request",'
                                '"status":"$status",'
                                '"method":"$request_method",'
                                '"uri":"$request_uri",'
                                '"body_bytes":$body_bytes_sent,'
                                '"referer":"$http_referer",'
                                '"ua":"$http_user_agent",'
                                '"request_time":$request_time,'
                                '"cache":"$upstream_http_x_cache",'
                                '"runtime":"$upstream_http_x_runtime",'
                                '"response_time":"$upstream_response_time",'
                                '"vhost":"$host"}';
```

#### H2O

```
access-log:
  escape: json
  format: '{"time":"%t","host":"%h","ua":"%{User-agent}i","status":%s,"req":"%r","uri":"%U","response_time":%{duration}x,"body_bytes":%b,"method":"%m"}'
```

## regexp

- Parses the log to match the regular expression
- By default, the following named capture groups are parsed:
    - `time`
        - datetime
    - `method`
        - HTTP Method
    - `uri`
        - URI
    - `status`
        - HTTP Status Code
    - `response_time`
        - Response time from the upstream server
    - `request_time`
        - Request Processing Time (Response time after receiving a request)
- The `--xxx-subexp` option can you change the name to any named capture groups

```console
$ cat example/logs/combined_access.log
127.0.0.1 - - [06/Sep/2015:05:58:05 +0900] "POST /foo/bar?token=xxx&uuid=1234 HTTP/1.1" 200 12 "-" "curl/7.54.0" "-" 0.057
127.0.0.1 - - [06/Sep/2015:05:58:41 +0900] "POST /foo/bar?token=yyy HTTP/1.1" 200 34 "-" "curl/7.54.0" "-" 0.100
127.0.0.1 - - [06/Sep/2015:06:00:42 +0900] "GET /foo/bar?token=zzz HTTP/1.1" 200 56 "-" "curl/7.54.0" "-" 0.123
127.0.0.1 - - [06/Sep/2015:06:00:43 +0900] "GET /foo/bar HTTP/1.1" 400 15 "-" "curl/7.54.0" "-" -
127.0.0.1 - - [06/Sep/2015:05:58:44 +0900] "POST /foo/bar?token=yyy HTTP/1.1" 200 34 "-" "curl/7.54.0" "-" 0.234
127.0.0.1 - - [06/Sep/2015:05:58:44 +0900] "POST /hoge/piyo?id=yyy HTTP/1.1" 200 34 "-" "curl/7.54.0" "-" 0.234
127.0.0.1 - - [06/Sep/2015:05:58:05 +0900] "POST foo/bar?token=xxx&uuid=1234 HTTP/1.1" 200 12 "-" "curl/7.54.0" "-" 0.057
127.0.0.1 - - [06/Sep/2015:05:58:41 +0900] "POST /foo/bar?token=yyy HTTP/1.1" 200 34 "-" "curl/7.54.0" "-" 0.100
127.0.0.1 - - [06/Sep/2015:06:00:42 +0900] "GET /foo/bar?token=zzz HTTP/1.1" 200 56 "-" "curl/7.54.0" "-" 0.123
127.0.0.1 - - [06/Sep/2015:06:00:43 +0900] "GET /foo/bar HTTP/1.1" 400 15 "-" "curl/7.54.0" "-" -
127.0.0.1 - - [06/Sep/2015:06:00:43 +0900] "GET /diary/entry/1234 HTTP/1.1" 200 15 "-" "curl/7.54.0" "-" 0.135
127.0.0.1 - - [06/Sep/2015:06:00:43 +0900] "GET /diary/entry/5678 HTTP/1.1" 200 30 "-" "curl/7.54.0" "-" 0.432
127.0.0.1 - - [06/Sep/2015:06:00:43 +0900] "GET /foo/bar/5xx HTTP/1.1" 504 15 "-" "curl/7.54.0" "-" 60.000

$ cat example/logs/combined_access.log | alp regexp
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN   |  MAX   |  SUM   |  AVG   |  P90   |  P95   |  P99   | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /req              |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | POST   | /hoge/piyo        |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.000 |    34.000 |    34.000 |    34.000 |    34.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
|     1 |   0 |   0 |   0 |   0 |   1 | GET    | /foo/bar/5xx      | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /foo/bar          |  0.123 |  0.123 |  0.246 |  0.123 |  0.123 |  0.123 |  0.123 |  0.000 |    56.000 |    56.000 |   112.000 |    56.000 |
|     5 |   0 |   5 |   0 |   0 |   0 | POST   | /foo/bar          |  0.057 |  0.234 |  0.548 |  0.110 |  0.234 |  0.234 |  0.234 |  0.065 |    12.000 |    34.000 |   126.000 |    25.200 |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
```

### Log format

#### Apache

```
LogFormat "%h %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-agent}i\" %D" combined_plus
```

#### Nginx

```
    log_format combined_plus '$remote_addr - $remote_user [$time_local] '
                             '"$request" $status $body_bytes_sent '
                             '"$http_referer" "$http_user_agent" $upstream_response_time $request_time';
```

#### H2O

```
access-log:
  format: "%h %l %u %t \"%r\" %s %b \"%{Referer}i\" \"%{User-agent}i\" %{duration}x"
```

## pcap

- Parses the pcap file to extract HTTP request/response packets to analyze the stats
  - Note that the actual response time may time duration from the actual response time because the difference in the timestamp of packet capturing is regarded as the real response time
  - The IP address and TCP port number of the server are required to distinguish between HTTP requests/responses to the server
- Able to specify the IP address of the HTTP server with the `--pcap-server-ip` option
  - This option can be specified more than once
  - By default, it automatically obtains the IP address from the network interfaces of its own host and uses it
  - However, `127.0.0.1` and `::1` will be the defaults in environments where permissions to retrieve network interface information are restricted.
- Able to specify the TCP port of the HTTP server with the `--pcap-server-port` option
  - The default server port number is 80.
- Cannot be used with `--pos`. (not yet supported)

```console
$ sudo tcpdump -i lo port 5000 -s0 -w http.cap -Z $USER
tcpdump: listening on lo, link-type EN10MB (Ethernet), capture size 262144 bytes
10000 packets captured
20000 packets received by filter
0 packets dropped by kernel

$ alp pcap --file=http.cap --pcap-server-port=5000
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN   |  MAX   |  SUM   |  AVG   |  P90   |  P95   |  P99   | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /req              |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | POST   | /hoge/piyo        |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.000 |    34.000 |    34.000 |    34.000 |    34.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
|     1 |   0 |   0 |   0 |   0 |   1 | GET    | /foo/bar/5xx      | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /foo/bar          |  0.123 |  0.123 |  0.246 |  0.123 |  0.123 |  0.123 |  0.123 |  0.000 |    56.000 |    56.000 |   112.000 |    56.000 |
|     5 |   0 |   5 |   0 |   0 |   0 | POST   | /foo/bar          |  0.057 |  0.234 |  0.548 |  0.110 |  0.234 |  0.234 |  0.234 |  0.065 |    12.000 |    34.000 |   126.000 |    25.200 |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
```

## diff

- Show the difference between the two profile results
- `+` means an increasing number of requests and body size, and response time is slower
- `-` means a decreasing number of requests and body size, and response time is faster

```console
$ cat /path/to/access.log | alp json --dump dumpfile1.yaml

$ cat /path/to/access.log | alp json --dump dumpfile2.yaml

$ alp diff dumpfile1.yaml dumpfile2.yaml -o count,2xx,method,uri,min,max,sum,avg,p90 --show-footers
+---------+---------+--------+-------------------+-----------------+-----------------+-----------------+-----------------+-----------------+
|  COUNT  |   2XX   | METHOD |        URI        |       MIN       |       MAX       |       SUM       |       AVG       |       P90       |
+---------+---------+--------+-------------------+-----------------+-----------------+-----------------+-----------------+-----------------+
| 1       | 1       | GET    | /req              | 0.221 (-0.100)  | 0.221 (-0.100)  | 0.221 (-0.100)  | 0.221 (-0.100)  | 0.221 (-0.100)  |
| 1       | 1       | GET    | /new              | 0.221           | 0.221           | 0.221           | 0.221           | 0.221           |
| 1       | 1       | POST   | /hoge/piyo        | 0.134 (-0.100)  | 0.134 (-0.100)  | 0.134 (-0.100)  | 0.134 (-0.100)  | 0.134 (-0.100)  |
| 1       | 1       | GET    | /diary/entry/1234 | 0.035 (-0.100)  | 0.035 (-0.100)  | 0.035 (-0.100)  | 0.035 (-0.100)  | 0.035 (-0.100)  |
| 1       | 0       | GET    | /foo/bar/5xx      | 59.000 (-1.000) | 59.000 (-1.000) | 59.000 (-1.000) | 59.000 (-1.000) | 59.000 (-1.000) |
| 2 (+1)  | 2 (+1)  | GET    | /diary/entry/5678 | 0.332 (-0.100)  | 2.332 (+1.900)  | 2.664 (+2.232)  | 1.332 (+0.900)  | 2.332 (+1.900)  |
| 2       | 2       | GET    | /foo/bar          | 0.023 (-0.100)  | 0.023 (-0.100)  | 0.046 (-0.200)  | 0.023 (-0.100)  | 0.023 (-0.100)  |
| 5       | 5       | POST   | /foo/bar          | 0.047 (-0.010)  | 0.134 (-0.100)  | 0.378 (-0.170)  | 0.076 (-0.034)  | 0.134 (-0.100)  |
+---------+---------+--------+-------------------+-----------------+-----------------+-----------------+-----------------+-----------------+
| 14 (+2) | 13 (+2) |
+---------+---------+--------+-------------------+-----------------+-----------------+-----------------+-----------------+-----------------+
```

## Global options

See: [Usage samples](./docs/usage_samples.md)

- `-c, --config`
    - The configuration file
    - YAML
- `--file=FILE` 
    - The access log file
- `-d, --dump=DUMP`
    - File path for creating the profile results to a file
- `-l, --load=LOAD`
    - File path to read the results of the profile created with the `-d, --dump` option
    - Can expect it to work fast if you change the `--sort` and `--reverse` options for the same profile results
- `--sort=count`
    - Output the results in sorted order
    - Sort in ascending order
    - `max`, `min`, `sum`, `avg`
    - `max-body`, `min-body`, `sum-body`, `avg-body`  
    - `p90`, `p95`, `p99`, `stddev`
    - `uri`
    - `method`
    - `count`
    - The default is `count`
    - `p90`, `p95`, and `p99` are modified by the values specified in `--percentiles`
- `-r, --reverse`
    - Sort in desecending order
- `-q, --query-string`
    - URIs up to and including query strings are included in the profile
- `--qs-ignore-values`
    - Ignore the value of the query string
    - It's not enabled unless use with `-q, --query-string`
- `--decode-uri`
    - Decode the URI
- `--format=table`
    - Print the profile results in a table, Markdown, TSV, CSV and HTML format
    - The default is table format
- `--noheaders`
    - Print no header when TSV and CSV format
- `--show-footers`
    - Print the total number of each 1xx ~ 5xx in the footer of the table or Markdown format
- `--limit=5000`
    - Maximum number of profile results to be printed
    - This setting is to avoid using too much memory
    - The default is 5000 lines
- `--location="Local"`
    - The timezone of the time specified in the filter condition.
    - Default is  localhost timezone
- `-o, --output="all"`
    - Specify the profile results to be print, separated by commas
    - `count`,`1xx`, `2xx`, `3xx`, `4xx`, `5xx`, `method`, `uri`, `min`, `max`, `sum`, `avg`, `p90`, `p95`, `p99`, `stddev`, `min_body`, `max_body`, `sum_body`, `avg_body`
        - `p90`, `p95`, and `p99` are modified by the values specified in `--percentiles`
    - The default is `all`
- `-m, --matching-groups=PATTERN,...`
    - Treat URIs that match regular expressions as the same URI
    - Evaluate in the specified order. If matched, no further evaluation is performed.
    - See [URI matching groups](#uri-matching-groups)
- `-f, --filters=FILTERS`
    - Filters the targets for profile
    - See [Filter](#filter)
- `--pos=POSITION_FILE`
    - Stores the number of bytes to which the file has been read.
    - If the number of bytes is stored in the POSITION_FILE, the data after that number of bytes will be profiled
    - You can profile without truncating the file
        - Also, it is expected to work fast because it seeks and skips files
- `--nosave-pos`
    - Data after the number of bytes specified by `--pos` is profiled, but the number of bytes reads is not stored
- `--percentiles`
    - Specifies the percentile values to output, separated by commas
    - The default is `90,95,99`
    
## URI matching groups

Consider the following cases like `/diary/entry/1234` and `/diary/entry/5678`.
If you simply profile URIs with different parameters on the same route, they will be profiled by parameter, but you may want to profile them by the route.

```console
$ cat example/logs/ltsv_access.log | alp ltsv --filters "Uri matches '^/diary/entry'"
+-------+-----+-----+-----+-----+-----+--------+-------------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN  |  MAX  |  SUM  |  AVG  |  P90  |  P95  |  P99  | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 | 0.135 | 0.135 | 0.135 | 0.135 | 0.135 | 0.135 | 0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 | 0.432 | 0.432 | 0.432 | 0.432 | 0.432 | 0.432 | 0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
+-------+-----+-----+-----+-----+-----+--------+-------------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
```

```console
$ cat example/logs/ltsv_access.log | alp ltsv --filters "Uri matches '^/diary/entry'" -m "/diary/entry/.+"
+-------+-----+-----+-----+-----+-----+--------+-----------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |       URI       |  MIN  |  MAX  |  SUM  |  AVG  |  P90  |  P95  |  P99  | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-----------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /diary/entry/.+ | 0.135 | 0.432 | 0.567 | 0.283 | 0.432 | 0.432 | 0.432 |  0.148 |    15.000 |    30.000 |    45.000 |    22.500 |
+-------+-----+-----+-----+-----+-----+--------+-----------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
```

In such a case, there is an option `-m, --matching-groups=PATTERN,...`.
You can also specify multiple items separated by commas.

## Filter

It is a function to include or exclude targets according to the conditions.

### Variables

Filter on the following variables:.

- `Uri`
    - URI
- `Method`
    - HTTP Method
- `Time`
    - Datetime string
    - See: https://github.com/tkuchiki/parsetime
- `ResponseTime`
    - Response time
- `BodyBytes`
    - Bytes of HTTP Body 
- `Status`
    - HTTP Status Code

### Operators

The following operators are available:.

- `+`, `-`, `*`, `/`, `%`, `**(pow)` 
- `==`, `!=`, `<`, `>`, `<=`, `>=`
- `not`, `!`
- `and`, `&&`
- `or`, `||`
- `matches`
    - e.g.
       - `Uri matches "PATTERN"`
       - `not(Uri matches "PATTERN")`
- `contains`
    - e.g.
        - `Uri contains "STRING"`
        - `not(Uri contains "STRING")`
- `startsWith`
    - e.g.
        - `Uri startsWith "PREFIX"`
        - `not(Uri startsWith "PREFIX")`
- `endsWith`
    - e.g.
        - `Uri endsWith "SUFFIX"`
        - `not(Uri endsWith "SUFFIX")`
- `in`
    - e.g.
        - `Method in ["GET", "POST"]`
        - `Method not in ["GET", "POST"]`

See: https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md  

### Functions

- `TimeAgo(duration)`
    - now - `duration`
    - `ns`, `us or µs`, `ms`, `s`, `m`, `h`
    - e.g.
        - `Time >= TimeAgo("5m")`
- `BetweenTime(val, start, end)`
    - Like SQL's `BETWEEN`, returns `start <= val && val <= end`
    - e.g.
        - `BetweenTime(Time, "2019-08-06T00:00:00", "2019-08-06T00:05:00")`

## Usage samples

See: [Usage samples](./docs/usage_samples.md)

## Donation

Donations are welcome as always!  
[:heart: Sponsor](https://github.com/sponsors/tkuchiki)
