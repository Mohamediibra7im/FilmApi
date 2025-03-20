package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupFiberApp() *fiber.App {
	app := fiber.New()
	app.Get("/films", getFilmsFiber)
	app.Get("/films/:id", getFilmByIDFiber)
	return app
}

func getFilmsFiber(c *fiber.Ctx) error {
	return c.JSON(films)
}

func getFilmByIDFiber(c *fiber.Ctx) error {
	id := c.Params("id")
	for _, g := range films {
		if strconv.Itoa(g.Id) == id {
			return c.JSON(g)
		}
	}
	return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "Film not found"})
}

func TestFiberGetFilm(t *testing.T) {
	app := setupFiberApp()
	req := httptest.NewRequest("GET", "/films", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberGetFilmByID(t *testing.T) {
	app := setupFiberApp()
	req := httptest.NewRequest("GET", "/films/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestFiberGetFilmByInvalidID(t *testing.T) {
	app := setupFiberApp()
	req := httptest.NewRequest("GET", "/films/999", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}