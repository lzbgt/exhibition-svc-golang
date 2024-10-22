package services

import (
	"go-http-svc/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetComments godoc
// @Summary Get all comments of an item, time desc
// @Description Get all comments of an item as an array
// @Tags comment
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.Comment ID"
// @Success 200 {array} models.ExComment
// @Router /api/{eid}/comments/{id} [get]
func GetExComments(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	id := c.Param("id")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatched eid"})
		return
	}

	var comments []models.ExComment
	if result := db.Where("eid=? and iid=?", eid, id).Order("create_time desc").Find(&comments); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}
	c.JSON(http.StatusOK, comments)
}

// CreateComment godoc
// @Summary Create/Update a comment an item
// @Tags comment
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param item body models.ExCommentInput true "ExComment Input"
// @Success 200 {object} models.ExComment
// @Router /api/{eid}/comments/ [put]
func CreateExComment(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatched eid"})
		return
	}

	var input models.ExCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.Uid = claims.UserId
	input.Eid = eid

	var comment models.ExComment
	if result := db.Where("uid=? and iid=? and eid=?", input.Uid, input.Iid, input.Eid).First(&comment); result.Error != nil {
		comment = models.ExComment{ExCommentInput: input}
		if result := db.Create(&comment); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		c.JSON(http.StatusOK, comment)
		return
	}

	db.Model(&comment).Updates(input)
	c.JSON(http.StatusOK, comment)
}
