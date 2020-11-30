package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/datatypes"
	"models"
	"net/http"
	"strconv"
	"testing"
)

type AuthorResponse struct {
	Id int  					`json:id`
	Name string                `json:name gorm:"unique;not null"`
	Genre string               `json:genre`
	DateOfBirth datatypes.Date `json:dob`
	Books []models.Book         `gorm:"many2many:author_book;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
var authorId int
var url = "http://localhost:3000/api/v1/author/"
var detailUrl = ""

func TestCreateAuthor (t *testing.T){
	var jsonStr = []byte(`{
    "Name": "Test author",
    "Genre": "F",
    "DateOfBirth" : "1110-07-17T00:00:00Z"
	}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	if resp.StatusCode != 201{
		fmt.Println("Error: Got status code "  )
		fmt.Println(resp.StatusCode)
		t.Fail()
	}
	if !t.Failed(){

		defer resp.Body.Close()
		decoder := json.NewDecoder(resp.Body)
		var author AuthorResponse
		err = decoder.Decode(&author)
		authorId = author.Id
		detailUrl = url + strconv.Itoa(authorId)
		fmt.Println("OK author created")
	}
}

func TestGetAllAuthors(t *testing.T){
	resp, err := http.Get("http://localhost:3000/api/v1/author")
	if err != nil  {
		fmt.Println("Error Get all authors " )
		fmt.Println(err.Error())
		t.FailNow()
	}
	if resp.StatusCode != 200{
		fmt.Println("Wrong status code get all authors, got")
		fmt.Println(resp.StatusCode)
		t.FailNow()
	}
	fmt.Println("OK get all authors")
}

func TestGetOneAuthor(t *testing.T){
	resp, err := http.Get(detailUrl)
	if err != nil  {
		fmt.Println("Error Get one authors " )
		fmt.Println(err.Error())
		t.FailNow()
	}
	if resp.StatusCode != 200{
		fmt.Println("Wrong status code get one author, got")
		fmt.Println(resp.StatusCode)
		t.FailNow()
	}
	fmt.Println("OK get one author")
}

func TestUpdateAuthor(t *testing.T){
	var author AuthorResponse
	var jsonStr = []byte(`{
    "Name": "updated",
    "DateOfBirth" : "1110-07-17T00:00:00Z",
    "Genre": "M"}`)

	req, err := http.NewRequest("PATCH", detailUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Patch author failed attempting request: " + err.Error())
		t.FailNow()
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&author)
	if err != nil{
		fmt.Println("Patch author failed decoding response: " + err.Error())
		t.FailNow()
	}
	if author.Name != "updated" || author.Genre != "M"{
				fmt.Println("Patch author failed: params dont match")
				t.FailNow()
	}

	fmt.Println("OK author Updated")
	}

func TestDeleteAuthor(t *testing.T)  {
	req, _ := http.NewRequest("DELETE", detailUrl,nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		fmt.Println("Delete author failed attempting request: " + err.Error())
		t.FailNow()
	}
	if resp.StatusCode != 202	{
		fmt.Println("Delete author got wrong status code: " + resp.Status)
		t.FailNow()
	}
	resp , err = http.Get(detailUrl)
	if resp.StatusCode != 204 {
		fmt.Println("Delete author did not delete the author in " + detailUrl + " , " + resp.Status)
		t.FailNow()
	}
	fmt.Println("OK author Deleted")

}