package services

import (
	"go-http-svc/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetExRate godoc
// @Summary create/update rate of a pku, with current user
// @Tags rate
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param eid path int true "models.Exibition ID"
// @Param rate body models.ExRateInput true "ExRate Input"
// @Success 200 {object} models.ExRate
// @Router /api/{eid}/rates [put]
func CreateExRate(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatched eid"})
		return
	}

	var input models.ExRateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.Uid = claims.UserId
	input.Eid = eid

	var rate models.ExRate
	if result := db.Where("uid=? and iid=? and eid=?", input.Uid, input.Iid, input.Eid).First(&rate); result.Error != nil {
		rate = models.ExRate{
			ExRateInput: input,
		}
		if result := db.Create(&rate); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		c.JSON(http.StatusOK, rate)
		return
	}

	db.Model(&rate).Updates(input)
	c.JSON(http.StatusOK, rate)
}

// GetExRate godoc
// @Summary Get an rate by ID
// @Tags rate
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.ExRate ID"
// @Success 200 {object} models.ExRate
// @Router /api/{eid}/rates/{id} [get]
func GetExRate(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var rate models.ExRate
	if result := db.First(&rate, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rate not found"})
		return
	}
	c.JSON(http.StatusOK, rate)
}

// GetMyRateByItemID godoc
// @Summary Get my rate of an Item
// @Tags rate
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.ExRate ID"
// @Success 200 {object} models.ExRate
// @Router /api/{eid}/rate_item/{id} [get]
func GetMyRateByItemID(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}

	id := c.Param("id")
	var rate models.ExRate
	if result := db.Where("iid=? and uid=?", id, claims.UserId).First(&rate); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rate not found"})
		return
	}
	c.JSON(http.StatusOK, rate)
}

// GetMyRateByItemID godoc
// @Summary Get my rate of an Item
// @Tags rate
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param rate body []int true "item IDs"
// @Success 200 {array} models.ExRate
// @Router /api/{eid}/my_rates_items [post]
func GetMyRatesByItemIDs(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}

	var input []int
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(input) < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	var rates []models.ExRate
	if result := db.Where("eid=? and iid in ? and uid=? ", eid, input, claims.UserId).Find(&rates); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rate not found"})
		return
	}
	c.JSON(http.StatusOK, rates)
}

// GetTotalRatesByItemID godoc
// @Summary Get total rate of an Item
// @Tags rate
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.ExRate ID"
// @Success 200 {object} models.ExRate
// @Router /api/{eid}/rates_item/{id} [get]
func GetTotalRatesByItemID(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}

	id := c.Param("id")
	var rates []models.ExRate
	if result := db.Where("iid=?", id).Find(&rates); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rate not found"})
		return
	}
	sum := 0.0
	for i := range rates {
		sum += rates[i].Rate
	}

	c.JSON(http.StatusOK, map[string]float64{"data": sum})
}
