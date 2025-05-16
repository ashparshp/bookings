package forms

import "net/url"

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

