package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/gofiber/cors"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	db, err := gorm.Open("sqlite3", "./urls.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&ShortUrl{})

	app := fiber.New()

	app.Use(middleware.Logger())

	app.Use(cors.New())

	if os.Getenv("NO_PUBLIC") != "true" {
		app.Static("/", "./public/")
	}

	api := app.Group("/")

	setupUrlRoutes(api, db)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "3000"
	}

	app.Listen(PORT)
}
