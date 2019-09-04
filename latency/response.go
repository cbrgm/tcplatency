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

import "fmt"

type TCPResponse struct {
	Host     string
	Port     int
	Sequence int
	Timeout  int
	Latency  float64
	Failed   bool
}

func (tcp *TCPResponse) String() string {
	if tcp.Latency == 0.00 {
		return fmt.Sprintf("%s via tcp seq=%d port=%d timeout=%d Failed \n", tcp.Host, tcp.Sequence, tcp.Port, tcp.Timeout)
	}
	return fmt.Sprintf("%s via tcp seq=%d port=%d timeout=%d time=%.2f ms \n", tcp.Host, tcp.Sequence, tcp.Port, tcp.Timeout, tcp.Latency)
}
