package models

import (
	"dbase"
	"github.com/jinzhu/gorm"
	"gorm.io/datatypes"
)


type Author struct {
	gorm.Model
	Name string            `json:name gorm:"unique;not null"`
	Genre string           `json:genre`
	DateOfBirth datatypes.Date  `json:dob`
	Books []Book           `gorm:"many2many:author_book;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}


func GetAuthorFromDb(id int) Author {
	db := dbase.Conn
	var author Author
	db.Find(&author, id)
	return author
}

func (a Author) AddRelatedBook(bookId int){
	db := dbase.Conn
	bookToAdd := GetBookFromDB(bookId)
	db.Model(&a).Association("Books").Append(bookToAdd)
}

func (a Author) DeleteRelatedBook(bookId int) {
	db := dbase.Conn
	bookToRemove := GetBookFromDB(bookId)
	db.Model(&a).Association("Books").Delete(bookToRemove)
	db.Save(&a)
}