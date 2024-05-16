package database

import (
	"fmt"
	"log"
	"os"

	"github.com/almirpernen/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func ConnectDb() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_ADDRESS"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("Connected")

	log.Println("Running migrations")
	db.AutoMigrate(&models.Post{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Comment{})
	db.AutoMigrate(&models.PostLike{})
	db.AutoMigrate(&models.CommentLike{})

	DB = Dbinstance{
		Db: db,
	}
}
