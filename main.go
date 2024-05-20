package main

import (
	"fmt"
	"os"

	"github.com/AlphaaaDev/GoGonicEcommerceApi/controllers"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/infrastructure"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/middlewares"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/seeds"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func drop(db *gorm.DB) {
	db.DropTableIfExists(
		&models.FileUpload{},
		&models.OrderItem{}, &models.Order{}, &models.Address{},
		&models.ProductCategory{},
		&models.Product{},
		&models.UserRole{}, &models.Role{}, &models.User{})
}

func migrate(database *gorm.DB) {

	database.AutoMigrate(&models.Address{})

	database.AutoMigrate(&models.Category{})

	database.AutoMigrate(&models.Order{})
	database.AutoMigrate(&models.OrderItem{})

	database.AutoMigrate(&models.Product{})
	database.AutoMigrate(&models.ProductCategory{})

	database.AutoMigrate(&models.Role{})
	database.AutoMigrate(&models.UserRole{})

	database.AutoMigrate(&models.User{})

	database.AutoMigrate(&models.FileUpload{})
}

func addDbConstraints(database *gorm.DB) {
	// TODO: it is well known GORM does not add foreign keys even after using ForeignKey in struct, but, why manually does not work neither ?

	dialect := database.Dialect().GetName() // mysql, sqlite3
	if dialect != "sqlite3" {
		database.Model(&models.Order{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
		database.Model(&models.Order{}).AddForeignKey("address_id", "addresses(id)", "CASCADE", "CASCADE")
		database.Model(&models.OrderItem{}).AddForeignKey("order_id", "orders(id)", "CASCADE", "CASCADE")
		database.Model(&models.OrderItem{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

		database.Model(&models.Address{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

		database.Model(&models.UserRole{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
		database.Model(&models.UserRole{}).AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")

		database.Model(models.ProductCategory{}).AddForeignKey("product_id", "products(id)", "CASCADE", "CASCADE")
		database.Model(models.ProductCategory{}).AddForeignKey("category_id", "categories(id)", "CASCADE", "CASCADE")
	}

	database.Model(&models.UserRole{}).AddIndex("user_roles__idx_user_id", "user_id")
}
func create(database *gorm.DB) {
	drop(database)
	migrate(database)
	addDbConstraints(database)
}

func main() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}
	println(os.Getenv("DB_DIALECT"))

	database := infrastructure.OpenDbConnection()

	defer database.Close()
	args := os.Args
	if len(args) > 1 {
		first := args[1]
		second := ""
		if len(args) > 2 {
			second = args[2]
		}

		if first == "create" {
			create(database)
		} else if first == "seed" {
			seeds.Seed()
			os.Exit(0)
		} else if first == "migrate" {
			migrate(database)
		}

		if second == "seed" {
			seeds.Seed()
			os.Exit(0)
		} else if first == "migrate" {
			migrate(database)
		}

		if first != "" && second == "" {
			os.Exit(0)
		}
	}

	migrate(database)

	// gin.New() - new gin Instance with no middlewares
	// goGonicEngine.Use(gin.Logger())
	// goGonicEngine.Use(gin.Recovery())
	goGonicEngine := gin.Default() // gin with the Logger and Recovery Middlewares attached
	// Allow all Origins
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "Authorization")
	goGonicEngine.Use(cors.New(config))
	//goGonicEngine.Use(cors.Default())

	goGonicEngine.Use(middlewares.Benchmark())

	// goGonicEngine.Use(middlewares.Cors())

	goGonicEngine.Use(middlewares.UserLoaderMiddleware())
	goGonicEngine.Static("/static", "./static")
	apiRouteGroup := goGonicEngine.Group("/api")

	controllers.RegisterUserRoutes(apiRouteGroup.Group("/users"))
	controllers.RegisterProductRoutes(apiRouteGroup.Group("/products"))
	controllers.RegisterAddressesRoutes(apiRouteGroup.Group("/addresses"))
	controllers.RegisterCategoryRoutes(apiRouteGroup.Group("/categories"))
	controllers.RegisterOrderRoutes(apiRouteGroup.Group("/orders"))

	goGonicEngine.Run(":8080") // listen and serve on 0.0.0.0:8080
}
