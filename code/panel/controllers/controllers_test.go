package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/CSPF-Founder/api-scanner/code/panel/config"
	"github.com/CSPF-Founder/api-scanner/code/panel/logger"
	"github.com/CSPF-Founder/api-scanner/code/panel/models"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// testContext is the data required to test API related functions
type testContext struct {
	app      *App
	server   *httptest.Server
	origPath string
	mock     sqlmock.Sqlmock
}

type MockHttpClient struct {
	// Add fields to simulate responses or errors
	Response *http.Response
	Err      error
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	if m.Err != nil {
		return nil, m.Err
	}

	return m.Response, nil
}

func setupConfig(t *testing.T) *config.Config {
	os.Setenv("PRODUCT_TITLE", "API Scanner")
	os.Setenv("SERVER_ADDRESS", "0.0.0.0:8080")
	os.Setenv("DATABASE_URI", "root:@(:3306)/api_db?charset=utf8&parseTime=True&loc=Local")
	os.Setenv("DBMS_TYPE", "mysql")
	os.Setenv("COPYRIGHT_FOOTER_COMPANY", "Cyber Security & Privacy Foundation")
	os.Setenv("WORK_DIR", "/app/data/")
	os.Setenv("TEMP_UPLOADS_DIR", "/app/data/temp_uploads/")
	os.Setenv("MIGRATIONS_PREFIX", "db")
	os.Setenv("LOG_FILENAME", "logfile")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("USE_TLS", "true")
	os.Setenv("CERT_PATH", "/app/certs/panel.crt")
	os.Setenv("KEY_PATH", "/app/certs/panel.key")

	os.Setenv("USE_DOTENV", "false")
	config, err := config.LoadConfig()
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}
	return config
}

func NewMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return gormDB, mock
}

func setupTest(t *testing.T) *testContext {
	wd, _ := os.Getwd()
	fmt.Println(wd)
	conf := setupConfig(t)
	appLogger, err := logger.NewLogger(conf.Logging)
	if err != nil {
		t.Fatal("Error setting up logging: ", err)
	}

	app := NewApp(conf, appLogger)

	gormDB, mock := NewMockDB()

	models.SetupDB(gormDB)
	app.SetupRoutes()

	ctx := &testContext{}
	ctx.app = app
	ctx.server = httptest.NewUnstartedServer(app.server.Handler)
	ctx.server.Config.Addr = ctx.app.config.ServerConf.ServerAddress
	ctx.mock = mock

	ctx.server.Start()

	origPath, _ := os.Getwd()
	ctx.origPath = origPath
	err = os.Chdir("../")
	if err != nil {
		t.Fatalf("error changing directories to setup asset discovery: %v", err)
	}
	return ctx
}

// SQL mocking helper functions
func (ctx *testContext) mockUserCount(count int) {
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(count)
	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users`")).
		WillReturnRows(rows)
}

func (ctx *testContext) mockGetByUsername(user models.User) {
	userRow := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(user.ID, user.Username, user.Password)

	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT 1")).
		WithArgs(user.Username).
		WillReturnRows(userRow)
}

func (ctx *testContext) mockGetByUserID(user models.User) {
	userRow := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(user.ID, user.Username, user.Password)

	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ? ORDER BY `users`.`id` LIMIT 1")).
		WithArgs(user.ID).
		WillReturnRows(userRow)
}
