package ui

import (
	"github.com/go-ui/ui/drivers"
	"github.com/go-ui/ui/events"
)

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
// properties for the window
type WindowConfig func(*windowConfig)

// WithSize sets the initial size for a window
func WithSize(w, h int) WindowConfig {
	return func(c *windowConfig) {
		c.width = w
		c.height = h
	}
}

// WithTitle sets the initial title for a window
func WithTitle(t string) WindowConfig {
	return func(c *windowConfig) {
		c.title = t
	}
}

// WithView sets the initial view for a window
func WithView(v View) WindowConfig {
	return func(c *windowConfig) {
		c.view = v
	}
}

func WithType(t WindowType) WindowConfig {
	return func(c *windowConfig) {
		c.typ = t
	}
}

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

type WindowType byte

const (
	NormalWindow WindowType = iota
	DialogWindow
	SplashWindow
	MenuWindow
)
