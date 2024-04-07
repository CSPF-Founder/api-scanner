package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/CSPF-Founder/api-scanner/code/panel/auth"
	"github.com/CSPF-Founder/api-scanner/code/panel/models"
	"github.com/PuerkitoBio/goquery"
)

const (
	InvalidCSRFTokenError = "Invalid CSRF token"
)

// Tests for the Users controller

func testLoginAttempt(t *testing.T, ctx *testContext, client *http.Client, username, password string) *http.Response {

	resp, err := http.Get(fmt.Sprintf("%s/users/login", ctx.server.URL))
	if err != nil {
		t.Fatalf("error requesting the /users/login endpoint: %v", err)
	}
	got := resp.StatusCode
	expected := http.StatusOK
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatalf("error parsing /login response body")
	}

	var csrfToken string
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		scriptContent := s.Text()
		// Adjust the regular expression according to how the token is set in the script
		re := regexp.MustCompile(`const CSRF_TOKEN = '([0-9a-f]+)'`)
		matches := re.FindStringSubmatch(scriptContent)
		if len(matches) > 1 {
			csrfToken = matches[1]
			return
		}
	})

	if client == nil {
		client = &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/login", ctx.server.URL), strings.NewReader(url.Values{
		"username":   {username},
		"password":   {password},
		"csrf_token": {csrfToken},
	}.Encode()))
	if err != nil {
		t.Fatalf("error creating new /users/login request: %v", err)
	}

	req.Header.Set("Cookie", resp.Header.Get("Set-Cookie"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("error requesting the /users/login endpoint: %v", err)
	}
	return resp

}

func TestLoginCSRF(t *testing.T) {
	ctx := setupTest(t)

	resp, err := http.PostForm(ctx.server.URL+"/users/login",
		url.Values{
			"username": {"test"},
			"password": {"test"},
		})

	if err != nil {
		t.Fatalf("error requesting the /login endpoint: %v", err)
	}

	got := resp.StatusCode
	expected := http.StatusForbidden
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}
}

func TestInvalidCredentials(t *testing.T) {
	ctx := setupTest(t)

	testUser := models.User{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}

	hash, err := auth.GeneratePasswordHash("test")
	if err != nil {
		t.Fatalf("error generating password hash: %v", err)
	}
	testUser.Password = hash

	// Mock SQL - series of expectations for the mock DB
	ctx.mockUserCount(1)
	ctx.mockGetByUsername(testUser)
	ctx.mockUserCount(1)

	resp := testLoginAttempt(t, ctx, &http.Client{}, "test", "invalidpass")
	got := resp.StatusCode
	expected := http.StatusOK
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Fatalf("error parsing /login response body")
	}

	if !strings.Contains(doc.Text(), "Invalid Username/Password") {
		t.Fatalf("Not found Invalid Username/Password")
	}
}

func TestSuccessfulLogin(t *testing.T) {
	testUser := models.User{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}
	_, _ = loggedSessionForTest(t, testUser)
}

func loggedSessionForTest(t *testing.T, testUser models.User) (*testContext, *http.Response) {
	ctx := setupTest(t)

	hash, err := auth.GeneratePasswordHash("test")
	if err != nil {
		t.Fatalf("error generating password hash: %v", err)
	}
	testUser.Password = hash

	// Mock SQL - series of expectations for the mock DB
	ctx.mockUserCount(1)
	ctx.mockGetByUsername(testUser)
	ctx.mockUserCount(1)

	resp := testLoginAttempt(t, ctx, nil, "test", "test")
	got := resp.StatusCode
	expected := http.StatusFound
	if got != expected {
		t.Fatalf("invalid status code received. expected %d got %d", expected, got)
	}

	// get redirect url path
	redirectURL, err := resp.Location()
	if err != nil {
		t.Fatalf("error getting redirect url: %v", err)
	}

	if redirectURL.Path != "/" {
		t.Fatalf("invalid redirect url path. expected %s got %s", "/", redirectURL.Path)
	}
	return ctx, resp
}
