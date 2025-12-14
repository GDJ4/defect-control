package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"defect-tracker/internal/domain"
	"defect-tracker/internal/pkg/storage"
	"defect-tracker/internal/service/defect"
	"defect-tracker/internal/transport/http/middleware"
)

type DefectHandler struct {
	service *defect.Service
	storage storage.Provider
}

func NewDefectHandler(service *defect.Service, storage storage.Provider) *DefectHandler {
	return &DefectHandler{service: service, storage: storage}
}

func (h *DefectHandler) Register(rg *gin.RouterGroup) {
	rg.GET("/defects", h.list)
	rg.POST("/defects", h.create)
	rg.GET("/defects/:id", h.get)
	rg.GET("/defects/:id/comments", h.listComments)
	rg.POST("/defects/:id/comments", h.addComment)
	rg.POST("/defects/:id/attachments", h.addAttachment)
	rg.PATCH("/defects/:id/status", h.updateStatus)
	rg.GET("/defects/:id/attachments/:attachmentId", h.downloadAttachment)
}

func (h *DefectHandler) list(c *gin.Context) {
	limit := parseLimit(c.DefaultQuery("limit", "20"))
	filter := domain.DefectFilter{
		Status:   c.Query("status"),
		Priority: c.Query("priority"),
		Project:  c.Query("projectId"),
		Limit:    limit,
	}

	items, err := h.service.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось загрузить дефекты"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"items": mapListResponse(items)})
}

func (h *DefectHandler) create(c *gin.Context) {
	var payload struct {
		ProjectID   string `json:"projectId"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		Severity    string `json:"severity"`
		AssigneeID  string `json:"assigneeId"`
		DueDate     string `json:"dueDate"`
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
		c.JSON(http.StatusForbidden, gin.H{"message": "Создавать дефекты может только менеджер"})
		return
	}

	due, err := parseDate(payload.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректная дата срока"})
		return
	}

	defectEntity, err := h.service.Create(c.Request.Context(), domain.DefectCreate{
		ProjectID:   payload.ProjectID,
		Title:       payload.Title,
		Description: payload.Description,
		Priority:    payload.Priority,
		Severity:    payload.Severity,
		AssigneeID:  payload.AssigneeID,
		DueDate:     due,
		CreatedBy:   user.ID,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, h.mapDefect(c, defectEntity))
}

func (h *DefectHandler) get(c *gin.Context) {
	defectEntity, err := h.service.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Дефект не найден"})
		return
	}
	c.JSON(http.StatusOK, h.mapDefect(c, defectEntity))
}

func (h *DefectHandler) listComments(c *gin.Context) {
	comments, err := h.service.ListComments(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось получить комментарии"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": mapComments(comments)})
}

func (h *DefectHandler) addComment(c *gin.Context) {
	var payload struct {
		Body string `json:"body"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Некорректные данные"})
		return
	}

	user, ok := middleware.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Необходима авторизация"})
		return
	}
	comment, err := h.service.AddComment(c.Request.Context(), domain.CommentCreate{
		DefectID: c.Param("id"),
		AuthorID: user.ID,
		Body:     payload.Body,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, mapComment(comment))
}

func (h *DefectHandler) addAttachment(c *gin.Context) {
	formFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Файл обязателен"})
		return
	}
	file, err := formFile.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Не удалось прочитать файл"})
		return
	}
	defer file.Close()

	storageKey, size, err := h.storage.Save(c.Request.Context(), file, formFile.Filename, formFile.Size, formFile.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось сохранить файл"})
		return
	}

	attachment, err := h.service.AddAttachment(c.Request.Context(), domain.AttachmentCreate{
		DefectID:    c.Param("id"),
		Filename:    formFile.Filename,
		ContentType: formFile.Header.Get("Content-Type"),
		SizeBytes:   size,
		StorageKey:  storageKey,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, h.mapSingleAttachment(c, c.Param("id"), attachment))
}

func (h *DefectHandler) updateStatus(c *gin.Context) {
	var payload struct {
		Status string `json:"status"`
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

	defect, err := h.service.UpdateStatus(c.Request.Context(), c.Param("id"), user, payload.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, h.mapDefect(c, defect))
}

func (h *DefectHandler) downloadAttachment(c *gin.Context) {
	attachment, err := h.service.GetAttachment(c.Request.Context(), c.Param("id"), c.Param("attachmentId"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Вложение не найдено"})
		return
	}

	local, ok := h.storage.(interface {
		PathFor(string) string
	})
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Загрузка через API недоступна для выбранного хранилища"})
		return
	}

	fullPath := local.PathFor(attachment.StorageKey)
	c.FileAttachment(fullPath, attachment.Filename)
}

func parseLimit(value string) int {
	limit, err := strconv.Atoi(value)
	if err != nil {
		return 20
	}
	return limit
}

func parseDate(value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	t, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func mapListResponse(items []domain.DefectListItem) []gin.H {
	result := make([]gin.H, 0, len(items))
	for _, item := range items {
		var due *string
		if item.DueDate != nil {
			formatted := item.DueDate.Format(time.DateOnly)
			due = &formatted
		}
		result = append(result, gin.H{
			"id":         item.ID,
			"projectId":  item.ProjectID,
			"project":    item.ProjectName,
			"title":      item.Title,
			"priority":   item.Priority,
			"status":     item.Status,
			"assigneeId": item.AssigneeID,
			"assignee":   item.AssigneeName,
			"dueDate":    due,
			"updatedAt":  item.UpdatedAt,
		})
	}
	return result
}

func (h *DefectHandler) mapDefect(c *gin.Context, d domain.Defect) gin.H {
	var due *string
	if d.DueDate != nil {
		formatted := d.DueDate.Format(time.DateOnly)
		due = &formatted
	}

	attachments := h.mapAttachments(c, d)

	return gin.H{
		"id":          d.ID,
		"projectId":   d.ProjectID,
		"project":     d.ProjectName,
		"title":       d.Title,
		"description": d.Description,
		"priority":    d.Priority,
		"severity":    d.Severity,
		"status":      d.Status,
		"assigneeId":  d.AssigneeID,
		"assignee":    d.Assignee,
		"dueDate":     due,
		"createdBy":   d.CreatedBy,
		"createdAt":   d.CreatedAt,
		"updatedAt":   d.UpdatedAt,
		"attachments": attachments,
		"comments":    mapComments(d.Comments),
	}
}

func mapComments(comments []domain.Comment) []gin.H {
	result := make([]gin.H, 0, len(comments))
	for _, comment := range comments {
		result = append(result, mapComment(comment))
	}
	return result
}

func mapComment(comment domain.Comment) gin.H {
	return gin.H{
		"id":        comment.ID,
		"authorId":  comment.AuthorID,
		"author":    comment.AuthorName,
		"body":      comment.Body,
		"createdAt": comment.CreatedAt,
	}
}

func (h *DefectHandler) mapAttachments(c *gin.Context, defect domain.Defect) []gin.H {
	result := make([]gin.H, 0, len(defect.Attachments))
	for _, att := range defect.Attachments {
		result = append(result, h.mapSingleAttachment(c, defect.ID, att))
	}
	return result
}

func (h *DefectHandler) mapSingleAttachment(c *gin.Context, defectID string, att domain.Attachment) gin.H {
	url := h.buildDownloadURL(c, defectID, att)
	return gin.H{
		"id":          att.ID,
		"filename":    att.Filename,
		"contentType": att.ContentType,
		"sizeBytes":   att.SizeBytes,
		"storageKey":  att.StorageKey,
		"uploadedAt":  att.UploadedAt,
		"downloadUrl": url,
	}
}

func (h *DefectHandler) buildDownloadURL(c *gin.Context, defectID string, att domain.Attachment) string {
	if url, err := h.storage.Presign(c.Request.Context(), att.StorageKey); err == nil && url != "" {
		return url
	}
	return fmt.Sprintf("/api/v1/defects/%s/attachments/%s", defectID, att.ID)
}
