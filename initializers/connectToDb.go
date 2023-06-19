package initializers

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func ConnectToDB() {
	dsn := os.Getenv("DB_URL")

	// https://github.com/go-gorm/postgres
	fmt.Print("Connecting to DB")

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("failed to connect database")
	}
	fmt.Printf("DB: %v\n", DB)
	fmt.Printf("DB: %v\n", err)

}
