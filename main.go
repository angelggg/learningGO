package main

import (
	"dbase"
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	"models"
	"routes"
)


func main(){
	InitDatabase()
	defer dbase.Conn.Close()
	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen("localhost:3000")
}


func InitDatabase() {
	var err error
	dbase.Conn, err = gorm.Open("mysql",
		"username:pwd@tcp(localhost:3306)/test?parseTime=true") // Fill this prop
	if err != nil {
		panic("No DB")
	}
	dbase.Conn.AutoMigrate(&models.Book{}, &models.Author{})
	fmt.Print("Connected")
}
