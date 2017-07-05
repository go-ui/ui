package drivers

import (
	"image/draw"
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

	Render(Surface, int, int)

	Close()
}

type Surface interface {
	// Implement the draw interface to enable other libraries to interact
	// with this surface.
	draw.Image
}

// Driver defines an interface for the drivers to implement.
type Driver interface {
	CreateWindow(string, int, int, WindowType, Window, func(events.Event)) Window
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

type WindowType byte

const (
	NormalWindow WindowType = iota
	DialogWindow
	SplashWindow
	MenuWindow
)
