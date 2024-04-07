package sessions

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/CSPF-Founder/api-scanner/code/panel/config"
	"github.com/CSPF-Founder/api-scanner/code/panel/models"
	"github.com/alexedwards/scs/v2"
)

// SetupSession configures and returns a session manager.

type SessionManager struct {
	*scs.SessionManager
}

// Flash is used to hold flash information for use in templates.
type SessionFlash struct {
	Type     string
	Message  string
	Closable bool
}

// Flash is used to hold flash information for use in templates.
const (
	FlashSuccess string = "success"
	FlashInfo    string = "info"
	FlashWarning string = "warning"
	FlashDanger  string = "danger"
)

const defaultFlashKey = "_flashes"

func init() {
	gob.Register(&models.User{})
	gob.Register([]SessionFlash{})
}

func SetupSession(config *config.Config) *SessionManager {

	scsSession := scs.New()
	scsSession.Lifetime = 24 * time.Hour
	scsSession.Cookie.Persist = true
	scsSession.Cookie.SameSite = http.SameSiteLaxMode
	scsSession.Cookie.HttpOnly = true
	if config.ServerConf.UseTLS {
		scsSession.Cookie.Secure = true
	} else {
		scsSession.Cookie.Secure = false
	}

	return &SessionManager{scsSession}
}

// Flashes retrieves flash messages from the session.
// An optional flash key can be provided, otherwise "_flash" is used by default.
func (session *SessionManager) Flashes(ctx context.Context, vars ...string) []SessionFlash {
	key := defaultFlashKey
	if len(vars) > 0 {
		key = vars[0]
	}

	flashes, _ := session.Pop(ctx, key).([]SessionFlash)
	return flashes
}

// AddFlash adds a flash message to the session.
// An optional flash key can be provided, otherwise "_flash" is used by default.
func (session *SessionManager) AddFlash(ctx context.Context, value SessionFlash, vars ...string) {
	key := defaultFlashKey
	if len(vars) > 0 {
		key = vars[0]
	}

	flashes := session.Flashes(ctx, key)
	flashes = append(flashes, value)
	session.Put(ctx, key, flashes)
}
