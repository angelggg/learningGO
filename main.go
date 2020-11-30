package main

import (
	"dbase"
	"fmt"
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	"models"
	"os"
	"routes"
)


func main(){
	// go run main.go test -> test DB
	isTest := false
	if  len(os.Args) > 1 && os.Args[1] == "test"{
		fmt.Println("Went test mode...")
		isTest=true
	}
	InitDatabase(isTest)
	defer dbase.Conn.Close()
	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen("localhost:3000")
}


func InitDatabase(isTest bool) {
	var err error
	dbName := "project"

	if isTest{
		dbName = "test"
	}
	dbase.Conn, err = gorm.Open("mysql",
		"nombre_usuario:tu_contrasena@tcp(localhost:3306)/"+dbName+"?parseTime=true") // Fill this prop
	if err != nil {
		panic("No DB")
	}
	dbase.Conn.AutoMigrate(&models.Book{}, &models.Author{})
	fmt.Print("Connected")
}
