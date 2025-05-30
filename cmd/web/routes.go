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

	mux.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.Get("/", handlers.Repo.HomePage)
	mux.Get("/about", handlers.Repo.AboutPage)
	mux.Get("/generals-quarters", handlers.Repo.GeneralsPage)
	mux.Get("/majors-suite", handlers.Repo.MajorsPage)
	mux.Get("/search-availability", handlers.Repo.AvailabilityPage)
	mux.Post("/search-availability", handlers.Repo.PostAvailabilityPage)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoomPage)
	mux.Get("/book-room", handlers.Repo.BookRoomPage)
	mux.Get("/contact", handlers.Repo.ContactPage)
	mux.Get("/make-reservation", handlers.Repo.ReservationPage)
	mux.Post("/make-reservation", handlers.Repo.PostReservationPage)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummaryPage)
	mux.Post("/reservation-summary", handlers.Repo.ReservationSummaryPage)
	mux.Get("/user/login", handlers.Repo.LoginPage)
	mux.Post("/user/login", handlers.Repo.PostLoginPage)
	mux.Get("/user/logout", handlers.Repo.LogoutPage)
	mux.Route("/admin", func(mux chi.Router) {
		// mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboardPage)
		mux.Get("/reservations-all", handlers.Repo.AdminAllReservationsPage)
		mux.Get("/reservations-new", handlers.Repo.AdminNewReservationPage)
		mux.Get("/reservations-calendar", handlers.Repo.AdminReservationCalendarPage)
		mux.Post("/reservations-calendar", handlers.Repo.AdminPostReservationCalendarPage)
		mux.Get("/process-reservation/{src}/{id}/do", handlers.Repo.AdminProcessReservationPage)
		mux.Get("/delete-reservation/{src}/{id}/do", handlers.Repo.AdminDeleteReservationPage)

		mux.Get("/reservations/{src}/{id}/show", handlers.Repo.AdminShowReservationPage)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservationPage)
		
	})
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
