package services

import (
	"go-http-svc/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetItems godoc
// @Summary Get all items
// @Description Get all items as an array
// @Tags item
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param q query string false "full body search string in name/description"
// @Param cid query int false "catalog id"
// @Success 200 {array} map[string]interface{}
// @Router /api/{eid}/items [get]
func GetExItems(c *gin.Context, db *gorm.DB) {
	eid := c.Param("eid")
	q := c.Query("q")
	if q != "" {
		q = "%" + q + "%"
	}

	cid, _ := strconv.Atoi(c.Query("cid"))

	var childIDs []int

	db.Raw(`
		WITH RECURSIVE child_tree AS (
			SELECT id FROM ex_catalogs WHERE pid = ?
			UNION ALL
			SELECT n.id FROM ex_catalogs n
			INNER JOIN child_tree ct ON n.pid = ct.id
		)
		SELECT id FROM child_tree;
	`, cid).Debug().Scan(&childIDs)
	childIDs = append(childIDs, cid)

	var results []struct {
		models.ExItem
		Cname     string  `json:"cname"`
		AvgRate   float64 `json:"avg_rate"`
		SumAmount int     `json:"sum_amount"`
	}

	query := db.Table("ex_items").
		Select("ex_items.*, ex_catalogs.name as cname, AVG(ex_rates.rate) as avg_rate, SUM(ex_amounts.amount) as sum_amount").
		Joins("LEFT JOIN ex_rates ON ex_rates.iid = ex_items.id and ex_rates.eid=ex_items.eid").
		Joins("LEFT JOIN ex_amounts ON ex_amounts.iid = ex_items.id and ex_amounts.eid=ex_items.eid").
		Joins("LEFT JOIN ex_catalogs ON ex_items.cid = ex_catalogs.id and ex_catalogs.eid=ex_items.eid").
		Group("ex_items.id").
		Where("ex_items.eid=? ", eid).Order("ex_items.id desc")
	if cid > 0 {
		query = query.Where("ex_items.cid in ?", childIDs)
	}
	if q != "" {
		query = query.Where("(ex_items.name like ? or ex_items.description like ?)", q, q)
	}
	query.Debug().
		Find(&results)
	c.JSON(http.StatusOK, results)
}

// CreateExItem godoc
// @Summary Create new item
// @Tags item
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param eid path int true "models.Exibition ID"
// @Param item body models.ExItemInput true "ExItem Input"
// @Success 200 {object} models.ExItem
// @Router /api/{eid}/items [put]
func CreateExItem(c *gin.Context, db *gorm.DB) {
	eid := c.Param("eid")
	var input models.ExItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.Eid, _ = strconv.Atoi(eid)
	item := models.ExItem{
		ExItemInput: input,
	}

	// TODO: NEED PRD check catalog is leaf
	// if !IsLeafCatalog(db, input.Cid) {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "必须选择叶子类目"})
	// 	return
	// }

	if result := db.Create(&item); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}
	c.JSON(http.StatusOK, item)
}

// GetExItem godoc
// @Summary Get an item by ID
// @Tags item
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.ExItem ID"
// @Success 200 {object} models.ExItem
// @Router /api/{eid}/items/{id} [get]
func GetExItem(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var item models.ExItem
	if result := db.First(&item, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// UpdateExItem godoc
// @Summary Update an item by ID
// @Tags item
// @Security BearerAuth
// @Accept  json
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "ExItem ID"
// @Param item body models.ExItemInput true "User Input"
// @Success 200 {object} models.ExItem
// @Router /api/{eid}/items/{id} [patch]
func UpdateExItem(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	id := c.Param("id")
	eid := c.Param("eid")
	var item models.ExItem
	if result := db.First(&item, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}
	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input["eid"] = eid

	var fieldsToUpdate []string
	for key := range input {
		fieldsToUpdate = append(fieldsToUpdate, key)
	}

	// Use GORM’s Updates method to perform a partial update
	if err := db.Model(&item).Select(fieldsToUpdate).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated user
	c.JSON(http.StatusOK, item)
}

// DeleteExItem godoc
// @Summary Delete an item by ID
// @Tags item
// @Security BearerAuth
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "ExItem ID"
// @Success 200
// @Router /api/{eid}/items/{id} [delete]
func DeleteExItem(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	id := c.Param("id")
	var item models.ExItem
	if result := db.First(&item, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}
	db.Delete(&item)
	c.JSON(http.StatusOK, gin.H{"message": "item deleted successfully"})
}
