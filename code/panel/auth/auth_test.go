package auth

import (
	"testing"
)

func TestCheckPasswordPolicy(t *testing.T) {
	testCases := []struct {
		name     string
		password string
		wantErr  error
	}{
		{"TooShort", "short", ErrPasswordTooShort},
		{"ValidPassword", "validPassword123", nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := CheckPasswordPolicy(tc.password); got != tc.wantErr {
				t.Errorf("CheckPasswordPolicy(%q) = %v, want %v", tc.password, got, tc.wantErr)
			}
		})
	}

}

func TestValidatePasswordChange(t *testing.T) {

	testCases := []struct {
		name            string
		currentPassword string
		newPassword     string
		confirmPassword string
		wantErr         error
	}{
		{"MismatchPassword", "currentPassword", "validNewPassword", "invalid", ErrPasswordMismatch},
		{"ReusedPassword", "currentPassword", "currentPassword", "currentPassword", ErrReusedPassword},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			currentHash, err := GeneratePasswordHash(tc.currentPassword)
			if err != nil {
				t.Fatalf("unexpected error generating password hash: %v", err)
			}

			_, got := ValidatePasswordChange(currentHash, tc.newPassword, tc.confirmPassword)
			if got != tc.wantErr {
				t.Errorf("ValidatePasswordChange() = %v, want %v", got, tc.wantErr)
			}
		})
	}
}
