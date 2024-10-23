package main

import (
	"encoding/json"
	"fmt"
	"go-http-svc/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with a username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param sec  query string true  "query param"
// @Param user body models.ExUserInput true "User registration details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /register [post]
func Register(c *gin.Context, db *gorm.DB) {
	var input models.ExUserInput
	// Bind the JSON input to the user struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sec := c.Query("sec")
	if sec != "Hz20012056" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "access denined"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	input.Password = string(hashedPassword)
	// Create the user record
	user := models.ExUser{ExUserInput: input}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

// Login godoc
// @Summary Log in a user and get a JWT
// @Description Log in with username and password to receive a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param user body models.UserLogin true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
func Login(c *gin.Context, db *gorm.DB) {
	var input models.UserLogin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists
	var user models.ExUser
	if err := db.Where("uname = ?", input.Name).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not exist"})
		return
	}
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user disabled"})
		return
	}

	// Compare password
	if user.Eid == 0 {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
	} else {
		if user.Password != input.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
	}

	// Generate JWT token
	token, err := GenerateToken(user.Name, user.ID, user.Eid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// JWT Middleware to protect routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Validate the token
		claims, err := ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set the username in the context
		c.Set("user", claims)
		bs, _ := json.Marshal(&claims)
		fmt.Println(string(bs))

		c.Next()
	}
}
