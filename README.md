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

# How to difference between v0.4.0 and v1.0.0

TBW

# Usage

```console
$ alp --help
usage: alp [<flags>] <command> [<args> ...]

alp is the access log profiler for LTSV, JSON, and others.

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
  -c, --config=CONFIG      The configuration file
      --file=FILE          The access log file
  -d, --dump=DUMP          Dump profiled data as YAML
  -l, --load=LOAD          Load the profiled YAML data
      --sort=count         Output the results in sorted order
  -r, --reverse            Sort results in reverse order
  -q, --query-string       Include the URI query string.
      --format=table       The output format (table, markdown, tsv and csv)
      --noheaders          Output no header line at all (only --format=tsv, csv)
      --show-footers       Output footer line at all (only --format=table, markdown)
      --limit=5000         The maximum number of results to display.
      --location=Local     Location name for the timezone
  -o, --output=all         Specifies the results to display, separated by commas
  -m, --matching-groups=PATTERN,...
                           Specifies URI matching groups separated by commas
  -f, --filters=FILTERS    Only the logs are profiled that match the conditions
      --pos=POSITION_FILE  The position file
      --nosave-pos         Do not save position file
      --version            Show application version.

Commands:
  help [<command>...]
    Show help.

  ltsv [<flags>]
    Profile the logs for LTSV

  json [<flags>]
    Profile the logs for JSON

  regexp [<flags>]
    Profile the logs that match a regular expression
```

## ltsv

TBW 

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
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN   |  MAX   |  SUM   |  AVG   |   P1   |  P50   |  P99   | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /foo/bar          |  0.123 |  0.123 |  0.246 |  0.123 |  0.123 |  0.123 |  0.123 |  0.000 |    56.000 |    56.000 |   112.000 |    56.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     5 |   0 |   5 |   0 |   0 |   0 | POST   | /foo/bar          |  0.057 |  0.234 |  0.548 |  0.110 |  0.057 |  0.100 |  0.057 |  0.065 |    12.000 |    34.000 |   126.000 |    25.200 |
|     1 |   0 |   1 |   0 |   0 |   0 | POST   | /hoge/piyo        |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.000 |    34.000 |    34.000 |    34.000 |    34.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
|     1 |   0 |   0 |   0 |   0 |   1 | GET    | /foo/bar/5xx      | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|  11   |  0  | 10  |  0  |  0  |  1  |                                                                                                                                                     
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
```

## json

TBW

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
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN   |  MAX   |  SUM   |  AVG   |   P1   |  P50   |  P99   | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /foo/bar          |  0.123 |  0.123 |  0.246 |  0.123 |  0.123 |  0.123 |  0.123 |  0.000 |    56.000 |    56.000 |   112.000 |    56.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     5 |   0 |   5 |   0 |   0 |   0 | POST   | /foo/bar          |  0.057 |  0.234 |  0.548 |  0.110 |  0.057 |  0.100 |  0.057 |  0.065 |    12.000 |    34.000 |   126.000 |    25.200 |
|     1 |   0 |   1 |   0 |   0 |   0 | POST   | /hoge/piyo        |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.000 |    34.000 |    34.000 |    34.000 |    34.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
|     1 |   0 |   0 |   0 |   0 |   1 | GET    | /foo/bar/5xx      | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|  11   |  0  | 10  |  0  |  0  |  1  |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
```

## regexp

TBW

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
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN   |  MAX   |  SUM   |  AVG   |   P1   |  P50   |  P99   | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|     1 |   0 |   1 |   0 |   0 |   0 | POST   | foo/bar           |  0.057 |  0.057 |  0.057 |  0.057 |  0.057 |  0.057 |  0.057 |  0.000 |    12.000 |    12.000 |    12.000 |    12.000 |
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /foo/bar          |  0.123 |  0.123 |  0.246 |  0.123 |  0.123 |  0.123 |  0.123 |  0.000 |    56.000 |    56.000 |   112.000 |    56.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     4 |   0 |   4 |   0 |   0 |   0 | POST   | /foo/bar          |  0.057 |  0.234 |  0.491 |  0.123 |  0.057 |  0.100 |  0.234 |  0.067 |    12.000 |    34.000 |   114.000 |    28.500 |
|     1 |   0 |   1 |   0 |   0 |   0 | POST   | /hoge/piyo        |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.000 |    34.000 |    34.000 |    34.000 |    34.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
|     1 |   0 |   0 |   0 |   0 |   1 | GET    | /foo/bar/5xx      | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|  11   |  0  | 10  |  0  |  0  |  1  |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
```

## Global options

TBW
    
## URI matching groups

TBW

```console
$ cat example/logs/ltsv_access.log | alp ltsv --filters "Uri matches '^/diary/entry'"
+-------+-----+-----+-----+-----+-----+--------+-------------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN  |  MAX  |  SUM  |  AVG  |  P1   |  P50  |  P99  | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 | 0.135 | 0.135 | 0.135 | 0.135 | 0.135 | 0.135 | 0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 | 0.432 | 0.432 | 0.432 | 0.432 | 0.432 | 0.432 | 0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
+-------+-----+-----+-----+-----+-----+--------+-------------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
|   2   |  0  |  2  |  0  |  0  |  0  |
+-------+-----+-----+-----+-----+-----+--------+-------------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
```

```console
$ cat example/logs/ltsv_access.log | alp ltsv --filters "Uri matches '^/diary/entry'" -m "/diary/entry/.+"
+-------+-----+-----+-----+-----+-----+--------+-----------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |       URI       |  MIN  |  MAX  |  SUM  |  AVG  |  P1   |  P50  |  P99  | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-----------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /diary/entry/.+ | 0.135 | 0.432 | 0.567 | 0.283 | 0.135 | 0.135 | 0.135 |  0.148 |    15.000 |    30.000 |    45.000 |    22.500 |
+-------+-----+-----+-----+-----+-----+--------+-----------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
|   2   |  0  |  2  |  0  |  0  |  0  |
+-------+-----+-----+-----+-----+-----+--------+-----------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
```

## Filter

TBW

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
My Kyash account is `tkuchiki`.

![kyash](./docs/imgs/kyash.jpg "kyash")
