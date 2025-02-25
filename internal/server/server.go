package server

import (
	"context"
	"net/http"
	"time"

	"guardian/api"
	"guardian/configs"
	guardianMiddleware "guardian/internal/middleware"
	"guardian/internal/mongodb"
	"guardian/internal/ratelimit"
	redisClient "guardian/internal/redis"
	"guardian/internal/setup"
	"guardian/utlis/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	swagger "github.com/swaggo/http-swagger"
	"golang.org/x/sync/errgroup"
)

func StartServer() error {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			writer.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(writer, request)
		})
	})

	setupRoutes(router)
	router.Get("/swagger/*", swagger.Handler(
		swagger.URL("/swagger/doc.json"),
	))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.GetLogger().Info("Server starting on port 8080")
		return http.ListenAndServe(":8080", router)
	})

	if err := g.Wait(); err != nil {
		logger.GetLogger().Errorf("Server error occurred:%v\n", err)
		return err
	}

	return nil
}

func setupRoutes(router *chi.Mux) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	authController := setup.InitializeAuthController(mongodb.Database)
	sendController := setup.InitializeSendHandlerController(mongodb.Database)

	addAuthRoutes(router, authController)
	setupRateLimiter(router)

	router.Group(func(protected chi.Router) {
		protected.Use(guardianMiddleware.VerifyJWT)
		addProtectedRoutes(protected, authController, sendController)
	})
}

func setupRateLimiter(router *chi.Mux) {
	if configs.GlobalConfig.EnableRateLimiter {
		router.Use(ratelimit.RateLimiterMiddleware(redisClient.Client))
	}
}

func addAuthRoutes(router *chi.Mux, authController *api.AuthController) {
	router.Route("/user", func(r chi.Router) {
		r.Post("/login", authController.Login)
		r.Post("/sign-up", authController.SignUp)
	})
}

func addProtectedRoutes(protected chi.Router, authController *api.AuthController,
	controller *api.SendHandlerController,
) {
	protected.Put("/user/update", authController.UpdateUser)
	protected.Patch("/user/activate", authController.ActivateUser)
	protected.Delete("/user/delete", authController.DeleteUser)

	protected.Post("/send", controller.SendHandler)
}
