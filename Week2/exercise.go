package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Mongo struct {
	client *mongo.Client
	url    string
}

type Students struct {
	Students []Student `json:"students"`
}

type Student struct {
	Name        string   `json:"name"`
	Age         int      `json: "age"`
	CompanyList []string `json:"companyList"`
}

var s Students
var m Mongo

func (s *Students) parse(file string) ([]byte, error) {
	// os.Open read file
	jsonFile, err := os.Open(file)
	// close file
	defer jsonFile.Close()
	// os.Open return error
	if err != nil {
		fmt.Println("error in read file: ", err)
		return nil, err
	}
	// parse io.Reader to []byte
	bs, _ := ioutil.ReadAll(jsonFile)
	// save data of json file to s' struct
	err = json.Unmarshal(bs, &s)
	if err != nil {
		fmt.Println("parse file to struct error: ", err)
		return bs, err
	}
	return bs, nil
}

func saveStudentsToDb(w http.ResponseWriter, r *http.Request) {
	// initialization collection on database
	collection := m.client.Database("test").Collection("trainers")
	// insert s to db
	insertResult, err := collection.InsertOne(context.TODO(), s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult)

	bs, error := json.Marshal(map[string]interface{}{
		"status": "OK",
	})

	if error != nil {
		fmt.Println("error: ", err)
	}

	w.Header().Set("Author", "NDLPhat")
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	bs, err := s.parse("students.json")
	if err != nil {
		fmt.Println("error: ", err)
	}
	w.Header().Set("Author", "NDLPhat")
	w.Header().Set("Content-Type", "application/json")
	w.Write(bs)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Read port from env
	port := os.Getenv("PORT")
	mongo_url := os.Getenv("MONGO_URL")

	// Route
	http.HandleFunc("/students", getStudents)
	http.HandleFunc("/save-student-to-db", saveStudentsToDb)

	// Mongo initialization
	clientOptions := options.Client().ApplyURI(mongo_url)
	client, MongoErr := mongo.Connect(context.TODO(), clientOptions)
	if MongoErr != nil {
		log.Fatal("error Mongo", err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	// save info to mongo struct
	m.url = mongo_url
	m.client = client

	// Start on local on port
	serverErr := http.ListenAndServe(":"+port, nil)
	if serverErr != nil {
		log.Println("Error starting server")
	} else {
		fmt.Println("Started on server localhost port: ", port)
	}
}
