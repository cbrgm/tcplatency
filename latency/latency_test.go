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
	"testing"
)

func TestNewMeasurement(t *testing.T) {
	want := &Measurement{
		Host:    "google.de",
		Port:    443,
		Timeout: 5,
		Runs:    5,
		Wait:    1,
	}

	got := NewMeasurement("google.de", 443, 5, 5, 1)
	if want.Host != got.Host {
		t.Fatalf("want: %s, got %s", want.Host, got.Host)
	}
	if want.Port != got.Port {
		t.Fatalf("want: %d, got %d", want.Port, got.Port)
	}
	if want.Timeout != got.Timeout {
		t.Fatalf("want: %d, got %d", want.Timeout, got.Timeout)
	}
	if want.Runs != got.Runs {
		t.Fatalf("want: %d, got %d", want.Runs, got.Runs)
	}
	if want.Wait != got.Wait {
		t.Fatalf("want: %d, got %d", want.Wait, got.Wait)
	}
}

func TestMeasurement_DialTCP(t *testing.T) {
	m := NewMeasurement("google.de", 443, 5, 5, 1)

	tcp := m.DialTCP()
	if tcp.Sequence != 0 {
		t.Fatalf("want: %d, got: %d", 0, tcp.Sequence)
	}

	tcp = m.DialTCP()
	if tcp.Sequence != 1 {
		t.Fatalf("want: %d, got: %d", 1, tcp.Sequence)
	}

}

func TestMeasurement_IsFinished(t *testing.T) {
	got := NewMeasurement("google.de", 443, 5, 2, 0)

	_ = got.DialTCP()
	if got.IsFinished() {
		t.Fatalf("want: false, got %t", got.IsFinished())
	}

	_ = got.DialTCP()
	if !got.IsFinished() {
		t.Fatalf("want: true, got %t", got.IsFinished())
	}

}

func TestMeasurement_AvgLatency(t *testing.T) {
	m := NewMeasurementWithTestData()

	got := m.AvgLatency()
	want := 23.00

	if got != want {
		t.Fatalf("got %.2f, want %.2f", got, want)
	}

}

func TestMeasurement_MaxLatency(t *testing.T) {
	m := NewMeasurementWithTestData()

	got := m.MaxLatency()
	want := 24.50

	if got != want {
		t.Fatalf("got %.2f, want %.2f", got, want)
	}
}

func TestMeasurement_MinLatency(t *testing.T) {
	m := NewMeasurementWithTestData()

	got := m.MinLatency()
	want := 21.50

	if got != want {
		t.Fatalf("got %.2f, want %.2f", got, want)
	}
}

func TestMeasurement_StdDevLatency(t *testing.T) {
	m := NewMeasurementWithTestData()

	got := m.StdDevLatency()
	want := 10.46

	if got == want {
		t.Fatalf("got %.2f, want %.2f", got, want)
	}
}

func TestAppendSuccessful(t *testing.T) {
	m := NewMeasurement("google.de", 443, 5, 3, 1)

	m.appendSuccessful(&TCPResponse{
		Host:     m.Host,
		Port:     m.Port,
		Latency:  21.50,
		Timeout:  m.Timeout,
		Sequence: 0,
	})

	if m.count != 1 {
		t.Fatalf("want: %d, got %d", 1, m.count)
	}

	if m.successful != 1 {
		t.Fatalf("want: %d, got %d", 1, m.successful)
	}

	if m.failed != 0 {
		t.Fatalf("want: %d, got %d", 0, m.failed)
	}

	if len(m.data) != 1 {
		t.Fatalf("want: %d, got %d", 1, len(m.data))
	}
}

func TestAppendFailed(t *testing.T) {
	m := NewMeasurement("google.de", 443, 5, 3, 1)

	m.appendFailed(&TCPResponse{
		Host:     m.Host,
		Port:     m.Port,
		Latency:  21.50,
		Timeout:  m.Timeout,
		Sequence: 0,
	})

	if m.count != 1 {
		t.Fatalf("want: %d, got %d", 1, m.count)
	}

	if m.successful != 0 {
		t.Fatalf("want: %d, got %d", 0, m.successful)
	}

	if m.failed != 1 {
		t.Fatalf("want: %d, got %d", 1, m.failed)
	}

	if len(m.data) != 1 {
		t.Fatalf("want: %d, got %d", 1, len(m.data))
	}

}

func NewMeasurementWithTestData() *Measurement {
	m := NewMeasurement("google.de", 443, 5, 3, 1)

	m.appendSuccessful(&TCPResponse{
		Host:     m.Host,
		Port:     m.Port,
		Latency:  21.50,
		Timeout:  m.Timeout,
		Sequence: 0,
	})
	m.appendSuccessful(&TCPResponse{
		Host:     m.Host,
		Port:     m.Port,
		Latency:  24.50,
		Timeout:  m.Timeout,
		Sequence: 1,
	})
	m.appendSuccessful(&TCPResponse{
		Host:     m.Host,
		Port:     m.Port,
		Latency:  23.00,
		Timeout:  m.Timeout,
		Sequence: 2,
	})
	m.appendFailed(
		&TCPResponse{
			Host:     m.Host,
			Port:     m.Port,
			Latency:  5.00,
			Timeout:  m.Timeout,
			Sequence: 2,
			Failed:   true,
		},
	)

	return m
}
