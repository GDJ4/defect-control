package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"defect-tracker/internal/domain"
	"defect-tracker/internal/pkg/auth"
	"defect-tracker/internal/service/user"
)

const userContextKey = "currentUser"

type AuthMiddleware struct {
	manager *auth.Manager
	userSrv *user.Service
}

func NewAuthMiddleware(manager *auth.Manager, userSrv *user.Service) *AuthMiddleware {
	return &AuthMiddleware{
		manager: manager,
		userSrv: userSrv,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := parseBearer(c.GetHeader("Authorization"))
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Требуется авторизация"})
			return
		}

		claims, err := m.manager.Parse(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Недействительный токен"})
			return
		}

		userEntity, err := m.userSrv.GetByID(c.Request.Context(), claims.Subject)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Пользователь не найден"})
			return
		}

		c.Set(userContextKey, userEntity)
		c.Next()
	}
}

func (m *AuthMiddleware) RequireRoles(roles ...string) gin.HandlerFunc {
	roleSet := make(map[string]struct{}, len(roles))
	for _, role := range roles {
		roleSet[role] = struct{}{}
	}

	return func(c *gin.Context) {
		user, ok := CurrentUser(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Требуется авторизация"})
			return
		}

		if len(roleSet) == 0 {
			c.Next()
			return
		}

		if _, allowed := roleSet[user.Role]; !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Недостаточно прав"})
			return
		}

		c.Next()
	}
}

func CurrentUser(c *gin.Context) (domain.User, bool) {
	user, ok := c.Get(userContextKey)
	if !ok {
		return domain.User{}, false
	}
	u, ok := user.(domain.User)
	return u, ok
}

func parseBearer(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}
	if strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}
