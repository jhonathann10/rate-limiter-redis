package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/jhonathann10/rate-limiter-redis/configs"
	middleware2 "github.com/jhonathann10/rate-limiter-redis/internal/infra/api/middleware"
	"github.com/jhonathann10/rate-limiter-redis/internal/infra/api/web"
	"github.com/jhonathann10/rate-limiter-redis/internal/infra/database"
	"github.com/jhonathann10/rate-limiter-redis/internal/infra/usecase"
	"github.com/redis/go-redis/v9"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(
		&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
	)

	redisDB := database.NewUserRepository(client)

	userUseCase := usecase.NewUserUseCase(redisDB)
	handler := web.NewHandler(userUseCase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.TokenAuth))
	r.Use(middleware.WithValue("JwtExperesIn", config.JwtExperesIn))

	rateLimiter := middleware2.NewRateLimiter(
		time.Duration(config.RateLimitTokenTime)*time.Second, time.Duration(config.RateLimitIpTime)*time.Second,
		config.RateLimitToken, config.RateLimitIp,
	)

	r.Route("/user", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Use(rateLimiter.Limit)
		r.Post("/{username}", handler.CreateUserController)
		r.Get("/", handler.GetUserController)
	})
	r.Post("/generate_token", handler.GetJWT)

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", r)
}
