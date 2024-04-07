package models

import (
	"regexp"
	"testing"

	"github.com/CSPF-Founder/api-scanner/code/panel/auth"
	"github.com/DATA-DOG/go-sqlmock"
)

// TODO: Avoid duplication of code between controller and models test

// SQL mocking helper functions
func (ctx *testContext) mockUserCount(count int) {
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(count)
	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users`")).
		WillReturnRows(rows)
}

func (ctx *testContext) mockGetByUsername(user User) {
	userRow := sqlmock.NewRows([]string{"id", "username", "password"}).
		AddRow(user.ID, user.Username, user.Password)

	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? ORDER BY `users`.`id` LIMIT 1")).
		WithArgs(user.Username).
		WillReturnRows(userRow)
}

func TestValidUserCount(t *testing.T) {
	ctx := setupTest(t)
	ctx.mockUserCount(1)
	var expected int64 = 1
	got, err := GetNumberOfUsers()
	if err != nil {
		t.Fatalf("error getting number of users: %v", err)
	}

	if got != expected {
		t.Fatalf("invalid number of users. expected %d got %d", expected, got)
	}

}

func TestInvalidUserCount(t *testing.T) {
	ctx := setupTest(t)
	ctx.mockUserCount(0)
	var expected int64 = 0
	got, err := GetNumberOfUsers()
	if err != nil {
		t.Fatalf("error getting number of users: %v", err)
	}

	if got != expected {
		t.Fatalf("invalid number of users. expected %d got %d", expected, got)
	}

}

func TestGetByUsername(t *testing.T) {
	ctx := setupTest(t)
	testUser := User{
		ID:       1,
		Username: "test",
		Email:    "test@example.com",
	}

	hash, err := auth.GeneratePasswordHash("test")
	if err != nil {
		t.Fatalf("error generating password hash: %v", err)
	}
	testUser.Password = hash

	ctx.mockGetByUsername(testUser)

	got, err := GetUserByUsername(testUser.Username)
	if err != nil {
		t.Fatalf("error getting user by username: %v", err)
	}

	if got.Username != testUser.Username {
		t.Fatalf("invalid username. expected %s got %s", testUser.Username, got.Username)
	}

}
