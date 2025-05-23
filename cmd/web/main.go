package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ashparshp/bookings/internal/driver"
	"github.com/ashparshp/bookings/internal/handlers"
	"github.com/ashparshp/bookings/internal/helpers"
	"github.com/ashparshp/bookings/internal/models"
	"github.com/ashparshp/bookings/internal/render"

	"github.com/ashparshp/bookings/internal/config"

	"github.com/alexedwards/scs/v2"
)
const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	/*
	http.HandleFunc("/", handlers.Repo.HomePage)
	http.HandleFunc("/about", handlers.Repo.AboutPage)
	*/
	
	fmt.Println("Server running on port", portNumber)
	/*
	err = http.ListenAndServe(portNumber, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
	*/

	srv := &http.Server{
		Addr: portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=ashparsh password=")
	if err != nil {
		log.Fatal("Cannot connect to database")
	}
	log.Println("Connected to database")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCahce = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)

	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}