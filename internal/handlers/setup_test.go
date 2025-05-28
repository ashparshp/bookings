package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ashparshp/bookings/internal/config"
	"github.com/ashparshp/bookings/internal/models"
	"github.com/ashparshp/bookings/internal/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
)

var functions = template.FuncMap{
	"humanDate": render.HumanDate,
	"formatDate": render.FormatDate,
	"iterate": render.Iterate,
	"add": render.Add,
}
var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

func TestMain(m *testing.M) {
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

    // change this to true when in production
    app.InProduction = false

    infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
    app.InfoLog = infoLog

    // Use a null writer for error logs during tests
    errorLog := log.New(io.Discard, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
    app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	app.Session = session

	// create a channel for mail
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)

	// listen for mail
	listenForMail()


	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCahce = true

	repo := NewTestRepo(&app)
	NewHandler(repo)

	render.NewRenderer(&app)
	app.Session = session

	os.Exit(m.Run())
}

func listenForMail() {
	go func() {
		for {
			_ = <-app.MailChan
		}
	}()
}

func getRoutes() http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.HomePage)
	mux.Get("/about", Repo.AboutPage)
	mux.Get("/generals-quarters", Repo.GeneralsPage)
	mux.Get("/majors-suite", Repo.MajorsPage)
	mux.Get("/search-availability", Repo.AvailabilityPage)
	mux.Post("/search-availability", Repo.PostAvailabilityPage)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)
	mux.Get("/contact", Repo.ContactPage)
	mux.Get("/make-reservation", Repo.ReservationPage)
	mux.Post("/make-reservation", Repo.PostReservationPage)
	mux.Get("/reservation-summary", Repo.ReservationSummaryPage)


	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

// NoSurf adds CSFR protection to POST request
func NoSurf(next http.Handler) http.Handler {
	csfrHandler := nosurf.New(next)

	csfrHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path: "/",
		Secure: app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csfrHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTestTemplateCache creates a template cache as a map
func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))

	if err != nil {
		return myCache, err
	}

	// range through all files with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}