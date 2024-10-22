package services

import (
	"go-http-svc/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetExibitions godoc
// @Summary Get all exibitions
// @Description Get all exibitions as an array
// @Tags exibition
// @Security BearerAuth
// @Produce json
// @Param q query string false "full body search string in name/description"
// @Success 200 {array} models.Exibition
// @Router /api/exibitions [get]
func GetExibitions(c *gin.Context, db *gorm.DB) {
	q := c.Query("q")
	if q != "" {
		q = "%" + q + "%"
	}

	var exibitions []models.Exibition
	query := db.Model(&models.Exibition{})
	if q != "" {
		query = query.Where("title like ? or description like ? or location like ?", q, q, q)
	}
	query.Order("create_time desc").Find(&exibitions)
	c.JSON(http.StatusOK, exibitions)
}

// CreateExibition godoc
// @Summary Create new exibition
// @Tags exibition
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param exibition body models.ExibitionInput true "Exibition Input"
// @Success 200 {object} models.Exibition
// @Router /api/exibitions [put]
func CreateExibition(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	var input models.ExibitionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exibition := models.Exibition{
		ExibitionInput: models.ExibitionInput{
			Description: input.Description,
			Images:      input.Images,
			Location:    input.Location,
			Sponsors:    input.Sponsors,
			Title:       input.Title,
			Videos:      input.Videos,
			Creator:     input.Creator,
		},
	}

	if result := db.Create(&exibition); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}
	c.JSON(http.StatusOK, exibition)
}

// GetExibition godoc
// @Summary Get an exibition by ID
// @Tags exibition
// @Security BearerAuth
// @Produce json
// @Param id path int true "models.Exibition ID"
// @Success 200 {object} models.Exibition
// @Router /api/exibitions/{id} [get]
func GetExibition(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var exibition models.Exibition
	if result := db.First(&exibition, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exibition not found"})
		return
	}
	c.JSON(http.StatusOK, exibition)
}

// UpdateExibition godoc
// @Summary Update an exibition by ID
// @Tags exibition
// @Security BearerAuth
// @Accept  json
// @Produce json
// @Param id path int true "Exibition ID"
// @Param exibition body models.ExibitionInput true "User Input"
// @Success 200 {object} models.Exibition
// @Router /api/exibitions/{id} [patch]
func UpdateExibition(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	id := c.Param("id")
	var exibition models.Exibition
	if result := db.First(&exibition, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exibition not found"})
		return
	}
	// Bind the incoming JSON payload to the input struct
	var input models.ExibitionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use GORM’s Updates method to perform a partial update
	if err := db.Model(&exibition).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated user
	c.JSON(http.StatusOK, exibition)
}

// DeleteExibition godoc
// @Summary Delete an exibition by ID
// @Tags exibition
// @Security BearerAuth
// @Param id path int true "Exibition ID"
// @Success 200
// @Router /api/exibitions/{id} [delete]
func DeleteExibition(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	id := c.Param("id")
	var exibition models.Exibition
	if result := db.First(&exibition, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exibition not found"})
		return
	}
	db.Delete(&exibition)
	c.JSON(http.StatusOK, gin.H{"message": "exibition deleted successfully"})
}

// GetExibition godoc
// @Summary Get active exibition
// @Tags exibition
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.Exibition
// @Router /ex_active [get]
func GetActiveExibition(c *gin.Context, db *gorm.DB) {
	var exibition models.Exibition
	if result := db.Where("is_active=?", true).Order("update_time desc").First(&exibition); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exibition not found"})
		return
	}
	c.JSON(http.StatusOK, exibition)
}

// SetActiveExibition godoc
// @Summary Update an exibition by ID
// @Tags exibition
// @Security BearerAuth
// @Produce json
// @Param eid path int true "Exibition ID"
// @Success 200 {object} models.Exibition
// @Router /api/ex_active/{eid} [post]
func SetActiveExibition(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	eid := c.Param("eid")
	var exibition models.Exibition
	if result := db.Where("id=?", eid).Order("update_time desc").First(&exibition); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "exibition not found"})
		return
	}
	if result := db.Model(&exibition).Update("is_active", true); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, exibition)
}
