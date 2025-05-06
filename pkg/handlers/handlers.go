package handlers

import (
	"net/http"

	"github.com/ashparshp/bookings/pkg/config"
	"github.com/ashparshp/bookings/pkg/models"
	"github.com/ashparshp/bookings/pkg/render"
)

// Repo is the repository used by the handler
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository {
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandler(r * Repository) {
	Repo = r
}

// HomePage is the handler for the home page
func (m *Repository) HomePage (w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
}

// AboutPage is the handler for the about page
func (m *Repository) AboutPage (w http.ResponseWriter, r *http.Request) {
	// perform some logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	// send the data to the template
	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// ReservationPage renders the make a reservation page and displays form
func (m *Repository) ReservationPage (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "make-reservation.page.tmpl", &models.TemplateData{})
}

// GeneralsPage renders the room page
func (m *Repository) GeneralsPage (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.tmpl", &models.TemplateData{})
}

// MajorsPage renders the room page
func (m *Repository) MajorsPage (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.tmpl", &models.TemplateData{})
}

// AvailabilityPage renders the room page
func (m *Repository) AvailabilityPage (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailabilityPage renders the room page
func (m *Repository) PostAvailabilityPage (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "post-search-availability.page.tmpl", &models.TemplateData{})
}

// ContactPage renders the room page
func (m *Repository) ContactPage (w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl", &models.TemplateData{})
}
