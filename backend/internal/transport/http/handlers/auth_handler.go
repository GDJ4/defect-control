package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"defect-tracker/internal/domain"
	"defect-tracker/internal/pkg/auth"
	"defect-tracker/internal/service/token"
	"defect-tracker/internal/service/user"
	"defect-tracker/internal/transport/http/middleware"
)

type AuthHandler struct {
	users   *user.Service
	tokens  *token.Service
	manager *auth.Manager
}

func NewAuthHandler(users *user.Service, tokens *token.Service, manager *auth.Manager) *AuthHandler {
	return &AuthHandler{users: users, tokens: tokens, manager: manager}
}

func (h *AuthHandler) RegisterPublic(rg *gin.RouterGroup) {
	rg.POST("/auth/login", h.login)
	rg.POST("/auth/register", h.register)
	rg.POST("/auth/refresh", h.refresh)
}

func (h *AuthHandler) RegisterProtected(rg *gin.RouterGroup) {
	rg.POST("/auth/logout", h.logout)
	rg.POST("/auth/password", h.changePassword)
}

func (h *AuthHandler) login(c *gin.Context) {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректный формат данных"})
		return
	}

	userEntity, err := h.users.Authenticate(c.Request.Context(), payload.Email, payload.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Неверный логин или пароль"})
		return
	}

	refresh, err := h.tokens.Issue(c.Request.Context(), userEntity.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось выпустить refresh token"})
		return
	}

	h.respondWithTokens(c, userEntity, refresh.Token, refresh.ExpiresAt)
}

func (h *AuthHandler) register(c *gin.Context) {
	var payload struct {
		Email    string `json:"email"`
		FullName string `json:"fullName"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректный формат данных"})
		return
	}

	userEntity, err := h.users.Register(c.Request.Context(), domain.UserRegister{
		Email:    payload.Email,
		FullName: payload.FullName,
		Password: payload.Password,
		Role:     payload.Role,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	refresh, err := h.tokens.Issue(c.Request.Context(), userEntity.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось выпустить refresh token"})
		return
	}

	h.respondWithTokens(c, userEntity, refresh.Token, refresh.ExpiresAt)
}

func (h *AuthHandler) refresh(c *gin.Context) {
	var payload struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil || payload.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Токен обязателен"})
		return
	}

	newRefresh, err := h.tokens.Rotate(c.Request.Context(), payload.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Недействительный refresh token"})
		return
	}

	userEntity, err := h.users.GetByID(c.Request.Context(), newRefresh.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Пользователь не найден"})
		return
	}

	h.respondWithTokens(c, userEntity, newRefresh.Token, newRefresh.ExpiresAt)
}

func (h *AuthHandler) logout(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Необходима авторизация"})
		return
	}

	var payload struct {
		RefreshToken string `json:"refreshToken"`
	}
	_ = c.ShouldBindJSON(&payload)

	if payload.RefreshToken != "" {
		_ = h.tokens.Revoke(c.Request.Context(), payload.RefreshToken)
	} else {
		_ = h.tokens.RevokeUserTokens(c.Request.Context(), user.ID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Выход выполнен"})
}

func (h *AuthHandler) changePassword(c *gin.Context) {
	user, ok := middleware.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Необходима авторизация"})
		return
	}

	var payload struct {
		CurrentPassword string `json:"currentPassword"`
		NewPassword     string `json:"newPassword"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректные данные"})
		return
	}

	if err := h.users.ChangePassword(c.Request.Context(), user.ID, payload.CurrentPassword, payload.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	_ = h.tokens.RevokeUserTokens(c.Request.Context(), user.ID)
	c.JSON(http.StatusOK, gin.H{"message": "Пароль обновлён"})
}

func (h *AuthHandler) respondWithTokens(c *gin.Context, user domain.User, refreshToken string, refreshExpires time.Time) {
	accessToken, err := h.manager.Generate(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось выпустить access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":      accessToken,
		"refreshToken":     refreshToken,
		"refreshExpiresAt": refreshExpires,
		"user": gin.H{
			"id":       user.ID,
			"email":    user.Email,
			"fullName": user.FullName,
			"role":     user.Role,
		},
	})
}
