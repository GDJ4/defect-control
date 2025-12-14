package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"defect-tracker/internal/domain"
	"defect-tracker/internal/service/project"
	"defect-tracker/internal/transport/http/middleware"
)

type ProjectHandler struct {
	service *project.Service
}

func NewProjectHandler(service *project.Service) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/projects", h.list)
	rg.POST("/projects", h.create)
}

func (h *ProjectHandler) list(c *gin.Context) {
	projects, err := h.service.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось получить проекты"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": mapProjects(projects)})
}

func (h *ProjectHandler) create(c *gin.Context) {
	var payload struct {
		Name        string `json:"name"`
		Stage       string `json:"stage"`
		Description string `json:"description"`
		StartDate   string `json:"startDate"`
		EndDate     string `json:"endDate"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректный формат данных"})
		return
	}

	user, ok := middleware.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Необходима авторизация"})
		return
	}

	if user.Role != "manager" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Создавать проекты может только менеджер"})
		return
	}

	start, err := parseDateValue(payload.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректная дата начала"})
		return
	}
	end, err := parseDateValue(payload.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректная дата окончания"})
		return
	}

	projectEntity, err := h.service.Create(c.Request.Context(), domain.ProjectCreate{
		Name:        payload.Name,
		Stage:       payload.Stage,
		Description: payload.Description,
		StartDate:   start,
		EndDate:     end,
		CreatedBy:   user.ID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapProject(projectEntity))
}

func mapProjects(projects []domain.Project) []gin.H {
	result := make([]gin.H, 0, len(projects))
	for _, project := range projects {
		result = append(result, mapProject(project))
	}
	return result
}

func mapProject(p domain.Project) gin.H {
	var start, end *string
	if p.StartDate != nil {
		v := p.StartDate.Format(time.DateOnly)
		start = &v
	}
	if p.EndDate != nil {
		v := p.EndDate.Format(time.DateOnly)
		end = &v
	}
	return gin.H{
		"id":          p.ID,
		"name":        p.Name,
		"stage":       p.Stage,
		"description": p.Description,
		"startDate":   start,
		"endDate":     end,
		"createdBy":   p.CreatedBy,
		"createdAt":   p.CreatedAt,
	}
}

func parseDateValue(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	t, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
