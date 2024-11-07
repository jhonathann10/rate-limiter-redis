package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jhonathann10/rate-limiter-redis/internal/infra/api/web"
	"github.com/jhonathann10/rate-limiter-redis/internal/infra/database"
	"github.com/jhonathann10/rate-limiter-redis/internal/infra/usecase"
	"github.com/redis/go-redis/v9"
)

func main() {
	// posso criar um endpoint para gerar um token e passar no header
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

	router := gin.Default()
	router.POST("/user/:username", handler.CreateUserController)
	router.GET("/user", handler.GetUserController)

	router.Run(":8080")
}
