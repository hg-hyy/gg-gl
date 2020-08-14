package model

import (
	"github.com/jinzhu/gorm"
)

// Product TEST
type Product struct {
	ID    int `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Code  string
	Price uint
}

// User   login user
type User struct {
	ID       int `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Username string
	Email    string
	Password string
}

// TestDB sqlite
func TestDB() {
	db, err := gorm.Open("sqlite3", "./model/hello.sqlite")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "L1212", Price: 1000})

	// Read
	var product Product
	db.First(&product, 1)                   // find product with id 1
	db.First(&product, "code = ?", "L1212") // find product with code l1212

	// Update - update product's price to 2000
	db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	// db.Delete(&product)
}

// UserDB sqlite
func UserDB() {
	db, err := gorm.Open("sqlite3", "./model/hello.sqlite")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&User{})

	// Read
	var user User
	db.First(&user, 1)
	db.First(&user, "Username = ?", "admin")
	db.Where("Username = ?", "admin").First(&user)

}
