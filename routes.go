package main

import (
	"fmt"
	"github.com/alioygur/is"
	"github.com/gofiber/fiber"
	"github.com/jinzhu/gorm"
	"strings"
)

func setupUrlRoutes(router fiber.Router, db *gorm.DB) {
	router.Get("/:short", func(c *fiber.Ctx) {
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

	router.Post("/", func(c *fiber.Ctx) {
		shortUrl := new(ShortUrl)
		if err := c.BodyParser(shortUrl); err != nil {
			c.Status(503).Send(err)
			return
		}

		// Make random short
		if shortUrl.Short == "" {
			shortUrl.Short = randShort(5)
			for findShortUrl(db, shortUrl.Short) != nil {
				shortUrl.Short = randShort(5)
			}
		}

		// Check short length
		if (len(shortUrl.Short) < 2) || len(shortUrl.Short) > 8 {
			c.Status(406).JSON(fiber.Map{
				"success": false,
				"message": "Short must be between 2 and 8 long",
			})
			return
		}

		// Check that short only contains alphanumeric characters
		if !is.Alphanumeric(shortUrl.Short) {
			c.Status(406).JSON(fiber.Map{
				"success": false,
				"message": "Short can only contain alphanumeric characters",
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
}
