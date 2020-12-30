package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/tempcke/rpm/api/rest"
	"github.com/tempcke/rpm/repository"
)

var propRepo = repository.NewInMemoryRepo()
var server = rest.NewServer(propRepo)

func main() {
	port := ":8080"
	fmt.Println("Listening on " + port)
	err := http.ListenAndServe(port, server)
	if err != nil {
		log.Fatal(err.Error())
	}
}
