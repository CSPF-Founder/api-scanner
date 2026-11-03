package controllers

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/CSPF-Founder/api-scanner/code/scanner/internal/schemas"
	"github.com/CSPF-Founder/api-scanner/code/scanner/logger"
)

func TestGetAuthHeaders(t *testing.T) {
	logFile, err := os.CreateTemp(t.TempDir(), "job.log")
	if err != nil {
		t.Fatalf("error creating log file: %v", err)
	}
	defer func() { _ = logFile.Close() }()
	fileLogger, err := logger.GetFileLogger(logFile)
	if err != nil {
		t.Fatalf("error creating logger: %v", err)
	}
	c := &JobController{logger: fileLogger}

	tests := []struct {
		name    string
		content string
		want    []schemas.AuthHeaderMap
		wantErr bool
	}{
		{
			name:    "no headers",
			content: "[AUTH_HEADERS]\n",
		},
		{
			name:    "single header",
			content: "[AUTH_HEADERS]\nAuthorization = Bearer abc\n",
			want:    []schemas.AuthHeaderMap{{Name: "Authorization", Value: "Bearer abc"}},
		},
		{
			name:    "comment and quote characters are kept",
			content: "[AUTH_HEADERS]\nx-api-key = \"k#1;2 %3\"\n",
			want:    []schemas.AuthHeaderMap{{Name: "x-api-key", Value: `"k#1;2 %3"`}},
		},
		{
			name:    "multiple headers",
			content: "[AUTH_HEADERS]\nAuthorization = a\nx-api-key = b\n",
			wantErr: true,
		},
		{
			name:    "malformed file",
			content: "[AUTH_HEADERS\n",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "auth_headers.conf")
			if err := os.WriteFile(path, []byte(tc.content), 0o600); err != nil {
				t.Fatalf("error writing auth header file: %v", err)
			}

			got, err := c.getAuthHeaders(path)
			if (err != nil) != tc.wantErr {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected headers %v, got %v", tc.want, got)
			}
		})
	}
}
