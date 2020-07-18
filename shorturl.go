package main

import (
	"github.com/jinzhu/gorm"
	"math/rand"
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

const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randShort(length int) string {
	res := make([]byte, length)
	numChars := len(allowedChars)
	for i := range res {
		res[i] = allowedChars[rand.Intn(numChars)]
	}
	return string(res)
}
