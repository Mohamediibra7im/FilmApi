# FilmApi

FilmApi is a simple API designed to manage and retrieve information about films. It provides endpoints for creating, reading, updating, and deleting film data.

## Features

- Add new films with details like title, director, release year, and genre.
- Retrieve a list of all films or search for specific films.
- Update film information.
- Delete films from the database.

## Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/your-username/FilmApi.git
    ```
2. Navigate to the project directory:
    ```bash
    cd FilmApi
    ```
3. Install dependencies:
    ```bash
    npm install
    ```

## Usage

1. Start the server:
    ```bash
    npm start
    ```
2. Access the API at `http://localhost:3000`.

## API Endpoints

### Films

- `GET /films` - Retrieve all films.
- `GET /films/:id` - Retrieve a specific film by ID.
- `POST /films` - Add a new film.
- `PUT /films/:id` - Update a film by ID.
- `DELETE /films/:id` - Delete a film by ID.