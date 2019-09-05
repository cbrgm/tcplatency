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
	"github.com/urfave/cli"
	"os"
	"time"
)

const (
	FlagHost    = "host"
	FlagPort    = "port,p"
	FlagTimeout = "timeout,t"
	FlagRuns    = "runs,r"
	FlagWait    = "wait,w"
)

type appConfig struct {
	Host    string
	Port    int
	Timeout int
	Runs    int
	Wait    int
}

var (
	conf = appConfig{}

	clientFlags = []cli.Flag{
		cli.StringFlag{
			Name:        FlagHost,
			Usage:       "the host address",
			Destination: &conf.Host,
		},
		cli.IntFlag{
			Name:        FlagPort,
			Usage:       "the host port",
			Value:       443,
			Destination: &conf.Port,
		},
		cli.IntFlag{
			Name:        FlagTimeout,
			Usage:       "timeout in seconds",
			Value:       5,
			Destination: &conf.Timeout,
		},
		cli.IntFlag{
			Name:        FlagRuns,
			Usage:       "number of latency points to return",
			Value:       5,
			Destination: &conf.Runs,
		},
		cli.IntFlag{
			Name:        FlagWait,
			Usage:       "seconds to wait between each run",
			Value:       1,
			Destination: &conf.Wait,
		},
	}
)

var (
	// Version of tcplatency
	Version string
	// Revision or Commit this binary was built from
	Revision string
	// BuildDate this binary was built
	BuildDate string
)

func main() {
	app := cli.NewApp()
	app.Name = "tcplatency"
	app.Usage = "tcplatency measures network latencies using tcp pings"
	app.Version = fmt.Sprintf("version: %s, revision: %s (%s)", Version, Revision, BuildDate)
	app.Author = "Christian Bargmann"
	app.Email = "chris@cbrgm.net"
	app.Action = appAction
	app.Flags = clientFlags

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("failed to run cli: %s", err)
		os.Exit(1)
	}
}

func appAction(c *cli.Context) error {

	// preflight checks
	if len(c.Args()) > 1 {
		fmt.Println("Incorrect number of arguments. Usage: `tcplatency <host>`, see tcplatency -h")
		return nil
	}

	if len(c.Args()) == 1 {
		conf.Host = c.Args()[0]
	}

	if conf.Host == "" {
		fmt.Println("Please provide a host either as argument or via --host flag, see tcplatency -h")
		return nil
	}

	var m = latency.NewMeasurement(conf.Host, conf.Port, conf.Timeout, conf.Runs, conf.Wait)

	for !m.IsFinished() {
		time.Sleep(time.Second * time.Duration(m.Wait))
		tcp := m.DialTCP()

		if tcp.Failed {
			fmt.Printf("%s via tcp seq=%d port=%d timeout=%d failed \n", tcp.Host, tcp.Sequence, tcp.Port, tcp.Timeout)
			continue
		}
		fmt.Printf("%s via tcp seq=%d port=%d timeout=%d time=%.2f ms \n", tcp.Host, tcp.Sequence, tcp.Port, tcp.Timeout, tcp.Latency)
	}

	printSummary(m.Host, m.Result())

	return nil
}

func printSummary(host string, result latency.MeasurementResult) {
	fmt.Printf("--- %s tcplatency statistics --- \n", host)
	fmt.Printf("%d packets transmitted, %d successful, %d failed  \n", result.Count, result.Successful, result.Failed)
	fmt.Printf("min/avg/max/mdev = %.2f/%.2f/%.2f/%.2f ms \n", result.Min, result.Average, result.Max, result.StdDev)
}
