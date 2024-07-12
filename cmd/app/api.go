package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/odev-swe/todoapp/docs"
	"github.com/odev-swe/todoapp/internal/handlers"
	"github.com/odev-swe/todoapp/internal/services"
	"github.com/odev-swe/todoapp/internal/store"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

//	@title			TodoApp API
//	@version		1.0
//	@description	This is a simple todo app API
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:3000
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				Description for what is this security definition being used

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func (app *application) Start() {
	// start the server
	// router
	router := chi.NewRouter()

	// swagger for chi

	// middlewares
	router.Use(app.LogMiddleware)
	router.Use(app.RateLimitMiddleware)

	router.Get("/swagger/*", httpSwagger.WrapHandler)

	// prefix
	router.Route("/api/v1", func(r chi.Router) {
		// ping route
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})

		// auth routes
		r.Route("/auth", func(r chi.Router) {
			authStore := store.NewAuthStore(app.db, app.config.JwtSecret)
			authService := services.NewAuthService(authStore)
			authHandler := handlers.NewAuthHandler(authService)
			authHandler.RegisterRoute(r)
		})

		// todos routes
		r.Route("/todos", func(r chi.Router) {
			r.Use(app.AuthMiddleware)
			todoStore := store.NewTodosStore(app.db)
			todoService := services.NewTodosService(todoStore)
			todoHandler := handlers.NewTodosHandler(todoService)
			todoHandler.RegisterRoute(r)

		})
	})

	zap.L().Info("Server started at", zap.String("port", app.config.Port))

	svr := &http.Server{
		Handler:  router,
		Addr:     ":" + app.config.Port,
		ErrorLog: zap.NewStdLog(zap.L()),
	}

	err := svr.ListenAndServe()

	if err != nil {
		zap.L().Fatal("Server error", zap.Error(err))
	}
}
