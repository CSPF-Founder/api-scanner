package openapivalidator

import (
	"testing"
)

func TestValidYaml(t *testing.T) {

	yaml_path := "../../test_data/yaml/valid.yml"

	serverURL, err := ValidateOpenapiFile(yaml_path)
	if err != nil {
		t.Errorf("Result was incorrect, got: %v, want: no error.", err)
	}

	expected := "http://192.168.56.1"
	if serverURL != expected {
		t.Errorf("Result was incorrect, got: %s, want: %s.", serverURL, expected)
	}

}

func TestInvalidYaml(t *testing.T) {

	yaml_path := "../../../../test_data/yaml/invalid.yml"

	serverURL, err := ValidateOpenapiFile(yaml_path)

	if err == nil {
		t.Errorf("Result was incorrect, got: no error, want: Invalid file")
	}

	if serverURL != "" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", serverURL, "empty string")
	}

}

func TestInvalidOpenApi(t *testing.T) {

	yaml_path := "../../../../test_data/yaml/invalid-openapi.yml"

	serverURL, err := ValidateOpenapiFile(yaml_path)

	if err == nil {
		t.Errorf("Result was incorrect, got: no error, want: Not valid spec")
	}

	if serverURL != "" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", serverURL, "empty string")
	}

}

func TestInvalidPath(t *testing.T) {

	yaml_path := "../../../../test_data/yaml/NON-Existing.yml"

	serverURL, err := ValidateOpenapiFile(yaml_path)

	if err == nil {
		t.Errorf("Result was incorrect, got: no error, want: Invalid file")
	}

	if serverURL != "" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", serverURL, "empty string")
	}

}

func TestMissingServer(t *testing.T) {

	yaml_path := "../../../../test_data/yaml/missing-server.yml"

	serverURL, err := ValidateOpenapiFile(yaml_path)

	if err == nil {
		t.Errorf("Result was incorrect, got: no error, want: No servers mentioned in spec")
	}

	if serverURL != "" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", serverURL, "empty string")
	}
}
