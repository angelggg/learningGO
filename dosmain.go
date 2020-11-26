package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/gofiber/fiber"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Article struct {
	Title string `json:"Title"`
	Description string `json:"Description"`
	Content string `json:"Content"`
	Price float32 `json:"Price"`

}

type Articles []Article

func getAllArticles(writer http.ResponseWriter, request *http.Request){
	fmt.Print("All")
	articles := Articles{
		Article{
			Title:       "Test",
			Description: "Descripto",
			Content:     "contnido",
			Price:       22.4,
		},
		Article{
			Title:       "seg",
			Description: "asdsa",
			Content:     "assad",
			Price:       20,
		},
	}
	json.NewEncoder(writer).Encode(articles)

}
func homePage(writer http.ResponseWriter, request *http.Request){
	fmt.Fprintf(writer, "Estamos")
}
func testPost(writer http.ResponseWriter, request *http.Request){
	fmt.Fprintf(writer, "Estamos")

}

func handleRequests(){
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.HandleFunc("/articles", getAllArticles).Methods("GET")
	router.HandleFunc("/articles", testPost).Methods("POST")

	log.Fatal(http.ListenAndServe(":8081", router))

}
func nomain(){
	db, err := sql.Open("mysql", "nombre_usuario:tu_contrasena@tcp(localhost:3306)/test")
	if err != nil{
		panic(err.Error())
	}
	fmt.Print("Conectado")
	query := "SELECT COUNT(*) FROM tasks;"
	fmt.Print( db.Exec(query) )
	db.Close()
	handleRequests()
}

