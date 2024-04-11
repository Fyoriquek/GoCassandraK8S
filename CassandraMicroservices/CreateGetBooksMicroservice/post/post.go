package Books

import (
	"net/http"
	"github.com/gocql/gocql"
	"encoding/json"
	"context"
	"insertget/Cassandra"
	"fmt"
	"io/ioutil"
)




func Post(w http.ResponseWriter, r *http.Request) {


	var NewBook Books

	var gocqlUUID gocql.UUID
        gocqlUUID = gocql.TimeUUID()
        json.NewEncoder(w).Encode(BookID{ID: gocqlUUID})


	reqBody, err := ioutil.ReadAll(r.Body)
        if err != nil {
                fmt.Fprintf(w, "wrong data")
        }



        json.Unmarshal(reqBody, &NewBook)


        if err := Cassandra.Session.Query(`INSERT INTO books(id, authorfirstname, authorlastname, bookname, price) VALUES(?, ?, ?, ?, ?);`,
	               		 	gocqlUUID,
       			         	NewBook.AuthorFirstName,
                	        	NewBook.AuthorLastName,
                			NewBook.BookName,
                			NewBook.Price,
                			).WithContext(context.Background()).Exec(); err != nil {
                			fmt.Println("insert error")

        		}

}




