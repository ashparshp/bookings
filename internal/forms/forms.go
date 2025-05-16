package forms

import "net/url"

// Valid returns true if there are no errors, false otherwise
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// errors is a custom form struct, embed url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New initializes a new Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

