package controllers

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/CSPF-Founder/api-scanner/code/panel/config"
	ctx "github.com/CSPF-Founder/api-scanner/code/panel/context"
	"github.com/CSPF-Founder/api-scanner/code/panel/db"
	"github.com/CSPF-Founder/api-scanner/code/panel/frontend"
	"github.com/CSPF-Founder/api-scanner/code/panel/internal/httpclient"
	"github.com/CSPF-Founder/api-scanner/code/panel/internal/sessions"
	"github.com/CSPF-Founder/api-scanner/code/panel/logger"
	mid "github.com/CSPF-Founder/api-scanner/code/panel/middlewares"
	"github.com/CSPF-Founder/api-scanner/code/panel/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	server     *http.Server
	config     *config.Config
	logger     *logger.Logger
	session    *sessions.SessionManager
	httpClient httpclient.HttpClient
}

// Change Configuration accordingly
var defaultTLSConfig = &tls.Config{
	PreferServerCipherSuites: true,
	CurvePreferences: []tls.CurveID{
		tls.X25519,
		tls.CurveP256,
	},
	MinVersion: tls.VersionTLS12,
	CipherSuites: []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,

		// Kept for backwards compatibility with some clients
		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
	},
}

// ServerOption is a functional option that is used to configure the
type ServerOption func(*App)

// NewApp returns a new instance of the app with
// provided options applied.
func NewApp(config *config.Config, appLogger *logger.Logger, options ...ServerOption) *App {
	defaultServer := &http.Server{
		ReadTimeout:  45 * time.Second,
		WriteTimeout: 45 * time.Second,
		Addr:         config.ServerConf.ServerAddress,
	}

	sessionManager := sessions.SetupSession(config)
	app := &App{
		server:  defaultServer,
		config:  config,
		logger:  appLogger,
		session: sessionManager,
	}
	for _, opt := range options {
		opt(app)
	}

	app.httpClient = &http.Client{
		Timeout: time.Second * 30,
	}
	return app
}

func (app *App) SetupDB() {
	err := db.Setup(app.config, app.session)
	if err != nil {
		app.logger.Fatal("Error setting up models", err)
	}
}

// Start launches the server, listening on the configured address.
func (app *App) StartServer() {
	// Use Tls if configured
	if app.config.ServerConf.UseTLS {
		app.server.TLSConfig = defaultTLSConfig

		app.logger.Info("Creating new self-signed certificate")
		err := utils.CheckAndCreateSSL(app.config.ServerConf.CertPath, app.config.ServerConf.KeyPath)
		if err != nil {
			app.logger.Fatal("Error creating SSL Certificates: ", err)
			return
		}

		app.logger.Info("TLS Certificate Generation complete")

		app.logger.Info(fmt.Sprintf("Starting server at https://%s", app.config.ServerConf.ServerAddress))
		err = app.server.ListenAndServeTLS(app.config.ServerConf.CertPath, app.config.ServerConf.KeyPath)
		if err != nil {
			app.logger.Fatal("Error starting server: ", err)
		}
	}
	// If TLS isn't configured, just listen on HTTP
	app.logger.Info(fmt.Sprintf("Starting server at http://%s", app.config.ServerConf.ServerAddress))
	err := app.server.ListenAndServe()
	if err != nil {
		app.logger.Fatal("Error starting server: ", err)
	}
}

// Shutdown attempts to gracefully shutdown the server.
func (app *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return app.server.Shutdown(ctx)
}

func (app *App) SetupRoutes() {

	router := chi.NewRouter() // Initialize Chi router

	// r.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5))

	errHandlers := newErrorController(app)

	// CSRF Validation Middleware
	router.Use(mid.CSRFMiddleware(
		app.session,
		http.HandlerFunc(errHandlers.csrfErrorHandler),
		[]string{"/api"}),
	)

	sc := newScansController(app)
	usec := newUserSetupController(app)
	uc := newUserController(app)

	// Middlewares
	router.NotFound(errHandlers.NotFoundHandler)
	// Setup logging
	router.Use(mid.LoggingMiddleware(app.logger))

	router.Get("/", mid.Use(app.HandleHomePage))
	router.Mount("/scans", sc.registerRoutes())
	router.Mount("/users", uc.registerRoutes())
	router.Mount("/setup", usec.registerRoutes())

	// Embedded static file serving
	fileServer := http.FileServer(http.FS(frontend.FileSystem))
	router.Handle("/static/*", fileServer)

	// External static files that will be later mounted via docker
	externalFS := http.FileServer(http.Dir("./frontend/external"))
	router.Handle("/external/*", http.StripPrefix("/external/", externalFS))

	routeHandler := mid.Use(
		router.ServeHTTP,
		mid.GetContext(app.session),
		mid.ApplySecurityHeaders,
	)
	app.server.Handler = app.session.LoadAndSave(routeHandler)
}

func (c *App) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	user := ctx.Get(r, "user")
	if user == nil {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/scans/add", http.StatusSeeOther)
}
