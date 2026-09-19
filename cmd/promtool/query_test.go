// Copyright The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseTime(t *testing.T) {
	ts, err := time.Parse(time.RFC3339Nano, "2015-06-03T13:21:58.555Z")
	require.NoError(t, err)

	for _, test := range []struct {
		input  string
		fail   bool
		result time.Time
	}{
		{
			input: "",
			fail:  true,
		}, {
			input: "abc",
			fail:  true,
		}, {
			input: "30s",
			fail:  true,
		}, {
			// strconv.ParseFloat accepts NaN, and converting it to int64 is
			// undefined in Go: it yields 0 on arm64 and math.MinInt64 on amd64.
			input: "NaN",
			fail:  true,
		}, {
			input: "Inf",
			fail:  true,
		}, {
			input: "+Inf",
			fail:  true,
		}, {
			input: "-Inf",
			fail:  true,
		}, {
			// Finite, but too large to convert to an int64 number of seconds.
			input: "1e300",
			fail:  true,
		}, {
			input:  "123",
			result: time.Unix(123, 0).UTC(),
		}, {
			input:  "123.123",
			result: time.Unix(123, 123000000).UTC(),
		}, {
			input:  "2015-06-03T13:21:58.555Z",
			result: ts,
		}, {
			input:  "2015-06-03T14:21:58.555+01:00",
			result: ts,
		},
	} {
		t.Run(test.input, func(t *testing.T) {
			result, err := parseTime(test.input)
			if test.fail {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.True(t, test.result.Equal(result), "expected %v, got %v", test.result, result)
		})
	}
}
