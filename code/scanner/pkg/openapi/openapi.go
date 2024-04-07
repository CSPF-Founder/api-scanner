package openapi

import (
	"errors"
	"os"

	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

// ValidateOpenapiFile validates the openapi file
// If it is valid, it returns true and the server URL
func GetServerURLFromOpenAPI(filePath string) (string, error) {

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	loader := openapi3.NewLoader()

	// Passing the filname
	// swagger, err := loader.LoadFromFile(file.Name())
	swagger, err := loader.LoadFromData(fileContent) // Passing the byte data, use this if using os.ReadFile()
	if err != nil {
		return "", err
	}

	// Passing the Context
	if err := swagger.Validate(loader.Context); err != nil {
		return "", err
	}
	servers := swagger.Servers
	if len(servers) == 0 {
		return "", errors.New("No servers mentioned in spec")
	}

	serverURL := servers[0].URL
	if serverURL == "" {
		return "", errors.New("No URL present")
	}

	parsedURL, err := url.Parse(serverURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", errors.New("Invalid URL")
	}

	return serverURL, nil
}
