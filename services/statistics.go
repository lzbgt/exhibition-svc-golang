package services

import (
	"encoding/json"
	"go-http-svc/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetTotoalUser godoc
// @Summary Get total users
// @Tags stats
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Success 200  {object} map[string]int
// @Router /api/{eid}/stats/num_users_total [get]
func GetTotoalUsers(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}

	var sum int64
	if result := db.Model(&models.ExUser{}).Where("eid = ?", eid).Count(&sum); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, map[string]int64{"data": sum})
}

// GetTotalAmount godoc
// @Summary Get total num of order
// @Tags stats
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Success 200  {object} map[string]int
// @Router /api/{eid}/stats/num_amount_total [get]
func GetTotalAmount(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}
	var sum struct {
		Sum int
	}
	if result := db.Model(&models.ExAmount{}).Select("sum(amount) as sum").Where("eid=?", eid).Scan(&sum); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "amount not found"})
		return
	}

	c.JSON(http.StatusOK, sum)
}

// GetTotalItems godoc
// @Summary Get total num of items
// @Tags stats
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Success 200  {object} map[string]int
// @Router /api/{eid}/stats/num_items_total [get]
func GetTotalItems(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}
	var sum int64
	if result := db.Model(&models.ExItem{}).Where("eid = ?", eid).Count(&sum); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, map[string]int64{"data": sum})
}

// GetExcellentItems godoc
// @Summary Get rates over 8 items
// @Tags stats
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param rate path int true "rate eg. 8.0"
// @Success 200  {object} map[string]int
// @Router /api/{eid}/stats/excellent_items/{rate} [get]
func GetExcellentItems(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	rate, _ := strconv.ParseFloat(c.Param("rate"), 64)
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}

	var results []struct {
		Iid     int     `json:"iid"`
		AvgRate float64 `json:"avg_rate"`
	}
	if result := db.Table("ex_rates").Select("ex_rates.iid, AVG(ex_rates.rate) as avg_rate").Where("eid = ?", eid).Having("avg_rate > ?", rate).Order("avg_rate desc").Group("iid").Scan(&results); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	num := len(results)

	c.JSON(http.StatusOK, map[string]int{"data": num})
}

// GetTopNAmountItems godoc
// @Summary Get top N amount items
// @Tags stats
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param topN path int true "top N"
// @Success 200  {array} map[string]interface{}
// @Router /api/{eid}/stats/topn_amount_items/{topN} [get]
func GetTopNAmountItems(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	topN, _ := strconv.Atoi(c.Param("topN"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}

	var results []struct {
		ID        int
		Name      string
		Images    json.RawMessage `json:"images" gorm:"type:json"`
		Cid       int             `json:"cid"`
		Cname     string
		AvgRate   float64 `json:"avg_rate"`
		SumAmount int     `json:"sum_amount"`
	}

	query := db.Table("ex_items").
		Select("ex_items.id, ex_items.name, ex_items.images,ex_items.cid, ex_catalogs.name as cname, AVG(ex_rates.rate) as avg_rate, SUM(ex_amounts.amount) as sum_amount").
		Joins("LEFT JOIN ex_rates ON ex_rates.iid = ex_items.id").
		Joins("LEFT JOIN ex_amounts ON ex_amounts.iid = ex_items.id").
		Joins("LEFT JOIN ex_catalogs ON ex_items.cid = ex_catalogs.id").
		Group("ex_items.id").
		Where("ex_items.eid=? ", eid).Order("sum_amount desc")

	if result := query.Limit(topN).Scan(&results); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "amount not found"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetTopNOrdersItems godoc
// @Summary Get top N amount items
// @Tags stats
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param topN path int true "top N"
// @Success 200  {array} map[string]interface{}
// @Router /api/{eid}/stats/topn_orders_items/{topN} [get]
func GetTopNOrdersItems(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	topN, _ := strconv.Atoi(c.Param("topN"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}

	var results []struct {
		ID     int
		Name   string
		Sum    int
		Orders int
	}

	// if result := db.Table("ex_amounts").Select("ex_amounts.iid, ex_items.name, ex_items.images, count(*) as sum").Joins("left join ex_items on ex_items.id=ex_amounts.iid").Group("ex_amounts.iid").Where("ex_amounts.eid=?", eid).Order("sum desc").Limit(topN).Scan(&results); result.Error != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "amount not found"})
	// 	return
	// }
	if result := db.Table("ex_users").Select("ex_users.id, ex_users.name, SUM(ex_amounts.amount) as sum, COUNT(ex_amounts.id) as orders").Joins("left join ex_amounts on ex_amounts.uid=ex_users.id").Group("ex_users.id").Where("ex_users.eid=?", eid).Order("sum desc").Limit(topN).Scan(&results); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "amount not found"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetTopNRateItems godoc
// @Summary Get top N rated items
// @Tags stats
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param topN path int true "top N"
// @Success 200  {array} map[string]interface{}
// @Router /api/{eid}/stats/topn_rate_items/{topN} [get]
func GetTopNRateItems(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	topN, _ := strconv.Atoi(c.Param("topN"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}

	var results []struct {
		Iid    int
		Name   string
		Images json.RawMessage `json:"images" gorm:"type:json"`
		Sum    float64
	}

	if result := db.Table("ex_rates").Select("ex_rates.iid, ex_items.name, SUM(rate) as sum").Joins("left join ex_items on ex_items.id=ex_rates.iid").Group("ex_rates.iid").Where("ex_rates.eid=?", eid).Order("sum desc").Limit(topN).Scan(&results); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rate not found"})
		return
	}

	c.JSON(http.StatusOK, results)
}

// GetTrending godoc
// @Summary Get catalog trending
// @Tags stats
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Success 200  {array} map[string]interface{}
// @Router /api/{eid}/stats/catalog_trending [get]
func GetCatalogTrending(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}

	type CatalogSummary struct {
		Catalog     int
		Name        string
		TotalOrders int
		TotalScores int
	}

	var result []CatalogSummary
	db.Table("ex_items").
		Select("ex_items.cid, ex_catalogs.name,SUM(ex_amounts.amount) as total_orders, SUM(ex_rates.rate) as total_scores").
		Joins("left join ex_amounts on ex_amounts.iid = ex_items.id").
		Joins("left join ex_rates on ex_rates.iid = ex_items.id").
		Joins("left join ex_catalogs on ex_catalogs.id = ex_items.cid").
		Group("ex_items.cid").
		Scan(&result)

	c.JSON(http.StatusOK, &result)
}

// GetItemsRateDistribution godoc
// @Summary Get items rate distribution
// @Tags stats
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Success 200  {array} models.RatingDistribution
// @Router /api/{eid}/stats/items_rate_distribution [get]
func GetItemsRateDistribution(c *gin.Context, db *gorm.DB) {
	user, _ := c.Get("user")
	eid, _ := strconv.Atoi(c.Param("eid"))
	claims, _ := user.(*models.Claims)
	if claims.Eid != 0 && claims.Eid != eid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "mismatch eid"})
		return
	}

	var totalItems int64
	db.Model(&models.ExItem{}).Count(&totalItems)
	if totalItems == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no items found"})
		return
	}

	var results []models.RatingDistribution

	// Query to get the count of ratings for each category
	db.Raw(`
		SELECT 
			CASE 
				WHEN avg_score BETWEEN 7 AND 10 THEN '优秀' 
				WHEN avg_score BETWEEN 4 AND 6 THEN '良好' 
				WHEN avg_score BETWEEN 1 AND 3 THEN '一般' 
			END AS category, 
			COUNT(*) AS count
		FROM (
			SELECT iid, AVG(rate) AS avg_score
			FROM ex_rates
			GROUP BY iid
		) AS avg_scores
		GROUP BY category
	`).Scan(&results)

	// Calculate the percentage for each category
	for i := range results {
		results[i].Percent = float64(results[i].Count) / float64(totalItems) * 100
	}

	c.JSON(http.StatusOK, results)
}
