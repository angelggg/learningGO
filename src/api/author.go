package api

import (
	"dbase"
	"github.com/gofiber/fiber"
	"gorm.io/datatypes"
	"log"
	"models"
	"utils"
)

type (
	SimpleAuthor struct {
		Name        string         `json:name gorm:"unique;not null"`
		Genre       string         `json:genre`
		DateOfBirth datatypes.Date `json:year`
	}

	createUpdateAuthor struct {
		Name string            `json:name gorm:"unique;not null"`
		Genre string           `json:genre`
		DateOfBirth datatypes.Date `json:year`
		Books []int `json:books`
	}
)


func GetAuthors(context *fiber.Ctx) error {
	var (
		authors       []models.Author
		simpleAuthors []SimpleAuthor
	)
	db := dbase.Conn
	sortableFields := []string{"date_of_birth", "genre", "name"}
	nameFilter := context.Query("name")
	sortBy := context.Query("sort_by")
	order := "asc"
	needsSorting := false
	if sortBy != "" {
		if sortBy[0:1] == "-"{
			sortBy = sortBy[1:]
			order = "desc"
		}
		needsSorting = utils.StringIsContained(sortableFields, sortBy)
	}

	if nameFilter != "" && needsSorting {
		db.Where("name LIKE ?", "%" + nameFilter + "%").Order(sortBy + " " + order).Find(&authors)
	} else if nameFilter != ""{
		db.Where("name LIKE ?", "%" + nameFilter + "%").Find(&authors)
	} else if needsSorting {
		db.Order(sortBy + " " + order).Find(&authors)
	} else {
		db.Find(&authors)
	}
	for _, author := range authors {
		simpleAuthors = append(simpleAuthors, getSimpleAuthorFields(author))
	}
	context.JSON(simpleAuthors)
	return nil
}

func GetOneAuthor(context *fiber.Ctx) error {
	id := context.Params("id")
	db := dbase.Conn
	var author models.Author
	db.Preload("Books").Find(&author, id)
	if author.Name == "" {
		context.Status(204)
	} else {
		context.Status(200)
		context.JSON(author)
	}
	return nil
}

func NewAuthor(context *fiber.Ctx) error {
	var (
		db       = dbase.Conn
		cuAuthor createUpdateAuthor
		books    []models.Book
	)

	if err := context.BodyParser(&cuAuthor); err != nil {
		log.Fatal(err)
	}

	author := getCreateUpdateToNormalAuthor(cuAuthor, models.Author{})
	for _, bookId := range cuAuthor.Books {
		books = append(books, models.GetBookFromDB(bookId))
	}
	author.Books = books
	db.Create(&author)
	context.Status(201)
	context.JSON(author)
	return nil
}


func UpdateAuthor(context *fiber.Ctx) error {
	db := dbase.Conn
	id := context.Params("id")
	var (
		udAuthor createUpdateAuthor
		author   models.Author
	)

	if err := db.First(&author, id).Error; err != nil {
		context.Status(204).SendString("Could not find Author")
		return err
	}

	if err := context.BodyParser(&udAuthor); err != nil {
		context.Status(400).SendString("Bad request")
		return err
	}

	author = getCreateUpdateToNormalAuthor(udAuthor, author)
	switch bookOp := context.Query("bookOp"); bookOp {
		case "addbook":
			// Adds a book or books to the author and nothing else w param ?op=addbook
			for _, bookId := range udAuthor.Books {
				author.AddRelatedBook(bookId)
			}
			return nil

		case "delbook":
			// Adds a book or books to the author and nothing else w param ?bookOp=delbook
			for _, bookId := range udAuthor.Books {
				author.DeleteRelatedBook(bookId)
			}
		default:
			// Any other autors update
			if err := context.BodyParser(&udAuthor); err != nil {
				context.Status(400).SendString("Bad request")
				return err
			}
	}
	db.Save(&author)
	context.Status(202)
	context.JSON(author)
	return nil
}


func DeleteAuthor(context *fiber.Ctx) error {
	var db = dbase.Conn
	var id = context.Params("id")
	var author models.Author
	if err := db.First(&author, id).Error; err != nil {
		context.Status(204)
		context.SendString("Could not find object")
		return err
	}
	db.Delete(&author)
	context.Status(202)
	return nil
}

func getSimpleAuthorFields(author models.Author) SimpleAuthor {
	var sa SimpleAuthor
	sa.Genre = author.Genre
	sa.Name = author.Name
	sa.DateOfBirth = author.DateOfBirth
	return sa
}

func getCreateUpdateToNormalAuthor(ucAuthor createUpdateAuthor, author models.Author) models.Author{
	author.Genre = ucAuthor.Genre
	author.Name = ucAuthor.Name
	author.DateOfBirth = ucAuthor.DateOfBirth
	return author
}


