package models

const (
	// Roles for users
	RoleAdmin = "admin"
	RoleUser  = "customer"

	// PermissionViewObjects determines if a role can view standard
	// objects such as jobs
	PermissionViewObjects = "view_objects"
	// PermissionModifyObjects determines if a role can create and modify
	// standard objects.
	PermissionModifyObjects = "modify_objects"
	// PermissionModifySystem determines if a role can manage system-level
	// configuration.
	PermissionModifySystem = "modify_system"
)

// Role represents a role that can be assigned to a user.
type Role struct {
	ID          int64        `json:"-"`
	Keyword     string       `json:"keyword"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"-" gorm:"many2many:role_permissions;"`
}

// Permission represents a permission that can be assigned to a role.
type Permission struct {
	ID          int64  `json:"id"`
	Keyword     string `json:"keyword"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// GetRoleByKeyword returns a role that can be assigned to a user.
func GetRoleByKeyword(keyword string) (Role, error) {
	var role Role
	err := db.Where("keyword = ?", keyword).First(&role).Error

	return role, err
}

// HasPermission determines if a role has a given permission.
func (u *User) HasPermission(keyword string) (bool, error) {
	var perm []Permission

	err := db.Model(Role{ID: u.RoleID}).Where("keyword=?", keyword).Association("Permissions").Find(&perm)
	if err != nil {
		return false, err
	}

	// Gorm doesn't return an ErrRecordNotFound whe scanning into a slice, so
	// You can check the length directly to determine if permissions exist
	// Reference (ref github.com/go-gorm/gorm/issues/228)
	if len(perm) == 0 {
		return false, nil
	}
	return true, nil
}
