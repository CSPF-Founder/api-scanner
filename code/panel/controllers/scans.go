package controllers

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	ctx "github.com/CSPF-Founder/api-scanner/code/panel/context"
	"github.com/CSPF-Founder/api-scanner/code/panel/enums/flashtypes"
	"github.com/CSPF-Founder/api-scanner/code/panel/enums/jobstatus"
	"github.com/CSPF-Founder/api-scanner/code/panel/internal/openapivalidator"
	mid "github.com/CSPF-Founder/api-scanner/code/panel/middlewares"
	"github.com/CSPF-Founder/api-scanner/code/panel/models"
	"github.com/CSPF-Founder/api-scanner/code/panel/utils"
	"github.com/CSPF-Founder/api-scanner/code/panel/views"
	"github.com/go-chi/chi/v5"
)

type scansController struct {
	*App
}

func newScansController(app *App) *scansController {
	return &scansController{
		App: app,
	}
}

func (c *scansController) registerRoutes() http.Handler {
	router := chi.NewRouter()

	// Authenticated Routes
	router.Group(func(r chi.Router) {
		r.Use(mid.RequireLogin)

		r.Get("/", c.List)          // List all scans
		r.Post("/", c.AddHandler)   // Add a new scan
		r.Get("/add", c.DisplayAdd) // Display add scan form

		// eg: DELETE /scans/1
		r.Route("/{scanID:[0-9]+}", func(r chi.Router) {
			r.Delete("/", c.Delete)
			r.Get("/report", c.DownloadReport)
			r.Get("/error-logs", c.DownloadErrorRequests)
		})

	})

	return router
}

func (c *scansController) DisplayAdd(w http.ResponseWriter, r *http.Request) {
	templateData := views.NewTemplateData(c.config, c.session, r)
	templateData.Title = "Scans Add"
	if err := views.RenderTemplate(w, "scans/add", templateData); err != nil {
		c.logger.Error("Error rendering template: ", err)
	}
}

func (c *scansController) AddHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		c.SendJSONError(w, "Invalid file")
		return
	}
	httpHeadersInput := r.Form["http_headers[][header_name]"]
	var customHeaderName string

	if len(httpHeadersInput) > 0 {
		for i, httpHeader := range httpHeadersInput {
			headerName := httpHeader
			headerValue := r.Form["http_headers[][header_value]"][i]

			if headerName == "" || headerValue == "" {
				// http.Error(w, "Please fill all the inputs", http.StatusBadRequest)
				c.SendJSONError(w, "Please fill all the inputs")
				return
			}

			if headerName == "custom" {
				customHeaderName = r.Form["http_headers[][custom_header_name]"][i]
				if customHeaderName == "" {
					c.SendJSONError(w, "Please fill all the inputs")
					return
				}
			}
		}
	}

	file, handler, err := r.FormFile("yaml_file")
	if err != nil {
		c.SendJSONError(w, "Invalid file")
		return
	}
	defer file.Close()

	fileName := handler.Filename
	fileType := filepath.Ext(fileName)[1:]

	allowedTypes := map[string]bool{"yml": true, "yaml": true}
	if !allowedTypes[fileType] {
		c.SendJSONError(w, "Invalid file format")
		return
	}

	fileSize := handler.Size
	// filesize max 10 mb
	if fileSize > 10000000 {
		c.SendJSONError(w, "File size should be less than 10 MB")
		return
	}

	tmpFileName := filepath.Join(c.config.TempUploadsDir, fmt.Sprintf("%s.yaml", utils.GetRandomString(32)))
	out, err := os.Create(tmpFileName)
	if err != nil {
		c.logger.Error("Error creating file:", err)
		c.SendJSONError(w, "There's some issue with file")
		return
	}

	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.logger.Error("Error copying file:", err)
		c.SendJSONError(w, "Invalid File")
		return
	}

	defer os.Remove(tmpFileName)

	c.handleJobCreation(w, r, tmpFileName)
}

func (c *scansController) handleJobCreation(w http.ResponseWriter, r *http.Request, tmpFileName string) {

	serverURL, err := openapivalidator.ValidateOpenapiFile(tmpFileName)
	if err != nil {
		c.SendJSONError(w, err.Error())
		return
	}

	err = IsServerReachable(serverURL)
	if err != nil {
		c.SendJSONError(w, err.Error())
		return
	}

	job := models.Job{
		Status:        jobstatus.Default,
		ApiURL:        serverURL,
		CreatedAt:     time.Now(),
		CompletedTime: time.Now(),
		UserID:        ctx.Get(r, "user").(models.User).ID,
	}

	if err := models.SaveJob(&job); err != nil {
		c.SendJSONError(w, "Unable to add the scan")
		return
	}

	targetDir := filepath.Join(c.config.WorkDir, fmt.Sprintf("user_%d/job_%d/", job.UserID, job.ID))
	destinationFile := filepath.Join(targetDir, "openapi.yaml")

	if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
		// if err := os.MkdirAll(targetDir, 0777); err != nil {
		c.logger.Error("HandleJobCreation:", err)
		c.SendJSONError(w, "Unable to create target directory")
		return
	}

	if err := os.Rename(tmpFileName, destinationFile); err != nil {
		c.logger.Error("HandleJobCreation:", err)
		c.SendJSONError(w, "Unable to upload the yaml file")
		return
	}

	authHeadersFileName := "auth_headers.conf"
	authHeadersPath := filepath.Join(targetDir, authHeadersFileName)

	fp, err := os.Create(authHeadersPath)
	if err != nil {
		c.logger.Error("HandleJobCreation:", err)
		c.SendJSONError(w, "Unable to create auth headers file")
		return
	}
	defer fp.Close()

	fmt.Fprintf(fp, "[AUTH_HEADERS]\n")
	httpHeadersInput := r.Form["http_headers[][header_name]"]

	if len(httpHeadersInput) > 0 {
		for i, httpHeader := range httpHeadersInput {
			headerName := httpHeader
			headerValue := r.Form["http_headers[][header_value]"][i]

			if headerName == "custom" {
				headerName = r.Form["http_headers[][custom_header_name]"][i]
			}
			fmt.Fprintf(fp, "%s = %s\n", headerName, headerValue)
		}
	}

	c.SendJSONSuccess(w, "Successfully added the scan")
}

func (c *scansController) List(w http.ResponseWriter, r *http.Request) {
	user := ctx.Get(r, "user").(models.User)
	data, err := models.GetJobs(&user)
	if err != nil {
		c.SendJSONError(w, err.Error())
	}
	templateData := views.NewTemplateData(c.config, c.session, r)
	templateData.Title = "Scans List"
	templateData.Data = data
	if err := views.RenderTemplate(w, "scans/list", templateData); err != nil {
		c.logger.Error("Error rendering template: ", err)
	}
}

func (c *scansController) Delete(w http.ResponseWriter, r *http.Request) {
	scanID := chi.URLParam(r, "scanID")

	// Convert the string to uint64
	jobID, err := strconv.ParseUint(scanID, 10, 64)
	if err != nil {
		c.SendJSONError(w, "Invalid Scan ID parameter")
		return
	}

	user := ctx.Get(r, "user").(models.User)

	job, err := models.GetByIDAndUser(jobID, user.ID)
	if err != nil {
		c.SendJSONError(w, "Invalid ID or already deleted")
		return
	}

	// Construct target directory path
	targetDirPath := fmt.Sprintf("%s/user_%d/job_%d", c.config.WorkDir, user.ID, jobID)

	// Check if the target directory exists before removing it
	if _, err := os.Stat(targetDirPath); err == nil {
		err := os.RemoveAll(targetDirPath)
		if err != nil {
			c.logger.Error("Unable to delete the scan", err)
			c.SendJSONError(w, "Unable to delete the scan folder")
			return
		}
	}

	delete := models.DeleteJob(job.ID)
	if delete != nil {
		c.SendJSONError(w, "Unable to delete the scan")
		return
	}

	c.SendJSONSuccess(w, "Successfully deleted the scan")
}

// Check Server is reachable
// Does not validate the status code
func IsServerReachable(serverURL string) error {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Head(serverURL)
	if err != nil {

		if urlErr, ok := err.(*url.Error); ok {
			if urlErr.Timeout() {
				return errors.New("Server connection timeout")
			}

			errMsg := urlErr.Err.Error()
			if strings.Contains(errMsg, "no such host") {
				return errors.New("Hostname not found")
			} else if strings.Contains(errMsg, "connection refused") {
				return errors.New("connection refused or host unreachable")
			} else if strings.Contains(errMsg, "network is unreachable") {
				return errors.New("network is unreachable")
			}
		}

		if _, ok := err.(*net.DNSError); ok {
			return errors.New("Hostname not found")
		}

		if netErr, ok := err.(net.Error); ok {
			if netErr.Timeout() {
				return errors.New("Server connection timeout")
			} else if _, ok := netErr.(*net.OpError); ok {
				return errors.New("connection refused or host unreachable")
			}
		}

		return err
	}

	defer resp.Body.Close()

	return nil
}

func (c *scansController) DownloadReport(w http.ResponseWriter, r *http.Request) {

	scanID := chi.URLParam(r, "scanID")

	jobID, _ := strconv.ParseUint(scanID, 10, 64)

	user := ctx.Get(r, "user").(models.User)

	job, err := models.GetByIDAndUser(jobID, user.ID)
	if err != nil {
		c.FlashAndGoBack(w, r, flashtypes.FlashDanger, "Invalid Scan Id")
		return
	}

	jobDir := filepath.Join(c.config.WorkDir, fmt.Sprintf("user_%d/job_%d/", user.ID, job.ID))
	reportDirPath := filepath.Join(jobDir, "report.docx")

	if fileInfo, err := os.Stat(reportDirPath); err == nil && !fileInfo.IsDir() {
		w.Header().Set("Content-Type", "application/txt")
		w.Header().Set("Content-Disposition", fmt.Sprintf("filename=report_%d.docx", job.ID))
		http.ServeFile(w, r, reportDirPath)
	} else {
		c.FlashAndGoBack(w, r, flashtypes.FlashWarning, "Report file does not exist!")
		return
	}
}

func (c *scansController) DownloadErrorRequests(w http.ResponseWriter, r *http.Request) {

	scanID := chi.URLParam(r, "scanID")

	// Convert the string to uint64
	jobID, _ := strconv.ParseUint(scanID, 10, 64)

	user := ctx.Get(r, "user").(models.User)

	job, err := models.GetByIDAndUser(jobID, user.ID)
	if err != nil {
		c.FlashAndGoBack(w, r, flashtypes.FlashDanger, "Invalid Scan Id")
		return
	}

	jobDir := filepath.Join(c.config.WorkDir, fmt.Sprintf("user_%d/job_%d/", user.ID, job.ID))

	reportDirPath := filepath.Join(jobDir, "error_requests.zip")

	if fileInfo, err := os.Stat(reportDirPath); err == nil && !fileInfo.IsDir() {
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", "attachment; filename=error_requests.zip")
		http.ServeFile(w, r, reportDirPath)
	} else {
		c.FlashAndGoBack(w, r, flashtypes.FlashWarning, "No error requests file to download")
		return
	}
}
