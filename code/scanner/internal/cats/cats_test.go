package cats

import (
	"reflect"
	"testing"

	"github.com/CSPF-Founder/api-scanner/code/scanner/internal/schemas"
)

func TestBuildArgs(t *testing.T) {
	c := &CatsModule{outputDir: "/out"}
	base := []string{"--contract", "spec.yaml", "--server", "https://api.test", "-b", "--output", "/out"}

	tests := []struct {
		name    string
		headers []schemas.AuthHeaderMap
		want    []string
	}{
		{
			name: "no auth header",
			want: base,
		},
		{
			name:    "auth header is passed",
			headers: []schemas.AuthHeaderMap{{Name: "x-api-key", Value: "k=1"}},
			want:    append(append([]string{}, base...), "-H", "x-api-key=k=1"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := c.buildArgs("spec.yaml", tc.headers, "https://api.test")
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected args %v, got %v", tc.want, got)
			}
		})
	}
}
