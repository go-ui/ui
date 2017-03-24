package drivers

import (
	"sync"

	"github.com/go-ui/ui/events"
)

type Window interface {
	Title() string
	SetTitle(string)

	Size() (int, int)
	SetSize(int, int)

	Position() (int, int)
	SetPosition(int, int)

	Close()
}

// Driver defines an interface for the drivers to implement.
type Driver interface {
	CreateWindow(string, int, int, func(events.Event)) Window
	Release() error
}

var (
	m       sync.Mutex
	drivers = make(map[string]func() Driver)
)

// Set sets the driver factory function for the given name/ID
func Set(name string, f func() Driver) {
	m.Lock()
	drivers[name] = f
	m.Unlock()
}

// Get returns the driver for the given name/ID
func Get(name string) Driver {
	m.Lock()
	f, ok := drivers[name]
	m.Unlock()
	if ok {
		return f()
	}
	return nil
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
