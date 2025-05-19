package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
			values := url.Values{}
			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			res, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal()
			}
			if res.StatusCode != e.expectedStatusCode {
				t.Errorf("For %s, expected %d, got %d", e.name, e.expectedStatusCode, res.StatusCode)
			}
			res.Body.Close()
		}

	}
}