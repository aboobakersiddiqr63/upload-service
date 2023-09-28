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
	fmt.Println("Starting the server on port 4001")
	helper.InitLogger()
	helper.InitDB()
	helper.InitStorageConnection()
	log.Fatal(http.ListenAndServe(":4001", app))
}
