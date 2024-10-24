package services

import (
	"fmt"
	"go-http-svc/models"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GetExUsers godoc
// @Summary Get all users
// @Description Get all users as an array
// @Tags user
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.ExUser EID"
// @Param q query string false "full body search string in name/title/mobile"
// @Success 200 {array} models.ExUser
// @Router /api/{eid}/users [get]
func GetExUsers(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}
	eid, err := strconv.Atoi(c.Param("eid"))
	fmt.Println(eid, err)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eid"})
		return
	}
	q := c.Query("q")
	if q != "" {
		q = "%" + q + "%"
	}

	var users []models.ExUser
	query := db.Where("eid=?", eid).Order("id desc")
	if q != "" {
		query = query.Where("name like ? or title like ? or uname like ? or mobile like ?", q, q, q, q)
	}
	if result := query.Find(&users); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error})
		return
	}
	c.JSON(http.StatusOK, users)
}

// CreateExUser godoc
// @Summary Create new user
// @Tags user
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param eid path int true "models.ExUser EID"
// @Param user body models.ExUserInput true "models.ExUser Input"
// @Success 200 {object} models.ExUser
// @Router /api/{eid}/users [post]
func CreateExUser(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	eid, err := strconv.Atoi(c.Param("eid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eid"})
		return
	}
	var input models.ExUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.Eid = eid
	user := models.ExUser{
		ExUserInput: input,
	}

	if result := db.Create(&user); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}
	c.JSON(http.StatusOK, user)
}

// BatchCreateExUser godoc
// @Summary Create new users by batch
// @Tags user
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param eid path int true "models.ExUser EID"
// @Param user body models.BatchUserInput true "BatchUserInput Input"
// @Success 200 {object} models.ExUser
// @Router /api/{eid}/users [put]
func BatchCreateExUser(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	eid, err := strconv.Atoi(c.Param("eid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eid"})
		return
	}
	var input models.BatchUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var users []models.ExUser
	for i := input.IndexRange[0]; i <= input.IndexRange[1]; i++ {
		name := input.NamePrefix + strconv.Itoa(i)
		user := models.ExUser{
			ExUserInput: models.ExUserInput{
				Eid:      eid,
				Uname:    name,
				Password: "0000",
				Name:     name,
			},
		}
		users = append(users, user)
	}

	if result := db.Create(&users); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetExUser godoc
// @Summary Get an user by ID
// @Tags user
// @Security BearerAuth
// @Produce json
// @Param eid path int true "models.ExUser EID"
// @Param id path int true "models.ExUser ID"
// @Success 200 {object} models.ExUser
// @Router /api/{eid}/users/{id} [get]
func GetExUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if !IsAdmin(c) && id != "0" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin/self_user only"})
		return
	}

	eid, err := strconv.Atoi(c.Param("eid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eid"})
		return
	}

	claim, _ := c.Get("user")
	u := claim.(*models.Claims)

	if id == "0" {
		id = strconv.Itoa(u.UserId)
	}

	fmt.Println("uid: ", id)

	var user models.ExUser
	if result := db.Where("eid=? and id=?", eid, id).First(&user); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateExUser godoc
// @Summary Update an user by ID
// @Tags user
// @Security BearerAuth
// @Accept  json
// @Produce json
// @Param eid path int true "models.ExUser EID"
// @Param id path int true "models.ExUser ID"
// @Param user body models.ExUser true "User Input"
// @Success 200 {object} models.ExUser
// @Router /api/{eid}/users/{id} [patch]
func UpdateExUser(c *gin.Context, db *gorm.DB) {
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
	var user models.ExUser
	if result := db.Where("eid=?", eid).First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var fieldsToUpdate []string
	for key := range input {
		fieldsToUpdate = append(fieldsToUpdate, key)
	}

	// Use GORM’s Updates method to perform a partial update
	if err := db.Model(&user).Select(fieldsToUpdate).Debug().Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated user
	c.JSON(http.StatusOK, user)
}

// DeleteExUser godoc
// @Summary Delete an user by ID
// @Tags user
// @Security BearerAuth
// @Param eid path int true "models.ExUser EID"
// @Param id path int true "ExUser ID"
// @Success 200
// @Router /api/{eid}/users/{id} [delete]
func DeleteExUser(c *gin.Context, db *gorm.DB) {
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
	var user models.ExUser
	if result := db.Where("eid=?", eid).First(&user, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	db.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// @Summary Create Users from template xlsx file
// @Tags user
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param file formData file true "user_template.xlsx to upload"
// @Router /api/{eid}/users_tmpl [put]
func CreateUserFromTemplate(c *gin.Context, db *gorm.DB) {
	// Get the file from the form
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	eid, _ := strconv.Atoi(c.Param("eid"))
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
		return
	}

	base := filepath.Base(file.Filename)
	ext := filepath.Ext(base) // Get the file extension (e.g. .txt)
	newUUID, _ := uuid.NewUUID()

	// Replace the stem with the UUID
	newFileName := "./uploads/" + newUUID.String() + ext

	// Save the file to a destination
	err = c.SaveUploadedFile(file, newFileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save file"})
		return
	}
	ret := ReadExcel(newFileName)
	users := []models.ExUser{}

	for i, row := range ret {

		user := models.ExUser{
			ExUserInput: models.ExUserInput{
				Eid:      eid,
				Title:    row["title"],
				Uname:    generateRandomUsername(5) + strconv.Itoa(i),
				Password: "0000",
				Name:     row["name"],
			},
		}
		users = append(users, user)
	}
	if len(users) > 0 {
		if result := db.Create(&users); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "users created successfully"})
	}
}

// ActiveUser godoc
// @Summary Active/Disable Users
// @Tags user
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param eid path int true "models.ExUser EID"
// @Param active query bool true "active/disable: true|false"
// @Param user body []int true "User IDs"
// @Success 200 {array} models.ExUser
// @Router /api/{eid}/active_users [post]
func ActiveExUser(c *gin.Context, db *gorm.DB) {
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}

	eid, err := strconv.Atoi(c.Param("eid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eid"})
		return
	}
	var input []int
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(input) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no input user ids"})
		return
	}
	active, _ := strconv.ParseBool(c.Query("active"))
	if result := db.Model(&models.ExUser{}).Where("eid=? and id in ?", eid, input).Update("is_active", active); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error})
		return
	}

	c.JSON(http.StatusOK, nil)
}
