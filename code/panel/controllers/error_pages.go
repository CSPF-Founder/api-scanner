package controllers

import (
	"net/http"

	"github.com/CSPF-Founder/api-scanner/code/panel/views"
)

type errorController struct {
	*App
}

func newErrorController(a *App) *errorController {
	return &errorController{a}
}

func (c *errorController) ForbiddenHandler(w http.ResponseWriter, r *http.Request) {
	isAjaxRequest := r.Header.Get("X-Requested-With") == "XMLHttpRequest"

	if isAjaxRequest {
		c.SendJSONError(w, "Access Denied", http.StatusForbidden)
		return
	}
	c.Flash(w, r, "danger", "You do not have permission to access this page.", true)
	templateData := views.NewTemplateData(c.config, c.session, r)
	templateData.Title = "Access Denied"
	templateData.HideHeaderAndFooter = true
	w.WriteHeader(http.StatusForbidden)
	if err := views.RenderTemplate(w, "errors/403", templateData); err != nil {
		c.logger.Error("Error rendering template: ", err)
	}
}

func (c *errorController) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	templateData := views.NewTemplateData(c.config, c.session, r)
	templateData.Title = "404 Not Found"
	templateData.HideHeaderAndFooter = true
	w.WriteHeader(http.StatusNotFound)
	if err := views.RenderTemplate(w, "errors/404", templateData); err != nil {
		c.logger.Error("Error rendering template: ", err)
	}
}

func (c *errorController) csrfErrorHandler(w http.ResponseWriter, r *http.Request) {
	isAjaxRequest := r.Header.Get("X-Requested-With") == "XMLHttpRequest"

	if isAjaxRequest {
		c.SendJSONError(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	c.Flash(w, r, "danger", "Invalid CSRF token", true)
	templateData := views.NewTemplateData(c.config, c.session, r)
	templateData.Title = "Access Denied"
	templateData.HideHeaderAndFooter = true

	w.WriteHeader(http.StatusForbidden)
	if err := views.RenderTemplate(w, "errors/403", templateData); err != nil {
		c.logger.Error("Error rendering template: ", err)
	}
}
