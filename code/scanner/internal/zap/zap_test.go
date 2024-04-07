package zap

import (
	"os"
	"strings"
	"testing"

	"github.com/CSPF-Founder/api-scanner/code/scanner/internal/schemas"
	"github.com/CSPF-Founder/api-scanner/code/scanner/models"
)

// func (z *ZapModule) WriteAuthFile(authHeaderData []schemas.AuthHeaderMap, siteNameForConfig string) error {
// 	if len(authHeaderData) == 0 {
// 		return fmt.Errorf("authHeaderData should not be empty")
// 	} else if len(authHeaderData) > 1 {
// 		return fmt.Errorf("authHeaderData should have only one element")
// 	}

// 	file, err := os.Create(z.authPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to create file: %v", err)
// 	}
// 	defer file.Close()

// 	firstHeader := authHeaderData[0]

// 	_, err = fmt.Fprintf(file, "ZAP_AUTH_HEADER=%s\n", firstHeader.Name)
// 	if err != nil {
// 		return fmt.Errorf("failed to write ZAP_AUTH_HEADER: %v", err)
// 	}

// 	_, err = fmt.Fprintf(file, "ZAP_AUTH_HEADER_VALUE=%s\n", firstHeader.Value)
// 	if err != nil {
// 		return fmt.Errorf("failed to write ZAP_AUTH_HEADER_VALUE: %v", err)
// 	}

// 	_, err = fmt.Fprintf(file, "ZAP_SITE_NAME=%s\n", siteNameForConfig)
// 	if err != nil {
// 		return fmt.Errorf("failed to write ZAP_SITE_NAME: %v", err)
// 	}

// 	return nil
// }

// write test for the above function

func TestWriteAuthFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "api-scanner-test")
	if err != nil {
		t.Errorf("Error creating temp dir: %v", err)
	}

	defer os.RemoveAll(tmpDir)

	outputDir := tmpDir + "/output"
	authPath := tmpDir + "/auth_headers.conf"

	z := &ZapModule{
		job: models.Job{
			ID:     10000030303003033,
			UserID: 1,
		},
		outputDir:    outputDir,
		authPath:     authPath,
		scannerImage: "scannerImage",
	}

	authHeaderData := []schemas.AuthHeaderMap{
		{
			Name:  "Authorization",
			Value: "Bearer DUMMY_TOKEN",
		},
	}

	siteNameForConfig := "example.com"

	err = z.WriteAuthFile(authHeaderData, siteNameForConfig)
	if err != nil {
		t.Errorf("WriteAuthFile() failed to write auth file: %v", err)
	}

	// Read the file
	data, err := os.ReadFile(z.authPath)
	if err != nil {
		t.Errorf("Error reading log file: %v", err)
	}

	if !strings.Contains(string(data), "ZAP_AUTH_HEADER=Authorization") {
		t.Errorf("Log message not found in file: %v", string(data))
	}

	if !strings.Contains(string(data), "ZAP_AUTH_HEADER_VALUE=Bearer DUMMY_TOKEN") {
		t.Errorf("Log message not found in file: %v", string(data))
	}

	if !strings.Contains(string(data), "ZAP_SITE_NAME=example.com") {
		t.Errorf("Log message not found in file: %v", string(data))
	}
}
