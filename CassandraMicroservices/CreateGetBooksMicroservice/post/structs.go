package Books

import (
	"github.com/gocql/gocql"
)


type Books struct {
	ID	        	gocql.UUID `json:"id"`
	AuthorFirstName 	string `json:"authorfirstname"`
	AuthorLastName 		string `json:"authorlastname"`
	BookName	     	string `json:"bookname"`
	Price       		int `json:"price"`
}


type BookID struct {
	ID gocql.UUID `json:"id"`
}

