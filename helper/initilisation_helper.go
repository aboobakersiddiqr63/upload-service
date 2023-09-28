package helper

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB
var Log = logrus.New()

func InitDB() {
	LoadEnv()
	dsn := os.Getenv("DB_URI")
	fmt.Println("Connected to authdb")

	var err error
	Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

func LoadEnv() {
	err := godotenv.Load(".env")
	HandleException(err, "Loading Env")
}

func InitLogger() {
	Log.SetLevel(logrus.InfoLevel)
	Log.SetFormatter(&logrus.JSONFormatter{})
}
