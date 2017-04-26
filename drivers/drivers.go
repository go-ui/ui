package drivers

import (
	"image"
	"image/color"
	"sync"

	"github.com/go-ui/ui/events"
)

type DriverNotFoundError struct {
	Driver string
}

func (d DriverNotFoundError) Error() string {
	return "Driver not found: " + d.Driver
}

type Window interface {
	Title() string
	SetTitle(string)

	Size() (int, int)
	SetSize(int, int)

	Position() (int, int)
	SetPosition(int, int)

	Close()
}

type Surface interface {
	// Changed returns iff the surface has been changed. This value is set to
	// true by the following functions and should be reset by the driver
	// implementing Surface.
	Changed() bool

	// PutPixel sets a single pixel to the specified color. Writes outside of
	// the surface area should be ignored.
	PutPixel(int, int, color.Color)

	// PutImage renders an image to at the specified position. The image should
	// be trimmed at the surface edges.
	PutImage(int, int, image.Image)

	// PutScaled renders an image to the surface while changing the dimensions,
	PutScaled(int, int, int, int, image.Image)
}

// Driver defines an interface for the drivers to implement.
type Driver interface {
	CreateWindow(string, int, int, func(events.Event)) Window
	CreateSurface(int, int) Surface
	Release() error
}

var (
	m       sync.Mutex
	drivers = make(map[string]func() (Driver, error))
)

// Set sets the driver factory function for the given name/ID
func Set(name string, f func() (Driver, error)) {
	m.Lock()
	drivers[name] = f
	m.Unlock()
}

// Get returns the driver for the given name/ID
func Get(name string) (Driver, error) {
	m.Lock()
	f, ok := drivers[name]
	m.Unlock()
	if ok {
		return f()
	}
	return nil, DriverNotFoundError{Driver: name}
}

// List returns a list of available drivers
func List() (l []string) {
	m.Lock()
	l = make([]string, 0, len(drivers))
	for name := range drivers {
		l = append(l, name)
	}
	m.Unlock()
	return
}
