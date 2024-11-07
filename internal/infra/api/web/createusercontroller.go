package web

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jhonathann10/rate-limiter-redis/internal/infra/usecase"
)

type Handler struct {
	userUseCase usecase.UserUseCaseInterface
}

func NewHandler(userUseCase usecase.UserUseCaseInterface) *Handler {
	return &Handler{
		userUseCase: userUseCase,
	}
}

var ctx = context.Background()

func (h *Handler) CreateUserController(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	err := h.userUseCase.SaveUser(ctx, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user created"})
}

func (h *Handler) GetUserController(c *gin.Context) {
	user, err := h.userUseCase.GetUser(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, user)
}
