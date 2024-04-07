package sessions

import (
	"context"

	"github.com/CSPF-Founder/api-scanner/code/panel/utils"
)

const (
	TokenName = "csrf_token"
)

func (session *SessionManager) GenerateCSRF(ctx context.Context) {
	token := utils.GetRandomHexString(32)
	session.Put(ctx, TokenName, token)
}

func (session *SessionManager) GetCSRF(ctx context.Context) string {
	token := session.GetString(ctx, TokenName)
	return token
}

func (session *SessionManager) ValidateCSRF(ctx context.Context, inputToken string) bool {
	sessionToken := session.GetString(ctx, TokenName)
	if sessionToken == "" {
		return false
	}

	return sessionToken == inputToken
}

// func ValidateJSONCSRF(r *http.Request, session *sessions.SessionManager) bool {

// 	var requestData struct {
// 		CSRFToken string `json:"csrf_token"`
// 	}

// 	err := json.NewDecoder(r.Body).Decode(&requestData)
// 	if err != nil {
// 		return false
// 	}

// 	sessionToken := session.GetString(r.Context(), TokenName)
// 	if sessionToken == "" || sessionToken != requestData.CSRFToken {
// 		return false
// 	}

// 	return true
// }
