package ui

import "github.com/go-ui/ui/drivers"

// NoMatchingDriverError is returned when New could not find a driver matching the name
type NoMatchingDriverError struct {
	name string
}

// Error returns a string with the description of the error
func (n NoMatchingDriverError) Error() string {
	return "No matching GUI driver found for \"" + n.name + "\""
}

// UI offers functionality to write UIs in Go.
type UI struct {
	d  drivers.Driver
	ws []*Window
}

// New creates a new UI instance with the driver that implements name
func New(name string) (*UI, error) {
	d := drivers.Get(name)
	if d == nil {
		return nil, NoMatchingDriverError{name: name}
	}
	return &UI{
		d: d,
	}, nil
}

// Release closes all windows that were opened by this instance of UI and
// releases all resouces allocated by it. This instance can not be used
// afterwards.
func (ui *UI) Release() error {
	return ui.d.Release()
}
