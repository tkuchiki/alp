# alp

alp はアクセスログ解析ツールです。
読み方は、[éɪélpíː] または [ælp] です(作者は [éɪélpíː] 派)。

# インストール

https://github.com/tkuchiki/alp/releases から環境にあったバイナリが含まれる zip ファイルをダウンロードして、解凍してください。

# v0.4.0 と v1.0.0 の違い

[v0.4.0 と v1.0.0 の違い](./docs/how_to_difference_between_v0_4_0_and_v1_0_0.ja.md) を参照してください。

# 使い方

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

- alp は ltsv, json, regexp の3つのサブコマンドで構成されています
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
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN   |  MAX   |  SUM   |  AVG   |   P1   |  P50   |  P99   | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /req              |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | POST   | /hoge/piyo        |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.000 |    34.000 |    34.000 |    34.000 |    34.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
|     1 |   0 |   0 |   0 |   0 |   1 | GET    | /foo/bar/5xx      | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /foo/bar          |  0.123 |  0.123 |  0.246 |  0.123 |  0.123 |  0.123 |  0.123 |  0.000 |    56.000 |    56.000 |   112.000 |    56.000 |
|     5 |   0 |   5 |   0 |   0 |   0 | POST   | /foo/bar          |  0.057 |  0.234 |  0.548 |  0.110 |  0.057 |  0.100 |  0.057 |  0.065 |    12.000 |    34.000 |   126.000 |    25.200 |
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
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN   |  MAX   |  SUM   |  AVG   |   P1   |  P50   |  P99   | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /req              |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | POST   | /hoge/piyo        |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.000 |    34.000 |    34.000 |    34.000 |    34.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
|     1 |   0 |   0 |   0 |   0 |   1 | GET    | /foo/bar/5xx      | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /foo/bar          |  0.123 |  0.123 |  0.246 |  0.123 |  0.123 |  0.123 |  0.123 |  0.000 |    56.000 |    56.000 |   112.000 |    56.000 |
|     5 |   0 |   5 |   0 |   0 |   0 | POST   | /foo/bar          |  0.057 |  0.234 |  0.548 |  0.110 |  0.057 |  0.100 |  0.057 |  0.065 |    12.000 |    34.000 |   126.000 |    25.200 |
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
                                '"rseponse_time":$upstream_response_time,'
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
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN   |  MAX   |  SUM   |  AVG   |   P1   |  P50   |  P99   | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-------------------+--------+--------+--------+--------+--------+--------+--------+--------+-----------+-----------+-----------+-----------+
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /req              |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.321 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | POST   | /hoge/piyo        |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.234 |  0.000 |    34.000 |    34.000 |    34.000 |    34.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/1234 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.135 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     1 |   0 |   1 |   0 |   0 |   0 | GET    | /diary/entry/5678 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.432 |  0.000 |    30.000 |    30.000 |    30.000 |    30.000 |
|     1 |   0 |   0 |   0 |   0 |   1 | GET    | /foo/bar/5xx      | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 | 60.000 |  0.000 |    15.000 |    15.000 |    15.000 |    15.000 |
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /foo/bar          |  0.123 |  0.123 |  0.246 |  0.123 |  0.123 |  0.123 |  0.123 |  0.000 |    56.000 |    56.000 |   112.000 |    56.000 |
|     5 |   0 |   5 |   0 |   0 |   0 | POST   | /foo/bar          |  0.057 |  0.234 |  0.548 |  0.110 |  0.057 |  0.100 |  0.057 |  0.065 |    12.000 |    34.000 |   126.000 |    25.200 |
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
    - `p1`, `p50`, `p99`, `stddev`
    - `uri`
    - `method`
    - `count`
    - デフォルトは `count`
- `-r, --reverse`
    - `--sort` オプションのソート結果を降順にします
- `-q, --query-string`
    - Query String までを含めた URI を集計対象にする
- `--format=table`
    - 解析結果を テーブル、Markdown, TSV, CSV 形式で出力する
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
    - `count`,`1xx`, `2xx`, `3xx`, `4xx`, `5xx`, `method`, `uri`, `min`, `max`, `sum`, `avg`, `p1`, `p50`, `p99`, `stddev`, `min_body`, `max_body`, `sum_body`, `avg_body`
    - デフォルトはすべて出力(`all`)
- `-m, --matching-groups=PATTERN,...`
    - 正規表現にマッチした URI を同じ集計対象として扱います
    - 後述の [URI matching groups](#URI matching groups) 参照
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
    
## URI matching groups

以下の `/diary/entry/1234` や `/diary/entry/5678` のように、同一のルーティングでパラメータが異なる URI を単純に集計すると、パラメータごとに集計されますが、ルーティングごとに集計したい場合もあるでしょう。

```console
$ cat example/logs/ltsv_access.log | alp ltsv --filters "Uri matches '^/diary/entry'"
+-------+-----+-----+-----+-----+-----+--------+-------------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |        URI        |  MIN  |  MAX  |  SUM  |  AVG  |  P1   |  P50  |  P99  | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
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
| COUNT | 1XX | 2XX | 3XX | 4XX | 5XX | METHOD |       URI       |  MIN  |  MAX  |  SUM  |  AVG  |  P1   |  P50  |  P99  | STDDEV | MIN(BODY) | MAX(BODY) | SUM(BODY) | AVG(BODY) |
+-------+-----+-----+-----+-----+-----+--------+-----------------+-------+-------+-------+-------+-------+-------+-------+--------+-----------+-----------+-----------+-----------+
|     2 |   0 |   2 |   0 |   0 |   0 | GET    | /diary/entry/.+ | 0.135 | 0.432 | 0.567 | 0.283 | 0.135 | 0.135 | 0.135 |  0.148 |    15.000 |    30.000 |    45.000 |    22.500 |
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
