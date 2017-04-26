package x11

import (
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/pkg/errors"

	"github.com/go-ui/ui/drivers"
	"github.com/go-ui/ui/events"
)

func init() {
	drivers.Set("x11", initX11)
}

// initX11
func initX11() (drivers.Driver, error) {
	x, err := xgb.NewConn()
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to X server")
	}

	r := &x11{
		conn: x,
		done: make(chan struct{}),
	}

	go r.processEvents()

	return r, nil
}

type x11 struct {
	conn *xgb.Conn
	done chan struct{}
}

func (x *x11) processEvents() {
	defer func() {
		close(x.done)
	}()

	for {
		ev, err := x.conn.WaitForEvent()
		if ev == nil && err == nil {
			return
		}
		if ev != nil {
			log.Printf("Event: %+v\n", ev)
		}
		if err != nil {
			log.Printf("Event error: %s\n", err)
		}
	}
}

type x11Window struct {
	x11      *x11
	windowID xproto.Window
}

func (x *x11) CreateWindow(title string, w, h int, ev func(events.Event)) drivers.Window {
	// Get a window ID
	windowID, err := xproto.NewWindowId(x.conn)
	if err != nil {
		return nil
	}

	// Select the default screen
	screen := xproto.Setup(x.conn).DefaultScreen(x.conn)

	// Create the window
	xproto.CreateWindow(
		x.conn,           // X11 connection
		screen.RootDepth, // Color depth
		windowID,         // Window ID
		screen.Root,      // Window parent
		0,                // X
		0,                // Y
		uint16(w),        // Width
		uint16(h),        // Height
		0,                // Border width
		xproto.WindowClassInputOutput, // Input class
		screen.RootVisual,             // Visual
		xproto.CwBackPixel|xproto.CwEventMask,
		[]uint32{
			0xffffffff,
			xproto.EventMaskStructureNotify |
				xproto.EventMaskKeyPress |
				xproto.EventMaskKeyRelease |
				xproto.EventMaskPointerMotion,
		},
	)

	// Display the window
	xproto.MapWindow(x.conn, windowID)

	// TODO Add to event functions list

	return &x11Window{
		x11:      x,
		windowID: windowID,
	}
}

// TODO
func (xw *x11Window) Title() string {
	return ""
}

// TODO
func (xw *x11Window) SetTitle(title string) {

}

// TODO
func (xw *x11Window) Size() (int, int) {
	return 0, 0
}

// TODO
func (xw *x11Window) SetSize(w int, h int) {

}

// TODO
func (xw *x11Window) Position() (int, int) {
	return 0, 0
}

// TODO
func (xw *x11Window) SetPosition(x int, y int) {

}

// TODO
func (xw *x11Window) Close() {
	// TODO remove from event functions list
	xproto.DestroyWindow(xw.x11.conn, xw.windowID)
}

func (x *x11) CreateSurface(w, h int) drivers.Surface {
	return nil
}

func (x *x11) Release() error {
	return nil
}
