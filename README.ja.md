# alp

![Test](https://github.com/tkuchiki/alp/workflows/Test/badge.svg?branch=master)

alp はアクセスログ解析ツールです。
読み方は、[éɪélpíː] または [ælp] です(作者は [éɪélpíː] 派)。

# インストール

### バイナリ配布

[ここ](https://github.com/tkuchiki/alp/releases)から任意のOS向けのバイナリをダウンロードすることができ、次のようにしてインストールすることが可能です。

```bash
sudo install <ダウンロードしたファイル> /usr/local/bin/alp
```

### ディストリビューションのパッケージシステムを使用

#### macOS (Homebrew)

Homebrewでalpをインストールします。

- `brew install alp`

### asdf

[asdf](https://github.com/asdf-vm/asdf)と[asdf-alp](https://github.com/asdf-community/asdf-alp)でalpをインストールします。

```bash
asdf plugin-add alp https://github.com/asdf-community/asdf-alp.git
asdf install alp <バージョン>
asdf global alp <バージョン>
```

# v0.4.0 と v1.0.0 の違い

[v0.4.0 と v1.0.0 の違い](./docs/the_difference_between_v0_4_0_and_v1_0_0.ja.md) を参照してください。

# 使い方

```console
$ alp --help
alp is the access log profiler for LTSV, JSON, Pcap, and others.

Usage:
  alp [flags]
  alp [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  count       Count by log entries
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

$ alp count --help
Count by log entries

Usage:
  alp count [flags]

Flags:
      --file string      The access log file
      --format string    Log format (json,ltsv,regexp) (default "json")
  -h, --help             help for count
      --keys string      Log key names (comma separated)
      --pattern string   Regular expressions pattern matching the log. (only use with --format=regexp) (default "^(\\S+)\\s\\S+\\s+(\\S+\\s+)+\\[(?P<time>[^]]+)\\]\\s\"(?P<method>\\S*)\\s?(?P<uri>(?:[^\"]*(?:\\\\\")?)*)\\s([^\"]*)\"\\s(?P<status>\\S+)\\s(?P<body_bytes>\\S+)\\s\"((?:[^\"]*(?:\\\\\")?)*)\"\\s\"(?:.+)\"\\s(?P<response_time>\\S+)(?:\\s(?P<request_time>\\S+))?$")
  -r, --reverse          Sort results in reverse order
```

- alp は ltsv, json, regexp, pcap, diff の5つのサブコマンドで構成されています
- `cat /path/to/log | alp` のようにパイプでデータを送るか、後述する `-f, --file` オプションでファイルを指定して解析します

## ltsv

- [LTSV](http://ltsv.org/) 形式のログを解析します
- デフォルトでは以下のラベルを抽出します
    - `time`
        - ログ時刻
    - `method`
        - HTTP Method
    - `uri`
        - URI
    - `status`
        - HTTP Status Code
    - `apptime`
        - Upstream server からのレスポンスタイム
    - `reqtime`
        - リクエスト処理時間(リクエストを受けてからレスポンスを返すまでの時間)
- `--xxx-label` オプションで、任意のラベル名に変更することができます 

```console
$ cat example/logs/ltsv_access.log
time:2015-09-06T05:58:05+09:00  method:POST     uri:/foo/bar?token=xxx&uuid=1234        status:200      size:12 apptime:0.057
time:2015-09-06T05:58:41+09:00  method:POST     uri:/foo/bar?token=yyy  status:200      size:34 apptime:0.100
time:2015-09-06T06:00:42+09:00  method:GET      uri:/foo/bar?token=zzz  status:200      size:56 apptime:0.123
time:2015-09-06T06:00:43+09:00  method:GET      uri:/foo/bar    status:400      size:15 apptime:-
time:2015-09-06T05:58:44+09:00  method:POST     uri:/foo/bar?token=yyy  status:200      size:34 apptime:0.234
time:2015-09-06T05:58:44+09:00  method:POST     uri:/hoge/piyo?id=yyy   status:200      size:34 apptime:0.234
time:2015-09-06T05:58:05+09:00  method:POST     uri:/foo/bar?token=xxx&uuid=1234        status:200      size:12 apptime:0.057
time:2015-09-06T05:58:41+09:00  method:POST     uri:/foo/bar?token=yyy  status:200      size:34 apptime:0.100
time:2015-09-06T06:00:42+09:00  method:GET      uri:/foo/bar?token=zzz  status:200      size:56 apptime:0.123
time:2015-09-06T06:00:43+09:00  method:GET      uri:/foo/bar    status:400      size:15 apptime:-
time:2015-09-06T06:00:43+09:00  method:GET      uri:/diary/entry/1234   status:200      size:15 apptime:0.135
time:2015-09-06T06:00:43+09:00  method:GET      uri:/diary/entry/5678   status:200      size:30 apptime:0.432
time:2015-09-06T06:00:43+09:00  method:GET      uri:/foo/bar/5xx        status:504      size:15 apptime:60.000
time:2015-09-06T06:00:43+09:00  method:GET      uri:/req        status:200      size:15 apptime:-       reqtime:0.321

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

- 1 行に 1 つの JSON が書かれているログを解析します
- デフォルトでは以下のキーを抽出します
    - `time`
        - ログ時刻
    - `method`
        - HTTP Method
    - `uri`
        - URI
    - `status`
        - HTTP Status Code
    - `response_time`
        - Upstream server からのレスポンスタイム
    - `request_time`
        - リクエスト処理時間(リクエストを受けてからレスポンスを返すまでの時間)
- `--xxx-key` オプションで、任意のキー名に変更することができます

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
{"time":"2015-09-06T06:00:43+09:00","method":"GET","uri":"/req","status":200,"body_bytes":15,"response_time":"-", "request_time":0.321}

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

- 正規表現にマッチするログを解析します
- デフォルトでは Apache combined log + ` "response time"` のログを、以下の名前付きキャプチャで抽出します
    - `time`
        - ログ時刻
    - `method`
        - HTTP Method
    - `uri`
        - URI
    - `status`
        - HTTP Status Code
    - `response_time`
        - Upstream server からのレスポンスタイム
    - `request_time`
        - リクエスト処理時間(リクエストを受けてからレスポンスを返すまでの時間)
- 以下の正規表現
    ```regexp
    ^(\S+)\s\S+\s+(\S+\s+)+\[(?P<time>[^]]+)\]\s"(?P<method>\S*)\s?(?P<uri>(?:[^"]*(?:\\")?)*)\s([^"]*)"\s(?P<status>\S+)\s(?P<body_bytes>\S+)\s"((?:[^"]*(?:\\")?)*)"\s"(?:.+)"\s(?P<response_time>\S+)(?:\s(?P<request_time>\S+))?$
    ```
- `--xxx-subexp` オプションで、任意の名前付きキャプチャに変更することができます

```console
$ cat example/logs/combined_access.log
127.0.0.1 - - [06/Sep/2015:05:58:05 +0900] "POST /foo/bar?token=xxx&uuid=1234 HTTP/1.1" 200 12 "-" "curl/7.54.0" "-" 0.057
127.0.0.1 - - [06/Sep/2015:05:58:41 +0900] "POST /foo/bar?token=yyy HTTP/1.1" 200 34 "-" "curl/7.54.0" "-" 0.100
127.0.0.1 - - [06/Sep/2015:06:00:42 +0900] "GET /foo/bar?token=zzz HTTP/1.1" 200 56 "-" "curl/7.54.0" "-" 0.123
127.0.0.1 - - [06/Sep/2015:06:00:43 +0900] "GET /foo/bar HTTP/1.1" 400 15 "-" "curl/7.54.0" "-" -
127.0.0.1 - - [06/Sep/2015:05:58:44 +0900] "POST /foo/bar?token=yyy HTTP/1.1" 200 34 "-" "curl/7.54.0" "-" 0.234
127.0.0.1 - - [06/Sep/2015:05:58:44 +0900] "POST /hoge/piyo?id=yyy HTTP/1.1" 200 34 "-" "curl/7.54.0" "-" 0.234
127.0.0.1 - - [06/Sep/2015:05:58:05 +0900] "POST /foo/bar?token=xxx&uuid=1234 HTTP/1.1" 200 12 "-" "curl/7.54.0" "-" 0.057
127.0.0.1 - - [06/Sep/2015:05:58:41 +0900] "POST /foo/bar?token=yyy HTTP/1.1" 200 34 "-" "curl/7.54.0" "-" 0.100
127.0.0.1 - - [06/Sep/2015:06:00:42 +0900] "GET /foo/bar?token=zzz HTTP/1.1" 200 56 "-" "curl/7.54.0" "-" 0.123
127.0.0.1 - - [06/Sep/2015:06:00:43 +0900] "GET /foo/bar HTTP/1.1" 400 15 "-" "curl/7.54.0" "-" -
127.0.0.1 - - [06/Sep/2015:06:00:43 +0900] "GET /diary/entry/1234 HTTP/1.1" 200 15 "-" "curl/7.54.0" "-" 0.135
127.0.0.1 - - [06/Sep/2015:06:00:43 +0900] "GET /diary/entry/5678 HTTP/1.1" 200 30 "-" "curl/7.54.0" "-" 0.432
127.0.0.1 - - [06/Sep/2015:06:00:43 +0900] "GET /foo/bar/5xx HTTP/1.1" 504 15 "-" "curl/7.54.0" "-" 60.000
127.0.0.1 - - [06/Sep/2015:06:00:43 +0900] "GET /req HTTP/1.1" 200 15 "-" "curl/7.54.0" "-" - 0.321

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

- pcap形式のファイルから生のHTTPのプロトコルとそのパケットのタイムスタンプをもとにその統計を分析します
  - パケットのタイムスタンプの差をレスポンスタイムを見なすために、実態と誤差が出る場合がある点に注意してください
  - サーバーへのリクエスト/レスポンスを区別するためにサーバーのIPアドレスとTCPポート番号が必要です
- `--pcap-server-ip` オプションでサーバーのIPアドレスを指定できます
  - このオプションは複数個指定することができます
  - デフォルトではローカルのネットワークインターフェースから自動で抽出します
  - ただし、ネットワークインターフェース情報の取得権限が制限されている環境下では `127.0.0.1` と `::1` がデフォルトになります
- `--pcap-server-port` オプションでサーバーのTCPポート番号を指定できます
  - デフォルトでは80になっています
- `--pos` オプションとの併用はできません

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

- 2つの解析結果のダンプファイルを比較します
- `+` はリクエスト数とボディサイズの増加、レスポンスタイムが遅くなったことを意味します
- `-` はリクエスト数とボディサイズの減少、レスポンスタイムが速くなったことを意味します

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

## グローバルオプション

sample は [Usage samples](./docs/usage_samples.ja.md) を参照してください。

- `-c, --config`
    - 各種オプションの設定ファイル
    - YAML
- `--file=FILE` 
    - 解析するファイルのパス
- `-d, --dump=DUMP`
    - 解析結果をファイルに書き出す際のファイルパス
- `-l, --load=LOAD`
    - `-d, --dump` オプションで書き出した解析結果を読み込む際のファイルパス
    - 同じ解析結果に対して、`--sort` や `--reverse` のオプションを変更したい場合に高速に動作することが期待できます
- `--sort=count`
    - 解析結果を表示する際にソートする条件
    - 昇順でソートする 
    - `max`, `min`, `sum`, `avg`
    - `max-body`, `min-body`, `sum-body`, `avg-body`  
    - `p90`, `p95`, `p99`, `stddev`
    - `uri`
    - `method`
    - `count`
    - デフォルトは `count`
    - `p90`, `p95`, `p99` は `--percentiles` で指定したパーセンタイル値によって変更されます
- `-r, --reverse`
    - `--sort` オプションのソート結果を降順にします
- `-q, --query-string`
    - Query String までを含めた URI を集計対象にする
- `--qs-ignore-values`
    - Query String の値を無視して集計します
    - `-q, --query-string` を指定しないと有効になりません
- `--decode-uri`
    - 解析結果の URI をデコードして表示します
- `--format=table`
    - 解析結果を テーブル、Markdown, TSV, CSV, HTML 形式で出力する
    - デフォルトはテーブル形式
- `--noheaders`
    - 解析結果を TSV, CSV で出力する際、header を表示しない
- `--show-footers`
    - 解析結果を テーブル, Markdown で出力する際、footer として 1xx ~ 5xx の合計数を表示する
- `--limit=5000`
    - 解析結果の表示上限数
    - 解析結果の表示数が想定より多かった場合でも、リソースを使いすぎないための設定です
    - デフォルトは 5000 行
- `--location="Local"`
    - フィルタ条件で指定する時刻の timezone
    - デフォルトは localhost に設定されている timezone
- `-o, --output="all"`
    - 出力する解析結果をカンマ区切りで指定する
    - `count`,`1xx`, `2xx`, `3xx`, `4xx`, `5xx`, `method`, `uri`, `min`, `max`, `sum`, `avg`, `p90`, `p95`, `p99`, `stddev`, `min_body`, `max_body`, `sum_body`, `avg_body`
        - `p90`, `p95`, `p99` は `--percentiles` で指定したパーセンタイル値によって変更されます
    - デフォルトはすべて出力(`all`)
- `-m, --matching-groups=PATTERN,...`
    - 正規表現にマッチした URI を同じ集計対象として扱います
    - 指定した順序で正規表現を評価します。マッチした場合、それ以降の正規表現を評価しません。
    - 後述の [URI matching groups](#uri-matching-groups) 参照
- `-f, --filters=FILTERS`
    - 集計対象をフィルタします
    - 後述の[フィルタ](#フィルタ)参照
- `--pos=POSITION_FILE`
    - ファイルをどこまで読み込んだかバイト数を記録します
    - POSITION_FILE にバイト数が書かれていた場合、そのバイト数以降のデータが解析対象になります
    - ファイルを truncate することなく前回解析後からの増分だけを解析することができます
        - また、ファイルを Seek して読み飛ばすので、高速に動作することが見込めます
- `--nosave-pos`
    - `--pos` で指定したバイト数以降のデータを解析対象としますが、読み込んだバイト数の記録はしないようにします
- `--percentiles`
    - 出力するパーセンタイル値をカンマ区切りで指定します
    - デフォルトは `90,95,99`
    
## URI matching groups

以下の `/diary/entry/1234` や `/diary/entry/5678` のように、同一のルーティングでパラメータが異なる URI を単純に集計すると、パラメータごとに集計されますが、ルーティングごとに集計したい場合もあるでしょう。

```console
$ cat example/logs/ltsv_access.log | alp ltsv --filters "Uri matches '^/diary/entry'"
+-------+-----+-----+-----+-----+-----+--------+-------------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN  |  MAX  |  SUM  |  AVG  |  P90  |  P95  |  P99  | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 | 0.135 | 0.135 | 0.135 | 0.135 | 0.135 | 0.135 | 0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 | 0.432 | 0.432 | 0.432 | 0.432 | 0.432 | 0.432 | 0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
+-------+-----+-----+-----+-----+-----+--------+-------------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
```

そのようなケースで、正規表現にマッチした URI を同じ集計対象とするオプションが `-m, --matching-groups=PATTERN,...` です。
カンマ区切りで複数指定することもできます。

```console
$ cat example/logs/ltsv_access.log | alp ltsv --filters "Uri matches '^/diary/entry'" -m "/diary/entry/.+"
+-------+-----+-----+-----+-----+-----+--------+-----------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |       URI       |  MIN  |  MAX  |  SUM  |  AVG  |  P90  |  P95  |  P99  | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-----------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /diary/entry/.+ | 0.135 | 0.432 | 0.567 | 0.283 | 0.432 | 0.432 | 0.432 |  0.148 |    15.000 |    30.000 |    45.000 |    22.500 |
+-------+-----+-----+-----+-----+-----+--------+-----------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
```

## フィルタ

集計対象を条件に応じて包含、除外する機能です。

### 変数

以下の変数に対してフィルタをかけることができます。

- `Uri`
    - URI
- `Method`
    - HTTP メソッド
- `Time`
    - 時刻文字列
    - https://github.com/tkuchiki/parsetime でパースできる時刻文字列に対応しています
- `ResponseTime`
    - レスポンスタイム
- `BodyBytes`
    - HTTP Body のバイト数
- `Status`
    - HTTP Status Code

### 演算子

以下の演算子を使用できます。

- `+`, `-`, `*`, `/`, `%`, `**(べき乗)` 
- `==`, `!=`, `<`, `>`, `<=`, `>=`
- `not`, `!`
- `and`, `&&`
- `or`, `||`
- `matches`
    - 正規表現(`PATTERN`)にマッチするか否か
    - e.g.
       - `Uri matches "PATTERN"`
       - `not(Uri matches "PATTERN")`
- `contains`
    - 文字列(`STRING`)を含むか否か
    - e.g.
        - `Uri contains "STRING"`
        - `not(Uri contains "STRING")`
- `startsWith`
    - 文字列に前方一致するか否か
    - e.g.
        - `Uri startsWith "PREFIX"`
        - `not(Uri startsWith "PREFIX")`
- `endsWith`
    - 文字列に後方一致するか否か
    - e.g.
        - `Uri endsWith "SUFFIX"`
        - `not(Uri endsWith "SUFFIX")`
- `in`
    - 配列の値を含むか否か
    - e.g.
        - `Method in ["GET", "POST"]`
        - `Method not in ["GET", "POST"]`

詳細は https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md を参照してください。  

### 関数

- `TimeAgo(duration)`
    - 現在時刻 - `duration` した時刻を返します
    - Go の time.Duration で使用できる時刻の単位を指定できます
    - `ns`, `us or µs`, `ms`, `s`, `m`, `h`
    - e.g.
        - `Time >= TimeAgo("5m")`
        - `Time` が現在時刻 -5分以上のログを集計対象とする
- `BetweenTime(val, start, end)`
    - SQL の `BETWEEN` のように、`start <= val && val <= end` の結果を返す   
    - e.g.
        - `BetweenTime(Time, "2019-08-06T00:00:00", "2019-08-06T00:05:00")`
        
## 利用例

[Usage samples](./docs/usage_samples.ja.md) を参照してください。

## 寄付

寄付はいつでも歓迎します！    
[:heart: Sponsor](https://github.com/sponsors/tkuchiki)
