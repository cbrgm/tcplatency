/*
 * Copyright 2019, Christian Bargmann <chris@cbrgm.net>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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
