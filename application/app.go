package application

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type App struct {
	router http.Handler
	DB     *sql.DB
}

func New() *App {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB URL is not found")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Started Database Connection...")

	app := &App{
		DB: db,
	}

	app.LoadRoutes()

	return app
}

func (app *App) Start() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT not found")
	}

	log.Println("Starting the server on port:", port)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: app.router,
	}

	defer app.DB.Close()

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
