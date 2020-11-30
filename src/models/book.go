package models

import (
	"dbase"
	"github.com/jinzhu/gorm"
)



type Book struct {
	gorm.Model
	Title string     `json:title`
	Year string      `json:year`
	Authors []Author `gorm:"many2many:author_book;"`
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Obj  Book   `json:"obj"`
}

func GetBookFromDB (id int) Book {
	var (
		db   = dbase.Conn
		book Book
	)
	db.Find(&book, id)
	return book
}

func (b Book) AddRelatedAuthor(authorId int){
	db := dbase.Conn
	authorToAdd := GetAuthorFromDb(authorId)
	db.Model(&b).Association("Authors").Append(authorToAdd)
}

func (b Book) DeleteRelatedAuthor(authorId int) {
	db := dbase.Conn
	authorToRemove := GetAuthorFromDb(authorId)
	db.Model(&b).Association("Authors").Delete(authorToRemove)
}