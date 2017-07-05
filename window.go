package ui

import (
	"github.com/go-ui/ui/drivers"
	"github.com/go-ui/ui/events"
)

// Window holds the information of a window that is displayed on the screen.
type Window struct {
	w drivers.Window
	v View

	OnResize  func()
	OnMove    func()
	OnKeyDown func(events.KeyEvent)
}

type windowConfig struct {
	title  string
	width  int
	height int
	view   View
	typ    WindowType
	parent *Window
}

// WindowConfig represents a configuration function that sets specific
// properties for the window.
type WindowConfig func(*windowConfig)

// WithSize sets the size for a new window.
func WithSize(w, h int) WindowConfig {
	return func(c *windowConfig) {
		c.width = w
		c.height = h
	}
}

// WithTitle sets the title for a new window.
func WithTitle(t string) WindowConfig {
	return func(c *windowConfig) {
		c.title = t
	}
}

// WithView sets the initial view for a new window.
func WithView(v View) WindowConfig {
	return func(c *windowConfig) {
		c.view = v
	}
}

// WithType sets the window type for a new window.
func WithType(t WindowType) WindowConfig {
	return func(c *windowConfig) {
		c.typ = t
	}
}

// WithParentWindow sets a parent window for a new window.
func WithParentWindow(w *Window) WindowConfig {
	return func(c *windowConfig) {
		c.parent = w
	}
}

// NewWindow creates a new window for the UI.
func (ui *UI) NewWindow(cf ...WindowConfig) *Window {
	c := windowConfig{
		title:  "UI",
		width:  400,
		height: 300,
		view:   DummyView{},
		typ:    NormalWindow,
		parent: nil,
	}

	for _, f := range cf {
		f(&c)
	}

	w := &Window{v: c.view}

	var p drivers.Window
	if c.parent != nil {
		p = c.parent.w
	}

	w.w = ui.d.CreateWindow(c.title, c.width, c.height, drivers.WindowType(c.typ), p, w.event)

	return w
}

func (w *Window) event(e events.Event) {
	switch ev := e.(type) {

	case events.KeyEvent:
		if w.OnKeyDown != nil {
			w.OnKeyDown(ev)
		}

	default:
		w.v.Event(ev)

	}
}

// WindowType indiates the type of the window, this can be used to create borderless windows.
type WindowType byte

const (
	// NormalWindow indicates that the window has normal window decorations, this is used for most main windows.
	NormalWindow WindowType = iota

	// DialogWindow indicates that the window is a dialog, like a message box.
	DialogWindow

	// SplashWindow indicates a borderless window that is often displayed while the main application is loading.
	SplashWindow

	// MenuWindow indicates that the window is a tooltip or pop-up menu.
	MenuWindow
)
