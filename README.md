# tcplatency [![Build Status](https://drone.cbrgm.net/api/badges/cbrgm/tcplatency/status.svg)](https://drone.cbrgm.net/cbrgm/tcplatency)

[![](https://img.shields.io/badge/license-Apache--2.0-blue.svg)](https://github.com/cbrgm/tcplatency/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/cbrgm/tcplatency)](https://goreportcard.com/report/github.com/cbrgm/tcplatency)
[![](https://img.shields.io/github/release/cbrgm/tcplatency.svg)](https://github.com/cbrgm/tcplatency/releases)

tcplatency provides an easy way to measure network latency using TCP.

tcplatency was created out of necessity to perform network diagnostics and troubleshooting on serverless infrastructure (but it can also be used in any other environment). Normal pinging is often not possible because many cloud providers do not offer ICMP support. tcplatency solves this problem and allows latency diagnostics in environments where pinging is not possible.
## Features
- Runs as a command line tool or can be used as a library in other projects
- Custom parameters for a port, runs, timeout and wait time between runs
- IPv4 (e.g 192.168.178.22) and dns (e.g google.com) host support
- Small and extensible

## Usage
`tcplatency` can be used both as a module and as a standalone script.

### Command-Line Tool

 Download the latest release here. Or build from source using `make build` (requires go 1.11+)
 
 ```
 AME:
    tcplatency - tcplatency measures network latencies using tcp pings
 
 USAGE:
    tcplatency [global options] command [command options] [arguments...]
 
 COMMANDS:
    help, h  Shows a list of commands or help for one command
 
 GLOBAL OPTIONS:
    --host value               the host address
    --port value, -p value     the host port (default: 443)
    --timeout value, -t value  timeout in seconds (default: 5)
    --runs value, -r value     number of latency points to return (default: 5)
    --wait value, -w value     seconds to wait between each run (default: 1)
    --help, -h                 show help
    --version, -v              print the version
 ```

Example: `$ tcplatency google.com`

```bash
google.com via tcp seq=0 port=443 timeout=5 time=20.63 ms 
google.com via tcp seq=1 port=443 timeout=5 time=14.10 ms 
google.com via tcp seq=2 port=443 timeout=5 time=8.99 ms 
google.com via tcp seq=3 port=443 timeout=5 time=8.41 ms 
google.com via tcp seq=4 port=443 timeout=5 time=8.57 ms 
--- google.com tcplatency statistics --- 
5 packets transmitted, 5 successful, 0 failed  
min/avg/max/mdev = 8.41/12.14/20.63/4.74 ms
````

Example: `$ tcplatency --port 80 --runs 3 --wait 1 52.26.14.11`

```bash
52.26.14.11 via tcp seq=0 port=80 timeout=5 time=224.45 ms 
52.26.14.11 via tcp seq=1 port=80 timeout=5 time=166.37 ms 
52.26.14.11 via tcp seq=2 port=80 timeout=5 time=187.80 ms 
--- 52.26.14.11 tcplatency statistics --- 
3 packets transmitted, 3 successful, 0 failed  
min/avg/max/mdev = 166.37/192.87/224.45/23.98 ms 
```

### Go Module

Use tcplatency as a library

```go
"github.com/cbrgm/tcplatency/latency"
```

Example: 
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

## Credit & License

tcplatency is open-source and is developed under the terms of the [Apache 2.0 License](https://github.com/cbrgm/tcplatency/blob/master/LICENSE).

Maintainer of this repository is:

-   [@cbrgm](https://github.com/cbrgm) | Christian Bargmann <mailto:chris@cbrgm.net>

Please refer to the git commit log for a complete list of contributors.

## Contributing

See the [Contributing Guide](https://github.com/cbrgm/contributing/blob/master/CONTRIBUTING.md).