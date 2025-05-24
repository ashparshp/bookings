package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ashparshp/bookings/internal/models"
)

var theTests = []struct {
    name               string
    method             string
    url                string
    expectedStatusCode int
}{
    {
        name:               "home",
        method:             "GET",
        url:                "/",
        expectedStatusCode: http.StatusOK,
    },
    {
        name:               "about",
        method:             "GET",
        url:                "/about",
        expectedStatusCode: http.StatusOK,
    },
    {
        name:               "gs",
        method:             "GET",
        url:                "/generals-quarters",
        expectedStatusCode: http.StatusOK,
    },
    {
        name:               "ms",
        method:             "GET",
        url:                "/majors-suite",
        expectedStatusCode: http.StatusOK,
    },
    {
        name:               "sa",
        method:             "GET",
        url:                "/search-availability",
        expectedStatusCode: http.StatusOK,
    },
    {
        name:               "contact",
        method:             "GET",
        url:                "/contact",
        expectedStatusCode: http.StatusOK,
    },
    /*
        {
            name: "psa",
            method: "POST",
            url: "/search-availability",
            params: []postData{
                {"start", "2025-10-01"},
                {"end", "2025-10-02"},
            },
            expectedStatusCode: http.StatusOK,
        },
        {
            name: "psa-json",
            method: "POST",
            url: "/search-availability-json",
            params: []postData{
                {key: "start", value: "2025-10-01"},
                {key: "end", value: "2025-10-02"},
            },
            expectedStatusCode: http.StatusOK,
        },
        {
            name: "mr",
            method: "POST",
            url: "/make-reservation",
            params: []postData{
                {key: "first_name", value: "ashparsh"},
                {key: "last_name", value: "pandey"},
                {key: "email", value: "ashparsh@gmail.com"},
                {key: "phone", value: "1234567890"},
            },
            expectedStatusCode: http.StatusOK,
        },
    */
}

func TestHandlers(t *testing.T) {
    routes := getRoutes()
    ts := httptest.NewTLSServer(routes)
    defer ts.Close()

    for _, e := range theTests {
        res, err := ts.Client().Get(ts.URL + e.url)
        if err != nil {
            t.Log(err)
            t.Fatal()
        }

        if res.StatusCode != e.expectedStatusCode {
            t.Errorf("For %s, expected %d, got %d", e.name, e.expectedStatusCode, res.StatusCode)
        }
    }
}

func TestRepository_Reservation(t *testing.T) {
    reservation := models.Reservation{
        RoomID: 1,
        Room: models.Room{
            ID:       1,
            RoomName: "General's Quarters",
        },
    }

    req, _ := http.NewRequest("GET", "/make-reservation", nil)
    ctx := getCtx(req)
    req = req.WithContext(ctx)

    rr := httptest.NewRecorder()
    session.Put(ctx, "reservation", reservation)

    handler := http.HandlerFunc(Repo.ReservationPage)

    handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusOK {
        t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
    }

    // test case where reservation is not in session (reset everything)
    req, _ = http.NewRequest("GET", "/make-reservation", nil)
    ctx = getCtx(req)
    req = req.WithContext(ctx)
    rr = httptest.NewRecorder()

    handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusTemporaryRedirect {
        t.Errorf("Expected status code %d, got %d", http.StatusTemporaryRedirect, rr.Code)
    }

    // test with non-existent room
    req, _ = http.NewRequest("GET", "/make-reservation", nil)
    ctx = getCtx(req)
    req = req.WithContext(ctx)
    rr = httptest.NewRecorder()
    reservation.RoomID = 999 // non-existent room ID

    session.Put(ctx, "reservation", reservation)

    handler.ServeHTTP(rr, req)
    if rr.Code != http.StatusTemporaryRedirect {
        t.Errorf("Expected status code %d, got %d", http.StatusTemporaryRedirect, rr.Code)
    }
}

func getCtx(req *http.Request) context.Context {
    ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
    if err != nil {
        log.Println("Error loading session:", err)
        return nil
    }
    return ctx
}

func TestRepository_PostReservation(t *testing.T) {
    // Create a reservation that we'll use in the session
    reservation := models.Reservation{
        RoomID:    1,
        StartDate: time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
        EndDate:   time.Date(2025, 01, 02, 0, 0, 0, 0, time.UTC),
        Room: models.Room{
            ID:       1,
            RoomName: "General's Quarters",
        },
    }

    // Case 1: valid form submission
    postedData := url.Values{}
    postedData.Add("first_name", "John")
    postedData.Add("last_name", "Smith")
    postedData.Add("email", "john@example.com")
    postedData.Add("phone", "123456789")

    req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
    ctx := getCtx(req)
    req = req.WithContext(ctx)

    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    rr := httptest.NewRecorder()

    // Put reservation in session
    session.Put(ctx, "reservation", reservation)

    handler := http.HandlerFunc(Repo.PostReservationPage)
    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusTemporaryRedirect {
        t.Errorf("PostReservationPage handler returned wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
    }

    // Case 2: missing session data
    req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
    ctx = getCtx(req)
    req = req.WithContext(ctx)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    rr = httptest.NewRecorder()

    handler = http.HandlerFunc(Repo.PostReservationPage)
    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusTemporaryRedirect {
        t.Errorf("PostReservationPage handler returned wrong status code for missing session: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
    }

    // Case 3: invalid form data
    postedData = url.Values{}
    postedData.Add("first_name", "J") // Too short, minimum length is 3
    postedData.Add("last_name", "Smith")
    postedData.Add("email", "invalid-email")
    postedData.Add("phone", "123456789")

    req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
    ctx = getCtx(req)
    req = req.WithContext(ctx)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    rr = httptest.NewRecorder()

    // Put reservation in session
    session.Put(ctx, "reservation", reservation)

    handler = http.HandlerFunc(Repo.PostReservationPage)
    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("PostReservationPage handler returned wrong status code for invalid form: got %d, wanted %d", rr.Code, http.StatusOK)
    }
}

func TestRepository_ReservationSummary(t *testing.T) {
    // Case 1: reservation in session
    reservation := models.Reservation{
        RoomID:    1,
        StartDate: time.Date(2025, 01, 01, 0, 0, 0, 0, time.UTC),
        EndDate:   time.Date(2025, 01, 02, 0, 0, 0, 0, time.UTC),
        Room: models.Room{
            ID:       1,
            RoomName: "General's Quarters",
        },
    }

    req, _ := http.NewRequest("GET", "/reservation-summary", nil)
    ctx := getCtx(req)
    req = req.WithContext(ctx)

    rr := httptest.NewRecorder()

    session.Put(ctx, "reservation", reservation)

    handler := http.HandlerFunc(Repo.ReservationSummaryPage)
    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusOK {
        t.Errorf("ReservationSummaryPage handler returned wrong status code: got %d, wanted %d", rr.Code, http.StatusOK)
    }

    // Case 2: reservation not in session
    req, _ = http.NewRequest("GET", "/reservation-summary", nil)
    ctx = getCtx(req)
    req = req.WithContext(ctx)

    rr = httptest.NewRecorder()

    handler = http.HandlerFunc(Repo.ReservationSummaryPage)
    handler.ServeHTTP(rr, req)

    if rr.Code != http.StatusTemporaryRedirect {
        t.Errorf("ReservationSummaryPage handler returned wrong status code when no reservation in session: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
    }
}
