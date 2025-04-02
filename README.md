# FilmApi

FilmApi is a robust and user-friendly API designed to manage and retrieve information about films. It offers a comprehensive set of endpoints for creating, reading, updating, and deleting film data, making it an ideal solution for film cataloging applications.

## Features

- Add new films with details such as title, director, release year, and genre.
- Retrieve a list of all films or search for specific films using filters.
- Update film information seamlessly.
- Delete films from the database with ease.
- Supports pagination and sorting for large datasets.
- Provides detailed error messages for better debugging.

## API Endpoints

### Films

- `GET /films` - Retrieve all films with optional query parameters for filtering, sorting, and pagination.
- `GET /films/:id` - Retrieve detailed information about a specific film by its ID.
- `POST /films` - Add a new film to the database. Requires a JSON payload with film details.
- `PUT /films/:id` - Update an existing film's information by its ID. Requires a JSON payload with updated details.
- `DELETE /films/:id` - Remove a film from the database by its ID.

### Example Usage

#### Retrieve All Films
```bash
GET /films
```

#### Add a New Film
```bash
POST /films
Content-Type: application/json

{
    "title": "Inception",
    "director": "Christopher Nolan",
    "releaseYear": 2010,
    "genre": "Science Fiction"
}
```

#### Update a Film
```bash
PUT /films/1
Content-Type: application/json

{
    "title": "Inception",
    "director": "Christopher Nolan",
    "releaseYear": 2010,
    "genre": "Sci-Fi"
}
```

#### Delete a Film
```bash
DELETE /films/1
```

## Getting Started

To get started with FilmApi, clone the repository and follow the setup instructions in the documentation.
