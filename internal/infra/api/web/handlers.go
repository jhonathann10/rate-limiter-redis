package web

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/jhonathann10/rate-limiter-redis/internal/infra/usecase"
)

type Message struct {
	Message string `json:"message"`
}

type Handler struct {
	userUseCase usecase.UserUseCaseInterface
}

func NewHandler(userUseCase usecase.UserUseCaseInterface) *Handler {
	return &Handler{
		userUseCase: userUseCase,
	}
}

var ctx = context.Background()

func (h *Handler) CreateUserController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	username := chi.URLParam(r, "username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		msg := Message{Message: "username is required"}
		json.NewEncoder(w).Encode(msg)
		return
	}

	err := h.userUseCase.SaveUser(ctx, username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := Message{Message: err.Message}
		json.NewEncoder(w).Encode(msg)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Message{Message: "user created"})
}

func (h *Handler) GetUserController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := h.userUseCase.GetUser(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := Message{Message: err.Message}
		json.NewEncoder(w).Encode(msg)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwtToken := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExperesIn").(int)

	_, tokenString, _ := jwtToken.Encode(map[string]interface{}{
		"sub": "user",
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})
	accessToken := map[string]string{"access_token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}
