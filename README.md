# Build a Fullstack App in Vanilla JS & Go Course
This is a companion repository for the [Build a Fullstack App in Vanilla JS & Go](https://frotnendmasters.com/courses/vanilla-js-go/) course on Frontend Masters.
[![Frontend Masters](https://static.frontendmasters.com/assets/brand/logos/full.png)](https://frotnendmasters.com/courses/vanilla-js-go/)

The code snippets below are referenced throughout the course so you can either code along with Max Firtman or copy/paste. In the assets folders, you will find a copy of the slides and the final project.

## A-Backend

### A1 - Init

Run `go mod init frontendmasters.com/reelingit` to create the module.

Install the dependencies

```
go get github.com/joho/godotenv
go get github.com/lib/pq
```

Create the *main.go* file.

```go
package main 

func main() {
    // Serve static files
    http.Handle("/", http.FileServer(http.Dir("public")))

    // Start server
    const addr = ":8080"
    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalf("Server failed: %v", err)
    }
}
```

Create a test *index.html*.

### A2 - Logger

Create a logger package with a *logger.go* file

```go
package logger

import (
	"log"
	"os"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	file        *os.File
}

// NewLogger creates a new logger with output to both file and stdout
func NewLogger(logFilePath string) (*Logger, error) {
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		file:        file,
	}, nil
}

// Info logs informational messages to stdout
func (l *Logger) Info(msg string) {
	l.infoLogger.Printf("%s", msg)
}

// Error logs error messages to file
func (l *Logger) Error(msg string, err error) {
	l.errorLogger.Printf("%s: %v", msg, err)
}

// Close closes the log file
func (l *Logger) Close() {
	l.file.Close()
}
```

Now, change *main.go* with:

```go
func main() {
	// Initialize logger
	logInstance := initializeLogger()

	http.Handle("/", http.FileServer(http.Dir("public")))

	// Start server
	const addr = ":8080"
	logInstance.Info("Server starting on " + addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		logInstance.Error("Server failed to start", err)
		log.Fatalf("Server failed: %v", err)
	}
}

func initializeLogger() *logger.Logger {
	logInstance, err := logger.NewLogger("movie-service.log")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	return logInstance
}

```

### A3 - Models

Create the *models* package with the following files:

*genre.go*
```go
package models

type Genre struct {
	ID   int   
	Name string 
}
```

*actor.go*
```go
package models

type Actor struct {
	ID        int     
	FirstName string  
	LastName  string  
	ImageURL  *string 
}

```

*movie.go*
```go
package models

type Movie struct {
	ID          int      
	TMDB_ID     int      
	Title       string   
	Tagline     *string  
	ReleaseYear int      
	Genres      []Genre  
	Overview    *string  
	Score       *float32 
	Popularity  *float32 
	Keywords    []string 
	Language    *string  
	PosterURL   *string  
	TrailerURL  *string  
	Casting     []Actor  
}

```

### A4 - Basic Handlers

Create the *handlers* package with a *movies_handlers.go* file.

```go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"frontendmasters.com/movies/logger"
	"frontendmasters.com/movies/models"
)

type MovieHandler struct {
}

// Utility functions
func (h *MovieHandler) writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *MovieHandler) GetTopMovies(w http.ResponseWriter, r *http.Request) {
    movies := []models.Movie{
		{
			ID:          1,
			TMDB_ID:     101,
			Title:       "The Hacker",
			ReleaseYear: 2022,
			Genres:      []models.Genre{{ID: 1, Name: "Thriller"}},
			Keywords:    []string{"hacking", "cybercrime"},
			Casting:     []models.Actor{{ID: 1, Name: "Jane Doe"}},
		},
		{
			ID:          2,
			TMDB_ID:     102,
			Title:       "Space Dreams",
			ReleaseYear: 2020,
			Genres:      []models.Genre{{ID: 2, Name: "Sci-Fi"}},
			Keywords:    []string{"space", "exploration"},
			Casting:     []models.Actor{{ID: 2, Name: "John Star"}},
		},
		{
			ID:          3,
			TMDB_ID:     103,
			Title:       "The Lost City",
			ReleaseYear: 2019,
			Genres:      []models.Genre{{ID: 3, Name: "Adventure"}},
			Keywords:    []string{"jungle", "treasure"},
			Casting:     []models.Actor{{ID: 3, Name: "Lara Hunt"}},
		},
	}

	if h.writeJSONResponse(w, movies) == nil {
		h.logger.Info("Successfully served top movies")
	}
}

```

Now setup the handler in *main.go*

```go
    // ...
    movieHandler := handlers.NewMovieHandler {}
	// Set up routes
	http.HandleFunc("/api/movies/top", movieHandler.GetTopMovies)
    // ...
```

### A5 - Install AIR

Check instructions at [https://github.com/air-verse/air](https://github.com/air-verse/air), such as executing

```
go install github.com/cosmtrek/air@latest
```

To customize it, create a *.air.toml* file. On Linux and macOS, you can use
```toml
# .air.toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./main.go"
bin = "./tmp/main"
include_ext = ["go"]  # Only watch .go files
exclude_dir = ["tmp", "vendor", "node_modules"]
delay = 1000  # ms

[log]
time = true

[misc]
clean_on_exit = true
```

On Windows, `cmd` & `bin` needs `.exe` after main script name, so it should be:

```toml
# .air.toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main.exe ./main.go"
bin = "./tmp/main.exe"
include_ext = ["go"]  # Only watch .go files
exclude_dir = ["tmp", "vendor", "node_modules"]
delay = 1000  # ms

[log]
time = true

[misc]
clean_on_exit = true
```

## B-Database

### B1 - Import Data

Set up a Postgres database and get a connection string, then, go to *import/install.go* and insert the string there.

Get into the *import* folder, and run `go run install.go`. That should populate your database with all the data.

### B2 - Create the data interface

Create the *data* package and the *interfaces.go* file

```go
package data

import "frontendmasters.com/movies/models"

type MovieStorage interface {
	GetTopMovies() ([]models.Movie, error)
	GetRandomMovies() ([]models.Movie, error)
	GetMovieByID(id int) (models.Movie, error)
	SearchMoviesByName(name string, order string, genre *int) ([]models.Movie, error)
	GetAllGenres() ([]models.Genre, error)
}
```

### B3 - Create the DB connection

Create a *.env* file in the root folder and add the connection string

```
DATABASE_URL=""
```

Open *main.go* and add this in the *main* function

```go
    // Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or failed to load: %v", err)
	}

	// Initialize logger
	logInstance := initializeLogger()

	// Database connection
	dbConnStr := os.Getenv("DATABASE_URL")
	if dbConnStr == "" {
		log.Fatalf("DATABASE_URL not set in environment")
	}
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
```

### B4 - Add field metadata to models

Modify the models to add metadata such as with *models/movie.go*

```go
package models

type Movie struct {
	ID          int      `json:"id"`
	TMDB_ID     int      `json:"tmdb_id,omitempty"`
	Title       string   `json:"title"`
	Tagline     *string  `json:"tagline,omitempty"`
	ReleaseYear int      `json:"release_year"`
	Genres      []Genre  `json:"genres"`
	Overview    *string  `json:"overview,omitempty"`
	Score       *float32 `json:"score,omitempty"`
	Popularity  *float32 `json:"popularity,omitempty"`
	Keywords    []string `json:"keywords"`
	Language    *string  `json:"language,omitempty"`
	PosterURL   *string  `json:"poster_url,omitempty"`
	TrailerURL  *string  `json:"trailer_url,omitempty"`
	Casting     []Actor  `json:"casting"`
}
```

### B5 - Create the Movie Repository

Create *data/movie_repository.go

```go
package data

import (
	"database/sql"
	"errors"
	"strconv"

	"frontendmasters.com/movies/logger"
	"frontendmasters.com/movies/models"
	_ "github.com/lib/pq"
)

type MovieRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewMovieRepository(db *sql.DB, log *logger.Logger) (*MovieRepository, error) {
	return &MovieRepository{
		db:     db,
		logger: log,
	}, nil
}

const defaultLimit = 20

func (r *MovieRepository) GetTopMovies() ([]models.Movie, error) {
	// Fetch movies
	query := `
		SELECT id, tmdb_id, title, tagline, release_year, overview, score, 
		       popularity, language, poster_url, trailer_url
		FROM movies
		ORDER BY popularity DESC
		LIMIT $1
	`
	return r.getMovies(query)
}


func (r *MovieRepository) getMovies(query string) ([]models.Movie, error) {
	rows, err := r.db.Query(query, defaultLimit)
	if err != nil {
		r.logger.Error("Failed to query movies", err)
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan movie row", err)
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

var (
	ErrMovieNotFound = errors.New("movie not found")
)
```

Back in *main.go*, initialize the repository after the database creation

```go
	// Initialize repositories
	movieRepo, err := data.NewMovieRepository(db, logInstance)
	if err != nil {
		log.Fatalf("Failed to initialize movie repository: %v", err)
	}
```

Update handlers

```go
type MovieHandler struct {
	storage data.MovieStorage
	logger  *logger.Logger
}

func (h *MovieHandler) handleStorageError(w http.ResponseWriter, err error, context string) bool {
	if err != nil {
		if err == data.ErrMovieNotFound {
			http.Error(w, context, http.StatusNotFound)
			return true
		}
		h.logger.Error(context, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}
	return false
}
```

Update handler instance in *main.go* to use the new structure:

```go
movieHandler := handlers.NewMovieHandler(movieRepo, logInstance)	
```

### B6 - Finish the Movie Repository

The final *movie_repository.go* should look like

```go
package data

import (
	"database/sql"
	"errors"
	"strconv"

	"frontendmasters.com/movies/logger"
	"frontendmasters.com/movies/models"
	_ "github.com/lib/pq"
)

type MovieRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewMovieRepository(db *sql.DB, log *logger.Logger) (*MovieRepository, error) {
	return &MovieRepository{
		db:     db,
		logger: log,
	}, nil
}

const defaultLimit = 20

func (r *MovieRepository) GetTopMovies() ([]models.Movie, error) {
	// Fetch movies
	query := `
		SELECT id, tmdb_id, title, tagline, release_year, overview, score, 
		       popularity, language, poster_url, trailer_url
		FROM movies
		ORDER BY popularity DESC
		LIMIT $1
	`
	return r.getMovies(query)
}

func (r *MovieRepository) GetRandomMovies() ([]models.Movie, error) {
	// Fetch movies
	randomQuery := `
		SELECT id, tmdb_id, title, tagline, release_year, overview, score, 
		       popularity, language, poster_url, trailer_url
		FROM movies
		ORDER BY random()
		LIMIT $1
	`
	return r.getMovies(randomQuery)
}

func (r *MovieRepository) getMovies(query string) ([]models.Movie, error) {
	rows, err := r.db.Query(query, defaultLimit)
	if err != nil {
		r.logger.Error("Failed to query movies", err)
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan movie row", err)
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

func (r *MovieRepository) GetMovieByID(id int) (models.Movie, error) {
	// Fetch movie
	query := `
		SELECT id, tmdb_id, title, tagline, release_year, overview, score, 
		       popularity, language, poster_url, trailer_url
		FROM movies
		WHERE id = $1
	`
	row := r.db.QueryRow(query, id)

	var m models.Movie
	err := row.Scan(
		&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
		&m.Overview, &m.Score, &m.Popularity, &m.Language,
		&m.PosterURL, &m.TrailerURL,
	)
	if err == sql.ErrNoRows {
		r.logger.Error("Movie not found", ErrMovieNotFound)
		return models.Movie{}, ErrMovieNotFound
	}
	if err != nil {
		r.logger.Error("Failed to query movie by ID", err)
		return models.Movie{}, err
	}

	// Fetch related data
	if err := r.fetchMovieRelations(&m); err != nil {
		return models.Movie{}, err
	}

	return m, nil
}

func (r *MovieRepository) SearchMoviesByName(name string, order string, genre *int) ([]models.Movie, error) {
	orderBy := "popularity DESC"
	switch order {
	case "score":
		orderBy = "score DESC"
	case "name":
		orderBy = "title"
	case "date":
		orderBy = "release_year DESC"
	}

	genreFilter := ""
	if genre != nil {
		genreFilter = ` AND ((SELECT COUNT(*) FROM movie_genres 
								WHERE movie_id=movies.id 
								AND genre_id=` + strconv.Itoa(*genre) + `) = 1) `
	}

	// Fetch movies by name
	query := `
		SELECT id, tmdb_id, title, tagline, release_year, overview, score, 
		       popularity, language, poster_url, trailer_url
		FROM movies
		WHERE (title ILIKE $1 OR overview ILIKE $1) ` + genreFilter + `
		ORDER BY ` + orderBy + `
		LIMIT $2
	`
	rows, err := r.db.Query(query, "%"+name+"%", defaultLimit)
	if err != nil {
		r.logger.Error("Failed to search movies by name", err)
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan movie row", err)
			return nil, err
		}
		movies = append(movies, m)
	}

	return movies, nil
}

func (r *MovieRepository) GetAllGenres() ([]models.Genre, error) {
	query := `SELECT id, name FROM genres ORDER BY id`
	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Error("Failed to query all genres", err)
		return nil, err
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			r.logger.Error("Failed to scan genre row", err)
			return nil, err
		}
		genres = append(genres, g)
	}
	return genres, nil
}

// fetchMovieRelations fetches genres, actors, and keywords for a movie
func (r *MovieRepository) fetchMovieRelations(m *models.Movie) error {
	// Fetch genres
	genreQuery := `
		SELECT g.id, g.name 
		FROM genres g
		JOIN movie_genres mg ON g.id = mg.genre_id
		WHERE mg.movie_id = $1
	`
	genreRows, err := r.db.Query(genreQuery, m.ID)
	if err != nil {
		r.logger.Error("Failed to query genres for movie "+strconv.Itoa(m.ID), err)
		return err
	}
	defer genreRows.Close()
	for genreRows.Next() {
		var g models.Genre
		if err := genreRows.Scan(&g.ID, &g.Name); err != nil {
			r.logger.Error("Failed to scan genre row", err)
			return err
		}
		m.Genres = append(m.Genres, g)
	}

	// Fetch actors
	actorQuery := `
		SELECT a.id, a.first_name, a.last_name, a.image_url
		FROM actors a
		JOIN movie_cast mc ON a.id = mc.actor_id
		WHERE mc.movie_id = $1
	`
	actorRows, err := r.db.Query(actorQuery, m.ID)
	if err != nil {
		r.logger.Error("Failed to query actors for movie "+strconv.Itoa(m.ID), err)
		return err
	}
	defer actorRows.Close()
	for actorRows.Next() {
		var a models.Actor
		if err := actorRows.Scan(&a.ID, &a.FirstName, &a.LastName, &a.ImageURL); err != nil {
			r.logger.Error("Failed to scan actor row", err)
			return err
		}
		m.Casting = append(m.Casting, a)
	}

	// Fetch keywords
	keywordQuery := `
		SELECT k.word
		FROM keywords k
		JOIN movie_keywords mk ON k.id = mk.keyword_id
		WHERE mk.movie_id = $1
	`
	keywordRows, err := r.db.Query(keywordQuery, m.ID)
	if err != nil {
		r.logger.Error("Failed to query keywords for movie "+strconv.Itoa(m.ID), err)
		return err
	}
	defer keywordRows.Close()
	for keywordRows.Next() {
		var k string
		if err := keywordRows.Scan(&k); err != nil {
			r.logger.Error("Failed to scan keyword row", err)
			return err
		}
		m.Keywords = append(m.Keywords, k)
	}

	return nil
}

var (
	ErrMovieNotFound = errors.New("movie not found")
)

```

### B7 - Finish the handlers

The final *movies_handler.go* file should look like

```go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"frontendmasters.com/movies/data"
	"frontendmasters.com/movies/logger"
	"frontendmasters.com/movies/models"
)

type MovieHandler struct {
	storage data.MovieStorage
	logger  *logger.Logger
}

// Utility functions
func (h *MovieHandler) writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *MovieHandler) handleStorageError(w http.ResponseWriter, err error, context string) bool {
	if err != nil {
		if err == data.ErrMovieNotFound {
			http.Error(w, context, http.StatusNotFound)
			return true
		}
		h.logger.Error(context, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}
	return false
}

func (h *MovieHandler) parseID(w http.ResponseWriter, idStr string) (int, bool) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Invalid ID format", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func (h *MovieHandler) GetTopMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.storage.GetTopMovies()

	if h.handleStorageError(w, err, "Failed to get movies") {
		return
	}
	if h.writeJSONResponse(w, movies) == nil {
		h.logger.Info("Successfully served top movies")
	}
}

func (h *MovieHandler) GetRandomMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.storage.GetRandomMovies()
	if h.handleStorageError(w, err, "Failed to get movies") {
		return
	}
	if h.writeJSONResponse(w, movies) == nil {
		h.logger.Info("Successfully served random movies")
	}
}

func (h *MovieHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	order := r.URL.Query().Get("order")
	genreStr := r.URL.Query().Get("genre")

	var genre *int
	if genreStr != "" {
		genreInt, ok := h.parseID(w, genreStr)
		if !ok {
			return
		}
		genre = &genreInt
	}

	var movies []models.Movie
	var err error
	if query != "" {
		movies, err = h.storage.SearchMoviesByName(query, order, genre)
	}
	if h.handleStorageError(w, err, "Failed to get movies") {
		return
	}
	if h.writeJSONResponse(w, movies) == nil {
		h.logger.Info("Successfully served movies")
	}
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/movies/"):]
	id, ok := h.parseID(w, idStr)
	if !ok {
		return
	}

	movie, err := h.storage.GetMovieByID(id)
	if h.handleStorageError(w, err, "Failed to get movie by ID") {
		return
	}
	if h.writeJSONResponse(w, movie) == nil {
		h.logger.Info("Successfully served movie with ID: " + idStr)
	}
}

func (h *MovieHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.storage.GetAllGenres()
	if h.handleStorageError(w, err, "Failed to get genres") {
		return
	}
	if h.writeJSONResponse(w, genres) == nil {
		h.logger.Info("Successfully served genres")
	}
}

func NewMovieHandler(storage data.MovieStorage, log *logger.Logger) *MovieHandler {
	return &MovieHandler{
		storage: storage,
		logger:  log,
	}
}
```

### B7 - Update the handlers

In *main.go* all the handlers for the API should look like:

```go
	// Initialize handlers
	movieHandler := handlers.NewMovieHandler(movieRepo, logInstance)
	// authHandler := handlers.NewAuthHandler(userStorage, jwt, logInstance)

	// Set up routes
	http.HandleFunc("/api/movies/random", movieHandler.GetRandomMovies)
	http.HandleFunc("/api/movies/top", movieHandler.GetTopMovies)
	http.HandleFunc("/api/movies/search", movieHandler.SearchMovies)
	http.HandleFunc("/api/movies/", movieHandler.GetMovie)
	http.HandleFunc("/api/genres", movieHandler.GetGenres)
	http.HandleFunc("/api/account/register", movieHandler.GetGenres)
	http.HandleFunc("/api/account/authenticate", movieHandler.GetGenres)
```

## C-Frontend

### C1 - Create the HTML

Create `public.html`

```html
<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>ReelingIt - Movies</title>
    <link href="https://fonts.googleapis.com/css2?family=Open+Sans:wght@300;400;700&display=swap" rel="stylesheet">
    <link rel="stylesheet"
        href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@48,400,0,0" />
    <link rel="stylesheet" href="/styles.css">
    <meta name="theme-color" content="#56bce8">
    <link rel="manifest" href="app.webmanifest">
    <link rel="icon" href="/images/icon.png" type="image/png">
    <script src="/app.js" type="module" defer></script>
    <base href="/">
</head>

<body>
    <header>
        <h1>
            <a href="/" class="navlink"><img src="/images/logo.png" height="35" alt="ReelingIt"></a>
        </h1>
        <nav>
            <ul>
                <li><a href="/" class="navlink">Movies</a></li>
                <li><a href="/account/favorites" class="navlink">Favorites</a></li>
                <li><a href="/account/watchlist" class="navlink">Watchlist</a></li>
                <li><a href="/account/" class="navlink">My Account</a></li>
            </ul>
        </nav>
        <div>
            <form onsubmit="app.search(event)">
                <input type="search" placeholder="Search movies">
            </form>
        </div>
    </header>

    <main>
    </main>

    <footer>
        <p>© ReelingIt - FrontendMasters.com</p>
    </footer>
         
</body>

</html>
```

### C2 - Initialize the client-side app

Create `app.js`

```js
window.app = { 
    search: (event) => {
        event.preventDefault();
        const keywords = document.querySelector("input[type=search]").value;
        
    },    
}

window.addEventListener("DOMContentLoaded", () => {

})
```

### C3 - Add a Manifest file

Add the *app.webmanifest* file to the project

```json
{
    "name": "ReelingIt",
    "short_name": "ReelingIt",
    "theme_color": "#43281C",
    "display": "browser",
    "background_color": "#56bce8",
    "description": "The ultimate app for movie lovers: discover trailers, reviews, showtimes, and more. Experience cinema like never before!",    "icons": [
        {
            "src": "images/icon.png",
            "sizes": "1024x1024",
            "type": "image/png"
        }
    ]
}
```

### C4 - Create the API Service

Create *services/API.js* file:

```js
export const API = {
    baseURL: '/api/',
    getTopMovies: async () => {
        return await API.fetch("movies/top");
    },
    getRandomMovies: async () => {
        return await API.fetch("movies/random");
    },
    getMovieById: async (id) => {
        return await API.fetch(`/movies/${id}`);
    },
    searchMovies: async (q, order, genre) => {
        return await API.fetch(`/movies/search`, {q, order, genre})
    },
    getGenres: async () => {
        return await API.fetch("genres");
    },
    fetch: async (service, args) => {
        try {
            const queryString = args ? new URLSearchParams(args).toString() : "";
            const response = await fetch(API.baseURL + service + '?' + queryString);
            const result = await response.json();
            return result;
        } catch (e) {
            console.error(e);
            app.showError();
        }
    }
}

export default API;
```

## D-Web Components

### D1 - HomePage template

Create a template in the *index.html*

```html
<template id="template-home">
    <section class="vertical-scroll" id="top-10">
        <h2>This Week's Top 10</h2>
        <ul>
            <animated-loading data-elements="5"
                data-width="150px" data-height="220px">
            </animated-loading> 
        </ul>
    </section>
    <section class="vertical-scroll" id="random">
        <h2>Something to watch today</h2>
        <ul>
            <animated-loading data-elements="5"
                data-width="150px" data-height="220px">
            </animated-loading> 
        </ul>
    </section>
</template>
```

### D2 - MovieItem Component

Create the *components* folder and *MovieItem.js* file

```js
export class MovieItemComponent extends HTMLElement {
    constructor(movie) {
        super();
        this.movie = movie;
    }

    connectedCallback() {
        this.innerHTML = `
                <a href="#">
                    <article>
                        <img src="${this.movie.poster_url}" alt="${this.movie.title} Poster">
                        <p>${this.movie.title} (${this.movie.release_year})</p>
                    </article>
                </a>
            `
    }
}

customElements.define("movie-item", MovieItemComponent);
```

### D3 - HomePage Component

Create components/*HomePage.js*

```js
import API from "../services/API.js";
import { MovieItemComponent } from "./MovieItem.js";

export default class HomePage extends HTMLElement {

    async render() {
        const topMovies = await API.getTopMovies();
        renderMoviesInList(topMovies, this.querySelector("#top-10 ul"));

        const randomMovies = await API.getRandomMovies();
        renderMoviesInList(randomMovies, this.querySelector("#random ul"));

        function renderMoviesInList(movies, ul) {
            ul.innerHTML = "";
            movies.forEach(movie => {
                const li = document.createElement("li");
                li.appendChild(new MovieItemComponent(movie));
                ul.appendChild(li);
            });    
        }
    }

    connectedCallback() {
        const template = document.getElementById("template-home");
        const content = template.content.cloneNode(true);
        this.appendChild(content);  

        this.render();
    }
}
customElements.define("home-page", HomePage);
```

### D4 - Animated Loading

Create the *components/AnimatedLoading.js* file:

```js
class AnimatedLoading extends HTMLElement {
    constructor() {
        super();
    }

    connectedCallback() {
        let qty = this.dataset.elements ?? 1;
        let width = this.dataset.width ?? "100px";
        let height = this.dataset.height ?? "10px";
        for (let i=0; i<qty; i++) {
            const wrapper = document.createElement('div');
            wrapper.setAttribute('class', 'loading-wave');    
            wrapper.style.width = width;
            wrapper.style.height = height;        
            wrapper.style.margin = "10px";
            wrapper.style.display = "inline-block";
            this.appendChild(wrapper);
        }
    }
}

customElements.define('animated-loading', AnimatedLoading);
```

### D5 - Movie Details

Add a new template to *index.html*

```html
    <template id="template-movie-details">
        <article id="movie">
            <h2><animated-loading elements="2"></animated-loading></h2>
            <h3></h3>
            <header>
                <img src="" alt="Poster">
                <youtube-embed id="trailer" data-url=""></youtube-embed>
                <section id="actions">
                    <dl id="metadata">
                    </dl>
                    <button>Add to Favorites</button>
                    <button>Add to Watchlist</button>    
                </section>
            </header>
            <ul id="genres"></ul>
            <p id="overview"></p>
            <ul id="cast"></ul>
        </article>
    </template>
```

Create the *components/MovieDetailsPage.js* file:

```js
import API from "../services/API.js";

export default class MovieDetailsPage extends HTMLElement {
    
    movie = null;

    async render(id) {
        try {
            this.movie = await API.getMovieById(id);
        } catch (e) {
            app.showError();
            return;
        }
        const template = document.getElementById("template-movie-details");
        const content = template.content.cloneNode(true);
        this.appendChild(content);  

        this.querySelector("h2").textContent = this.movie.title;
        this.querySelector("h3").textContent = this.movie.tagline;
        this.querySelector("img").src = this.movie.poster_url;
        this.querySelector("#trailer").dataset.url = this.movie.trailer_url;
        this.querySelector("#overview").textContent = this.movie.overview;
        this.querySelector("#metadata").innerHTML = `                        
            <dt>Release Date</dt>
            <dd>${this.movie.release_year}</dd>                        
            <dt>Score</dt>
            <dd>${this.movie.score} / 10</dd>                        
            <dt>Original languae</dt>
            <dd>${this.movie.language}</dd>                        
        `;

        const ulGenres = this.querySelector("#genres");
        ulGenres.innerHTML = "";
        this.movie.genres.forEach(genre => {
            const li = document.createElement("li");
            li.textContent = genre.name;
            ulGenres.appendChild(li);
        });

        const ulCast = this.querySelector("#cast");
        ulCast.innerHTML = "";
        this.movie.casting.forEach(actor => {
            const li = document.createElement("li");
            li.innerHTML = `
                <img src="${actor.image_url ?? '/images/generic_actor.jpg'}" alt="Picture of ${actor.last_name}">
                <p>${actor.first_name} ${actor.last_name}</p>
            `;
            ulCast.appendChild(li);
        });
    }

    connectedCallback() {
        const id = this.params[0];

        this.render(id);

    }
}
customElements.define("movie-details-page", MovieDetailsPage);
```

### D6 - YouTube Embed Component

Create the *components/YouTubeEmbed.js* file:

```js
export class YouTubeEmbed extends HTMLElement {
    
    static get observedAttributes() {
        return ['data-url'];
    }

    attributeChangedCallback(prop, value) {
        if (prop === 'data-url') {
            const url = this.dataset.url;
            const videoId = url.substring(url.indexOf("?v")+3);
            console.log(videoId);

            this.innerHTML = `
                <iframe width="100%" height="300" src="https://www.youtube.com/embed/${videoId}" title="YouTube video player" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>
            `;
        }
    }

}

customElements.define("youtube-embed", YouTubeEmbed);
```

## E-Client Side Routing

### E1 - Routes

Let's create our routes for the client in */services/Routes.js* and boilerplate code for all the Web Components.

```js
export const routes = [
    {
        path: "/",
        component: HomePage
    },
    {
        path: "/movies",
        component: MoviesPage
    },
    {
        path: /\/movies\/(\d+)/,
        component: MovieDetailsPage
    },
    {
        path: "/account/register",
        component: RegisterPage
    },
    {
        path: "/account/login",
        component: LoginPage
    },     
    {
        path: "/account/",
        component: AccountPage
    },
    {
        path: "/account/favorites",
        component: FavoritesPage
    },	
{
        path: "/account/watchlist",
        component: WatchlistPage
    },	
]
```

### E2 - The Router

Now let's make the router in *services/Router.js*

```js
import { routes } from "./Routes.js";

const Router = {
    init: () => {
        document.querySelectorAll("a.navlink").forEach(a => {
            a.addEventListener("click", event => {
                event.preventDefault();
                const href = a.getAttribute("href");
                Router.go(href);
            });
        });  
        window.addEventListener("popstate", () => {
            Router.go(location.pathname, false);
        });      
        // Process initial URL   
        Router.go(location.pathname + location.search);
    },
    go: (route, addToHistory=true) => {
        if (addToHistory) {
            history.pushState(null, "", route);
        }
        const routePath = route.includes('?') ? route.split('?')[0] : route;
        let pageElement = null;
        for (const r of routes) {
            if (typeof r.path === "string" && r.path === routePath) {
                pageElement = new r.component();
                break;
            } else if (r.path instanceof RegExp) {
                const match = r.path.exec(route);
                if (match) {
                    const params = match.slice(1);
                    pageElement = new r.component();
                    pageElement.params = params;                    
                    break;
                }
            }
        }
        if (pageElement==null) {
            pageElement = document.createElement("h1");
            pageElement.textContent = "Page not found";
        }       

		document.querySelector("main").innerHTML = "";
		document.querySelector("main").appendChild(pageElement); 
    }

}

export default Router;
```

Now let's add the Router to app in `app.js` and call init to enhance our links

```js
window.app = { 
    API,
    Router,
}
window.addEventListener("DOMContentLoaded", () => {
    app.Router.init();
})
```

### E3 - Server-side dynamic routes

When we refresh the page on dynamic routes, we get a 404, to solve the problem, let's add some new Handlers in our backend, adding them before the file serving at *main.go*:

```go
	// Handler catch-all
	catchAllHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	}
	http.HandleFunc("/movies", catchAllHandler)
	http.HandleFunc("/movies/", catchAllHandler)
	http.HandleFunc("/account/", catchAllHandler)
```
### E4 - Animating the transition

At *Router.js* we can remove the last two code lines with this:

```js
	function updatePage() {
		document.querySelector("main").innerHTML = "";
		document.querySelector("main").appendChild(pageElement); 
	}

	if (!document.startViewTransition) {
		updatePage();
	} else {
		const oldPage = document.querySelector("main").firstElementChild;
		if (oldPage) oldPage.style.viewTransitionName = "old";
		pageElement.style.viewTransitionName = "new";
		document.startViewTransition( () => updatePage() );
	}
```

### E5 - Error Messages

Add in *index.html*

```html
 <!-- Alert Modal -->
<dialog id="alert-modal">
	<h3>Error</h3>
	<p>There was an error loading the page</p>
	<button class="action-btn" onclick="app.closeError()">OK</button>
</dialog>
```

Then, in *app.js*

```js
window.app = { 
	// ...
    showError: (message = 'There was an error loading the page', goToHome=true) => {
        document.querySelector("#alert-modal").showModal()
		document.querySelector("#alert-modal p").textContents = message;
        if (goToHome) app.Router.go("/");
        return;
    },
    closeError: () => {
        document.getElementById('alert-modal').close()        
    },
	// ...
}
```

## F-The Search Section

### F1 - The MoviesPage

Add the template in *index.html*

```html
<template id="template-movies">
        <section>
            <div id="search-header">
                <h2></h2>
                <section id="filters">
                    <select id="filter" onchange="app.searchFilterChange(this.value)">
                        <option>Filter by Genre</option>                        
                    </select>
                    <select id="order" onchange="app.searchOrderChange(this.value)">
                        <option value="popularity">Sort by Popularity</option>
                        <option value="score">Sort by Score</option>
                        <option value="date">Sort by Release Date</option>
                        <option value="name">Sort by Name</option>
                    </select>
                </section>
            </div>
            <ul id="movies-result">
                <animated-loading data-elements="5"
                    data-width="150px" data-height="220px">
                </animated-loading> 
            </ul>
        </section>
    </template> 
```

And then for the Web Component *MoviesPage*

```js
import API from "../services/API.js";
import { MovieItemComponent } from "./MovieItem.js";

export default class MoviesPage extends HTMLElement {
    
    async render(query) {
        const urlParams = new URLSearchParams(window.location.search);
        const order = urlParams.get("order") ?? "";
        const genre = urlParams.get("genre") ?? "";

        const movies = await API.searchMovies(query, order, genre);
        
        const ulMovies = this.querySelector("ul");
        ulMovies.innerHTML = "";
        if (movies && movies.length>0) {
            movies.forEach(movie => {
                const li = document.createElement("li");
                li.appendChild(new MovieItemComponent(movie));
                ulMovies.appendChild(li);
            });    
        } else {
            ulMovies.innerHTML = "<h3>There are no movies with your search</h3>";
        }        

        //await this.loadGenres();

        if (order) this.querySelector("#order").value = order;
        if (genre) this.querySelector("#filter").value = genre;

    }
    
   
    connectedCallback() {
        const template = document.getElementById("template-movies");
        const content = template.content.cloneNode(true);
        this.appendChild(content);  

        const urlParams = new URLSearchParams(window.location.search);
        const query = urlParams.get('q');
        if (query) {
            this.querySelector("h2").textContent = `'${query}' movies`;
            this.render(query);
        } else {
            app.showError();
        }
    }
}
customElements.define("movies-page", MoviesPage);
```

At *app.js* add the following function to app

```js
// ...
    search: (event) => {
        event.preventDefault();
        const keywords = document.querySelector("input[type=search]").value;
        if (keywords.length>1) {
            app.Router.go(`/movies?q=${keywords}`)
        }
    },
// ...	
```

### F2 - Connecting the Filter with Genre

Add in *MoviesComponent.js*

```js
 async loadGenres() {
	const genres = await API.getGenres();
	const select = this.querySelector("#filter");
	select.innerHTML = `
		<option value=''>Filter by Genre</option>
	`;
	genres.forEach(genre => {
		var option = document.createElement("option");
		option.value = genre.id;
		option.textContent = genre.name;
		select.appendChild(option);
	})
}
```

Call it from *connectedCallback*

### F3 - Adding Support for Order and Filter

In *app.js* add the new functions to the app object:

```js
// ...
    searchOrderChange: (order) => {
        const urlParams = new URLSearchParams(window.location.search);
        const q = urlParams.get("q");
        const genre = urlParams.get("genre") ?? "";
        app.Router.go(`/movies?q=${q}&order=${order}&genre=${genre}`);
    },
    searchFilterChange: (genre) => {
        const urlParams = new URLSearchParams(window.location.search);
        const q = urlParams.get("q");
        const order = urlParams.get("order") ?? "";
        app.Router.go(`/movies?q=${q}&order=${order}&genre=${genre}`);
// ...
    }
```


## G-Authentication

### G1 - Adding the new Storage

Let's add a new dependency, executing:

```
go get "golang.org/x/crypto/bcrypt"
```


Let's add a new interface in *data/interfaces.go*

```go
type AccountStorage interface {
	Authenticate(string, string) (bool, error)
	Register(string, string, string) (bool, error)
	GetAccountDetails(string) (models.User, error)
	SaveCollection(models.User, int, string) (bool, error)
}
```
And we create a new implementation at *data/account_repository.go*

```go
package data

import (
	"database/sql"
	"errors"
	"time"

	"frontendmasters.com/movies/logger"
	"frontendmasters.com/movies/models"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type AccountRepository struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewAccountRepository(db *sql.DB, log *logger.Logger) (*AccountRepository, error) {
	return &AccountRepository{
		db:     db,
		logger: log,
	}, nil
}

func (r *AccountRepository) Register(name, email, password string) (bool, error) {
	// Validate basic requirements
	if name == "" || email == "" || password == "" {
		r.logger.Error("Registration validation failed: missing required fields", nil)
		return false, ErrRegistrationValidation
	}

	// Check if user already exists
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`, email).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check existing user", err)
		return false, err
	}
	if exists {
		r.logger.Error("User already exists with email: "+email, ErrUserAlreadyExists)
		return false, ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		r.logger.Error("Failed to hash password", err)
		return false, err
	}

	// Insert new user
	query := `
		INSERT INTO users (name, email, password_hashed, time_created)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	var userID int
	err = r.db.QueryRow(
		query,
		name,
		email,
		string(hashedPassword),
		time.Now(),
	).Scan(&userID)
	if err != nil {
		r.logger.Error("Failed to register user", err)
		return false, err
	}

	return true, nil
}

func (r *AccountRepository) Authenticate(email string, password string) (bool, error) {
	if email == "" || password == "" {
		r.logger.Error("Authentication validation failed: missing credentials", nil)
		return false, ErrAuthenticationValidation
	}

	// Fetch user by email
	var user models.User
	query := `
		SELECT id, name, email, password_hashed
		FROM users 
		WHERE email = $1 AND time_deleted IS NULL
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHashed,
	)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found for email: "+email, nil)
		return false, ErrAuthenticationValidation
	}
	if err != nil {
		r.logger.Error("Failed to query user for authentication", err)
		return false, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHashed), []byte(password))
	if err != nil {
		r.logger.Error("Password mismatch for email: "+email, nil)
		return false, ErrAuthenticationValidation
	}

	// Update last login time
	updateQuery := `
		UPDATE users 
		SET last_login = $1
		WHERE id = $2
	`
	_, err = r.db.Exec(updateQuery, time.Now(), user.ID)
	if err != nil {
		r.logger.Error("Failed to update last login", err)
		// Don't fail authentication just because last login update failed
	}

	return true, nil
}

func (r *AccountRepository) GetAccountDetails(email string) (models.User, error) {
	var user models.User
	query := `
		SELECT id, name, email
		FROM users 
		WHERE email = $1 AND time_deleted IS NULL
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found for email: "+email, nil)
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("Failed to query user by email", err)
		return models.User{}, err
	}

	// Fetch favorites
	favoritesQuery := `
		SELECT m.id, m.tmdb_id, m.title, m.tagline, m.release_year, 
		       m.overview, m.score, m.popularity, m.language, 
		       m.poster_url, m.trailer_url
		FROM movies m
		JOIN user_movies um ON m.id = um.movie_id
		WHERE um.user_id = $1 AND um.relation_type = 'favorite'
	`
	favoriteRows, err := r.db.Query(favoritesQuery, user.ID)
	if err != nil {
		r.logger.Error("Failed to query user favorites", err)
		return user, err
	}
	defer favoriteRows.Close()

	for favoriteRows.Next() {
		var m models.Movie
		if err := favoriteRows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan favorite movie row", err)
			return user, err
		}
		user.Favorites = append(user.Favorites, m)
	}

	// Fetch watchlist
	watchlistQuery := `
		SELECT m.id, m.tmdb_id, m.title, m.tagline, m.release_year, 
		       m.overview, m.score, m.popularity, m.language, 
		       m.poster_url, m.trailer_url
		FROM movies m
		JOIN user_movies um ON m.id = um.movie_id
		WHERE um.user_id = $1 AND um.relation_type = 'watchlist'
	`
	watchlistRows, err := r.db.Query(watchlistQuery, user.ID)
	if err != nil {
		r.logger.Error("Failed to query user watchlist", err)
		return user, err
	}
	defer watchlistRows.Close()

	for watchlistRows.Next() {
		var m models.Movie
		if err := watchlistRows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			r.logger.Error("Failed to scan watchlist movie row", err)
			return user, err
		}
		user.Watchlist = append(user.Watchlist, m)
	}

	return user, nil
}

func (r *AccountRepository) SaveCollection(user models.User, movieID int, collection string) (bool, error) {
	// Validate inputs
	if movieID <= 0 {
		r.logger.Error("SaveCollection failed: invalid movie ID", nil)
		return false, errors.New("invalid movie ID")
	}
	if collection != "favorite" && collection != "watchlist" {
		r.logger.Error("SaveCollection failed: invalid collection type", nil)
		return false, errors.New("collection must be 'favorite' or 'watchlist'")
	}

	// Get user ID from email
	var userID int
	err := r.db.QueryRow(`
		SELECT id 
		FROM users 
		WHERE email = $1 AND time_deleted IS NULL
	`, user.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		r.logger.Error("User not found", nil)
		return false, ErrUserNotFound
	}
	if err != nil {
		r.logger.Error("Failed to query user ID", err)
		return false, err
	}

	// Check if the relationship already exists
	var exists bool
	err = r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM user_movies 
			WHERE user_id = $1 
			AND movie_id = $2 
			AND relation_type = $3
		)
	`, userID, movieID, collection).Scan(&exists)
	if err != nil {
		r.logger.Error("Failed to check existing collection entry", err)
		return false, err
	}
	if exists {
		r.logger.Info("Movie already in " + collection + " for user")
		return true, nil // Return true since the movie is already in the collection
	}

	// Insert the new relationship
	query := `
		INSERT INTO user_movies (user_id, movie_id, relation_type, time_added)
		VALUES ($1, $2, $3, $4)
	`
	_, err = r.db.Exec(query, userID, movieID, collection, time.Now())
	if err != nil {
		r.logger.Error("Failed to save movie to "+collection, err)
		return false, err
	}

	r.logger.Info("Successfully added movie " + string(movieID) + " to " + collection + " for user")
	return true, nil
}

var (
	ErrRegistrationValidation   = errors.New("registration failed")
	ErrAuthenticationValidation = errors.New("authentication failed")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrUserNotFound             = errors.New("user not found")
)

```

### G2 - Add the handlers

Now we create `handlers/account_handlers.go`

```go
package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"frontendmasters.com/movies/data"
	"frontendmasters.com/movies/logger"
	"frontendmasters.com/movies/models"
	"frontendmasters.com/movies/token"
)

// Define request structure
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Define request structure
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type AccountHandler struct {
	storage data.AccountStorage
	logger  *logger.Logger
}

// Utility functions
func (h *AccountHandler) writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *AccountHandler) handleStorageError(w http.ResponseWriter, err error, context string) bool {
	if err != nil {
		switch err {
		case data.ErrAuthenticationValidation, data.ErrUserAlreadyExists, data.ErrRegistrationValidation:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: err.Error()})
			return true
		case data.ErrUserNotFound:
			http.Error(w, "User not found", http.StatusNotFound)
			return true
		default:
			h.logger.Error(context, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return true
		}
	}
	return false
}

func (h *AccountHandler) Register(w http.ResponseWriter, r *http.Request) {

	// Parse request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode registration request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Register the user
	success, err := h.storage.Register(req.Name, req.Email, req.Password)
	if h.handleStorageError(w, err, "Failed to register user") {
		return
	}

	// Return success response
	response := AuthResponse{
		Success: success,
		Message: "User registered successfully",
	}

	if err := h.writeJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully registered user with email: " + req.Email)
	}
}

func (h *AccountHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode authentication request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Authenticate the user
	success, err := h.storage.Authenticate(req.Email, req.Password)
	if h.handleStorageError(w, err, "Failed to authenticate user") {
		return
	}

	// Return success response
	response := AuthResponse{
		Success: success,
		Message: "User registered successfully",
	}

	if err := h.writeJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully authenticated user with email: " + req.Email)
	}
}


func NewAccountHandler(storage data.AccountStorage, log *logger.Logger) *AccountHandler {
	return &AccountHandler{
		storage: storage,
		logger:  log,
	}
}


```

Finally, we register the handlers in *main.go*

```go
	accountRepo, err := data.NewAccountRepository(db, logInstance)
	if err != nil {
		log.Fatalf("Failed to initialize account repository: %v", err)
	}

	// ...

	accountHandler := handlers.NewAccountHandler(accountRepo, logInstance)
	http.HandleFunc("/api/account/register/", accountHandler.Register)
	http.HandleFunc("/api/account/authenticate/", accountHandler.Authenticate)

```


### G3 - Adding the Forms

Add the new templates and populate the Web Component objects for LoginPage and RegisterPage

```html
<template id="template-register">
        <section>
            <h2>Register a New Account</h2>
            <form onsubmit="app.register(event)">
                <label for="register-name">Name</label>
                <input type="text" id="register-name" placeholder="Name" required autocomplete="name">
                <label for="register-email">Email</label>
                <input type="email" id="register-email" placeholder="Email" required autocomplete="email">
                <label for="register-password">Password</label>
                <input type="password" id="register-password" placeholder="Password" required autocomplete="new-password">
                <label for="register-password-confirm">Confirm Password</label>
                <input type="password" id="register-password-confirm" placeholder="Confirm Password" required autocomplete="new-password">                 
                <button>Register</button>
                <p>If you already have an account, please <a href="/account/login">login</a>.</p>
            </form>
        </section>
    </template>
    <template id="template-login">
        <section>
            <h2>Login into Your Account</h2>
            <form onsubmit="app.login(event)">
                <label for="login-email">Email</label>
                <input type="email" id="login-email" placeholder="Email" required autocomplete="email">
                <label for="login-password">Password</label>
                <input type="password" id="login-password" placeholder="Password" required autocomplete="current-password">
                <button>Log In</button>
                <p>If you don't have an account, please <a href="/account/register">register</a>.</p>
            </form>
        </section>
    </template> 
```

### G4 - Connecting the API

Add to *API.js*

```js
    register: async (name, email, password) => {
        return await API.send("account/register/", {name, email, password})
    },
    authenticate: async (email, password) => {
        return await API.send("account/authenticate/", {email, password})
    },   
    send: async (service, args) => {
        try {
            const response = await fetch(API.baseURL + service, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(args)
            });
            const result = await response.json();
            return result;
        } catch (e) {
            console.error(e);
            app.showError();
        }
    },   

```

### G5 - Calling the APIs

In *app.js* we will now add:

```js
    register: async (event) => {
        event.preventDefault();
        let errors = [];
        const name = document.getElementById("register-name").value;
        const email = document.getElementById("register-email").value;
        const password = document.getElementById("register-password").value;
        const passwordConfirm = document.getElementById("register-password-confirm").value;

        if (name.length < 4) errors.push("Enter your complete name");
        if (email.length < 8) errors.push("Enter your complete email");
        if (password.length < 6) errors.push("Enter a password with 6 characters");
        if (password != passwordConfirm) errors.push("Passwords don't match");
        if (errors.length==0) {
            const response = await API.register(name, email, password);
            if (response.success) {
                app.Router.go("/account/")
            } else {
                app.showError(response.message, false);
            }        
        } else {
            app.showError(errors.join(". "), false);
        }
    },
    login: async (event) => {
        event.preventDefault();
        let errors = [];
        const email = document.getElementById("login-email").value;
        const password = document.getElementById("login-password").value;
 
        if (email.length < 8) errors.push("Enter your complete email");
        if (password.length < 6) errors.push("Enter a password with 6 characters");
        if (errors.length==0) {
            const response = await API.authenticate(email, password);
            if (response.success) {
                app.Router.go("/account/")
            } else {
                app.showError(response.message, false);
            }
        } else {
            app.showError(errors.join(". "), false);
        }
    },
```

### G6 - Sending the Token

To implement JWT, we need to install a package

```
go get -u github.com/golang-jwt/jwt/v5
```

Then, we create the *token* package and inside two file. Starting with *getsecret.go*

```go
package token

import (
	"os"
	"frontendmasters.com/movies/logger"
)

func GetJWTSecret(logger logger.Logger) string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-for-dev"
		logger.Info("JWT_SECRET not set, using default development secret")
	} else {
		logger.Info("Using JWT_SECRET from environment")
	}
	return jwtSecret
}
```

Then, *creation.go*

```go
package token

import (
	"time"

	"frontendmasters.com/movies/logger"
	"frontendmasters.com/movies/models"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(user models.User, logger logger.Logger) string {
	jwtSecret := GetJWTSecret(logger)

	// Create a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	})

	// Sign the token with the secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		logger.Error("Failed to sign JWT", err)
		return ""
	}

	return tokenString
}
```

And finally, *validation.go*

```go
package token

import (
	"frontendmasters.com/movies/logger"
	"github.com/golang-jwt/jwt/v5"
)

func ValidateJWT(tokenString string, logger logger.Logger) (*jwt.Token, error) {
	jwtSecret := GetJWTSecret(logger)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Error("Unexpected signing method", nil)
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		logger.Error("Failed to validate JWT", err)
		return nil, err
	}

	if !token.Valid {
		logger.Error("Invalid JWT token", nil)
		return nil, jwt.ErrTokenInvalidId
	}

	return token, nil
}
```

Finally, we add it to our AuthResponse struct in *account_handlers.go* and return it after a successful registration or authentication.

```go
type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	JWT     string `json:"jwt"`
}

// ...

// in register
	response := AuthResponse{
		Success: success,
		Message: "User registered successfully",
		JWT:     token.CreateJWT(models.User{Email: req.Email, Name: req.Name}, *h.logger),
	}

// ...

// in authenticate
	response := AuthResponse{
		Success: success,
		Message: "User registered successfully",
		JWT:     token.CreateJWT(models.User{Email: req.Email}, *h.logger),
	}
```

### G7 - Working with the Token Client-Side

We first create a *services/Store.js*

```js

const Store = {
    jwt: null,
    get loggedIn() {
        return this.jwt !== null;
    }
}

if (localStorage.getItem("jwt")) {
    Store.jwt = localStorage.getItem("jwt");
}

const proxiedStore = new Proxy(Store, {
    set: (target, prop, value) => {
        switch (prop) {
            case "jwt":
                target[prop] = value;
                localStorage.setItem("jwt", value)
                break;
        }
        return true;
    }
});


export default proxiedStore;
```

Then, in *app.js* we save the jwt after a successful login or registration

```js
// ...
if (response.success) {
	app.Store.jwt = response.jwt;
	app.Router.go("/account/")
} else {
	app.showError(response.message, false);
}
// ...			
```

Finally, we update our router, so it can detect and work with URLs needing authentication

```js
// ...
        for (const r of routes) {
            if (typeof r.path === "string" && r.path === routePath) {
                pageElement = new r.component();
                pageElement.loggedIn = r.loggedIn;
            } else if (r.path instanceof RegExp) {
                const match = r.path.exec(route);
                if (match) {
                    const params = match.slice(1);
                    pageElement = new r.component();
                    pageElement.loggedIn = r.loggedIn;

                    pageElement.params = params;                    
                }
            }
            if (pageElement) {
                // A page was found, we checked if we have access to it.
                if (pageElement.loggedIn && app.Store.loggedIn==false) {
                    app.Router.go("/account/login");
                    return;
                }
                break;
            }
        }
// ...
```

### G8 - Creating the My Account page

We start by adding a template to our *index.html*

```html
 <template id="template-account">
	<section id="account">
		<h2>You are Logged In</h2>
		<button onclick="app.logout()">Log out</button>
		<button onclick="app.Router.go('/account/favorites')">Your Favorites</button>
		<button onclick="app.Router.go('/account/watchlist')">Your Watchlist</button>

	</section>
</template>        
```

We create the component *AccountPage.js*

```js
import API from "../services/API.js";
import { MovieItemComponent } from "./MovieItem.js";

export default class AccountPage extends HTMLElement {


    connectedCallback() {
        const template = document.getElementById("template-account");
        const content = template.content.cloneNode(true);
        this.appendChild(content);  
    }
}
customElements.define("account-page", AccountPage);
```

We finally define the route in *Router.js* and implement app.logout in *main.js*

```js
// ...
    {
        path: "/account/",
        component: AccountPage,
        loggedIn: true
    },
// ...
```



## H-Favorites and Watchlist

### H1 - Creating a middleware for authentication check

At *AccountHandlers.go* add the following code:

```go
func (h *AccountHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// Parse and validate the token

		token, err := jwt.Parse(tokenStr,
			func(t *jwt.Token) (interface{}, error) {
				// Ensure the signing method is HMAC
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(token.GetJWTSecret(*h.logger)), nil
			},
		)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Get the email from claims
		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "Email not found in token", http.StatusUnauthorized)
			return
		}

		// Inject email into the request context
		ctx := context.WithValue(r.Context(), "email", email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

```

### H2 - Adding Web Services

In  *AccountHandlers.go* add the new handlers

```go


func (h *AccountHandler) SaveToCollection(w http.ResponseWriter, r *http.Request) {
	type CollectionRequest struct {
		MovieID    int    `json:"movie_id"`
		Collection string `json:"collection"`
	}

	var req CollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode collection request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}

	success, err := h.storage.SaveCollection(models.User{Email: email},
		req.MovieID, req.Collection)
	if h.handleStorageError(w, err, "Failed to save to collection") {
		return
	}

	response := AuthResponse{
		Success: success,
		Message: "Movie added to " + req.Collection + " successfully",
	}

	if err := h.writeJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully saved movie to " + req.Collection)
	}
}

func (h *AccountHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}
	details, err := h.storage.GetAccountDetails(email)
	if err != nil {
		http.Error(w, "Unable to retrieve collections", http.StatusInternalServerError)
		return
	}
	if err := h.writeJSONResponse(w, details.Favorites); err == nil {
		h.logger.Info("Successfully sent favorites")
	}
}

func (h *AccountHandler) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}
	details, err := h.storage.GetAccountDetails(email)
	if err != nil {
		http.Error(w, "Unable to retrieve collections", http.StatusInternalServerError)
		return
	}
	if err := h.writeJSONResponse(w, details.Watchlist); err == nil {
		h.logger.Info("Successfully sent favorites")
	}
}
```

Then register them in *main.go*

```go
	http.Handle("/api/account/favorites/",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.GetFavorites)))

	http.Handle("/api/account/watchlist/",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.GetWatchlist)))

	http.Handle("/api/account/save-to-collection/",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.SaveToCollection)))
```

### H3 - Connecting with the Client

We need to change *API.js* to send the token if it's available

```js
    send: async (service, args) => {
        const response = await fetch(API.baseURL + service, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": app.Store.jwt ? `Bearer ${app.Store.jwt}` : null
            },
            body: JSON.stringify(args)
        });
        const result = await response.json();
        return result;
    },
    fetch: async (service, args) => {
        const queryString = args ? new URLSearchParams(args).toString() : "";
        const response = await fetch(API.baseURL + service + '?' + queryString, {
            headers: {
                "Authorization": app.Store.jwt ? `Bearer ${app.Store.jwt}` : null
            
            }
        });
        const result = await response.json();
        return result;
    }
```

Then, on the same file we will add the new services:

```js
    getFavorites: async () => {
        try {
            return await API.fetch("account/favorites");
        } catch (e) {
            app.Router.go("/account/")
        }
    },     
    getWatchlist: async () => {
        try {
            return await API.fetch("account/watchlist");
        } catch (e) {
            app.Router.go("/account/")
        }
    
    },     
    saveToCollection: async (movie_id, collection) => {
        return await API.send("account/save-to-collection/", {
            movie_id, collection
        });
    },
```

### H4 - Creating the New Pages

Add this template to *index.html*

```html
    <template id="template-collection">
        <section>
            <ul id="movies-result">
                <animated-loading data-elements="5"
                    data-width="150px" data-height="220px">
                </animated-loading> 
            </ul>
        </section>
    </template>  
```
Then, create *CollectionPage.js*

```js
import { MovieItemComponent } from "./MovieItem.js";

export class CollectionPage extends HTMLElement {

    constructor(endpoint, title) {
        super();
        this.endpoint = endpoint;
        this.title = title;   
    }

    async render() {
        const movies = await this.endpoint()
        const ulMovies = this.querySelector("ul");
        ulMovies.innerHTML = "";
        if (movies && movies.length>0) {
            movies.forEach(movie => {
                const li = document.createElement("li");
                li.appendChild(new MovieItemComponent(movie));
                ulMovies.appendChild(li);
            });    
        } else {
            ulMovies.innerHTML = "<h3>There are no movies</h3>";
        }        
            ;
    }

    connectedCallback() {
        const template = document.getElementById("template-collection");
        const content = template.content.cloneNode(true);
        this.appendChild(content);  

        this.render();
    }
}
```

And two other components, *FavoritePage.js*

```js
import API from "../services/API.js";
import { CollectionPage } from "./CollectionPage.js";

export default class FavoritePage extends CollectionPage {

    constructor() {
        super(API.getFavorites, "Favorite Movies")
    }

}
customElements.define("favorite-page", FavoritePage);
```

And then *WatchlistPage.js*

```js
import API from "../services/API.js";
import { CollectionPage } from "./CollectionPage.js";

export default class WatchlistPage extends CollectionPage {

    constructor() {
        super(API.getWatchlist, "Movie Watchlist")
    }

}
customElements.define("watchlist-page", WatchlistPage);
```

And finally define the routes in *Routes.js*

```js
// ...

    {
        path: "/account/favorites",
        component: FavoritePage,
        loggedIn: true
    },    
    {
        path: "/account/watchlist",
        component: WatchlistPage,
        loggedIn: true
    }, 
```

### H5 - Adding to Collections

In *app.js* lets add the new two features:

```js
    saveToCollection: async (movie_id, collection) => {
        if (app.Store.loggedIn) {
            try {
                const response = await API.saveToCollection(movie_id, collection);
                if (response.success) {
                    switch(collection) {
                        case "favorite":
                            app.Router.go("/account/favorites")
                        break;
                        case "watchlist":
                            app.Router.go("/account/watchlist")
                    }
                } else {
                    app.showError("We couldn't save the movie.")
                }
            } catch (e) {
                console.log(e)
            }
        } else {
            app.Router.go("/account/");
        }
    }
```	

Back in *MovieDetailsPage.js*, add the following calls

```js
        this.querySelector("#btnFavorites").addEventListener("click", () => {
            app.saveToCollection(this.movie.id, "favorite")
        })
        this.querySelector("#btnWatchlist").addEventListener("click", () => {
            app.saveToCollection(this.movie.id, "watchlist")
        })

```

# Intermediate Course

## I-Passkeys

### I1 - Adding new dependences

Run in the console

```
go get "github.com/go-webauthn/webauthn/webauthn"
```

### I2 - Update our models and the database

Run in the database the following query:

```sql
CREATE SEQUENCE IF NOT EXISTS passkeys_id_seq;

CREATE TABLE "public"."passkeys" (
    "id" int4 NOT NULL DEFAULT nextval('passkeys_id_seq'::regclass),
    "user_id" int4,
    "keys" text,
    PRIMARY KEY ("id")
);
```

Then let's create a new model as *passkeyuser.go*

```go
package models

import "github.com/go-webauthn/webauthn/webauthn"

type PasskeyUser struct {
	ID          []byte
	DisplayName string
	Name        string

	Credentials []webauthn.Credential
}

func (u *PasskeyUser) WebAuthnID() []byte {
	return u.ID
}

func (u *PasskeyUser) WebAuthnName() string {
	return u.Name
}

func (u *PasskeyUser) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u *PasskeyUser) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func (u PasskeyUser) WebAuthnIcon() string {
	return ""
}

func (u *PasskeyUser) PutCredential(credential webauthn.Credential) {
	u.Credentials = append(u.Credentials, credential)
}

func (u *PasskeyUser) AddCredential(credential *webauthn.Credential) {
	u.Credentials = append(u.Credentials, *credential)
}

func (u *PasskeyUser) UpdateCredential(credential *webauthn.Credential) {
	for i, c := range u.Credentials {
		if string(c.ID) == string(credential.ID) {
			u.Credentials[i] = *credential
		}
	}
}

```

### I3 - Create the Passkey Repository

Let's update our *interfaces.go*

```go
type PasskeyStore interface {
	GetUserByEmail(userName string) (*models.PasskeyUser, error)
	GetUserByID(ID int) (*models.PasskeyUser, error)
	SaveUser(models.PasskeyUser)
	GenSessionID() (string, error)
	GetSession(token string) (webauthn.SessionData, bool)
	SaveSession(token string, data webauthn.SessionData)
	DeleteSession(token string)
}
```

Now let's create the *data/passkey_repository.go*

```go
package data

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"frontendmasters.com/movies/logger"
	"frontendmasters.com/movies/models"
	"github.com/go-webauthn/webauthn/webauthn"
)

// PasskeyRepository manages WebAuthn passkey data using a database.
type PasskeyRepository struct {
	db       *sql.DB                         // Database connection
	sessions map[string]webauthn.SessionData // In-memory session storage
	log      logger.Logger                   // Logger for debugging and errors
}

// NewPasskeyRepository initializes a new PasskeyRepository with a database connection.
func NewPasskeyRepository(db *sql.DB, log logger.Logger) *PasskeyRepository {
	return &PasskeyRepository{
		db:       db,
		sessions: make(map[string]webauthn.SessionData),
		log:      log,
	}
}

// GenSessionID generates a random session ID.
func (r *PasskeyRepository) GenSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetSession retrieves session data from the in-memory map.
func (r *PasskeyRepository) GetSession(token string) (webauthn.SessionData, bool) {
	r.log.Info(fmt.Sprintf("GetSession: %v", r.sessions[token]))
	val, ok := r.sessions[token]
	return val, ok
}

// SaveSession stores session data in the in-memory map.
func (r *PasskeyRepository) SaveSession(token string, data webauthn.SessionData) {
	r.log.Info(fmt.Sprintf("SaveSession: %s - %v", token, data))
	r.sessions[token] = data
}

// DeleteSession removes session data from the in-memory map.
func (r *PasskeyRepository) DeleteSession(token string) {
	r.log.Info(fmt.Sprintf("DeleteSession: %v", token))
	delete(r.sessions, token)
}

func (r *PasskeyRepository) GetUserByEmail(email string) (*models.PasskeyUser, error) {
	r.log.Info(fmt.Sprintf("Get User: %v", email))

	// Check if user exists by email
	var userID int
	err := r.db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&userID)
	if err == sql.ErrNoRows {
		r.log.Error("Failed to find new user", err)
		return nil, err
	} else if err != nil {
		r.log.Error("Failed to query user", err)
		return nil, err
	}

	// Fetch user credentials from passkeys table
	rows, err := r.db.Query("SELECT keys FROM passkeys WHERE user_id = $1", userID)
	if err != nil {
		r.log.Error("Failed to query passkeys", err)
		return nil, err
	}
	defer rows.Close()

	var credentials []webauthn.Credential
	for rows.Next() {
		var keys string
		if err := rows.Scan(&keys); err != nil {
			r.log.Error("Failed to scan passkey row", err)
			return nil, err
		}
		cred, err := deserializeCredential(keys)
		if err != nil {
			r.log.Error("Failed to deserialize credential", err)
			continue // Skip invalid credentials
		}
		credentials = append(credentials, cred)
	}

	// Construct and return PasskeyUser
	user := models.PasskeyUser{
		ID:          []byte(strconv.Itoa(userID)), // Convert int ID to byte slice
		Name:        email,
		DisplayName: email,
		Credentials: credentials,
	}
	return &user, nil
}

func (r *PasskeyRepository) GetUserByID(id int) (*models.PasskeyUser, error) {
	r.log.Info(fmt.Sprintf("Get User: %v", id))

	// Check if user exists by id
	var userID int
	err := r.db.QueryRow("SELECT id FROM users WHERE id = $1", id).Scan(&userID)
	if err == sql.ErrNoRows {
		r.log.Error("Failed to find new user", err)
		return nil, err
	} else if err != nil {
		r.log.Error("Failed to query user", err)
		return nil, err
	}

	// Fetch user credentials from passkeys table
	rows, err := r.db.Query("SELECT keys FROM passkeys WHERE user_id = $1", userID)
	if err != nil {
		r.log.Error("Failed to query passkeys", err)
		return nil, err
	}
	defer rows.Close()

	// Fetch the email of the user
	var email string
	err = r.db.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&email)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.Error("Failed to find user email", err)
			return nil, err
		}
		r.log.Error("Failed to query user email", err)
		return nil, err
	}

	var credentials []webauthn.Credential
	for rows.Next() {
		var keys string
		if err := rows.Scan(&keys); err != nil {
			r.log.Error("Failed to scan passkey row", err)
			return nil, err
		}
		cred, err := deserializeCredential(keys)
		if err != nil {
			r.log.Error("Failed to deserialize credential", err)
			continue // Skip invalid credentials
		}
		credentials = append(credentials, cred)
	}

	// Construct and return PasskeyUser
	user := models.PasskeyUser{
		ID:          []byte(strconv.Itoa(userID)), // Convert int ID to byte slice
		Name:        email,
		DisplayName: email,
		Credentials: credentials,
	}
	return &user, nil
}

// SaveUser updates the user's credentials in the database.
func (r *PasskeyRepository) SaveUser(user models.PasskeyUser) {
	r.log.Info(fmt.Sprintf("SaveUser: %v", user.WebAuthnName()))

	// Convert user ID from byte slice to integer
	userID, err := strconv.Atoi(string(user.ID))
	if err != nil {
		r.log.Error("Invalid user ID", err)
		return
	}

	// Insert new credentials
	for _, cred := range user.Credentials {
		keys, err := serializeCredential(cred)
		if err != nil {
			r.log.Error("Failed to serialize credential", err)
			continue
		}
		// Check if the key already exists in the database
		var exists bool
		err = r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM passkeys WHERE user_id = $1 AND keys = $2)", userID, keys).Scan(&exists)
		if err != nil {
			r.log.Error("Failed to check if passkey exists", err)
			continue
		}

		// Insert the key only if it does not already exist
		if !exists {
			_, err = r.db.Exec("INSERT INTO passkeys (user_id, keys) VALUES ($1, $2)", userID, keys)
			if err != nil {
				r.log.Error("Failed to insert passkey", err)
			}
		} else {
			r.log.Info(fmt.Sprintf("Passkey already exists for user_id: %d", userID))
		}
	}
}

// serializeCredential converts a WebAuthn credential to a JSON string.
func serializeCredential(cred webauthn.Credential) (string, error) {
	data, err := json.Marshal(cred)
	if err != nil {
		return "", fmt.Errorf("failed to marshal credential: %w", err)
	}
	return string(data), nil
}

// deserializeCredential converts a JSON string back to a WebAuthn credential.
func deserializeCredential(data string) (webauthn.Credential, error) {
	var cred webauthn.Credential
	err := json.Unmarshal([]byte(data), &cred)
	if err != nil {
		return webauthn.Credential{}, fmt.Errorf("failed to unmarshal credential: %w", err)
	}
	return cred, nil
}

```

### I4 - Create the Passkey handlers

Create *handlers/passkey_handler.go*

```go
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"frontendmasters.com/movies/data"
	"frontendmasters.com/movies/logger"
	"frontendmasters.com/movies/models"
	"frontendmasters.com/movies/token"
	"github.com/go-webauthn/webauthn/webauthn"
)

type WebAuthnHandler struct {
	storage  data.PasskeyStore
	logger   *logger.Logger
	webauthn *webauthn.WebAuthn
}

func NewWebAuthnHandler(storage data.PasskeyStore, logger *logger.Logger, webauthn *webauthn.WebAuthn) *WebAuthnHandler {
	return &WebAuthnHandler{
		storage:  storage,
		logger:   logger,
		webauthn: webauthn,
	}
}

func (h *WebAuthnHandler) writeJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *WebAuthnHandler) WebAuthnRegistrationBeginHandler(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		h.logger.Error("Unable to retrieve email", nil)
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}
	user, err := h.storage.GetUserByEmail(email)
	if err != nil {
		h.logger.Error("Failed to find user", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	options, session, err := h.webauthn.BeginRegistration(user)
	if err != nil {
		h.logger.Error("Unable to retrieve email", err)
		http.Error(w, "Can't begin WebAuthn Registration", http.StatusInternalServerError)

		return
	}

	// Make a session key and store the sessionData values
	t, err := h.storage.GenSessionID()
	if err != nil {
		h.logger.Error("Can't generate session id: %s", err)
	}

	h.storage.SaveSession(t, *session)

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    t,
		Path:     "api/passkey/registerStart",
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	h.writeJSONResponse(w, options)
}

func (h *WebAuthnHandler) WebAuthnRegistrationEndHandler(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		h.logger.Error("Unable to retrieve email", nil)
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}

	// Get the session key from cookie
	sid, err := r.Cookie("sid")
	if err != nil {
		h.logger.Error("Couldn't get the cookie for the session", err)
	}

	// Get the session data
	session, _ := h.storage.GetSession(sid.Value)

	user, err := h.storage.GetUserByEmail(email)
	if err != nil {
		h.logger.Error("Failed to find user", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	credential, err := h.webauthn.FinishRegistration(user, session, r)
	if err != nil {
		h.logger.Error("Coudln't finish the WebAuthn Registration", err)
		// clean up sid cookie
		http.SetCookie(w, &http.Cookie{
			Name:  "sid",
			Value: "",
		})
		http.Error(w, "Couldn't finish registration", http.StatusBadRequest)
		return
	}

	// Store the credential object
	user.AddCredential(credential)
	h.storage.SaveUser(*user)
	// Delete the session data
	h.storage.DeleteSession(sid.Value)
	http.SetCookie(w, &http.Cookie{
		Name:  "sid",
		Value: "",
	})

	h.writeJSONResponse(w, "{'success': true}")

}

func (h *WebAuthnHandler) WebAuthnAuthenticationBeginHandler(w http.ResponseWriter, r *http.Request) {
	type CollectionRequest struct {
		Email string `json:"email"`
	}
	var req CollectionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode collection request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := req.Email

	h.logger.Info("Finding user " + email)

	user, err := h.storage.GetUserByEmail(email) // Find the user

	if err != nil {
		h.logger.Error("Failed to find user", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	options, session, err := h.webauthn.BeginLogin(user)
	if err != nil {
		h.logger.Error("Coudln't start a login", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Make a session key and store the sessionData values
	t, err := h.storage.GenSessionID()
	if err != nil {
		h.logger.Error("Coudln't create a session ID", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	h.storage.SaveSession(t, *session)

	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    t,
		Path:     "api/passkey/loginStart",
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // TODO: SameSiteStrictMode maybe?
	})

	h.writeJSONResponse(w, options)
}

func (h *WebAuthnHandler) WebAuthnAuthenticationEndHandler(w http.ResponseWriter, r *http.Request) {
	// Get the session key from cookie
	sid, err := r.Cookie("sid")
	if err != nil {
		h.logger.Error("Coudln't get a session", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Get the session data stored from the function above
	session, _ := h.storage.GetSession(sid.Value)

	userID, err := strconv.Atoi(string(session.UserID)) // Convert []byte to int
	if err != nil {
		h.logger.Error("Failed to convert UserID to int", err)
		http.Error(w, "Invalid session data", http.StatusBadRequest)
		return
	}
	user, err := h.storage.GetUserByID(userID) // Get the user
	if err != nil {
		h.logger.Error("Failed to find user", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	credential, err := h.webauthn.FinishLogin(user, session, r)
	if err != nil {
		h.logger.Error("Coudln't finish login", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Handle credential.Authenticator.CloneWarning
	if credential.Authenticator.CloneWarning {
		h.logger.Error("Couldn't finish login", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// If login was successful
	user.UpdateCredential(credential)
	h.storage.SaveUser(*user)

	// Delete the session data
	h.storage.DeleteSession(sid.Value)
	http.SetCookie(w, &http.Cookie{
		Name:  "sid",
		Value: "",
	})

	// Add the new session cookie
	t, err := h.storage.GenSessionID()
	if err != nil {
		h.logger.Error("Couldn't generate session", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	h.storage.SaveSession(t, webauthn.SessionData{
		Expires: time.Now().Add(time.Hour),
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "sid",
		Value:    t,
		Path:     "/",
		MaxAge:   3600,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // TODO: SameSiteStrictMode maybe?
	})

	type PasskeyResponse struct {
		Success bool   `json:"success"`
		JWT     string `json:"jwt"`
	}
	h.logger.Info("Sending JWT for " + user.Name)
	// Return success response
	response := PasskeyResponse{
		Success: true,
		JWT:     token.CreateJWT(models.User{Email: user.Name}, *h.logger),
	}

	h.writeJSONResponse(w, response)
}

```

### I5 - Update the Main

Update `main.go` to support our new handlers

```go
	// WebAuthn Handlers
	wconfig := &webauthn.Config{
		RPDisplayName: "ReelingIt",
		RPID:          "localhost",
		RPOrigins:     []string{"http://localhost:8080"},
	}

	var webAuthnManager *webauthn.WebAuthn

	if webAuthnManager, err = webauthn.New(wconfig); err != nil {
		logInstance.Error("Error creating WebAuthn", err)
	}

	if err != nil {
		logInstance.Error("Error initialing Passkey engine", err)
	}

	passkeyRepo := data.NewPasskeyRepository(db, *logInstance)
	webAuthnHandler := handlers.NewWebAuthnHandler(passkeyRepo, logInstance, webAuthnManager)
	http.Handle("/api/passkey/registration-begin",
		accountHandler.AuthMiddleware(http.HandlerFunc(webAuthnHandler.WebAuthnRegistrationBeginHandler)))
	http.Handle("/api/passkey/registration-end",
		accountHandler.AuthMiddleware(http.HandlerFunc(webAuthnHandler.WebAuthnRegistrationEndHandler)))
	http.HandleFunc("/api/passkey/authentication-begin", webAuthnHandler.WebAuthnAuthenticationBeginHandler)
	http.HandleFunc("/api/passkey/authentication-end", webAuthnHandler.WebAuthnAuthenticationEndHandler)

```

### I6 - Add a dependency client side

In our HTML add the SimpleWebAuthn library script before app.js loading:

```html
    <script src="https://unpkg.com/@simplewebauthn/browser/dist/bundle/index.umd.min.js" defer></script>
```

### I7 - Update the UI

Change in *index.html* the template-account 

```html
 <template id="template-account">
        <section id="account">
            <h2>You are Logged In</h2>
            <button onclick="app.logout()">Log out</button>
            <button onclick="app.Router.go('/account/favorites')">Your Favorites</button>
            <button onclick="app.Router.go('/account/watchlist')">Your Watchlist</button>
            <!-- NEW -->
            <button onclick="app.addPasskey()">Add a Passkey for faster login</button>
        </section>
    </template>      
```

And also the template-login removing the `required` attribute for the password:

```html
 <template id="template-login">
        <section>
            <h2>Login into Your Account</h2>
            <form onsubmit="app.login(event)">
                <label for="login-email">Email</label>
                <input type="email" id="login-email" placeholder="Email" required autocomplete="email">
                <label for="login-password">Password</label>
                <input type="password" id="login-password" placeholder="Password" autocomplete="current-password">
                <button>Log In</button>

                <!-- NEW -->
                <button type="button" onclick="app.loginWithPasskey()">Log In with a Passkey</button>

                <p>If you don't have an account, please <a href="/account/register">register</a>.</p>
            </form>
        </section>
    </template>    
```

### I8 - Passkey Services

Create a new Service *services/Passkey.js*

```js
export const Passkeys = {
    register: async (username) => {
        try {
            // Get registration options with the challenge.
            const response = await fetch('/api/passkey/registration-begin', {
                method: 'POST', 
                headers: {
                    'Content-Type': 'application/json',
                    "Authorization": app.Store.jwt ? `Bearer ${app.Store.jwt}` : null
                },                
                body: JSON.stringify({username: username})
            });
    
            // Check if the options are ok.
            if (!response.ok) {
                const err = await response.json();
                app.showError('Failed to get registration options from server.' + err)
            }
    
            const options = await response.json();
    
            // This triggers the browser to display the passkey modal 
            // A new public-private-key pair is created.
            const attestationResponse = await SimpleWebAuthnBrowser.startRegistration({optionsJSON: options.publicKey});
    
            // Send attestationResponse back to server for verification and storage.
            const verificationResponse = await fetch('/api/passkey/registration-end', {
                method: 'POST',
                credentials: 'same-origin',
                headers: {
                    'Content-Type': 'application/json',
                    "Authorization": app.Store.jwt ? `Bearer ${app.Store.jwt}` : null
                },
                body: JSON.stringify(attestationResponse)
            });
    
            const msg = await verificationResponse.json();
            if (verificationResponse.ok) {
                app.showError("Your passkey was saved. You can use it next time to login")
            } else {
                app.showError(msg, false);
            }
        } catch (e) {
            app.showError('Error: ' + e.message, false);
        }        
    },
    authenticate: async (email) => {
        try {
            // Get login options from your server with the challenge
            const response = await fetch('/api/passkey/authentication-begin', {
                method: 'POST', headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({email})
            });
            const options = await response.json();
    
            // This triggers the browser to display the passkey / WebAuthn modal 
            // The challenge has been signed after this.
            const assertionResponse = await SimpleWebAuthnBrowser.startAuthentication({optionsJSON: options.publicKey});
    
            // Send assertionResponse back to server for verification.
            const verificationResponse = await fetch('/api/passkey/authentication-end', {
                method: 'POST',
                credentials: 'same-origin',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(assertionResponse)
            });
    
            const serverResponse = await verificationResponse.json();
            if (serverResponse.success) {
                app.Store.jwt = serverResponse.jwt;
                app.Router.go("/account/")
            } else {
                app.showError(msg, false);
            }
        } catch (e) {
            console.log(e)
            app.showError('We couldn\'t authenticate you using a Passkey', false);
        }        
    }
}
```

### I9 - App Controller

In *app.js* add

```js
    addPasskey: async () => {
        const username = "testuser";
        await Passkeys.register(username);
    },
    loginWithPasskey: async () => {
        const username = document.getElementById("login-email").value;
        if (username.length < 4) {
            app.showError("To use a passkey, enter your email address first")
        } else {
            await Passkeys.authenticate(username);
        }
    }    
```

### J-Server-Side render

### J1 - Create a new handler

Create *handlers/ssr_handler.go*

```go
package handlers

import (
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"frontendmasters.com/movies/data"
	"frontendmasters.com/movies/logger"
	"frontendmasters.com/movies/models"
)

// In main.go, add this new handler function before the main function
func SSRMovieDetailsHandler(movieRepo *data.MovieRepository, logInstance *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract movie ID from URL
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 {
			http.Error(w, "Movie ID required", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(pathParts[2])
		if err != nil {
			http.Error(w, "Invalid movie ID", http.StatusBadRequest)
			return
		}

		// Get movie from repository
		movie, err := movieRepo.GetMovieByID(id)
		if err != nil {
			if errors.Is(err, data.ErrMovieNotFound) {
				http.Error(w, "Movie not found", http.StatusNotFound)
			} else {
				logInstance.Error("Error fetching movie", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		// Serve the HTML with movie data
		w.Header().Set("Content-Type", "text/html")
		err = renderMovieDetails(w, movie)
		if err != nil {
			logInstance.Error("Error rendering movie details", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

// Add this function to render the HTML
func renderMovieDetails(w io.Writer, movie models.Movie) error {
	// Read the index.html file
	htmlContent, err := os.ReadFile("./public/index.html")
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Convert movie data to HTML
	genresHTML := ""
	for _, genre := range movie.Genres {
		genresHTML += fmt.Sprintf(`<li>%s</li>`, html.EscapeString(genre.Name))
	}

	castHTML := ""
	for _, actor := range movie.Casting {
		imageURL := "/images/generic_actor.jpg"
		if actor.ImageURL != nil {
			imageURL = *actor.ImageURL
		}
		castHTML += fmt.Sprintf(`
            <li>
                <img src="%s" alt="Picture of %s">
                <p>%s %s</p>
            </li>`,
			html.EscapeString(imageURL),
			html.EscapeString(actor.LastName),
			html.EscapeString(actor.FirstName),
			html.EscapeString(actor.LastName))
	}

	// Replace the main content
	mainContent := fmt.Sprintf(`
        <main>
            <article id="movie">
                <h2>%s</h2>
                <h3>%s</h3>
                <header>
                    <img src="%s" alt="Poster">
                    <youtube-embed id="trailer" data-url="%s"></youtube-embed>
                    <section id="actions">
                        <dl id="metadata">
                            <dt>Release Date</dt>
                            <dd>%d</dd>
                            <dt>Score</dt>
                            <dd>%.1f / 10</dd>
                            <dt>Original language</dt>
                            <dd>%s</dd>
                        </dl>
                        <button id="btnFavorites">Add to Favorites</button>
                        <button id="btnWatchlist">Add to Watchlist</button>
                    </section>
                </header>
                <ul id="genres">%s</ul>
                <p id="overview">%s</p>
                <ul id="cast">%s</ul>
            </article>
        </main>`,
		html.EscapeString(movie.Title),
		html.EscapeString(*movie.Tagline),
		html.EscapeString(*movie.PosterURL),
		html.EscapeString(*movie.TrailerURL),
		movie.ReleaseYear,
		*movie.Score,
		html.EscapeString(*movie.Language),
		genresHTML,
		html.EscapeString(*movie.Overview),
		castHTML)

	// Replace the main tag content in the HTML
	htmlStr := string(htmlContent)
	htmlStr = strings.Replace(htmlStr, "<main></main>", mainContent, 1)
	fmt.Println(htmlStr)

	// Write the response
	_, err = w.Write([]byte(htmlStr))
	return err
}
```

### J2 - Register the handler for the route

At *main.go* replace the current `/movies/` handler with:

```go
http.HandleFunc("/movies/", func(w http.ResponseWriter, r *http.Request) {
	if strings.Count(r.URL.Path, "/") == 2 && strings.HasPrefix(r.URL.Path, "/movies/") {
		handlers.SSRMovieDetailsHandler(movieRepo, logInstance)(w, r)
	} else {
		catchAllHandler(w, r)
	}
})
```

## K-Offline

### K1 - Update the Web App Manifest

Add to `app.webmanifest` the `display: standalone`, `scope` and `start_url` attributes.

### K2 - Create a Service Worker

Create *sw.js* in the root of the public folder:

```js
// service-worker.js

const CACHE_NAME = 'my-cache-v1';

// Install event - precache any initial resources if needed
self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(() => {
        // Skip waiting to activate immediately
        self.skipWaiting();
      })
  );
});

// Activate event - clean up old caches
self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          if (cacheName !== CACHE_NAME) {
            return caches.delete(cacheName);
          }
        })
      );
    }).then(() => {
      // Take control of clients immediately
      return self.clients.claim();
    })
  );
});

// Fetch event - handle caching strategies
self.addEventListener('fetch', (event) => {
  const requestUrl = new URL(event.request.url);

  // Handle /api/ requests (network first, cache fallback)
  if (requestUrl.pathname.startsWith('/api/')) {
    event.respondWith(
      fetch(event.request)
        .then((networkResponse) => {
          // Cache successful network response
          return caches.open(CACHE_NAME).then((cache) => {
            cache.put(event.request, networkResponse.clone());
            return networkResponse;
          });
        })
        .catch(() => {
          // If network fails, try cache
          return caches.match(event.request)
            .then((cachedResponse) => {
              return cachedResponse || Promise.reject('No network or cache available');
            });
        })
    );
  } 
  // Handle all other requests (stale-while-revalidate)
  else {
    event.respondWith(
      caches.open(CACHE_NAME).then((cache) => {
        return cache.match(event.request).then((cachedResponse) => {
          // Start fetching new version in background
          const fetchPromise = fetch(event.request)
            .then((networkResponse) => {
              // Update cache with new response
              cache.put(event.request, networkResponse.clone());
              return networkResponse;
            })
            .catch((error) => {
              console.error('Fetch failed:', error);
            });

          // Return cached version if available, otherwise wait for network
          return cachedResponse || fetchPromise;
        });
      })
    );
  }
});
```

### K3 - Register the Service Worker

In *app.js* change the `DOMContentLoad` event with:

```js
    navigator.serviceWorker.register("/sw.js")
```
