package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"insertget/post"
	"insertget/Cassandra"
	"insertget/get"
)



func main() {


	cassSession := Cassandra.Session


        defer cassSession.Close()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/books/new", Books.Post)
	router.HandleFunc("/books/all", AllBooks.Get)

	log.Fatal(http.ListenAndServe(":8083", router))



}


