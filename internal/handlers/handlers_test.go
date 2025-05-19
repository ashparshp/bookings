package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key  string
	value string
}

var theTests = []struct {
	name       string
	method     string
	url        string
	params      []postData
	expectedStatusCode int
}{
	{
		name: "home",
		method: "GET",
		url: "/",
		params: []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name: "about",
		method: "GET",
		url: "/about",
		params: []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name: "gs",
		method: "GET",
		url: "/generals-quarters",
		params: []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name: "ms",
		method: "GET",
		url: "/majors-suite",
		params: []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name: "sa",
		method: "GET",
		url: "/search-availability",
		params: []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name: "contact",
		method: "GET",
		url: "/contact",
		params: []postData{},
		expectedStatusCode: http.StatusOK,
	},
	{
		name: "ms",
		method: "GET",
		url: "/make-reservation",
		params: []postData{},
		expectedStatusCode: http.StatusOK,
	},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()
	
	for _, e := range theTests {
		if e.method == "GET" {
			res, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal()
			}

			if res.StatusCode != e.expectedStatusCode {
				t.Errorf("For %s, expected %d, got %d", e.name, e.expectedStatusCode, res.StatusCode)
			}
		} else {

		}

	}
}