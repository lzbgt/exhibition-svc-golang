package services

import (
	"encoding/json"
	"fmt"
	"go-http-svc/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetExCatalog godoc
// @Summary Get all catalogs
// @Description Get all catalogs as an array
// @Tags catalog
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Success 200 {array} models.ExCatalog
// @Router /api/{eid}/catalogs [get]
func GetExCatalogs(c *gin.Context, db *gorm.DB) {
	eid, err := strconv.Atoi(c.Param("eid"))
	fmt.Println(eid, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eid"})
		return
	}
	var catalogs []models.ExCatalog
	if result := db.Where("eid=?", eid).Order("create_time desc").Find(&catalogs); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}
	c.JSON(http.StatusOK, catalogs)
}

func getRootCatalog(db *gorm.DB, pid int, eid int) models.ExCatalog {
	var catalog models.ExCatalog
	query := `
		WITH RECURSIVE parent_tree AS (
			SELECT *
			FROM ex_catalogs
			WHERE id = ? and eid = ?
			UNION ALL
			SELECT t.*
			FROM ex_catalogs t
			INNER JOIN parent_tree pt ON t.id = pt.pid
		)
		SELECT * FROM parent_tree ORDER BY id limit 1;
	`
	db.Raw(query, pid, eid).Scan(&catalog)
	return catalog
}

func IsLeafCatalog(db *gorm.DB, cid int) bool {
	var count int64
	db.Model(&models.ExCatalog{}).Where("pid=?", cid).Count(&count)
	return count == 0
}

// CreateExCatalog godoc
// @Summary Create new catalog
// @Tags catalog
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param eid path int true "models.Exibition ID"
// @Param catalog body models.ExCatalogInput true "models.ExCatalog Input"
// @Success 200 {object} models.ExCatalog
// @Router /api/{eid}/catalogs [put]
func CreateExCatalog(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	eid, err := strconv.Atoi(c.Param("eid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eid"})
		return
	}
	var input models.ExCatalogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.Eid = eid

	rid := 0

	if input.Pid != 0 {
		root := getRootCatalog(db, input.Pid, input.Eid)
		if root.Pid != 0 {
			rid = root.ID
		}
	}

	catalog := models.ExCatalog{
		ExCatalogInput: input,
		RootId:         rid,
	}

	if result := db.Create(&catalog); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}
	c.JSON(http.StatusOK, catalog)
}

// GetExCatalog godoc
// @Summary Get an catalog by ID
// @Tags catalog
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.ExCatalog ID"
// @Success 200 {object} models.ExCatalog
// @Router /api/{eid}/catalogs/{id} [get]
func GetExCatalog(c *gin.Context, db *gorm.DB) {
	eid, err := strconv.Atoi(c.Param("eid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eid"})
		return
	}

	id := c.Param("id")
	var catalog models.ExCatalog
	if result := db.Where("eid=?", eid).First(&catalog, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "catalog not found"})
		return
	}
	c.JSON(http.StatusOK, catalog)
}

// UpdateExCatalog godoc
// @Summary Update an catalog by ID
// @Tags catalog
// @Security BearerAuth
// @Accept  json
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.ExCatalog ID"
// @Param catalog body models.ExCatalogInput true "models.ExCatalogInput Input"
// @Success 200 {object} models.ExCatalog
// @Router /api/{eid}/catalogs/{id} [patch]
func UpdateExCatalog(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	eid, err := strconv.Atoi(c.Param("eid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eid"})
		return
	}
	id := c.Param("id")
	var catalog models.ExCatalog
	if result := db.Where("eid=?", eid).First(&catalog, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "catalog not found"})
		return
	}
	// Bind the incoming JSON payload to the input struct
	var input models.ExCatalogInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.Eid = eid

	// Use GORM’s Updates method to perform a partial update
	if err := db.Model(&catalog).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated catalog
	c.JSON(http.StatusOK, catalog)
}

// DeleteExCatalog godoc
// @Summary Delete an catalog by ID
// @Tags catalog
// @Security BearerAuth
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "ExCatalog ID"
// @Success 200
// @Router /api/{eid}/catalogs/{id} [delete]
func DeleteExCatalog(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	eid, err := strconv.Atoi(c.Param("eid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eid"})
		return
	}
	id := c.Param("id")
	var catalog models.ExCatalog
	if result := db.Where("eid=?", eid).First(&catalog, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "catalog not found"})
		return
	}
	db.Delete(&catalog)
	c.JSON(http.StatusOK, gin.H{"message": "catalog deleted successfully"})
}

// GetExCatalogsRoot godoc
// @Summary Get root of a catelag
// @Description Get all catalogs as an array
// @Tags catalog
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.ExCatalog ID"
// @Success 200 {obj} models.ExCatalog
// @Router /api/{eid}/catalogs_root/{id} [get]
func GetExCatalogsRoot(c *gin.Context, db *gorm.DB) {
	eid, err := strconv.Atoi(c.Param("eid"))
	fmt.Println(eid, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eid"})
		return
	}
	id := c.Param("id")
	var catalog models.ExCatalog
	if result := db.Where("eid=? AND id=?", eid, id).Order("id desc").First(&catalog); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	if catalog.Pid == 0 {
		c.JSON(http.StatusOK, catalog)
		return
	}

	r1, _ := json.Marshal(catalog)
	fmt.Println("catalog", string(r1))
	pid := catalog.Pid
	catalog.ID = 0
	if result := db.Where("id=?", pid).Order("id desc").First(&catalog); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, catalog)
}

func GetParentsUsingCTE(db *gorm.DB, id string, eid int) ([]models.ExCatalog, error) {
	var parents []models.ExCatalog

	// Write the recursive CTE query
	query := `
		WITH RECURSIVE parent_tree AS (
			SELECT *
			FROM ex_catalogs
			WHERE id = ? and eid = ?
			UNION ALL
			SELECT t.*
			FROM ex_catalogs t
			INNER JOIN parent_tree pt ON t.id = pt.pid
		)
		SELECT * FROM parent_tree ORDER BY id;
	`

	// Execute the query and scan the result into the parents slice
	if err := db.Raw(query, id, eid).Scan(&parents).Error; err != nil {
		return nil, err
	}

	return parents, nil
}

// GetExCatalogsPath godoc
// @Summary Get full path to a catalog
// @Description Get path catalogs as an array
// @Tags catalog
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.ExCatalog ID"
// @Success 200 {array} models.ExCatalog
// @Router /api/{eid}/catalogs_path/{id} [get]
func GetExCatalogsPath(c *gin.Context, db *gorm.DB) {
	eid, err := strconv.Atoi(c.Param("eid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	// Get the "id" from the URL parameters
	id := c.Param("id")

	// Get the parents from the database using CTE
	parents, err := GetParentsUsingCTE(db, id, eid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the parents in JSON format
	c.JSON(http.StatusOK, parents)
}

// Query the node and its children up to a certain depth using a recursive CTE
func GetNodeAndChildren(db *gorm.DB, eid, id, depth string) ([]models.ExCatalog, error) {
	var nodes []models.ExCatalog

	query := `
		WITH RECURSIVE children_tree AS (
			SELECT id, pid, eid, name, description,images, videos, 0 AS depth
			FROM ex_catalogs
			WHERE id = ? and eid = ?
			UNION ALL
			-- Recursive step
			SELECT t.id, t.pid, t.eid, t.name, t.description,t.images, t.videos, ct.depth + 1
			FROM ex_catalogs t
			INNER JOIN children_tree ct ON t.pid = ct.id
			WHERE ct.depth < ?
		)
		SELECT * FROM children_tree ORDER BY id;
	`

	// Execute the query with the given id and depth
	if err := db.Raw(query, id, eid, depth).Scan(&nodes).Error; err != nil {
		return nil, err
	}

	return nodes, nil
}

// GetExCatalogsChildren godoc
// @Summary Get all catalogs
// @Description Get all catalogs as an array
// @Tags catalog
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.Exibition ID"
// @Param id path int true "models.ExCatalog ID"
// @Param depth  query int true  "Depth of children to fetch"
// @Success 200 {array} models.ExCatalog
// @Router /api/api/{eid}/sub_catalogs/{id} [get]
func GetExCatalogsChildren(c *gin.Context, db *gorm.DB) {
	eid := c.Param("eid")
	// Get the "id" from the URL parameters
	id := c.Param("id")
	depth := c.Query("depth")

	fmt.Println(eid, id, depth)

	// Query the node and its children using the service function
	nodes, err := GetNodeAndChildren(db, eid, id, depth)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the result in JSON format
	c.JSON(http.StatusOK, nodes)
}
