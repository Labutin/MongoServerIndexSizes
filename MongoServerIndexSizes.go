package main

import (
	"flag"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"fmt"
)

// getMongoDBConnection init MongoDB connection
func getMongoDBConnection(connectionString string) (*mgo.Session, error)  {
	session, err := mgo.Dial(connectionString)
	return session, err
}

func collectAddIndexSize(session *mgo.Session) float64 {
	indexSize := float64(0)
	databases, _ := session.DatabaseNames()
	for _, db := range databases {
		result := bson.M{}
		session.DB(db).Run("dbstats", &result)
		t := result["indexSize"]
		fmt.Printf("%s: %d\n", db, int64(t.(float64)))
		indexSize += t.(float64)
	}
	return indexSize
}

func main() {
	flagMongoConnectString := flag.String("mongoConnectionString", "127.0.0.1", "Mongo connection string")
	flag.Parse()
	session, err := getMongoDBConnection(*flagMongoConnectString)
	if err != nil {
		log.Fatalf("Can't connect to MongoDB (%s): %v", *flagMongoConnectString, err)
	}
	indexSize := collectAddIndexSize(session)
	fmt.Printf("----------------------\nTotal: %d (%.2f GB)\n", int64(indexSize), (indexSize / 1024 / 1024 / 1024))
}
