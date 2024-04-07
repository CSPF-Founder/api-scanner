package db

import (
	"fmt"
	"testing"

	"github.com/ory/dockertest/v3"
)

func TestConnectDB(t *testing.T) {
	// Create a new Docker pool
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Could not construct pool: %s", err)
	}

	// uses pool to try to connect to Docker
	err = pool.Client.Ping()
	if err != nil {
		t.Fatalf("Could not connect to Docker: %s", err)
	}

	// Start a MySQL container
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mariadb",
		Tag:        "11.2.2",
		Env:        []string{"MYSQL_ROOT_PASSWORD=secret"},
	})
	if err != nil {
		t.Fatalf("Could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		dbURI := fmt.Sprintf("root:secret@(localhost:%s)/mysql", resource.GetPort("3306/tcp"))
		db, err := ConnectDB(dbURI)
		if err != nil {
			return err
		}
		defer db.Close()

		return nil
	}); err != nil {
		t.Fatalf("Could not connect to docker: %s", err)
	}

	// Clean up
	if err := pool.Purge(resource); err != nil {
		t.Fatalf("Could not purge resource: %s", err)
	}
}
