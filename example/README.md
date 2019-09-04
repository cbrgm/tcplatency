# Go Module

Use tcplatency as a library

```
"github.com/cbrgm/tcplatency/latency"
```

Example (see main.go): 
```go
package main

import (
	"fmt"
	"github.com/cbrgm/tcplatency/latency"
)

func main() {
	var host = "google"
	var port = 443
	var timeout = 5
	var runs = 5
	var wait = 1

	result := latency.Measure(host, port, timeout, runs, wait)
	printSummary(host, result)

}

func printSummary(host string, result latency.MeasurementResult) {
	fmt.Printf("--- %s tcplatency statistics --- \n", host)
	fmt.Printf("%d packets transmitted, %d successful, %d failed  \n", result.Count, result.Successful, result.Failed)
	fmt.Printf("min/avg/max/mdev = %.2f/%.2f/%.2f/%.2f ms \n", result.Min, result.Average, result.Max, result.StdDev)
}

```