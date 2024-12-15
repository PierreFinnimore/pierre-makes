package app

import (
	"log/slog"
	"pierre/app/handlers"
	"pierre/app/views/errors"
	"pierre/plugins/auth"

	"pierre/kit"
	"pierre/kit/middleware"

	"github.com/go-chi/chi/v5"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// Define your global middleware
func InitializeMiddleware(router *chi.Mux) {
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Compress(5))
	router.Use(middleware.WithRequest)
}

// Define your routes in here
func InitializeRoutes(router *chi.Mux) {
	// Authentication plugin:
	auth.InitializeRoutes(router)

	authConfig := kit.AuthenticationConfig{
		AuthFunc:    auth.AuthenticateUser,
		RedirectURL: "/login",
	}

	// Routes that "might" have an authenticated user
	router.Group(func(app chi.Router) {
		app.Use(kit.WithAuthentication(authConfig, false)) // strict set to false

		// Routes
		app.Get("/", kit.Handler(handlers.HandleLandingIndex))
		app.Get("/art", kit.Handler(handlers.HandleArtIndex))
		app.Get("/thoughts", kit.Handler(handlers.HandleThoughtsIndex))
		app.Get("/tools", kit.Handler(handlers.HandleToolsIndex))
		app.Get("/robots.txt", kit.Handler(handlers.HandleRobotsTxt))

		app.Get("/poetry", kit.Handler(handlers.HandleConsequencesIndex))

		app.Post("/poetry/room", kit.Handler(handlers.HandleCreateRoom))
		app.Get("/poetry/room/{code}", kit.Handler(handlers.HandleGetRoomIndex))
		app.Post("/poetry/room/join", kit.Handler(handlers.HandleJoinRoom))
		app.Get("/poetry/room/{code}/poem", kit.Handler(handlers.HandleGetAvailablePoem))
		app.Post("/poetry/room/{code}/poem/{poemid}", kit.Handler(handlers.HandlePoemSubmission))

		app.Post("/poetry/poet", kit.Handler(handlers.HandleChoosePoet))
		app.Get("/poetry/poet", kit.Handler(handlers.HandleGetPoet))
		app.Get("/poetry/poet/auth", kit.Handler(handlers.HandleGetPoetAuth))

	})

	// Authenticated routes
	//
	// Routes that "must" have an authenticated user or else they
	// will be redirected to the configured redirectURL, set in the
	// AuthenticationConfig.
	router.Group(func(app chi.Router) {
		app.Use(kit.WithAuthentication(authConfig, true)) // strict set to true

		// Routes
		// app.Get("/path", kit.Handler(myHandler.HandleIndex))
	})
}

// NotFoundHandler that will be called when the requested path could
// not be found.
func NotFoundHandler(kit *kit.Kit) error {
	return kit.Render(errors.Error404())
}

// ErrorHandler that will be called on errors return from application handlers.
func ErrorHandler(kit *kit.Kit, err error) {
	slog.Error("internal server error", "err", err.Error(), "path", kit.Request.URL.Path)
	kit.Render(errors.Error500())
}
