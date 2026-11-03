package cats

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/CSPF-Founder/api-scanner/code/scanner/config"
	"github.com/CSPF-Founder/api-scanner/code/scanner/internal/schemas"
	"github.com/CSPF-Founder/api-scanner/code/scanner/models"
)

type CatsModule struct {
	binPath     string
	job         models.Job
	outputDir   string
	requestsDir string
}

func NewCatsModule(cfg config.Config, job models.Job) *CatsModule {
	return &CatsModule{
		job:         job,
		binPath:     "/app/bin/qats",
		outputDir:   filepath.Join(job.GetLocalWorkDir(cfg), "catsoutput/"),
		requestsDir: filepath.Join(job.GetLocalWorkDir(cfg), "requests/"),
	}
}

func (c *CatsModule) Run(ctx context.Context, openApiFile string, authHeaderData []schemas.AuthHeaderMap, serverUrl string) error {
	cmd := exec.CommandContext(ctx, c.binPath, c.buildArgs(openApiFile, authHeaderData, serverUrl)...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error in running catz: %v", string(out))
	}
	return nil
}

// buildArgs returns the CATS command line. Only the first auth header is
// passed, matching what the ZAP scan supports.
func (c *CatsModule) buildArgs(openApiFile string, authHeaderData []schemas.AuthHeaderMap, serverUrl string) []string {
	args := []string{"--contract", openApiFile, "--server", serverUrl, "-b", "--output", c.outputDir}
	if len(authHeaderData) > 0 {
		header := authHeaderData[0]
		args = append(args, "-H", header.Name+"="+header.Value)
	}
	return args
}
