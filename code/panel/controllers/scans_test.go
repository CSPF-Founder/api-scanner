package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
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

	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `jobs` WHERE id = ? ORDER BY `jobs`.`id` LIMIT ?")).
		WithArgs(job.ID, 1).
		WillReturnRows(jobRow)
}

func (ctx *testContext) mockEmptyGetJobByID(job models.Job) {
	jobRow := sqlmock.NewRows([]string{"id", "status", "created_at", "user_id"})

	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `jobs` WHERE id = ? ORDER BY `jobs`.`id` LIMIT ?")).
		WithArgs(job.ID, 1).
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

func TestParseAuthHeaders(t *testing.T) {
	const (
		nameKey   = "http_headers[][header_name]"
		valueKey  = "http_headers[][header_value]"
		customKey = "http_headers[][custom_header_name]"
	)

	tests := []struct {
		name    string
		form    url.Values
		want    []authHeader
		wantErr error
	}{
		{
			name: "no headers",
			form: url.Values{},
			want: []authHeader{},
		},
		{
			name: "predefined header",
			form: url.Values{
				nameKey:   {"Authorization"},
				valueKey:  {"Bearer abc"},
				customKey: {""},
			},
			want: []authHeader{{Name: "Authorization", Value: "Bearer abc"}},
		},
		{
			name: "custom header",
			form: url.Values{
				nameKey:   {"custom"},
				valueKey:  {"secret"},
				customKey: {"X-Tenant"},
			},
			want: []authHeader{{Name: "X-Tenant", Value: "secret"}},
		},
		{
			name: "custom name omitted for non-custom rows",
			form: url.Values{
				nameKey:  {"x-api-key"},
				valueKey: {"key"},
			},
			want: []authHeader{{Name: "x-api-key", Value: "key"}},
		},
		{
			name: "more than one header",
			form: url.Values{
				nameKey:   {"Authorization", "x-api-key"},
				valueKey:  {"Bearer abc", "key"},
				customKey: {"", ""},
			},
			wantErr: errTooManyHeaders,
		},
		{
			name:    "name without value",
			form:    url.Values{nameKey: {"Authorization"}},
			wantErr: errIncompleteHeader,
		},
		{
			name:    "empty value",
			form:    url.Values{nameKey: {"Authorization"}, valueKey: {""}},
			wantErr: errIncompleteHeader,
		},
		{
			name:    "custom row without custom name",
			form:    url.Values{nameKey: {"custom"}, valueKey: {"secret"}},
			wantErr: errIncompleteHeader,
		},
		{
			name: "empty custom name",
			form: url.Values{
				nameKey:   {"custom"},
				valueKey:  {"secret"},
				customKey: {""},
			},
			wantErr: errIncompleteHeader,
		},
		{
			name:    "newline in value",
			form:    url.Values{nameKey: {"Authorization"}, valueKey: {"abc\n[OTHER]\nkey = value"}},
			wantErr: errInvalidHeader,
		},
		{
			name: "INI syntax in custom name",
			form: url.Values{
				nameKey:   {"custom"},
				valueKey:  {"secret"},
				customKey: {"a = b"},
			},
			wantErr: errInvalidHeader,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseAuthHeaders(tc.form)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("expected error %v, got %v", tc.wantErr, err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("expected headers %v, got %v", tc.want, got)
			}
		})
	}
}

func TestAddScanRejectsMalformedInput(t *testing.T) {
	testUser := models.User{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}
	ctx, loginResp := loggedSessionForTest(t, testUser)
	cookie := loginResp.Header.Get("Set-Cookie")

	ctx.mockGetByUserID(testUser)
	req, err := http.NewRequest("GET", ctx.server.URL+"/scans/add", nil)
	if err != nil {
		t.Fatalf("error creating /scans/add request: %v", err)
	}
	req.Header.Set("Cookie", cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error requesting /scans/add: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		t.Fatalf("error reading /scans/add response: %v", err)
	}
	matches := regexp.MustCompile(`const CSRF_TOKEN = '([0-9a-f]+)'`).FindSubmatch(body)
	if matches == nil {
		t.Fatal("CSRF token not found on /scans/add")
	}
	csrfToken := string(matches[1])

	tests := []struct {
		name      string
		fields    map[string]string
		fileName  string
		wantError string
	}{
		{
			name:      "header name without value",
			fields:    map[string]string{"http_headers[][header_name]": "Authorization"},
			fileName:  "spec.yaml",
			wantError: errIncompleteHeader.Error(),
		},
		{
			name:      "file name without extension",
			fileName:  "spec",
			wantError: "Invalid file format",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			mw := multipart.NewWriter(&buf)
			_ = mw.WriteField("csrf_token", csrfToken)
			for k, v := range tc.fields {
				_ = mw.WriteField(k, v)
			}
			fw, err := mw.CreateFormFile("yaml_file", tc.fileName)
			if err != nil {
				t.Fatalf("error creating form file: %v", err)
			}
			_, _ = fw.Write([]byte("openapi: 3.0.0\n"))
			if err := mw.Close(); err != nil {
				t.Fatalf("error closing multipart writer: %v", err)
			}

			ctx.mockGetByUserID(testUser)
			req, err := http.NewRequest("POST", ctx.server.URL+"/scans/", &buf)
			if err != nil {
				t.Fatalf("error creating /scans/ request: %v", err)
			}
			req.Header.Set("Cookie", cookie)
			req.Header.Set("Content-Type", mw.FormDataContentType())

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("error requesting /scans/: %v", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusUnprocessableEntity {
				t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
			}
			var got map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Fatalf("error decoding response: %v", err)
			}
			if got["error"] != tc.wantError {
				t.Fatalf("expected error %q, got %q", tc.wantError, got["error"])
			}
		})
	}
}
