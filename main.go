package main

import (
	"fmt"
	"github.com/alioygur/is"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"strings"
)

type ShortUrl struct {
	gorm.Model
	Short string `json:"short" gorm:"unique_index:short"`
	Url   string `json:"url"`
}

func findShortUrl(db *gorm.DB, short string) *ShortUrl {
	res := new(ShortUrl)
	db.Where("short = ?", short).First(res)

	if res.Short == "" || res.Url == "" {
		return nil
	}

	return res
}

func main() {
	db, err := gorm.Open("sqlite3", "./urls.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&ShortUrl{})

	app := fiber.New()

	if os.Getenv("GO_ENV") != "production" {
		app.Use(middleware.Logger())
	}

	api := app.Group("/")

	api.Get("/:short", func(c *fiber.Ctx) {

		res := findShortUrl(db, c.Params("short"))

		if res == nil {
			c.Status(404).JSON(fiber.Map{
				"success": false,
				"message": "Short not found",
			})
			return
		}

		c.Redirect(res.Url)
	})

	api.Post("/", func(c *fiber.Ctx) {
		shortUrl := new(ShortUrl)
		if err := c.BodyParser(shortUrl); err != nil {
			c.Status(503).Send(err)
			return
		}

		// Check short length
		if len(shortUrl.Short) < 2 || len(shortUrl.Short) > 8 {
			c.Status(406).JSON(fiber.Map{
				"success": false,
				"message": "Short must be between 2 and 8 long.",
			})
			return
		}

		// Check that url is valid
		if !is.URL(shortUrl.Url) {
			c.Status(406).JSON(fiber.Map{
				"success": false,
				"message": "Url must be a valid url",
			})
			return
		}

		// Add protocol if there isn't one
		if !strings.HasPrefix(shortUrl.Url, "http://") || !strings.HasPrefix(shortUrl.Url, "https://") {
			shortUrl.Url = fmt.Sprintf("%s%s", "https://", shortUrl.Url)
		}

		if findShortUrl(db, shortUrl.Short) != nil {
			c.Status(400).JSON(fiber.Map{
				"success": false,
				"message": "Short already in use",
			})
			return
		}

		db.Create(shortUrl)

		c.JSON(fiber.Map{
			"success": true,
			"short":   shortUrl.Short,
			"url":     shortUrl.Url,
		})
	})

	app.Static("/", "./public/")

	app.Listen(3000)
}
