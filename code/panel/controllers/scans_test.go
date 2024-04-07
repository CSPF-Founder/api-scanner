package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/CSPF-Founder/api-scanner/code/panel/enums/jobstatus"
	"github.com/CSPF-Founder/api-scanner/code/panel/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PuerkitoBio/goquery"
)

// Tests for the Scams controller

func (ctx *testContext) mockGetJobByID(job models.Job) {
	jobRow := sqlmock.NewRows([]string{"id", "status", "created_at", "user_id"}).
		AddRow(job.ID, job.Status, job.CreatedAt, job.UserID)

	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `jobs` WHERE id = ? ORDER BY `jobs`.`id` LIMIT 1")).
		WithArgs(job.ID).
		WillReturnRows(jobRow)
}

func (ctx *testContext) mockEmptyGetJobByID(job models.Job) {
	jobRow := sqlmock.NewRows([]string{"id", "status", "created_at", "user_id"})

	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `jobs` WHERE id = ? ORDER BY `jobs`.`id` LIMIT 1")).
		WithArgs(job.ID).
		WillReturnRows(jobRow)
}

func TestDeleteScanWithMissingCSRF(t *testing.T) {
	testUser := models.User{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}
	ctx, resp := loggedSessionForTest(t, testUser)

	testJob := models.Job{
		ID:     1,
		Status: jobstatus.ScanStarted,
		UserID: 1,
	}

	ctx.mockGetByUserID(testUser)

	deleteURL := fmt.Sprintf("%s/scans/%d", ctx.server.URL, testJob.ID)
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		t.Fatalf("error creating new /users/login request: %v", err)
	}

	req.Header.Set("Cookie", resp.Header.Get("Set-Cookie"))

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the /users/login endpoint: %v", err)
	}
	got := resp.StatusCode
	expected := http.StatusForbidden
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatalf("error parsing /login response body")
	}

	if !strings.Contains(doc.Text(), InvalidCSRFTokenError) {
		t.Fatalf("Expected %s in response body but not found", InvalidCSRFTokenError)
	}

}

func TestDeleteInvalidScanID(t *testing.T) {
	testUser := models.User{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}
	ctx, resp := loggedSessionForTest(t, testUser)

	testJob := models.Job{
		ID:     1,
		Status: jobstatus.ScanStarted,
		UserID: 1,
	}

	ctx.mockGetByUserID(testUser)
	ctx.mockEmptyGetJobByID(testJob)

	deleteURL := fmt.Sprintf("%s/scans/%d", ctx.server.URL, testJob.ID)
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		t.Fatalf("error creating new /users/login request: %v", err)
	}

	req.Header.Set("Cookie", resp.Header.Get("Set-Cookie"))

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the /users/login endpoint: %v", err)
	}
	got := resp.StatusCode
	expected := http.StatusForbidden
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatalf("error parsing /login response body")
	}

	if !strings.Contains(doc.Text(), InvalidCSRFTokenError) {
		t.Fatalf("Expected %s in response body but not found", InvalidCSRFTokenError)
	}

}

func TestDeleteValidScanID(t *testing.T) {
	testUser := models.User{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}
	ctx, resp := loggedSessionForTest(t, testUser)

	testJob := models.Job{
		ID:     1,
		Status: jobstatus.ScanStarted,
		UserID: 1,
	}

	ctx.mockGetByUserID(testUser)
	ctx.mockGetJobByID(testJob)

	deleteURL := fmt.Sprintf("%s/scans/%d", ctx.server.URL, testJob.ID)
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		t.Fatalf("error creating new /users/login request: %v", err)
	}

	req.Header.Set("Cookie", resp.Header.Get("Set-Cookie"))

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the /users/login endpoint: %v", err)
	}
	got := resp.StatusCode
	expected := http.StatusForbidden
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatalf("error parsing /login response body")
	}

	if !strings.Contains(doc.Text(), InvalidCSRFTokenError) {
		t.Fatalf("Expected %s in response body but not found", InvalidCSRFTokenError)
	}

}

func TestDeleteScanWithInvalidSession(t *testing.T) {
	testUser := models.User{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}
	// ctx, resp := loggedSessionForTest(t, testUser)
	ctx := setupTest(t)

	testJob := models.Job{
		ID:     1,
		Status: jobstatus.ScanStarted,
		UserID: 1,
	}

	ctx.mockGetByUserID(testUser)
	ctx.mockGetJobByID(testJob)

	deleteURL := fmt.Sprintf("%s/scans/%d", ctx.server.URL, testJob.ID)
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		t.Fatalf("error creating new /users/login request: %v", err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the /users/login endpoint: %v", err)
	}
	got := resp.StatusCode
	expected := http.StatusForbidden
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatalf("error parsing /login response body")
	}

	if !strings.Contains(doc.Text(), InvalidCSRFTokenError) {
		t.Fatalf("Expected %s in response body but not found", InvalidCSRFTokenError)
	}

}
