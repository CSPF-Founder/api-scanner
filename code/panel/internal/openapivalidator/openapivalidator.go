package openapivalidator

import (
	"errors"
	"net/url"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

// function to validate an OpenAPI3 File Spec
func ValidateOpenapiFile(fileNameOpenApi string) (string, error) {

	// file, err := os.Stat(fileNameOpenApi)
	fileContent, err := os.ReadFile(fileNameOpenApi) // this reads the byte, can be used as well
	if err != nil {
		return "", errors.New("Invalid file")
	}

	loader := openapi3.NewLoader()

	// Passing the filname
	// swagger, err := loader.LoadFromFile(file.Name())
	swagger, err := loader.LoadFromData(fileContent) // Passing the byte data, use this if using os.ReadFile()
	if err != nil {
		return "", errors.New("Invalid file")
	}

	// Passing the Context
	if err := swagger.Validate(loader.Context); err != nil {
		return "", errors.New("Not valid spec")
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
		return "", errors.New("No URL present")
	}

	return serverURL, nil
}
