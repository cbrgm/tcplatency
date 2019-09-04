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

package latency

import (
	"fmt"
	"math"
	"net"
	"time"
)

// Measurements represents a new tcplatency measurement
type Measurement struct {
	Host    string
	Port    int
	Timeout int
	Runs    int
	Wait    int

	// internal vars
	data       []*TCPResponse
	count      int
	failed     int
	successful int
}

// MeasurementResult represents a new measurement result
type MeasurementResult struct {
	Successful int
	Failed     int
	Count      int
	Average    float64
	Max        float64
	Min        float64
	StdDev     float64
}

func NewMeasurement(host string, port int, timeout int, runs int, wait int) *Measurement {
	return &Measurement{
		Host:    host,
		Port:    port,
		Timeout: timeout,
		Runs:    runs,
		Wait:    wait,

		data:       []*TCPResponse{},
		count:      0,
		failed:     0,
		successful: 0,
	}
}

func Measure(host string, port int, timeout int, runs int, wait int) MeasurementResult {
	var m = NewMeasurement(host, port, timeout, runs, wait)
	for !m.IsFinished() {
		time.Sleep(time.Second * time.Duration(m.Wait))
		_ = m.DialTCP()
	}
	return m.Result()
}

// DialTCP calculates a latency point using sockets.
func (m *Measurement) DialTCP() *TCPResponse {

	// Simply check that the server is up and can  accept connections.
	result := &TCPResponse{
		Host:     m.Host,
		Port:     m.Port,
		Latency:  0,
		Timeout:  m.Timeout,
		Sequence: m.count,
	}

	address := fmt.Sprintf("%s:%d", m.Host, m.Port)
	start := time.Now()

	conn, err := net.DialTimeout("tcp", address, time.Duration(m.Timeout)*time.Second)
	if err != nil {

		m.appendFailed(result)
		return result
	}
	defer conn.Close()

	result.Latency = time.Since(start).Seconds() * 1000

	m.appendSuccessful(result)
	return result
}

// IsFinished returns true when all tcp ping runs have been executed, false if not
func (m *Measurement) IsFinished() bool {
	if m.count >= m.Runs {
		return true
	}
	return false
}

func (m *Measurement) appendFailed(response *TCPResponse) {

	response.Failed = true

	m.count++
	m.failed++
	m.data = append(m.data, response)
}

func (m *Measurement) appendSuccessful(response *TCPResponse) {

	response.Failed = false

	m.count++
	m.successful++
	m.data = append(m.data, response)
}

// AvgLatency returns the average tcp latency
func (m *Measurement) AvgLatency() float64 {
	var count = 0
	var avg = 0.00

	for _, item := range m.data {
		if item.Failed {
			continue
		}
		count++
		avg += item.Latency
	}
	return avg / float64(count)
}

// MaxLatency returns the max latency
func (m *Measurement) MaxLatency() float64 {
	var max float64
	for i, e := range m.data {
		if e.Failed {
			continue
		}
		if i == 0 || max < e.Latency {
			max = e.Latency
		}
	}
	return max
}

// MinLatency returns the min latency
func (m *Measurement) MinLatency() float64 {
	var min float64
	for i, e := range m.data {
		if e.Failed {
			continue
		}
		if i == 0 || min > e.Latency {
			min = e.Latency
		}
	}
	return min
}

// StdDevLatency calculates an average of how far each latency is from the mean latency.
// The higher MedianDerivation  is, the more variable the latency is over time.
func (m *Measurement) StdDevLatency() float64 {
	avg := m.AvgLatency()
	dev := 0.
	for _, d := range m.data {
		diff := d.Latency - avg
		dev += diff * diff
	}
	dev /= float64(m.successful)
	dev = math.Sqrt(dev)
	return dev
}

// Result returns  a new measurement result.
// Can be used to summarize results of multiple measurements
func (m *Measurement) Result() MeasurementResult {
	return MeasurementResult{
		Successful: m.successful,
		Failed:     m.failed,
		Count:      m.count,
		Average:    m.AvgLatency(),
		Max:        m.MaxLatency(),
		Min:        m.MinLatency(),
		StdDev:     m.StdDevLatency(),
	}
}
