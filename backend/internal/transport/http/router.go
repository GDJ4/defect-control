package transporthttp

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"defect-tracker/internal/transport/http/handlers"
	"defect-tracker/internal/transport/http/middleware"
)

// NewRouter wires up the HTTP routes and shared middleware.
func NewRouter(
	appName string,
	authHandler *handlers.AuthHandler,
	authMW *middleware.AuthMiddleware,
	defectHandler *handlers.DefectHandler,
	projectHandler *handlers.ProjectHandler,
) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"app":   appName,
			"state": "ok",
		})
	})

	api := router.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})

		authHandler.RegisterPublic(api)

		secured := api.Group("/")
		secured.Use(authMW.RequireAuth())

		authHandler.RegisterProtected(secured)
		projectHandler.Register(secured)
		defectHandler.Register(secured)
	}

	return router
}
