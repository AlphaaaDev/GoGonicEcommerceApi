package seeds

import (
	"math/rand"
	"time"

	"github.com/AlphaaaDev/GoGonicEcommerceApi/infrastructure"
	"github.com/AlphaaaDev/GoGonicEcommerceApi/models"
	"github.com/icrowley/fake"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

func randomInt(min, max int) int {

	return rand.Intn(max-min) + min
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func seedAdmin(db *gorm.DB) {
	count := 0
	adminRole := models.Role{Name: "ROLE_ADMIN", Description: "Only for admin"}
	query := db.Model(&models.Role{}).Where("name = ?", "ROLE_ADMIN")
	query.Count(&count)

	if count == 0 {
		db.Create(&adminRole)
	} else {
		query.First(&adminRole)
	}

	adminRoleUsers := 0
	var adminUsers []models.User
	db.Model(&adminRole).Related(&adminUsers, "Users")

	db.Model(&models.User{}).Where("username = ?", "admin").Count(&adminRoleUsers)
	if adminRoleUsers == 0 {

		// query.First(&adminRole) // First would fetch the Role admin because the query status name='ROLE_ADMIN'
		password, _ := bcrypt.GenerateFromPassword([]byte("administrator"), bcrypt.DefaultCost)
		// Approach 1
		user := models.User{FirstName: "Admin", LastName: "Admin", Email: "admin@example.com", Username: "admin", Password: string(password)}
		user.Roles = append(user.Roles, adminRole)

		// Do not try to update the adminRole
		db.Set("gorm:association_autoupdate", false).Create(&user)

		// Approach 2
		// user := models.User{FirstName: "AdminFN", LastName: "AdminFN", Email: "admin@golang.com", Username: "admin", Password: "password"}
		// user.Roles = append(user.Roles, adminRole)
		// db.NewRecord(user)
		// db.Set("gorm:association_autoupdate", false).Save(&user)

		if db.Error != nil {
			print(db.Error)
		}
	}
}

func seedUserRole(db *gorm.DB) {
	count := 0
	role := models.Role{Name: "ROLE_USER", Description: "Only for standard users"}
	q := db.Model(&models.Role{}).Where("name = ?", "ROLE_USER")
	q.Count(&count)

	if count == 0 {
		db.Create(&role)
	}
}

func seedCategories(db *gorm.DB) {
	var categories [3]models.Category
	db.Where(models.Category{Name: "Women"}).Attrs(models.Category{Description: "Clothes for women", IsNewRecord: true}).FirstOrCreate(&categories[0])
	db.Where(models.Category{Name: "Men"}).Attrs(models.Category{Description: "Clothes for men", IsNewRecord: true}).FirstOrCreate(&categories[1])
	db.Where(models.Category{Name: "Kids"}).Attrs(models.Category{Description: "Clothes for kids", IsNewRecord: true}).FirstOrCreate(&categories[2])
}

func seedProducts(db *gorm.DB) {
	productsCount := 0
	productsToSeed := 20
	db.Model(&models.Product{}).Count(&productsCount)
	productsToSeed -= productsCount

	if productsToSeed > 0 {
		rand.Seed(time.Now().Unix())
		categories := []models.Category{}
		db.Find(&categories)
		for i := 0; i < productsToSeed; i++ {
			// add a tag and a category for each product
			categoryForProduct := categories[rand.Intn(len(categories))]

			product := &models.Product{Name: fake.ProductName(), Description: fake.Paragraph(),
				Stock: randomInt(100, 2000), Price: randomInt(50, 1000),
				Categories: []models.Category{categoryForProduct}}
			for i := 0; i < randomInt(1, 4); i++ {
				db.Set("gorm:association_autoupdate", false).Create(&product)
			}
		}
	}
}

func Seed() {
	db := infrastructure.GetDb()
	rand.Seed(time.Now().UnixNano())
	seedAdmin(db)
	seedUserRole(db)
	seedCategories(db)
	seedProducts(db)
}
