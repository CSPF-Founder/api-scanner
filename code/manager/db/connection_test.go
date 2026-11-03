package db

import (
	"fmt"
	"testing"

	"github.com/ory/dockertest/v4"
)

func TestConnectDB(t *testing.T) {
	// The pool removes its containers when the test finishes.
	pool := dockertest.NewPoolT(t, "")

	// Start a MySQL container
	resource := pool.RunT(t, "mariadb",
		dockertest.WithTag("11.2.2"),
		dockertest.WithEnv([]string{"MYSQL_ROOT_PASSWORD=secret"}),
		dockertest.WithoutReuse(),
	)

	if err := pool.Retry(t.Context(), 0, func() error {
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
}
