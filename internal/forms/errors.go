package forms

import "net/http"

type errors map[string][]string

// Add adds a new error message for a given form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get retrieves the first error message for the given field
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

// Has checks if form field is in post and not empty
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

