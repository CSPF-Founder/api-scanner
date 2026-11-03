package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/CSPF-Founder/api-scanner/code/manager/db"
	"github.com/ory/dockertest/v4"
	"github.com/pressly/goose"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	// os.Exit skips deferred calls, so cleanup lives in run.
	os.Exit(run(m))
}

func run(m *testing.M) int {
	ctx := context.Background()

	// Create a new Docker pool
	pool, err := dockertest.NewPool(ctx, "")
	if err != nil {
		log.Printf("Could not construct pool: %s", err)
		return 1
	}
	defer func() {
		if err := pool.Close(ctx); err != nil {
			log.Printf("Could not purge resources: %s", err)
		}
	}()

	// Start a MySQL container
	resource, err := pool.Run(ctx, "mariadb",
		dockertest.WithTag("11.2.2"),
		dockertest.WithEnv([]string{"MYSQL_ROOT_PASSWORD=secret"}),
		dockertest.WithoutReuse(),
	)
	if err != nil {
		log.Printf("Could not start resource: %s", err)
		return 1
	}

	if err := pool.Retry(ctx, 0, func() error {
		dbURI := fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp"))
		testDB, err = db.ConnectDB(dbURI)
		return err
	}); err != nil {
		log.Printf("Could not connect to docker: %s", err)
		return 1
	}

	// DB migrations
	// Run Goose migrations

	if err := goose.SetDialect("mysql"); err != nil {
		log.Print(err)
		return 1
	}

	if err := goose.Up(testDB, "../../db/migrations"); err != nil {
		log.Printf("Goose failed to run migrations: %s", err)
		return 1
	}

	return m.Run()
}
