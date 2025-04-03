package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ashparshp/bookings/pkg/handlers"
	"github.com/ashparshp/bookings/pkg/render"

	"github.com/ashparshp/bookings/pkg/config"

	"github.com/alexedwards/scs/v2"
)
const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager

func main() {
	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCahce = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)

	render.NewTemplates(&app)

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