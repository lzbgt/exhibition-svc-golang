// Author: Bruce Lu
// Email: lzbgt_AT_icloud.com

package services

import (
	"go-http-svc/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetExAmount godoc
// @Summary create/update amount of a pku, with current user
// @Tags amount
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param eid path int true "models.Exibition ID"
// @Param amount body models.ExAmountInput true "ExAmount Input"
// @Success 200 {object} models.ExAmount
// @Router /api/{eid}/amounts [put]
func CreateExAmount(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatched eid"})
		return
	}

	var input models.ExAmountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.Eid = eid
	input.Uid = claims.UserId

	var amount models.ExAmount
	if result := db.Where("uid=? and iid=? and eid=?", claims.UserId, input.Iid, input.Eid).First(&amount); result.Error != nil {
		amount = models.ExAmount{
			ExAmountInput: input,
		}
		if result := db.Create(&amount); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		c.JSON(http.StatusOK, amount)
		return
	}
	if input.Amount != 0 {
		db.Model(&amount).Updates(input)
	} else {
		db.Model(&amount).Update("amount", 0)
	}

	c.JSON(http.StatusOK, amount)
}

// GetExAmount godoc
// @Summary Get an amount by ID
// @Tags amount
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.ExAmount ID"
// @Success 200 {object} models.ExAmount
// @Router /api/{eid}/amounts/{id} [get]
func GetExAmount(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var amount models.ExAmount
	if result := db.First(&amount, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "amount not found"})
		return
	}
	c.JSON(http.StatusOK, amount)
}

// GetMyAmountByItemID godoc
// @Summary Get my amount by Item ID
// @Tags amount
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.ExAmount ID"
// @Success 200 {object} models.ExAmount
// @Router /api/{eid}/amount_item/{id} [get]
func GetMyAmountByItemID(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}

	id := c.Param("id")
	var amount models.ExAmount
	if result := db.Where("iid=? and uid=?", id, claims.UserId).First(&amount); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "amount not found"})
		return
	}
	c.JSON(http.StatusOK, amount)
}

// GetTotalAmountsByItemID godoc
// @Summary Get total amount by Item ID
// @Tags amount
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.ExAmount ID"
// @Success 200 {object} models.ExAmount
// @Router /api/{eid}/amounts_item/{id} [get]
func GetTotalAmountsByItemID(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}

	id := c.Param("id")
	var amounts []models.ExAmount
	if result := db.Where("iid=?", id).Find(&amounts); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "amount not found"})
		return
	}
	sum := 0
	for i := range amounts {
		sum += amounts[i].Amount
	}

	c.JSON(http.StatusOK, map[string]int{"data": sum})
}
