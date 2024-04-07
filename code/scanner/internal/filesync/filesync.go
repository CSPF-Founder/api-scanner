package filesync

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/CSPF-Founder/api-scanner/code/scanner/config"
	"github.com/CSPF-Founder/api-scanner/code/scanner/logger"
	"github.com/CSPF-Founder/api-scanner/code/scanner/models"
	"github.com/CSPF-Founder/api-scanner/code/scanner/pkg/fileutil"
)

// Copy from remote directory to local directory
func CopyFromRemote(cfg config.Config, l *logger.FileLogger, job *models.Job) (bool, error) {
	scanDir := job.GetRemoteWorkDir(cfg)
	destDir := job.GetLocalWorkDir(cfg)

	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		os.RemoveAll(destDir)
	}
	err := fileutil.CopyDir(destDir, scanDir)
	if err != nil {
		l.Error("failed to copy remote directory", err)
		return false, err
	} else {
		return true, nil
	}
}

// Copy log file from local directory to remote directory
func CopyToRemoteOnlyLog(cfg config.Config, job *models.Job) error {
	srcDir := job.GetLocalWorkDir(cfg)
	destDir := job.GetRemoteWorkDir(cfg)

	logPath := filepath.Join(srcDir, "logs")
	destPath := filepath.Join(destDir, "logs")

	srcFile, err := os.Open(logPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return err
}

// Copy from local directory to remote directory
func CopyToRemote(cfg config.Config, job *models.Job) (bool, error) {
	srcDir := job.GetLocalWorkDir(cfg)
	destDir := job.GetRemoteWorkDir(cfg)

	// Copy report to remote dir
	reportPath := filepath.Join(srcDir, "report/report.docx")
	destReportPath := filepath.Join(destDir, "report.docx")
	err := fileutil.CopyFile(destReportPath, reportPath)
	if err != nil {
		return false, err
	}

	// Remove local report file
	// _ = os.RemoveAll(filepath.Dir(reportPath))

	// Copy request archive to remote dir
	requestPath := filepath.Join(srcDir, "requests")

	createRequestArchive := false
	if stat, err := os.Stat(requestPath); err == nil && stat.IsDir() {
		files, err := os.ReadDir(requestPath)
		if err == nil && len(files) > 0 {
			createRequestArchive = true
		}
	}

	if createRequestArchive {
		destRequestPath := filepath.Join(destDir, "error_requests.zip")
		err := fileutil.ZipDir(requestPath, destRequestPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Copy logs to remote dir
	logPath := filepath.Join(srcDir, "logs")
	destLogPath := filepath.Join(destDir, "logs")
	logErr := fileutil.CopyFile(destLogPath, logPath)
	if logErr != nil {
		return false, logErr
	}

	// Remove local log file to save space
	// _ = os.Remove(logPath)

	zapDir := filepath.Join(srcDir, "zapoutput")
	destZapDir := filepath.Join(destDir, "zapoutput")
	_ = fileutil.CopyDir(destZapDir, zapDir)

	// Remove all local files
	_ = os.RemoveAll(srcDir)

	return true, nil
}
