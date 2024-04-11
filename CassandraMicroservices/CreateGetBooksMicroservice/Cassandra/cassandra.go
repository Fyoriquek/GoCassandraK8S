package Cassandra 


import (
	"github.com/gocql/gocql"
	"fmt"
	"os"
)


var (
        cassandraHost     string = os.Getenv("CASSANDRA_HOST")
        cassandraUser     string = os.Getenv("CASSANDRA_USER")
        cassandraPassword string = os.Getenv("CASSANDRA_PASSWORD")

)

var Session *gocql.Session

func init() {
	var err error

	cluster := gocql.NewCluster(cassandraHost)
	cluster.Keyspace = "stats"
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("Cassandra Session created...")
}


