package models

import "github.com/ashparshp/bookings/internal/forms"

// TemplateData holds data sent from the handlers to github.com/ashparshp/bookings
type TemplateData struct {
	StringMap map[string]string
	IntMap map[string]int
	FloatMap map[string]float32
	Data map[string]interface{}
	CSRFToken string
	Flash string
	Warning string
	Error string
	Form *forms.Form
}