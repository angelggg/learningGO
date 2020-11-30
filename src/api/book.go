package api

import (
	"dbase"
	"fmt"
	"github.com/gofiber/fiber"
	"models"
	"utils"
)

type (
	SimpleBook struct {
		Title        string         `json:title gorm:"unique;not null"`
		Year       string         `json:year`
	}

	CreateUpdateBook struct {
		Title string `json:title`
		Year string `json:year`
		Authors []int `json:authors`
	}
)



func GetBooks(context *fiber.Ctx) error {

	var (
		books []models.Book
		listedBooks []SimpleBook
	)
	db := dbase.Conn
	sortableFields := []string{"title", "year"}
	titleFilter := context.Query("title")
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

	if titleFilter != "" && needsSorting {
		db.Where("name LIKE ?", "%" + titleFilter + "%").Order(sortBy + " " + order).Find(&books)
	} else if titleFilter != ""{
		db.Where("name LIKE ?", "%" + titleFilter + "%").Find(&books)
	} else if needsSorting {
		db.Order(sortBy + " " + order).Find(&books)
	} else {
		db.Find(&books)
	}
	for _, book := range books {
		listedBooks = append(listedBooks, getSimpleBookFields(book))
	}
	context.Status(200).JSON(listedBooks)
	return nil
}

func GetOneBook(context *fiber.Ctx) error {
	id := context.Params("id")
	db := dbase.Conn
	var book models.Book
	db.Preload("Authors").Find(&book, id)
	context.Status(200).JSON(book)
	return nil
}

func NewBook(context *fiber.Ctx) error {

	var (
		db      = dbase.Conn
		cuBook  CreateUpdateBook
		authors []models.Author
	)

	if err := context.BodyParser(&cuBook); err != nil {
		context.Status(400).SendString("Bad request")
		return err
	}

	author := getCreateUpdateToNormalBook(cuBook, models.Book{})
	for _, authorId := range cuBook.Authors {
		authors = append(authors, models.GetAuthorFromDb(authorId))
	}
	db.Create(&author)
	context.Status(201)
	context.JSON(author)
	return nil
}

func DeleteBook(context *fiber.Ctx) error {
	db := dbase.Conn
	id := context.Params("id")
	var book models.Book
	if err := db.First(&book, id).Error; err != nil {
		fmt.Print(err)
		context.Status(204).SendString("Could not find object")
		return err
	} else {
		db.Delete(&book)
		context.Status(202).SendString("Book deleted") //Success
		return nil
	}
}

func UpdateBook(context *fiber.Ctx) error {
	db := dbase.Conn
	id := context.Params("id")
	var (
		udBook CreateUpdateBook
		book   models.Book
	)

	if err := db.First(&book, id).Error; err != nil {
		context.Status(204).SendString("Could not find Author")
		return err
	}

	if err := context.BodyParser(&udBook); err != nil {
		context.Status(400).SendString("Bad request")
		return err
	}

	book = getCreateUpdateToNormalBook(udBook, book)
	switch authorOp := context.Query("authorOp"); authorOp {
	case "addauthor":
		// Adds a book or books to the author and nothing else w param ?op=addbook
		for _, authorId := range udBook.Authors {
			fmt.Println(authorId)
			book.AddRelatedAuthor(authorId)
		}

	case "delauthor":
		// Adds a book or books to the author and nothing else w param ?bookOp=delbook
		for _, authorId := range udBook.Authors {
			book.DeleteRelatedAuthor(authorId)
		}
	default:
		// Any other autors update
		if err := context.BodyParser(&udBook); err != nil {
			context.Status(400).SendString("Bad request")
			return err
		}
	}
	db.Save(&book)
	context.Status(201).JSON(book)
	return nil
}

func getSimpleBookFields(book models.Book) SimpleBook {
	return SimpleBook{
		Title: book.Title,
		Year: book.Year,
	}
}

func getCreateUpdateToNormalBook(ucBook CreateUpdateBook, book models.Book) models.Book {
	book.Year = ucBook.Year
	book.Title = ucBook.Title
	return book
}

