package main

import (
	"DynamoDBLocal/pkg/server"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	ds, err := server.NewDynamoLocal("test.db", []byte("#"))
	if err != nil {
		log.Fatal(err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/", ds.RunDynamo)

	if err := http.ListenAndServe("localhost:8080", router); err != nil {
		log.Fatal(err)
	}
}
