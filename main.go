// package main

// import (
// 	"context"
// 	"craft/api"
// 	"craft/config"
// 	"craft/database"
// 	"fmt"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

// func main() {
// 	fmt.Println("hello mayank1")
// 	config := config.GetConfig()
// 	bgctx := context.Background()
// 	fmt.Println("config :", config)
// 	db := database.ConnectDB(bgctx, config.Mongo)
// 	collection := db.Collection(config.Mongo.Collection)
// 	client := &database.MongoUtilsClient{
// 		Ctx:  bgctx,
// 		Coll: collection,
// 	}

// 	fmt.Println("db :", db)
// 	r := mux.NewRouter()
// 	r.HandleFunc("crafts", api.InsertCraftConfig(client)).Methods("POST")
// 	http.ListenAndServe(":80", r)

// }

package main

import (
	"context"
	"craft/api"
	"craft/config"
	"craft/database"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	port = ":8080"
)

func main() {
	fmt.Println("Swrvice up and running on  http://localhost:2021 ...")
	fmt.Println("hello team")
	config := config.GetConfig()
	bgctx := context.Background()
	db := database.ConnectDB(bgctx, config.Mongo)
	collection := db.Collection(config.Mongo.Collection)
	client := &database.MongoUtilsClient{
		Ctx:  bgctx,
		Coll: collection,
	}
	router := mux.NewRouter()
	api.RegisterRoutes(router, client)
	log.Fatal(http.ListenAndServe(port, router))

}
