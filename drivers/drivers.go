package drivers

import (
	"image/draw"
	"sync"

	"github.com/go-ui/ui/events"
)

// DriverNotFoundError indicates that the chosen driver is not available.
type DriverNotFoundError struct {
	Driver string
}

func (d DriverNotFoundError) Error() string {
	return "Driver not found: " + d.Driver
}

// Window abstracts the window implementation from the driver.
type Window interface {
	// Title returns the window title, if the window does not have a title, the
	// result is an undefined UTF-8 string.
	Title() string

	// SetTitle sets the title that is shown in the window. When the window
	// does not have a title this is a no-op.
	SetTitle(string)

	// Size returns the size of the window in pixel.
	Size() (int, int)

	// SetSize sets the size of the window in pixel.
	SetSize(int, int)

	// Position returns the position of the window on the screen in pixel.
	Position() (int, int)

	// SetPosition sets the position of the window on the screen in pixel.
	SetPosition(int, int)

	// Render renders the provided surface at the provided position onto the window.
	Render(Surface, int, int)

	// Close closes the window and frees the resources allocated for it.
	Close()
}

// Surface abstracts a surface that can be drawn upon.
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

// WindowType indicates the type of the window. Some types may not be available
// on all systems.
type WindowType byte

const (
	// NormalWindow indicates a normal, decorated window.
	NormalWindow WindowType = iota

	// DialogWindow indicates a dialog window, often these have a smaller
	// border and may be attached to the parent window.
	DialogWindow

	// SplashWindow indicates a window without or with limited decorations.
	SplashWindow

	// MenuWindow indicates a window used for overlay menues, like context
	// menues.
	MenuWindow
)
