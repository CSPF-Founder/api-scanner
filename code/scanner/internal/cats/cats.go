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

	var cmd *exec.Cmd

	if len(authHeaderData) > 1 {
		firstHeader := authHeaderData[0]
		// if authHeaderData is not empty
		cmd = exec.CommandContext(
			ctx,
			c.binPath,
			"--contract",
			openApiFile,
			"--server",
			serverUrl,
			"-b",
			"--output",
			c.outputDir,
			"-H",
			firstHeader.Name+"="+firstHeader.Value,
		)
	} else {
		cmd = exec.CommandContext(
			ctx,
			c.binPath,
			"--contract",
			openApiFile,
			"--server",
			serverUrl,
			"-b",
			"--output",
			c.outputDir,
		)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error in running catz: %v", string(out))
	}
	return nil
}
