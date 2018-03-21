gologz helps you searching hits on Logz.io platform by simple or aggregated map query's

## How to install:

```bash
go get github.com/stacknowledge/gologz
```

## How to use:

### Simple string tracing
```golang
package main

import (
	"fmt"

	"github.com/stacknowledge/gologz"
)

func main() {
    client := gologz.New("YOUR-LOGZ.IO-API-TOKEN", "eu") // Accepts EU or US regions
    
    // query, start index, limit, you can also compose a query with AND/OR selector
    // like when_ts:1521504127731 AND commit_hash:70793d6aee43b22ba08c5b95
    traces, err := client.Trace("when_ts:1521504127731", 0, 5)  

    if err != nil {
        fmt.Errorf("%s", err)
    }

    for index := range traces.Logs {
		fmt.Printf("\n When:%v \n Service:%v \n Log:%v \n", traces.Logs[index].When, traces.Logs[index].Service, traces.Logs[index].Log)
	}
}
```

### Aggregate map tracing
```golang
package main

import (
	"fmt"

	"github.com/stacknowledge/gologz"
)

func main() {
    client := gologz.New("YOUR-LOGZ.IO-API-TOKEN", "eu")            // Accepts EU or US regions
    traces, _ = client.TracePairs(map[string]interface{}{
		"when_ts":     1521504127731,
		"commit_hash": "70793d6aee43b22ba08c5b956761bba93b85fa75",
	}, 0, 5)                                                        // Aggregate map, start index, limit

    if err != nil {
        fmt.Errorf("%s", err)
    }

    for index := range traces.Logs {
		fmt.Printf("\n When:%v \n Service:%v \n Log:%v \n", traces.Logs[index].When, traces.Logs[index].Service, traces.Logs[index].Log)
	}
}
```
