package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/avinash-gautam-ios/go-easy-mongo-rest-api/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const databaseName = "users-mongo-db-learn"
const usersCollection = "users"

type UserController struct {
	dbClient *mongo.Client
}

func NewUserController(s *mongo.Client) *UserController {
	uc := UserController{
		dbClient: s,
	}
	return &uc
}

func (uc UserController) GetAllUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	/// get the collection
	coll := uc.dbClient.Database(databaseName).Collection(usersCollection)

	/// create the context
	ctx := context.Background()

	/// query db
	filter := bson.M{} // way to tell mongo to get all the items
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	/// get all records
	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	/// preapre json data from bson
	usersJson, err := json.Marshal(users)
	if err != nil {
		log.Fatal(err)
	}

	/// write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "users: %s", string(usersJson))
}

func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	/// create the context
	ctx := context.Background()

	/// get the collection
	coll := uc.dbClient.Database(databaseName).Collection(usersCollection)
	filter := bson.M{"_id": userId}

	var user models.User
	result := coll.FindOne(ctx, filter)

	if err := result.Decode(&user); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	/// marshall to json
	userJson, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusNotFound)
	}

	/// send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", userJson)
}

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	userId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	/// create the context
	ctx := context.Background()

	/// get the collection
	coll := uc.dbClient.Database(databaseName).Collection(usersCollection)
	filter := bson.M{"_id": userId}

	result, err := coll.DeleteOne(ctx, filter, options.Delete())
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User Delete with id = %s, delete count = %d", id, result.DeletedCount)
}

func (uc UserController) DeleteAllUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := context.Background()
	coll := uc.dbClient.Database(databaseName).Collection(usersCollection)
	filter := bson.D{}
	result, err := coll.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err, "\n")
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "all users deleted, total = %d", result.DeletedCount)
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := models.User{}

	/// get the received body in user
	json.NewDecoder(r.Body).Decode(&user)

	/// create new object id
	oih := primitive.NewObjectID()
	user.Id = oih

	options := options.InsertOne()
	database := uc.dbClient.Database("users-mongo-db-learn")
	collection := database.Collection("users")
	_, err := collection.InsertOne(context.Background(), user, options)
	if err != nil {
		log.Fatal("Error inserting recording in the database")
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "User created successfully with id = %s \n", oih)
}

func (uc UserController) Ping(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintf(w, "Pinged! Serving is running")
}
