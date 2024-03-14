package main

import (
	"log"
	"main/db"
	"main/handler"
	"main/router"
	"net/http"
)

func main() {
	dbSettings := `host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable`
	repos, err := db.GetConnection(dbSettings)
	if err != nil {
		log.Println(err)
		return
	}
	responses := handler.GetHandler(repos)
	rout := router.GetRouter(responses)

	err = http.ListenAndServe("localhost:8080", rout)
	if err != nil {
		log.Println(err)
		return
	}
}
