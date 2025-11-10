package handlers

import (
	"net/http"
	"strconv"
	"time"
	"todo-api/internal/database"
	"todo-api/internal/dto"
	"todo-api/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateTodo(c *gin.Context){
	var req dto.CreateTodoRequest

	if err := c.ShouldBindJSON(&req); err != nil{
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	} 

	userID, _ := c.Get("user_id")
	
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	todo := models.Todo{
		UserID: userID.(uint),
		Title: req.Title,
		Description: req.Description,
		Priority: req.Priority,
		Category: req.Category,
		DueDate: req.DueDate,
	}
	if err := database.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to Create Todo",
		})
		return
	}
	c.JSON(http.StatusCreated, todo)

}

func GetTodos(c *gin.Context){
	userID, _ := c.Get("user_id")

	// Query parameters for filtering
	completed := c.Query("completed")
	priority := c.Query("priority")
	category := c.Query("category")
	search := c.Query("search")

	// Pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page","1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "1"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := database.DB.Where("user_id = ?", userID)

	if completed != "" {
		query = query.Where("completed = ?", completed == "true")
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	var total int64
	query.Model(&models.Todo{}).Count(&total)

	var todos []models.Todo
	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&todos).Error; err != nil{
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to Fetch Todos",
		})
		return
	}
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"todos": todos,
		"pagination" : gin.H{
			"page": page,
			"limit": limit,
			"total": total,
			"total_pages": totalPages,
		},
	})
}

func GetTodo(c *gin.Context){
	todoID := c.Param("id")
	userID, _ := c.Get("user_id")

	var todo models.Todo

	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Todo Not Found",
		})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func UpdateTodo(c *gin.Context){
	todoID := c.Param("id")
	userID, _ := c.Get("user_id")

	var req dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var todo models.Todo

	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "Todo Not Found",
		})
		return
	}
	if req.Title != "" {
		todo.Title = req.Title
	}
	if req.Description != nil {
		todo.Description = *req.Description
	}
	if req.Priority != "" {
		todo.Priority = req.Priority
	}
	if req.Category != nil {
		todo.Category = *req.Category
	}
	if req.DueDate != nil {
		todo.DueDate = req.DueDate
	}
	if req.Completed != nil {
		todo.Completed = *req.Completed
	}
	if err := database.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to Update Todo",
		})
		return
	}
	c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context){
	todoID := c.Param("id")
	userID, _ := c.Get("user_id")

	var todo models.Todo

	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Todo Not Found",
		})
		return
	}
	if err := database.DB.Delete(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to Delete Todo",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todo Deleted Successfully",
	})
}

func ToggleComplete(c *gin.Context){
	todoID := c.Param("id")
	userID, _ := c.Get("user_id")

	var todo models.Todo
	if err := database.DB.Where("id = ? AND user_id = ?", todoID, userID).First(&todo).Error; err != nil {
	c.JSON(http.StatusNotFound, dto.ErrorResponse{
		Error: "Todo not found",
		})
		return
	}
	todo.Completed = !todo.Completed

	if err := database.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to Toggle Complete",
		})
		return
	}
	c.JSON(http.StatusOK, todo)
}


func GetStats(c *gin.Context){
	userID, _ := c.Get("user_id")
	
	var stats struct {
		Total         int64 `json:"total"`
		Completed     int64 `json:"completed"`
		Pending       int64 `json:"pending"`
		HighPriority  int64 `json:"high_priority"`
		Overdue       int64 `json:"overdue"`
	}

	// Total todos
	database.DB.Model(&models.Todo{}).Where("user_id = ?", userID).Count(&stats.Total)

	// Completed todos
	database.DB.Model(&models.Todo{}).Where("user_id = ? AND completed = ?", userID, true).Count(&stats.Completed)

	// Pending todos
	database.DB.Model(&models.Todo{}).Where("user_id = ? AND completed = ?", userID, false).Count(&stats.Pending)

	// High priority todos
	database.DB.Model(&models.Todo{}).Where("user_id = ? AND priority = ? AND completed = ?", userID, "high", false).Count(&stats.HighPriority)

	// Overdue todos (due date passed and not completed)
	database.DB.Model(&models.Todo{}).
		Where("user_id = ? AND completed = ? AND due_date < ?", userID, false, time.Now()).
		Count(&stats.Overdue)

	c.JSON(http.StatusOK, stats)
}

func BulkDelete(c *gin.Context){
	userID, _ := c.Get("user_id")

	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	result := database.DB.Where("id IN ? AND user_id = ?", req.IDs, userID).Delete(&models.Todo{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to Delete Todos",
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"message": "Todos deleted successfully",
		"deleted": result.RowsAffected,
	})
}

// BulkComplete marks multiple todos as completed
func BulkComplete(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		IDs       []uint `json:"ids" binding:"required,min=1"`
		Completed bool   `json:"completed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Update todos that belong to the user
	result := database.DB.Model(&models.Todo{}).
		Where("id IN ? AND user_id = ?", req.IDs, userID).
		Update("completed", req.Completed)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to update todos",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Todos updated successfully",
		"updated": result.RowsAffected,
	})
}


// GetCategories returns all unique categories for the user
func GetCategories(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var categories []string
	database.DB.Model(&models.Todo{}).
		Where("user_id = ? AND category != ''", userID).
		Distinct("category").
		Pluck("category", &categories)

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}