package main

import (
	"net/http"

	"github.com/ashparshp/bookings/internal/config"
	"github.com/ashparshp/bookings/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(_ *config.AppConfig) http.Handler {
	
	/*
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Repo.HomePage))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.AboutPage))
	*/

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.HomePage)
	mux.Get("/about", handlers.Repo.AboutPage)
	mux.Get("/generals-quarters", handlers.Repo.GeneralsPage)
	mux.Get("/majors-suite", handlers.Repo.MajorsPage)
	mux.Get("/search-availibility", handlers.Repo.AvailabilityPage)
	mux.Post("/search-availibility", handlers.Repo.PostAvailabilityPage)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/contact", handlers.Repo.ContactPage)
	mux.Get("/make-reservation", handlers.Repo.ReservationPage)
	mux.Post("/make-reservation", handlers.Repo.PostReservationPage)


	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

