package ui

import (
	"github.com/go-ui/ui/drivers"
	"github.com/go-ui/ui/events"
)

// View is the abstraction of a thing(tm) that can render in a window and can
// reveive events from a window.
type View interface {
	// Dirty signals the window that this view is dirty and should be repainted.
	Dirty() bool

	// Render renders the view into an image, if the view has not been resized, this will reuse the image.
	Render(drivers.Surface)

	// Event received events from the window, so the view can react to them.
	Event(events.Event)
}

// DummyView implements a basic view that does nothing.
type DummyView struct{}

// Dirty always returns false.
func (DummyView) Dirty() bool { return false }

// Render does not render anything.
func (DummyView) Render(_ drivers.Surface) {}

// Event ignores all events.
func (DummyView) Event(_ events.Event) {}
