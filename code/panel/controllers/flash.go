package controllers

import (
	"net/http"

	"github.com/CSPF-Founder/api-scanner/code/panel/internal/sessions"
	"github.com/CSPF-Founder/api-scanner/code/panel/utils"
)

// Flash handles the rendering flash messages
func (app *App) Flash(_ http.ResponseWriter, r *http.Request, t string, m string, c bool) {
	app.session.AddFlash(r.Context(), sessions.SessionFlash{
		Type:     t,
		Message:  m,
		Closable: c,
	})
}

func (app *App) FlashAndGoBack(w http.ResponseWriter, r *http.Request, t string, message string) {
	app.Flash(w, r, t, message, true)
	utils.RedirectBack(w, r)
}
