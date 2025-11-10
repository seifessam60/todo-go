package handlers

import (
	"net/http"
	"todo-api/internal/database"
	"todo-api/internal/dto"
	"todo-api/internal/models"
	"todo-api/pkg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var req dto.ResgisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Email already exists",
		})
		return
	}
	if err := database.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: "Username already exists",
		})
		return
	}
	hashedPaswword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to hash password",
		})
		return
	}
	user := models.User{
		Email: req.Email,
		Username: req.Username,
		Password: string(hashedPaswword),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to create user",
		})
		return
	}
	c.JSON(http.StatusCreated, dto.ResgisterResponse{
		Message: "User Created Successfully",
		UserID: user.ID,
	})

}

func Login(c *gin.Context){
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Invalid Email or Password",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Invalid Email or Password",
		})
		return
	}
	token, err := utils.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: "Failed to generate token",
		})
		return
	}
	c.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
		User: dto.UserInfo{
			ID: user.ID,
			Username: user.Username,
			Email: user.Email,
		},
	})
}

// GetProfile returns the current user's profile
func GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	})
}
