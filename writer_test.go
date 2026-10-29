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
		"multi writes without newline": {
			input: [][]byte{
				[]byte("Hello "),
				[]byte("shitty "),
				[]byte("world"),
			},
			detector:       goaway.NewProfanityDetector(),
			expectedOutput: "Hello ****ty world",
		},
		"nil detector fallback": {
			input: [][]byte{
				[]byte("This is shitty"),
			},
			detector:       nil,
			expectedOutput: "This is ****ty",
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

func TestWriter_BufferReuse(t *testing.T) {
	buf := &bytes.Buffer{}
	writer := goaway.NewWriter(buf, nil)

	reusable := []byte("hello ")
	_, err := writer.Write(reusable)
	if err != nil {
		t.Fatalf("unexpected write error: %v", err)
	}

	// Mutate reusable buffer to ensure Writer did not retain caller slice
	copy(reusable, "XXXXX ")

	if err := writer.Flush(); err != nil {
		t.Fatalf("unexpected flush error: %v", err)
	}

	if got := buf.String(); got != "hello " {
		t.Errorf("expected buffer not to be corrupted by caller mutation, got %q, want %q", got, "hello ")
	}
}
