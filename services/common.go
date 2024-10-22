package services

import (
	"fmt"
	"go-http-svc/models"
	"log"
	"net/http"
	"time"

	"path/filepath"

	"math/rand"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

func IsAdmin(c *gin.Context) bool {
	claims, ok := c.Get("user")
	if !ok {
		fmt.Println("claims not found")
		return false
	}

	if claims == nil {
		fmt.Println("claims nil")
		return false
	}

	claim := claims.(*models.Claims)
	return claim.Eid == 0
}

func ReadExcel(file_path string) (ret []map[string]string) {
	f, err := excelize.OpenFile(file_path)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		// Close the spreadsheet after reading
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Get all the rows from the first sheet
	rows, err := f.GetRows("Sheet 1")
	if err != nil {
		log.Fatal(err)
	}

	// Loop through the rows and print each cell value
	for i, row := range rows {
		if i <= 1 {
			continue
		}
		m := make(map[string]string)
		m["name"] = row[0]
		m["title"] = row[1]
		m["mobile"] = row[2]

		ret = append(ret, m)
		fmt.Println(i, m)
	}

	return ret
}

// @Summary Upload a file
// @Description Upload a file to the server
// @Tags helpers
// @Security BearerAuth
// @Accept mpfd
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/file_upload [post]
func UploadFile(c *gin.Context) {
	// Get the file from the form
	if !IsAdmin(c) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "admin only"})
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
		return
	}

	base := filepath.Base(file.Filename)
	ext := filepath.Ext(base) // Get the file extension (e.g. .txt)
	newUUID, _ := uuid.NewUUID()

	// Replace the stem with the UUID
	newFileName := newUUID.String() + ext

	// Save the file to a destination
	err = c.SaveUploadedFile(file, "./uploads/"+newFileName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "url": "uploads/" + newFileName})
}

func generateRandomUsername(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random string of the specified length
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
