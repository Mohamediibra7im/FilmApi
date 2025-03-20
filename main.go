package main

import (
	"net/http"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type Films struct {
	Id	int
	Title string
	Year  int
	FilmType string
}

var films = []Films{
	{Id: 1, Title: "The Shawshank Redemption", Year: 1994, FilmType: "Drama"},
	{Id: 2, Title: "The Godfather", Year: 1972, FilmType: "Crime"},
	{Id: 3, Title: "The Dark Knight", Year: 2008, FilmType: "Action"},
	{Id: 4, Title: "Nun", Year: 2018, FilmType: "Horror"},
	{Id: 5, Title: "The Lord of the Rings: The Return of the King", Year: 2003, FilmType: "Fantasy"},
}

func getFilms(c *fiber.Ctx) error {
	return c.JSON(films)
}

func getFilmByID(c *fiber.Ctx) error {
	id := c.Params("id")
intID, err := strconv.Atoi(id)
if err != nil {
	return c.Status(http.StatusBadRequest).SendString("Invalid ID")
}
for _, film := range films {
	if film.Id == intID {
		return c.JSON(film)
	}
}
return c.Status(http.StatusNotFound).SendString("Film not found")
}

func addFilm(c *fiber.Ctx) error {
	film := new(Films)
	if err := c.BodyParser(film); err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid film")
	}
	films = append(films, *film)
	return c.JSON(film)
}

func updateFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid ID")
	}
	film := new(Films)
	if err := c.BodyParser(film); err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid film")
	}
	for i, f := range films {
		if f.Id == intID {
			films[i] = *film
			return c.JSON(film)
		}
	}
	return c.Status(http.StatusNotFound).SendString("Film not found")
}

func deleteFilm(c *fiber.Ctx) error {
	id := c.Params("id")
	for i, film := range films {
		intID, err := strconv.Atoi(id)
		if err != nil {
			return c.Status(http.StatusBadRequest).SendString("Invalid ID")
		}
		if film.Id == intID {
			films = append(films[:i], films[i+1:]...)
			return c.SendString("Film is deleted")
		}
	}
	return c.Status(http.StatusNotFound).SendString("Film not found")
}


func main() {
    app := fiber.New()

	app.Get("/films", getFilms)

	app.Get("/films/:id", getFilmByID)
	app.Post("/films", addFilm)

	app.Put("/films/:id", updateFilm)
	app.Delete("/films/:id", deleteFilm)

    app.Listen(":3000")
}