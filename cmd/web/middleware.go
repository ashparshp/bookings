package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

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