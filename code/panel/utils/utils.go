package utils

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

func CheckAllParamsExist(r *http.Request, requiredParams []string) bool {
	for _, param := range requiredParams {
		if r.FormValue(param) == "" {
			// Parameter is missing, handle the error
			return false
		}
	}

	// All parameters exist
	return true
}

func GetRandomHexString(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}

	randomHex := hex.EncodeToString(randomBytes)
	return randomHex
}

func GetRandomString(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return ""
	}

	randomHex := hex.EncodeToString(randomBytes)
	return randomHex
}
