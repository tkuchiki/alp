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

```
$ ./alp --help                                                                                        
usage: alp --file=FILE [<flags>]

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
  --restime-label="apptime"
                       apptime label
  --body-label="size"  size label
  --method-label="method"
                       method label
  --uri-label="uri"    uri label
  --limit=5000         set an upper limit of the target uri
  --include=PATTERN    don't exclude uri matching PATTERN
  --exclude=PATTERN    exclude uri matching PATTERN
  --version            Show application version.

```

## Log Format

See "Labels for Web server's Log" of http://ltsv.org .

## Sample

[sample log file](./access.log)

### Basic

```
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
```
