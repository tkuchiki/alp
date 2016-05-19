# go-timezone
Timezone utility for Golang

## Example

### Code

```
package main

import (
	"fmt"
	"github.com/tkuchiki/go-timezone"
    "time"
)

func main() {
	offset, err := timezone.GetOffset("JST")
	fmt.Println(offset, err)

	offset, err = timezone.GetOffset("hogehoge")
	fmt.Println(offset, err)

	var zones []string
	zones, err = timezone.GetTimezones("UTC")
	fmt.Println(zones, err)

	zones, err = timezone.GetTimezones("foobar")
	fmt.Println(zones, err)

	now := time.Now()

	fmt.Println("## current timezone")
	fmt.Println(now)

	var jst time.Time
	loc, _ := time.LoadLocation("UTC")
	utc := now.In(loc)

	jst, _ = timezone.FixedTimezone(utc, "")

	fmt.Println("## UTC")
	fmt.Println(utc)
	fmt.Println("## UTC -> JST (current timezone)")
	fmt.Println(jst)

	var est time.Time
	est, _ = timezone.FixedTimezone(now, "America/New_York")
	fmt.Println("## JST -> EST")
	fmt.Println(est)
}
```

### Result

```
32400 <nil>
0 Invalid short timezone: hogehoge
[Antarctica/Troll Etc/UTC Etc/Universal Etc/Zulu UTC Universal Zulu] <nil>
[] Invalid short timezone: foobar
## current timezone
2016-03-02 14:33:49.078798783 +0900 JST
## UTC
2016-03-02 05:33:49.078798783 +0000 UTC
## UTC -> JST (current timezone)
2016-03-02 14:33:49.078798783 +0900 JST
## JST -> EST
2016-03-02 00:33:49.078798783 -0500 EST
```
