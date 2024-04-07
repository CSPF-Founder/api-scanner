package scanner

import (
	"context"

	"github.com/CSPF-Founder/api-scanner/code/scanner/config"
	"github.com/CSPF-Founder/api-scanner/code/scanner/enums/jobstatus"
	"github.com/CSPF-Founder/api-scanner/code/scanner/internal/cats"
	"github.com/CSPF-Founder/api-scanner/code/scanner/internal/schemas"
	"github.com/CSPF-Founder/api-scanner/code/scanner/internal/zap"
	"github.com/CSPF-Founder/api-scanner/code/scanner/logger"
	"github.com/CSPF-Founder/api-scanner/code/scanner/models"
)

type ScannerModule struct {
	logger *logger.FileLogger
	Job    models.Job
}

func Scanner(job models.Job, logger *logger.FileLogger) *ScannerModule {
	return &ScannerModule{Job: job, logger: logger}
}

func (s *ScannerModule) UpdateJobStatus(ctx context.Context, model models.DBModel, jobStatus jobstatus.JobStatus) {
	isUpdated, err := model.UpdateJobStatus(ctx, s.Job.ID, jobStatus)
	if isUpdated {
		s.logger.Info("Job Status Updated")
	} else {
		s.logger.Error("Job Status Failed to Update", err)
	}
}

func (s *ScannerModule) Run(
	ctx context.Context,
	cfg config.Config,
	model models.DBModel,
	serverURL string,
	openApiFile string,
	authHeaderData []schemas.AuthHeaderMap,
) bool {

	// s.UpdateJobStatus(ctx, model, jobstatus.ScanInitiated)
	s.UpdateJobStatus(ctx, model, jobstatus.ScanStarted)

	s.logger.Info("Valid OpenAPI File")
	s.logger.Info("Scanning Started" + string(serverURL))

	catScanner := cats.NewCatsModule(cfg, s.Job)
	s.logger.Info("Catz scan start")
	err := catScanner.Run(ctx, openApiFile, authHeaderData, serverURL)

	catsCompleted := false
	if err != nil {
		s.logger.Error("Error running catz", err)
	} else {
		catsCompleted = true
		s.logger.Info("completed catz")
		s.UpdateJobStatus(ctx, model, jobstatus.CatzCompleted)

	}

	zapScanner := zap.NewZapModule(cfg, s.Job, cfg.ScannerImage)
	s.logger.Info("Zap scan start")
	err = zapScanner.Run(ctx, openApiFile, authHeaderData, serverURL)
	zapCompleted := false
	if err != nil {
		s.logger.Error("Error running zap", err)
	} else {
		zapCompleted = true
		s.logger.Info("completed zap")
		s.UpdateJobStatus(ctx, model, jobstatus.ZapCompleted)
	}

	if !zapCompleted && !catsCompleted {
		return false
	}

	// logger.info("Catz and Zap Success")
	//     self.update_job_status(JobStatus.MODULES_FINISHED)

	s.logger.Info("Catz and Zap Success")
	s.UpdateJobStatus(ctx, model, jobstatus.ModulesFinished)
	return true
}
