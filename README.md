# FilmApi

FilmApi is a simple API designed to manage and retrieve information about films. It provides endpoints for creating, reading, updating, and deleting film data.

## Features


- Add new films with details like title, director, release year, and genre.
- Retrieve a list of all films or search for specific films.
- Update film information.
- Delete films from the database.

## API Endpoints

### Films

- `GET /films` - Retrieve all films.
- `GET /films/:id` - Retrieve a specific film by ID.
- `POST /films` - Add a new film.
- `PUT /films/:id` - Update a film by ID.
- `DELETE /films/:id` - Delete a film by ID.