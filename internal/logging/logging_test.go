package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewWithWriter(t *testing.T) {
	tests := []struct {
		name       string
		config     Config
		wantOutput string
	}{
		{
			name: "json",
			config: Config{
				Level:  "info",
				Format: "json",
			},
			wantOutput: `"msg":"started"`,
		},
		{
			name: "text",
			config: Config{
				Level:  "debug",
				Format: "text",
			},
			wantOutput: "msg=started",
		},
		{
			name: "defaults",
			config: Config{
				Level:  "",
				Format: "",
			},
			wantOutput: `"msg":"started"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			logger, err := NewWithWriter(tt.config, &buf)
			if err != nil {
				t.Fatalf(
					"NewWithWriter() error = %v, want nil",
					err,
				)
			}

			logger.Info("started")

			if !strings.Contains(buf.String(), tt.wantOutput) {
				t.Errorf(
					"logged output = %q, want it to contain %q", buf.String(), tt.wantOutput,
				)
			}
		})
	}
}

func TestNewWithWriterRejectsInvalidLevel(t *testing.T) {
	var buf bytes.Buffer

	_, err := NewWithWriter(
		Config{
			Level:  "verbose",
			Format: "json",
		},
		&buf,
	)
	if err == nil {
		t.Fatal(
			"NewWithWriter() error = nil, want invalid level error",
		)
	}

	want := `unsupported log level "verbose"`
	if !strings.Contains(err.Error(), want) {
		t.Fatalf(
			"NewWithWriter() error = %q, want it to contain %q", err, want,
		)
	}
}

func TestNewWithWriterRejectsInvalidFormat(t *testing.T) {
	var buf bytes.Buffer

	_, err := NewWithWriter(
		Config{
			Level:  "info",
			Format: "xml",
		},
		&buf,
	)
	if err == nil {
		t.Fatal("NewWithWriter() error = nil, want invalid format error")
	}

	want := `unsupported log format "xml"`
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("NewWithWriter() error = %q, want it to contain %q", err, want)
	}
}
