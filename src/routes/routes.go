package routes

import (
	"api"
	"github.com/gofiber/fiber"
)

func SetupRoutes(app *fiber.App){
	// Set routes
	getBookRoutes(app)
	getAuthorRoutes(app)

}

func getBookRoutes(app *fiber.App){
	app.Get("/api/v1/book", api.GetBooks)
	app.Get("/api/v1/book/:id", api.GetOneBook)
	app.Post("/api/v1/book", api.NewBook)
	app.Delete("/api/v1/book/:id", api.DeleteBook)
	app.Patch("/api/v1/book/:id", api.UpdateBook)
}


func getAuthorRoutes(app *fiber.App){
	app.Get("/api/v1/author", api.GetAuthors)
	app.Get("/api/v1/author/:id", api.GetOneAuthor)
	app.Post("/api/v1/author", api.NewAuthor)
	app.Delete("/api/v1/author/:id", api.DeleteAuthor)
	app.Patch("/api/v1/author/:id", api.UpdateAuthor)
}

