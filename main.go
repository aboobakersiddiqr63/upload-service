package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aboobakersiddiqr63/upload-service/helper"
	router "github.com/aboobakersiddiqr63/upload-service/routes"
)

func main() {
	app := router.Router()
	fmt.Println("Starting the server on port 4000")
	helper.InitLogger()
	helper.InitDB()
	log.Fatal(http.ListenAndServe(":4001", app))
}
