package AllBooks


import (
        "net/http"
	"github.com/gocql/gocql"
        "encoding/json"
        "insertget/Cassandra"
        "fmt"
)




func Get(w http.ResponseWriter, r *http.Request) {

	var showbooks []Books

	b := map[string]interface{}{}

	iter := Cassandra.Session.Query("SELECT * FROM books").Iter()
	for iter.MapScan(b) {
		showbooks = append(showbooks, Books{
			ID:        b["id"].(gocql.UUID),
			AuthorFirstName: b["authorfirstname"].(string),
			AuthorLastName:  b["authorlastname"].(string),
			BookName: b["bookname"].(string),
			Price:       b["price"].(int),
		})
		b = map[string]interface{}{}
	}

	Conv, _ := json.MarshalIndent(showbooks, "", " ")
	fmt.Fprintf(w, "%s", string(Conv))

}
