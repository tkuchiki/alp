# alp

alp is Access Log Profiler for Labeled Tab-separated Values(LTSV).  
(See: [Labeled Tab-separated Values](http://ltsv.org))

# Installation

```
curl -sLO https://github.com/tkuchiki/alp/releases/download/VERSION/alp_linux_amd64.zip
unzip alp_linux_amd64.zip
mv alp_linux_amd64 /usr/local/bin/alp
```

# Usage

Read from stdin or an input file(`-f`).  

```
$ ./alp --help
usage: alp [<flags>]

Access Log Profiler for LTSV (read from file or stdin).

Flags:
  --help               Show context-sensitive help (also try --help-long and --help-man).
  -f, --file=FILE      access log file
  --max                sort by max response time
  --min                sort by min response time
  --avg                sort by avg response time
  --sum                sort by sum response time
  --cnt                sort by count
  --uri                sort by uri
  --method             sort by method
  --max-body           sort by max body size
  --min-body           sort by min body size
  --avg-body           sort by avg body size
  --sum-body           sort by sum body size
  -r, --reverse        reverse the result of comparisons
  -q, --query-string   include query string
  --tsv                tsv format (default: table)
  --apptime-label="apptime"
                       apptime label
  --size-label="size"  size label
  --method-label="method"
                       method label
  --uri-label="uri"    uri label
  --limit=5000         set an upper limit of the target uri
  --includes=PATTERN,...
                       don't exclude uri matching PATTERN (comma separated)
  --excludes=PATTERN,...
                       exclude uri matching PATTERN (comma separated)
  --noheaders          print no header line at all (only --tsv)
  --aggregates=PATTERN,...
                       aggregate uri matching PATTERN (comma separated)
  --version            Show application version.

```

## Log format

See "Labels for Web server's Log" of http://ltsv.org .

### Apache

```
LogFormat "time:%t\tforwardedfor:%{X-Forwarded-For}i\thost:%h\treq:%r\tstatus:%>s\tmethod:%m\turi:%U%q\tsize:%B\treferer:%{Referer}i\tua:%{User-Agent}i\treqtime_microsec:%D\tapptime:%D\tcache:%{X-Cache}o\truntime:%{X-Runtime}o\tvhost:%{Host}i" ltsv
```

### Nginx

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

## Sample

[sample log file](./access.log)

### Basic

```
$ cat access.log | ./alp
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+----------+
| COUNT |  MIN  |  MAX  |  SUM  |  AVG  | MAX(BODY) | MIN(BODY) | SUM(BODY) | AVG(BODY) | METHOD |   URI    |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+----------+
| 1     | 0.123 | 0.123 | 0.123 | 0.123 |    56.000 |    56.000 |    56.000 |    56.000 | GET    | /foo/bar |
| 3     | 0.057 | 0.234 | 0.391 | 0.130 |    12.000 |    34.000 |    80.000 |    26.667 | POST   | /foo/bar |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+----------+

$ ./alp -f access.log
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+----------+
| COUNT |  MIN  |  MAX  |  SUM  |  AVG  | MAX(BODY) | MIN(BODY) | SUM(BODY) | AVG(BODY) | METHOD |   URI    |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+----------+
| 1     | 0.123 | 0.123 | 0.123 | 0.123 |    56.000 |    56.000 |    56.000 |    56.000 | GET    | /foo/bar |
| 3     | 0.057 | 0.234 | 0.391 | 0.130 |    12.000 |    34.000 |    80.000 |    26.667 | POST   | /foo/bar |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+----------+

$ ./alp -f access.log -r
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+----------+
| COUNT |  MIN  |  MAX  |  SUM  |  AVG  | MAX(BODY) | MIN(BODY) | SUM(BODY) | AVG(BODY) | METHOD |   URI    |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+----------+
| 3     | 0.057 | 0.234 | 0.391 | 0.130 |    12.000 |    34.000 |    80.000 |    26.667 | POST   | /foo/bar |
| 1     | 0.123 | 0.123 | 0.123 | 0.123 |    56.000 |    56.000 |    56.000 |    56.000 | GET    | /foo/bar |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+----------+

$ ./alp -f access.log -q
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-----------------------------+
| COUNT |  MIN  |  MAX  |  SUM  |  AVG  | MAX(BODY) | MIN(BODY) | SUM(BODY) | AVG(BODY) | METHOD |             URI             |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-----------------------------+
| 1     | 0.057 | 0.057 | 0.057 | 0.057 |    12.000 |    12.000 |    12.000 |    12.000 | POST   | /foo/bar?token=xxx&uuid=xxx |
| 1     | 0.123 | 0.123 | 0.123 | 0.123 |    56.000 |    56.000 |    56.000 |    56.000 | GET    | /foo/bar?token=xxx          |
| 2     | 0.100 | 0.234 | 0.334 | 0.167 |    34.000 |    34.000 |    68.000 |    34.000 | POST   | /foo/bar?token=xxx          |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-----------------------------+

$ ./alp -f access.log -q -r
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-----------------------------+
| COUNT |  MIN  |  MAX  |  SUM  |  AVG  | MAX(BODY) | MIN(BODY) | SUM(BODY) | AVG(BODY) | METHOD |             URI             |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-----------------------------+
| 2     | 0.100 | 0.234 | 0.334 | 0.167 |    34.000 |    34.000 |    68.000 |    34.000 | POST   | /foo/bar?token=xxx          |
| 1     | 0.123 | 0.123 | 0.123 | 0.123 |    56.000 |    56.000 |    56.000 |    56.000 | GET    | /foo/bar?token=xxx          |
| 1     | 0.057 | 0.057 | 0.057 | 0.057 |    12.000 |    12.000 |    12.000 |    12.000 | POST   | /foo/bar?token=xxx&uuid=xxx |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-----------------------------+
```

### TSV

```
$ ./alp -f access.log --tsv
Count	Min	Max	Sum	Avg	Max(Body)	Min(Body)	Sum(Body)	Avg(Body)	Method	Uri
1	0.123	0.123	0.123	0.123	56	56	56	56	GET	/foo/bar
3	0.057	0.234	0.391	0.130	12	34	80	26.667	POST	/foo/bar

$ ./alp -f access.log --tsv --noheaders
1	0.123	0.123	0.123	0.123	56	56	56	56	GET	/foo/bar
3	0.057	0.234	0.391	0.130	12	34	80	26.667	POST	/foo/bar
```

## Include

```
$ ./alp -f access.log
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-------------------+
| COUNT |  MIN  |  MAX  |  SUM  |  AVG  | MAX(BODY) | MIN(BODY) | SUM(BODY) | AVG(BODY) | METHOD |        URI        |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-------------------+
| 2     | 0.123 | 0.123 | 0.246 | 0.123 |    56.000 |    56.000 |   112.000 |    56.000 | GET    | /foo/bar          |
| 1     | 0.135 | 0.135 | 0.135 | 0.135 |    15.000 |    15.000 |    15.000 |    15.000 | GET    | /diary/entry/1234 |
| 5     | 0.057 | 0.234 | 0.548 | 0.110 |    12.000 |    34.000 |   126.000 |    25.200 | POST   | /foo/bar          |
| 1     | 0.234 | 0.234 | 0.234 | 0.234 |    34.000 |    34.000 |    34.000 |    34.000 | POST   | /hoge/piyo        |
| 1     | 0.432 | 0.432 | 0.432 | 0.432 |    30.000 |    30.000 |    30.000 |    30.000 | GET    | /diary/entry/5678 |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-------------------+

$ ./alp -f access.log --includes "foo,\d+"
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-------------------+
| COUNT |  MIN  |  MAX  |  SUM  |  AVG  | MAX(BODY) | MIN(BODY) | SUM(BODY) | AVG(BODY) | METHOD |        URI        |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-------------------+
| 2     | 0.123 | 0.123 | 0.246 | 0.123 |    56.000 |    56.000 |   112.000 |    56.000 | GET    | /foo/bar          |
| 1     | 0.135 | 0.135 | 0.135 | 0.135 |    15.000 |    15.000 |    15.000 |    15.000 | GET    | /diary/entry/1234 |
| 5     | 0.057 | 0.234 | 0.548 | 0.110 |    12.000 |    34.000 |   126.000 |    25.200 | POST   | /foo/bar          |
| 1     | 0.432 | 0.432 | 0.432 | 0.432 |    30.000 |    30.000 |    30.000 |    30.000 | GET    | /diary/entry/5678 |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-------------------+
```

## Exclude

```
$ ./alp -f access.log
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-------------------+
| COUNT |  MIN  |  MAX  |  SUM  |  AVG  | MAX(BODY) | MIN(BODY) | SUM(BODY) | AVG(BODY) | METHOD |        URI        |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-------------------+
| 2     | 0.123 | 0.123 | 0.246 | 0.123 |    56.000 |    56.000 |   112.000 |    56.000 | GET    | /foo/bar          |
| 1     | 0.135 | 0.135 | 0.135 | 0.135 |    15.000 |    15.000 |    15.000 |    15.000 | GET    | /diary/entry/1234 |
| 5     | 0.057 | 0.234 | 0.548 | 0.110 |    12.000 |    34.000 |   126.000 |    25.200 | POST   | /foo/bar          |
| 1     | 0.234 | 0.234 | 0.234 | 0.234 |    34.000 |    34.000 |    34.000 |    34.000 | POST   | /hoge/piyo        |
| 1     | 0.432 | 0.432 | 0.432 | 0.432 |    30.000 |    30.000 |    30.000 |    30.000 | GET    | /diary/entry/5678 |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-------------------+

$ ./alp -f access.log --excludes "foo,\d+"
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+------------+
| COUNT |  MIN  |  MAX  |  SUM  |  AVG  | MAX(BODY) | MIN(BODY) | SUM(BODY) | AVG(BODY) | METHOD |    URI     |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+------------+
| 1     | 0.234 | 0.234 | 0.234 | 0.234 |    34.000 |    34.000 |    34.000 |    34.000 | POST   | /hoge/piyo |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+------------+
```

### Aggregate

```
$ ./alp -f access.log
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-------------------+
| COUNT |  MIN  |  MAX  |  SUM  |  AVG  | MAX(BODY) | MIN(BODY) | SUM(BODY) | AVG(BODY) | METHOD |        URI        |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-------------------+
| 2     | 0.123 | 0.123 | 0.246 | 0.123 |    56.000 |    56.000 |   112.000 |    56.000 | GET    | /foo/bar          |
| 1     | 0.135 | 0.135 | 0.135 | 0.135 |    15.000 |    15.000 |    15.000 |    15.000 | GET    | /diary/entry/1234 |
| 5     | 0.057 | 0.234 | 0.548 | 0.110 |    12.000 |    34.000 |   126.000 |    25.200 | POST   | /foo/bar          |
| 1     | 0.234 | 0.234 | 0.234 | 0.234 |    34.000 |    34.000 |    34.000 |    34.000 | POST   | /hoge/piyo        |
| 1     | 0.432 | 0.432 | 0.432 | 0.432 |    30.000 |    30.000 |    30.000 |    30.000 | GET    | /diary/entry/5678 |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+-------------------+

$ ./alp -f access.log --aggregates "/diary/entry/\d+"
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+------------------+
| COUNT |  MIN  |  MAX  |  SUM  |  AVG  | MAX(BODY) | MIN(BODY) | SUM(BODY) | AVG(BODY) | METHOD |       URI        |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+------------------+
| 2     | 0.123 | 0.123 | 0.246 | 0.123 |    56.000 |    56.000 |   112.000 |    56.000 | GET    | /foo/bar         |
| 5     | 0.057 | 0.234 | 0.548 | 0.110 |    12.000 |    34.000 |   126.000 |    25.200 | POST   | /foo/bar         |
| 1     | 0.234 | 0.234 | 0.234 | 0.234 |    34.000 |    34.000 |    34.000 |    34.000 | POST   | /hoge/piyo       |
| 2     | 0.135 | 0.432 | 0.567 | 0.283 |    15.000 |    30.000 |    45.000 |    22.500 | GET    | /diary/entry/\d+ |
+-------+-------+-------+-------+-------+-----------+-----------+-----------+-----------+--------+------------------+
```
