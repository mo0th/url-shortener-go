package main

import (
	"fmt"
	"github.com/alioygur/is"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"os"
	"strings"
)

type ShortUrl struct {
	Short string `json:"short"`
	Url   string `json:"url"`
}

func main() {
	app := fiber.New()

	urls := make(map[string]ShortUrl)

	if os.Getenv("GO_ENV") != "production" {
		app.Use(middleware.Logger())
	}

	api := app.Group("/")

	api.Get("/:short", func(c *fiber.Ctx) {
		val, ok := urls[c.Params("short")]

		if ok {
			c.Redirect(val.Url)
			return
		}

		c.Status(404).Send("Not Found")
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

		urls[shortUrl.Short] = *shortUrl
		c.JSON(shortUrl)
	})

	app.Static("/", "./public/")

	app.Listen(3000)
}
