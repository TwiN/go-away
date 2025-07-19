package goaway_test

import (
	"bytes"
	"testing"

	goaway "github.com/TwiN/go-away"
)

func TestWriter(t *testing.T) {
	tests := map[string]struct {
		input          [][]byte
		detector       *goaway.ProfanityDetector
		expectedOutput string
	}{
		"no writing, empty output": {
			input:          [][]byte{},
			detector:       goaway.NewProfanityDetector(),
			expectedOutput: "",
		},
		"single uncensored write": {
			input: [][]byte{
				[]byte("I'm just a normal line"),
			},
			detector:       goaway.NewProfanityDetector(),
			expectedOutput: "I'm just a normal line",
		},
		"single censored write": {
			input: [][]byte{
				[]byte("I'm just a shitty line"),
			},
			detector:       goaway.NewProfanityDetector(),
			expectedOutput: "I'm just a ****ty line",
		},
		"multi-line single write": {
			input: [][]byte{
				[]byte("I'm just a shitty line\nAnd I'm another line"),
			},
			detector:       goaway.NewProfanityDetector(),
			expectedOutput: "I'm just a ****ty line\nAnd I'm another line",
		},
		"single-line multi writes": {
			input: [][]byte{
				[]byte("I'm just a shitty line\n"),
				[]byte("And I'm another line"),
				[]byte("\nAnd I'm the final fucking line"),
			},
			detector:       goaway.NewProfanityDetector(),
			expectedOutput: "I'm just a ****ty line\nAnd I'm another line\nAnd I'm the final ****ing line",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			writer := goaway.NewWriter(buf, tc.detector)
			for _, write := range tc.input {
				n, err := writer.Write(write)
				if n != len(write) {
					t.Errorf("unexpected write count %d != %d", n, len(write))
				}

				if err != nil {
					t.Errorf("unexpected writing error %v", err)
				}
			}

			err := writer.Flush()
			if err != nil {
				t.Errorf("unexpected error flushing writer %v", err)
			}

			result := buf.String()
			if tc.expectedOutput != result {
				t.Errorf("expected %q but recieved %q", tc.expectedOutput, result)
			}
		})
	}
}
