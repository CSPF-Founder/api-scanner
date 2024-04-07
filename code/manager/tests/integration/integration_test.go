package integration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/CSPF-Founder/api-scanner/code/manager/db"
	"github.com/ory/dockertest/v3"
	"github.com/pressly/goose"
)

var testDB *sql.DB

func TestMain(m *testing.M) {

	// Create a new Docker pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// Start a MySQL container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mariadb",
		Tag:        "11.2.2",
		Env:        []string{"MYSQL_ROOT_PASSWORD=secret"},
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	if err := pool.Retry(func() error {
		dbURI := fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp"))
		testDB, err = db.ConnectDB(dbURI)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
		return
	}

	// DB migrations
	// Run Goose migrations

	if err := goose.SetDialect("mysql"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(testDB, "../../db/migrations"); err != nil {
		log.Fatalf("Goose failed to run migrations: %s", err)
		return
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)

}
