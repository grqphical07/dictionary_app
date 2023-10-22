package main

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/django/v3"
)

var cache Cache

const CACHE_FILE string = "cache.json"

// / Returns the homepage
func Index(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

// / calls the dictionary API and renders the content as HTML to be used by HTMX
func DefinitionRoute(c *fiber.Ctx) error {
	word := c.FormValue("word")
	cacheResult, inCache := cache[word]

	var result DictionaryResult
	var err error

	if inCache {
		result = *cacheResult
	} else {
		result, err = GetDefinition(word)
		if err != nil {
			return c.Render("errors/word_not_found", fiber.Map{"word": word})
		}

		cache.AddItem(word, &result)
		if err != nil {
			log.Println(err)
		}
	}

	return c.Render("definition", fiber.Map{"result": result})
}

func main() {
	// I prefer django's templating language over Go's
	engine := django.New("./templates", ".html")

	// Create a new instance of our cache
	var err error
	cache, err = CacheFromFile(CACHE_FILE)
	if err != nil {
		log.Println("Failed to load cache from disk. Resulting to in-memory cache")
	}

	app := fiber.New(
		fiber.Config{
			Views: engine,
		})
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${method} ${path}\n",
		TimeFormat: time.DateTime,
		TimeZone:   "Canada/Vancouver",
	}))

	app.Get("/", Index)
	app.Get("/word", DefinitionRoute)

	app.Static("/static", "./static")

	app.Listen("127.0.0.1:8000")
}
