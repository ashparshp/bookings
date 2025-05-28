package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ashparshp/bookings/internal/config"
	"github.com/ashparshp/bookings/internal/driver"
	"github.com/ashparshp/bookings/internal/forms"
	"github.com/ashparshp/bookings/internal/helpers"
	"github.com/ashparshp/bookings/internal/models"
	"github.com/ashparshp/bookings/internal/render"
	"github.com/ashparshp/bookings/internal/repository"
	"github.com/ashparshp/bookings/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

// Repo is the repository used by the handler
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository {
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository for testing
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository {
		App: a,
		DB:  dbrepo.NewTestRepo(a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandler(r * Repository) {
	Repo = r
}

// HomePage is the handler for the home page
func (m *Repository) HomePage (w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// AboutPage is the handler for the about page
func (m *Repository) AboutPage (w http.ResponseWriter, r *http.Request) {
	// send the data to the template
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// ReservationPage renders the make a reservation page and displays form
func (m *Repository) ReservationPage (w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")
	
	StringMap := (map[string]string{})
	StringMap["start_date"] = sd
	StringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: StringMap,
	})
}

// PostReservationPage handles the posting of a reservation form
func (m *Repository) PostReservationPage (w http.ResponseWriter, r *http.Request) {

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	/*
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	*/

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	/*
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}
	*/

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email", "phone")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert reservation into database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate: reservation.StartDate,
		EndDate:   reservation.EndDate,
		RoomID:    reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't insert room restriction")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// send an email to the user
	htmlMessage := fmt.Sprintf(`
	<strong>Reservation Confirmation</strong><br>
	Dear %s:, <br>
	Thank you for your reservation from %s to %s.<br>
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))


	msg := models.MailData{
		To:     reservation.Email,
		From:    "me@here.com",
		Subject: "Reservation Confirmation",
		Content: htmlMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	// send an email to the admin
	adminMessage := fmt.Sprintf(`
	<strong>New Reservation</strong><br>
	New reservation for %s %s from %s to %s.<br>
	Email: %s<br>
	Phone: %s<br>
	Room ID: %d<br>
	`, reservation.FirstName, reservation.LastName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"), reservation.Email, reservation.Phone, reservation.RoomID)
	
	adminMsg := models.MailData{
		To:     "me@here.com",
		From:    "me@here.com",
		Subject: "New Reservation",
		Content: adminMessage,
		Template: "basic.html",
	}
	
	m.App.MailChan <- adminMsg


	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusTemporaryRedirect)
}

// GeneralsPage renders the room page
func (m *Repository) GeneralsPage (w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// MajorsPage renders the room page
func (m *Repository) MajorsPage (w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// AvailabilityPage renders the room page
func (m *Repository) AvailabilityPage (w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailabilityPage handles post
func (m *Repository) PostAvailabilityPage (w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	RoomID  string    `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handles request for availability and sends JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	available, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	resp := jsonResponse{
		OK:      available,
		Message: "",
		StartDate: sd,
		EndDate:   ed,
		RoomID: strconv.Itoa(roomID),
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// ContactPage renders the room page
func (m *Repository) ContactPage (w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummaryPage renders the room page
func (m *Repository) ReservationSummaryPage (w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("Can't get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	StringMap := (map[string]string{})
	StringMap["start_date"] = sd
	StringMap["end_date"] = ed


	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
		StringMap: StringMap,
	})
}

// ChooseRoomPage renders the room page
func (m *Repository) ChooseRoomPage (w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoomPage takes URL parameters, builds a sessional variable and redirects to make reservation page
func (m *Repository) BookRoomPage (w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	
	var res models.Reservation
	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate

	room, err := m.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// LoginPage renders the login page
func (m *Repository) LoginPage (w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostLoginPage handles the login form submission
func (m *Repository) PostLoginPage (w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/user/login", http.StatusTemporaryRedirect)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.AuthenticateUser(email, password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogoutPage handles the logout process
func (m *Repository) LogoutPage (w http.ResponseWriter, r *http.Request) {
	m.App.Session.Destroy(r.Context())
	m.App.Session.RenewToken(r.Context())
	m.App.Session.Put(r.Context(), "flash", "Logged out successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AdminDashboardPage renders the admin dashboard page
func (m *Repository) AdminDashboardPage(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{
	})
}

// AdminNewReservationPage renders the admin new reservations page
func (m *Repository) AdminNewReservationPage(w http.ResponseWriter, r *http.Request) {
	newReservations, err := m.DB.AllNewReservations()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Unable to retrieve new reservations")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["new_reservations"] = newReservations
	render.Template(w, r, "admin-new-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminAllReservationsPage renders the admin all reservations page
func (m *Repository) AdminAllReservationsPage(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Unable to retrieve reservations")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations

	render.Template(w, r, "admin-all-reservations.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// AdminShowReservationPage renders the admin show reservation page
func (m *Repository) AdminShowReservationPage(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(pathSegments[4])
	if err != nil {
		helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "Invalid reservation ID")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	src := pathSegments[3]
	stringMap:= make(map[string]string)
	stringMap["src"] = src

	year := r.URL.Query().Get("y")
	if year != "" {
		stringMap["year"] = year
	}
	month := r.URL.Query().Get("m")
	if month != "" {
		stringMap["month"] = month
	}

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "Unable to retrieve reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	if res.Room.ID == 0 {
		helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "No room found for this reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "admin-show-reservation.page.tmpl", &models.TemplateData{
		Data: data,
		StringMap: stringMap,
		Form: forms.New(nil),
	})
}

// AdminPostShowReservationPage handles the post request for showing a reservation
func (m *Repository) AdminPostShowReservationPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	pathSegments := strings.Split(r.RequestURI, "/")

	id, err := strconv.Atoi(pathSegments[4])
	if err != nil {
		helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "Invalid reservation ID")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	src := pathSegments[3]
	stringMap:= make(map[string]string)
	stringMap["src"] = src

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "Unable to retrieve reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")
	
	err = m.DB.UpdateReservation(res, id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	month := r.Form.Get("month")
	year := r.Form.Get("year")

	m.App.Session.Put(r.Context(), "flash", "Reservation updated!")

	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

// AdminReservationCalendarPage renders the admin reservation calendar page
func (m *Repository) AdminReservationCalendarPage(w http.ResponseWriter, r *http.Request) {
	// assume that there is no month/year specified
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, err := strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		month, err := strconv.Atoi(r.URL.Query().Get("m"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]interface{})
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	previous := now.AddDate(0, -1, 0)
	
	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")
	previousMonth := previous.Format("01")
	previousMonthYear := previous.Format("2006")
	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["previous_month"] = previousMonth
	stringMap["previous_month_year"] = previousMonthYear
	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	// get the first and last day of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	for _, room := range rooms {
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; !d.After(lastOfMonth); d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		restrictions , err := m.DB.GetRestrictionsForRoomByDate(room.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, restriction := range restrictions {
			if restriction.ReservationID > 0 {
				for d := restriction.StartDate; !d.After(restriction.EndDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = restriction.ReservationID
				}
			} else {
				blockMap[restriction.StartDate.Format("2006-01-2")] = restriction.ID
			}
		}

		m.App.Session.Put(r.Context(), fmt.Sprintf("reservation_map_%d", room.ID), reservationMap)
		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", room.ID), blockMap)

		data[fmt.Sprintf("reservation_map_%d", room.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", room.ID)] = blockMap
	}

	render.Template(w, r, "admin-reservations-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data: data,
		IntMap: intMap,
	})
}

// AdminProcessReservationPage processes a reservation based on the source and ID
func (m *Repository) AdminProcessReservationPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "Invalid reservation ID")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	src := chi.URLParam(r, "src")

	err = m.DB.UpdateProcessedForReservation(id, 1)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "Unable to process reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	m.App.Session.Put(r.Context(), "flash", "Reservation processed!")

	if year == "" {
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
		return
	}

}

// AdminDeleteReservationPage deletes a reservation based on the source and ID
func (m *Repository) AdminDeleteReservationPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "Invalid reservation ID")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	src := chi.URLParam(r, "src")

	err = m.DB.DeleteReservation(id)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "Unable to delete reservation")
		http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
		return
	}
	
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation deleted!")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
		return
	}
}

// AdminPostReservationCalendarPage handles the post request for the admin reservation calendar
func (m *Repository) AdminPostReservationCalendarPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	// process blocks
	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	forms := forms.New(r.PostForm)

	for _, room := range rooms {
		curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", room.ID)).(map[string]int)
		for name, value := range curMap {
			if val, ok := curMap[name]; ok {
				if val > 0 {
					if !forms.Has(fmt.Sprintf("remove_block_%d_%s", room.ID, name)) {
						// delete the restriction by id
						err = m.DB.DeleteBlockByID(value)
						if err != nil {
							m.App.Session.Put(r.Context(), "error", "Unable to delete block")
							http.Redirect(w, r, "/admin/reservations-calendar", http.StatusSeeOther)
							return
						}
					}
				}
			}
		}
	}

	// now handle new blocks
	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block_") {
			parts := strings.Split(name, "_")

			roomID, err := strconv.Atoi(parts[2])
			if err != nil {
				m.App.Session.Put(r.Context(), "error", "Invalid room ID")
				http.Redirect(w, r, "/admin/reservations-calendar", http.StatusSeeOther)
				return
			}

			dateStr := parts[3]
			layout := "2006-01-2"
			blockDate, err := time.Parse(layout, dateStr)
			if err != nil {
				m.App.Session.Put(r.Context(), "error", "Invalid date format")
				http.Redirect(w, r, "/admin/reservations-calendar", http.StatusSeeOther)
				return
			}

			err = m.DB.InsertBlockForRoom(roomID, blockDate)
			if err != nil {
				m.App.ErrorLog.Println("Error inserting block for room:", err)
				m.App.Session.Put(r.Context(), "error", "Unable to insert block")
				http.Redirect(w, r, "/admin/reservations-calendar", http.StatusSeeOther)
				return
			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Calendar updated")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}