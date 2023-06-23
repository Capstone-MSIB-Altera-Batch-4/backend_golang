package config

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var (
	Db         *gorm.DB
	serverPort string
	JWT        string
	ServerKey  string
	ClientKey  string
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Read environment variables
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	port := os.Getenv("DB_PORT")
	data := os.Getenv("DB_NAME")

	JWT = os.Getenv("KEY_SECRET")
	ServerKey = os.Getenv("SERVER_KEY")
	ClientKey = os.Getenv("CLIENT_KEY")

	// Setup database connection
	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + data + "?charset=utf8mb4&parseTime=True&loc=Asia%2FJakarta"
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	Db = conn
	serverPort = os.Getenv("APP_PORT")
}

func GetServer() string {
	return serverPort
}
