package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/sirupsen/logrus"
)

type Films struct {
	Id       int
	Title    string
	Year     int
	FilmType string
}

var films = []Films{
	{Id: 1, Title: "The Shawshank Redemption", Year: 1994, FilmType: "Drama"},
	{Id: 2, Title: "The Godfather", Year: 1972, FilmType: "Crime"},
	{Id: 3, Title: "The Dark Knight", Year: 2008, FilmType: "Action"},
	{Id: 4, Title: "Nun", Year: 2018, FilmType: "Horror"},
	{Id: 5, Title: "The Lord of the Rings: The Return of the King", Year: 2003, FilmType: "Fantasy"},
}

var logger = logrus.New()

func getFilms(c *fiber.Ctx) error {
	logger.Info("Fetching all films")
	return c.JSON(films)
}

func getFilmByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Invalid ID")
		return c.Status(http.StatusBadRequest).SendString("Invalid ID")
	}
	for _, film := range films {
		if film.Id == intID {
			logger.WithField("id", intID).Info("Film found")
			return c.JSON(film)
		}
	}
	logger.WithField("id", intID).Warn("Film not found")
	return c.Status(http.StatusNotFound).SendString("Film not found")
}

func addFilm(c *fiber.Ctx) error {
	film := new(Films)
	if err := c.BodyParser(film); err != nil {
		logger.WithField("error", err.Error()).Error("Invalid film data")
		return c.Status(http.StatusBadRequest).SendString("Invalid film")
	}
	films = append(films, *film)
	logger.WithField("film", film).Info("Film added")
	return c.JSON(film)
}

func updateFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"id":    id,
		}).Error("Invalid ID")
		return c.Status(http.StatusBadRequest).SendString("Invalid ID")
	}
	film := new(Films)
	if err := c.BodyParser(film); err != nil {
		logger.WithField("error", err.Error()).Error("Invalid film data")
		return c.Status(http.StatusBadRequest).SendString("Invalid film")
	}
	for i, f := range films {
		if f.Id == intID {
			films[i] = *film
			logger.WithField("film", film).Info("Film updated")
			return c.JSON(film)
		}
	}
	logger.WithField("id", intID).Warn("Film not found for update")
	return c.Status(http.StatusNotFound).SendString("Film not found")
}

func deleteFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	for i, film := range films {
		intID, err := strconv.Atoi(id)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"id":    id,
			}).Error("Invalid ID")
			return c.Status(http.StatusBadRequest).SendString("Invalid ID")
		}
		if film.Id == intID {
			films = append(films[:i], films[i+1:]...)
			logger.WithField("id", intID).Info("Film deleted")
			return c.SendString("Film is deleted")
		}
	}
	logger.WithField("id", id).Warn("Film not found for deletion")
	return c.Status(http.StatusNotFound).SendString("Film not found")
}

func main() {
	app := fiber.New()

	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	/*
	using Sliding Window Rate Limiting instead of Fixed Window Rate Limiting
	-> becauese it allows for a more flexible and responsive rate limiting strategy.
	-> It can help to smooth out spikes in traffic and provide a more consistent user experience.
	
	---> The rate limit is set to 4 requests per 20 seconds.
	*/
	app.Use(limiter.New(limiter.Config{
		Max:        4,
		Expiration: 20 * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			logger.Warn("Rate limit exceeded")
			return c.Status(http.StatusTooManyRequests).SendString("Rate limit exceeded")
		},
	}))

	app.Get("/films", getFilms)
	app.Get("/films/:id", getFilmByID)
	app.Post("/films", addFilm)
	app.Put("/films/:id", updateFilm)
	app.Delete("/films/:id", deleteFilm)

	logger.Info("Starting server on :3000")
	app.Listen(":3000")
}