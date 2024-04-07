package zap

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/CSPF-Founder/api-scanner/code/scanner/config"
	"github.com/CSPF-Founder/api-scanner/code/scanner/internal/schemas"
	"github.com/CSPF-Founder/api-scanner/code/scanner/models"
	"github.com/CSPF-Founder/api-scanner/code/scanner/pkg/fileutil"
)

type ZapModule struct {
	job          models.Job
	outputDir    string
	authPath     string
	scannerImage string
}

func NewZapModule(cfg config.Config, job models.Job, scannerImage string) *ZapModule {
	return &ZapModule{
		job:          job,
		outputDir:    filepath.Join(job.GetLocalWorkDir(cfg), "zapoutput/"),
		authPath:     filepath.Join(filepath.Join(job.GetLocalWorkDir(cfg), "zapoutput/"), "zap_config"),
		scannerImage: scannerImage,
	}
}

func (z *ZapModule) WriteAuthFile(authHeaderData []schemas.AuthHeaderMap, siteNameForConfig string) error {
	if len(authHeaderData) == 0 {
		return fmt.Errorf("authHeaderData should not be empty")
	} else if len(authHeaderData) > 1 {
		return fmt.Errorf("authHeaderData should have only one element")
	}

	file, err := os.Create(z.authPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	firstHeader := authHeaderData[0]

	_, err = fmt.Fprintf(file, "ZAP_AUTH_HEADER=%s\n", firstHeader.Name)
	if err != nil {
		return fmt.Errorf("failed to write ZAP_AUTH_HEADER: %v", err)
	}

	_, err = fmt.Fprintf(file, "ZAP_AUTH_HEADER_VALUE=%s\n", firstHeader.Value)
	if err != nil {
		return fmt.Errorf("failed to write ZAP_AUTH_HEADER_VALUE: %v", err)
	}

	_, err = fmt.Fprintf(file, "ZAP_SITE_NAME=%s\n", siteNameForConfig)
	if err != nil {
		return fmt.Errorf("failed to write ZAP_SITE_NAME: %v", err)
	}

	return nil
}

func (z *ZapModule) Run(ctx context.Context, srcOpenAPI string, authHeaderData []schemas.AuthHeaderMap, serverUrl string) error {

	err := os.MkdirAll(z.outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error creating output directory: %v", err)
	}

	err = fileutil.CopyFile(filepath.Join(z.outputDir, filepath.Base(srcOpenAPI)), srcOpenAPI)
	if err != nil {
		return fmt.Errorf("Error copying openApiFile: %v", err)
	}

	var cmd *exec.Cmd
	if len(authHeaderData) > 0 {
		parsedURL, err := url.Parse(serverUrl)
		if err != nil {
			return fmt.Errorf("Error parsing serverUrl: %v", err)
		}
		siteNameForConfig := parsedURL.Hostname()
		err = z.WriteAuthFile(authHeaderData, siteNameForConfig)
		if err != nil {
			return fmt.Errorf("Error writing auth file: %v", err)
		}

		cmd = exec.CommandContext(
			ctx,
			"docker",
			"run",
			"--network",
			"host",
			"--rm",
			"--env-file",
			z.authPath,
			"-v",
			string(z.outputDir)+":/zap/wrk/:rw",
			"-t",
			z.scannerImage,
			"zap-api-scan.py",
			"-t",
			filepath.Base(srcOpenAPI),
			"-f",
			"openapi",
			"-r",
			"zapReport.html",
			"-x",
			"zapReport.xml",
			"-I",
		)

	} else {
		cmd = exec.CommandContext(
			ctx,
			"docker",
			"run",
			"--rm",
			"--network",
			"host",
			"-v",
			string(z.outputDir)+":/zap/wrk/:rw",
			"-t",
			z.scannerImage,
			"zap-api-scan.py",
			"-t",
			filepath.Base(srcOpenAPI),
			"-f",
			"openapi",
			"-r",
			"zapReport.html",
			"-x",
			"zapReport.xml",
			"-I",
		)
	}

	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Error in running ZAP: %v", err)
	}
	return nil
}
