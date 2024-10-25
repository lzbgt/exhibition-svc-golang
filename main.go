package main

import (
	_ "embed"
	"fmt"
	"go-http-svc/docs"
	"go-http-svc/models"
	"go-http-svc/services"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title BSD Exibition API
// @version 1.0
// @description XYZ
// @BasePath /

var db *gorm.DB
var err error

//go:embed user_template.xlsx
var userTemplate []byte

func extractTemplateIfNotExists() error {
	// Define the path where the file should be extracted
	uploadsDir := "uploads"

	// Step 4.1: Check if the file already exists
	filePath := filepath.Join(uploadsDir, "user_template.xlsx")
	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("File already exists, no extraction needed.")
		return nil
	}

	// Step 4.2: Create the uploads directory if it doesn't exist
	err := os.MkdirAll(uploadsDir, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create uploads directory: %w", err)
	}

	// Step 4.3: Write the embedded file to the uploads directory
	err = os.WriteFile(filePath, userTemplate, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	fmt.Println("File extracted successfully to", filePath)
	return nil
}

// Initialize the MySQL database
func initDB() {
	DB_URL := os.Getenv("DB_URL")
	db, err = gorm.Open(mysql.Open(DB_URL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Auto-migrate the User model
	db.AutoMigrate(&models.Exibition{}, &models.ExUser{}, &models.ExCatalog{},
		&models.ExItem{}, &models.ExRate{}, &models.ExAmount{}, &models.ExComment{})
}

// @title Exibition System API
// @version 1.0
// @description This is a sample server that demonstrates JWT with Swagger and Gin.
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @type apiKey
// @name Authorization
// @in header
// @description Provide your JWT token in the format: {your_token}
func main() {
	(&services.Window{}).DisableConsoleQuickEdit()
	extractTemplateIfNotExists()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Initialize the database
	initDB()

	// Create a new Gin router
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	docs.SwaggerInfo.BasePath = "/"
	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/register", func(c *gin.Context) {
		Register(c, db)
	})
	r.Static("/uploads", "./uploads")

	r.POST("/login", func(c *gin.Context) {
		Login(c, db)
	})
	r.GET("/ex_active", func(c *gin.Context) {
		services.GetActiveExibition(c, db)
	})

	router := r.Group("/api")
	router.Use(AuthMiddleware())

	// Set up routes
	router.POST("/file_upload", func(c *gin.Context) {
		services.UploadFile(c)
	})
	router.GET("/exibitions", func(c *gin.Context) {
		services.GetExibitions(c, db)
	})
	router.PUT("/exibitions", func(c *gin.Context) {
		services.CreateExibition(c, db)
	})
	router.GET("/exibitions/:id", func(c *gin.Context) {
		services.GetExibition(c, db)
	})
	router.PATCH("/exibitions/:id", func(c *gin.Context) {
		services.UpdateExibition(c, db)
	})
	router.DELETE("/exibitions/:id", func(c *gin.Context) {
		services.DeleteExibition(c, db)
	})
	router.POST("/ex_active/:eid", func(c *gin.Context) {
		services.SetActiveExibition(c, db)
	})
	//users

	router.PUT("/:eid/users_tmpl", func(c *gin.Context) {
		services.CreateUserFromTemplate(c, db)
	})

	router.GET("/:eid/users", func(c *gin.Context) {
		services.GetExUsers(c, db)
	})
	router.POST("/:eid/users", func(c *gin.Context) {
		services.CreateExUser(c, db)
	})
	router.PUT("/:eid/users", func(c *gin.Context) {
		services.BatchCreateExUser(c, db)
	})
	router.GET("/:eid/users/:id", func(c *gin.Context) {
		services.GetExUser(c, db)
	})
	router.PATCH("/:eid/users/:id", func(c *gin.Context) {
		services.UpdateExUser(c, db)
	})
	router.DELETE("/:eid/users/:id", func(c *gin.Context) {
		services.DeleteExUser(c, db)
	})

	router.POST("/:eid/active_users", func(c *gin.Context) {
		services.ActiveExUser(c, db)
	})

	//catalogs
	router.GET("/:eid/catalogs", func(c *gin.Context) {
		services.GetExCatalogs(c, db)
	})
	router.PUT("/:eid/catalogs", func(c *gin.Context) {
		services.CreateExCatalog(c, db)
	})
	router.GET("/:eid/catalogs/:id", func(c *gin.Context) {
		services.GetExCatalog(c, db)
	})
	router.PATCH("/:eid/catalogs/:id", func(c *gin.Context) {
		services.UpdateExCatalog(c, db)
	})
	router.DELETE("/:eid/catalogs/:id", func(c *gin.Context) {
		services.DeleteExCatalog(c, db)
	})

	router.GET("/:eid/catalogs_root/:id", func(c *gin.Context) {
		services.GetExCatalogsRoot(c, db)
	})

	router.GET("/:eid/catalogs_path/:id", func(c *gin.Context) {
		services.GetExCatalogsPath(c, db)
	})

	router.GET("/:eid/sub_catalogs/:id", func(c *gin.Context) {
		services.GetExCatalogsChildren(c, db)
	})

	// Item
	router.GET("/:eid/items", func(c *gin.Context) {
		services.GetExItems(c, db)
	})

	router.GET("/:eid/items/:id", func(c *gin.Context) {
		services.GetExItem(c, db)
	})

	router.PUT("/:eid/items", func(c *gin.Context) {
		services.CreateExItem(c, db)
	})
	router.PATCH("/:eid/items/:id", func(c *gin.Context) {
		services.UpdateExItem(c, db)
	})
	router.DELETE("/:eid/items/:id", func(c *gin.Context) {
		services.DeleteExItem(c, db)
	})

	// router.GET("/:eid/search_items", func(c *gin.Context) {
	// 	services.SearchExItems(c, db)
	// })

	// Rate
	router.GET("/:eid/rates/:id", func(c *gin.Context) {
		services.GetExRate(c, db)
	})
	router.PUT("/:eid/rates", func(c *gin.Context) {
		services.CreateExRate(c, db)
	})
	router.GET("/:eid/rate_item/:id/", func(c *gin.Context) {
		services.GetMyRateByItemID(c, db)
	})
	router.GET("/:eid/rates_item/:id/", func(c *gin.Context) {
		services.GetTotalRatesByItemID(c, db)
	})

	router.POST("/:eid/my_rates_items", func(c *gin.Context) {
		services.GetMyRatesByItemIDs(c, db)
	})

	//Amount
	router.GET("/:eid/amounts/:id", func(c *gin.Context) {
		services.GetExAmount(c, db)
	})
	router.PUT("/:eid/amounts", func(c *gin.Context) {
		services.CreateExAmount(c, db)
	})
	router.GET("/:eid/amount_item/:id/", func(c *gin.Context) {
		services.GetMyAmountByItemID(c, db)
	})
	router.GET("/:eid/amounts_item/:id/", func(c *gin.Context) {
		services.GetTotalAmountsByItemID(c, db)
	})

	// comment
	router.GET("/:eid/comments/:id", func(c *gin.Context) {
		services.GetExComments(c, db)
	})

	router.PUT("/:eid/comments", func(c *gin.Context) {
		services.CreateExComment(c, db)
	})

	// stats
	router.GET("/:eid/stats/topn_rate_items/:topN", func(c *gin.Context) {
		services.GetTopNRateItems(c, db)
	})
	router.GET("/:eid/stats/topn_amount_items/:topN", func(c *gin.Context) {
		services.GetTopNAmountItems(c, db)
	})
	router.GET("/:eid/stats/topn_orders_items/:topN", func(c *gin.Context) {
		services.GetTopNOrdersItems(c, db)
	})
	router.GET("/:eid/stats/items_rate_distribution", func(c *gin.Context) {
		services.GetItemsRateDistribution(c, db)
	})
	router.GET("/:eid/stats/num_amount_total", func(c *gin.Context) {
		services.GetTotalAmount(c, db)
	})
	router.GET("/:eid/stats/num_items_total", func(c *gin.Context) {
		services.GetTotalItems(c, db)
	})
	router.GET("/:eid/stats/num_users_total", func(c *gin.Context) {
		services.GetTotoalUsers(c, db)
	})
	router.GET("/:eid/stats/excellent_items/:rate", func(c *gin.Context) {
		services.GetExcellentItems(c, db)
	})
	router.GET("/:eid/stats/catalog_trending", func(c *gin.Context) {
		services.GetCatalogTrending(c, db)
	})

	router.GET("/:eid/stats/orders_users_rate", func(c *gin.Context) {
		services.GetOrdersRateOfUsers(c, db)
	})

	// Start the server
	PORT := os.Getenv("PORT")
	fmt.Println("Server is running on port ", PORT)
	r.Run(":" + PORT)
}
