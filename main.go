package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

type Film struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	FilmType string `json:"filmtype"`
}

var db *pgxpool.Pool
var logger = logrus.New()

func connectDB() {
	var err error
	db, err = pgxpool.Connect(context.Background(), "postgres://postgres:2003@localhost:5432/filmdb")
	if err != nil {
		logger.WithField("error", err.Error()).Fatal("Unable to connect to database")
	}
	logger.Info("Connected to database")
}

func getFilms(c *fiber.Ctx) error {
	rows, err := db.Query(context.Background(), "SELECT id, title, year, filmtype FROM films")
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to fetch films")
		return c.Status(http.StatusInternalServerError).SendString("Failed to fetch films")
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var film Film
		if err := rows.Scan(&film.Id, &film.Title, &film.Year, &film.FilmType); err != nil {
			logger.WithField("error", err.Error()).Error("Failed to scan film")
			return c.Status(http.StatusInternalServerError).SendString("Failed to fetch films")
		}
		films = append(films, film)
	}

	return c.JSON(films)
}

func getFilmByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Invalid ID")
		return c.Status(http.StatusBadRequest).SendString("Invalid ID")
	}

	var film Film
	err = db.QueryRow(context.Background(), "SELECT id, title, year, filmtype FROM films WHERE id=$1", intID).
		Scan(&film.Id, &film.Title, &film.Year, &film.FilmType)
	if err != nil {
		logger.WithField("error", err.Error()).Warn("Film not found")
		return c.Status(http.StatusNotFound).SendString("Film not found")
	}

	return c.JSON(film)
}

func addFilm(c *fiber.Ctx) error {
	film := new(Film)
	if err := c.BodyParser(film); err != nil {
		logger.WithField("error", err.Error()).Error("Invalid film data")
		return c.Status(http.StatusBadRequest).SendString("Invalid film data")
	}

	err := db.QueryRow(context.Background(),
		"INSERT INTO films (title, year, filmtype) VALUES ($1, $2, $3) RETURNING id",
		film.Title, film.Year, film.FilmType).Scan(&film.Id)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Failed to add film")
		return c.Status(http.StatusInternalServerError).SendString("Failed to add film")
	}

	logger.WithField("film", film).Info("Film added")
	return c.JSON(film)
}

func updateFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Invalid ID")
		return c.Status(http.StatusBadRequest).SendString("Invalid ID")
	}

	film := new(Film)
	if err := c.BodyParser(film); err != nil {
		logger.WithField("error", err.Error()).Error("Invalid film data")
		return c.Status(http.StatusBadRequest).SendString("Invalid film data")
	}

	commandTag, err := db.Exec(context.Background(),
		"UPDATE films SET title=$1, year=$2, filmtype=$3 WHERE id=$4",
		film.Title, film.Year, film.FilmType, intID)
	if err != nil || commandTag.RowsAffected() == 0 {
		logger.WithField("error", err.Error()).Warn("Film not found for update")
		return c.Status(http.StatusNotFound).SendString("Film not found")
	}

	logger.WithField("film", film).Info("Film updated")
	return c.JSON(film)
}

func deleteFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		logger.WithField("error", err.Error()).Error("Invalid ID")
		return c.Status(http.StatusBadRequest).SendString("Invalid ID")
	}

	commandTag, err := db.Exec(context.Background(), "DELETE FROM films WHERE id=$1", intID)
	if err != nil || commandTag.RowsAffected() == 0 {
		logger.WithField("error", err.Error()).Warn("Film not found for deletion")
		return c.Status(http.StatusNotFound).SendString("Film not found")
	}

	logger.WithField("id", intID).Info("Film deleted")
	return c.SendString("Film is deleted")
}

func main() {
	connectDB()
	defer db.Close()

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
