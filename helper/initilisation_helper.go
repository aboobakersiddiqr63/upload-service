package helper

import (
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB
var Log = logrus.New()
var Client *azblob.Client

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

func InitStorageConnection() {
	var connectionString string
	cloudProvider := os.Getenv("CLOUD_PROVIDER")

	connectionString = os.Getenv("STORAGE_ACCOUNT_CONN_STRING")

	switch cloudProvider {
	case "Azure":
		Client, _ = azblob.NewClientFromConnectionString(connectionString, nil)
	}

	Log.Infoln("Initialised Storage Connection")
}
