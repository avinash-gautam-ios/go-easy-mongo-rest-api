package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/avinash-gautam-ios/go-easy-mongo-rest-api/controllers"
	"github.com/julienschmidt/httprouter"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	/// create new router
	r := httprouter.New()

	/// create controller
	ctx := context.TODO()
	client := getMongoClient(ctx)
	uc := controllers.NewUserController(client)

	r.GET("/ping", uc.Ping)

	r.GET("/user/:id", uc.GetUser)
	r.GET("/user", uc.GetAllUser)

	r.POST("/user", uc.CreateUser)

	r.DELETE("/user/:id", uc.DeleteUser)
	r.DELETE("/user", uc.DeleteAllUsers)

	/// start the server
	port := ":8080"
	// start and listen for error if failed to start
	err := http.ListenAndServe(port, r)
	if err != nil {
		panic(err)
	}

	/// close the connection when main function finished
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func getMongoClient(ctx context.Context) *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("MongoDB: Connected \n")
	return client
}
