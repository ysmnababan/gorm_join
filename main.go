package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// Company has many Users
type Company struct {
	CompID uint   `gorm:"column:id;primaryKey" json:"id"`
	Name   string `gorm:"size:100" json:"name"`

	// Users []User `gorm:"foreignKey:IDPerusahaan;references:CompID" json:"users"`
	Users []User `gorm:"foreignKey:IDPerusahaan" json:"users"` // you don't have to add the references tag
}

// User belongs to a Company and has many Orders
type User struct {
	ID           uint     `gorm:"primaryKey" json:"id"`
	Name         string   `gorm:"size:100" json:"name"`
	Age          int      `json:"age"`
	IDPerusahaan uint     `json:"company_id" gorm:"column:company_id"`
	Company      *Company `gorm:"foreignKey:IDPerusahaan;references:CompID" json:"company"`

	Orders []Order `gorm:"foreignKey:UserID" json:"orders"`
}

// Order belongs to a User
type Order struct {
	ID     uint    `gorm:"primaryKey" json:"id"`
	UserID uint    `json:"user_id"`
	Amount float64 `json:"amount"`
}

func seedData(db *gorm.DB) error {
	// Check if data exists
	var count int64
	db.Model(&User{}).Count(&count)
	if count > 0 {
		return nil // Data already exists
	}

	// Create Companies
	company1 := Company{Name: "Tech Corp"}
	company2 := Company{Name: "Finance Inc."}
	db.Create(&company1)
	db.Create(&company2)

	// Create Users
	user1 := User{Name: "Alice", Age: 30, IDPerusahaan: company1.CompID}
	user2 := User{Name: "Bob", Age: 25, IDPerusahaan: company2.CompID}
	user3 := User{Name: "Charlie", Age: 35, IDPerusahaan: company1.CompID}
	db.Create(&user1)
	db.Create(&user2)
	db.Create(&user3)

	// Create Orders
	order1 := Order{UserID: user1.ID, Amount: 100.50}
	order2 := Order{UserID: user1.ID, Amount: 200.75}
	order3 := Order{UserID: user2.ID, Amount: 50.25}
	db.Create(&order1)
	db.Create(&order2)
	db.Create(&order3)

	fmt.Println("Database seeded successfully!")
	return nil
}
func main() {
	dsn := "host=localhost user=admin password=password dbname=mydatabase port=5433 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info), // Enable SQL logging
	})
	if err != nil {
		panic(err)
	}

	// db.AutoMigrate(&Company{}, &User{}, &Order{})
	// err = seedData(db)
	// if err != nil {
	// log.Fatal(err)
	// }

	users := []User{}
	res := db.
		Joins("Company").
		// Joins("LEFT JOIN orders ON orders.user_id = users.id AND orders.id is null ").
		// Select(`users.*, "Company"."id" AS "Company__id", "Company"."name" AS "Company__name", orders.id AS Order__id, orders.amount AS Order__amount`).
		Find(&users)
	if res.Error != nil {
		fmt.Println(res.Error)
		return
	}
	for _, val := range users {
		fmt.Println(val.ID, val.Name, val.Company, "->", val.Orders)
	}
}
