package models

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// func (ctx *testContext) mockGetRoleWithPermissions(role Role, permissions []Permission) {
// 	// Mocking the role query
// 	roleRow := sqlmock.NewRows([]string{"id", "keyword", "name", "description"}).
// 		AddRow(role.ID, role.Keyword, role.Name, role.Description)
// 	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `roles` WHERE keyword = ?")).
// 		WithArgs(role.Keyword).
// 		WillReturnRows(roleRow)

// 	// Mocking the permissions query with JOIN
// 	for _, perm := range permissions {
// 		permRows := sqlmock.NewRows([]string{"id", "keyword", "name", "description"}).
// 			AddRow(perm.ID, perm.Keyword, perm.Name, perm.Description)
// 		ctx.mock.ExpectQuery(regexp.QuoteMeta(
// 			"SELECT `permissions`.`id`,`permissions`.`keyword`,`permissions`.`name`,`permissions`.`description` FROM `permissions` JOIN `role_permissions` ON `role_permissions`.`permission_id` = `permissions`.`id` AND `role_permissions`.`role_id` = ? WHERE `permissions`.`keyword` = ?")).
// 			WithArgs(role.ID, perm.Keyword).
// 			WillReturnRows(permRows)
// 	}
// }

// func TestHasPermission(t *testing.T) {
// 	ctx := setupTest(t)
// 	adminRole := Role{
// 		ID:      1,
// 		Keyword: RoleAdmin,
// 		Name:    "Administrator",
// 		Permissions: []Permission{
// 			{Keyword: PermissionViewObjects},
// 			{Keyword: PermissionModifyObjects},
// 			{Keyword: PermissionModifySystem},
// 		},
// 	}

// 	ctx.mockGetRoleWithPermissions(adminRole, adminRole.Permissions)

// 	user := User{
// 		ID:       1,
// 		Username: "admin",
// 		RoleID:   adminRole.ID,
// 	}

// 	// Check if the user has specific permissions
// 	for _, perm := range adminRole.Permissions {
// 		hasPerm, err := user.HasPermission(perm.Keyword)
// 		if err != nil {
// 			t.Fatalf("error checking permission '%s': %v", perm.Keyword, err)
// 		}
// 		if !hasPerm {
// 			t.Errorf("expected user to have permission '%s', but they did not", perm.Keyword)
// 		}
// 	}
// }

func TestGetRoleByKeyword(t *testing.T) {
	ctx := setupTest(t)

	// Define roles
	roles := []Role{
		{ID: 1, Keyword: RoleAdmin},
		{ID: 2, Keyword: RoleUser},
	}

	// Mock database responses for each role
	for _, role := range roles {
		roleRow := sqlmock.NewRows([]string{"id", "keyword"}).
			AddRow(role.ID, role.Keyword)

		ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `roles` WHERE keyword = ?")).
			WithArgs(role.Keyword).
			WillReturnRows(roleRow)
	}

	// Test for each role
	for _, role := range roles {
		got, err := GetRoleByKeyword(role.Keyword)
		if err != nil {
			t.Fatalf("unexpected error when getting role by keyword '%s': %v", role.Keyword, err)
		}
		if got.Keyword != role.Keyword {
			t.Errorf("expected role keyword %s, got %s", role.Keyword, got.Keyword)
		}
	}

}

func TestGetNonExistentRoleByKeyword(t *testing.T) {
	ctx := setupTest(t)

	// Mock the database response for a non-existent role
	ctx.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `roles` WHERE keyword = ?")).
		WithArgs("bogus").
		WillReturnError(fmt.Errorf("role not found"))

	// Test for a non-existent role
	_, err := GetRoleByKeyword("bogus")
	if err == nil {
		t.Error("expected error for non-existent role keyword, got nil")
	}
}
