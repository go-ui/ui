package x11

import (
	"image"
	"image/color"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/pkg/errors"

	"github.com/go-ui/ui/drivers"
	"github.com/go-ui/ui/events"
)

func init() {
	drivers.Set("x11", initX11)
}

// initX11
func initX11() (drivers.Driver, error) {
	// Create a new connection to X
	x, err := xgbutil.NewConn()
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to X server")
	}

	// Prepare our structure
	r := &x11{
		conn: x,
		done: make(chan struct{}),
	}

	// Start the event processing
	go xevent.Main(x)

	return r, nil
}

type x11 struct {
	conn *xgbutil.XUtil
	done chan struct{}
}

type x11Window struct {
	x  *x11
	w  *xwindow.Window
	i  *xgraphics.Image
	ev func(events.Event)
}

type x11Surface struct {
	x *x11
	i *xgraphics.Image
}

func (x *x11) CreateWindow(title string, w, h int, t drivers.WindowType, p drivers.Window, ev func(events.Event)) drivers.Window {
	// Set the right parent
	var pw xproto.Window
	if p != nil {
		if xpw, ok := p.(*x11Window); ok {
			pw = xpw.w.Id
		} else {
			pw = x.conn.RootWin()
		}
	} else {
		pw = x.conn.RootWin()
	}

	// Create the window
	win, err := xwindow.Create(x.conn, pw)
	if err != nil {
		return nil
	}

	// Set the correct size
	win.Resize(w, h)

	// Set background pixmap
	i := xgraphics.New(x.conn, image.Rect(0, 0, w, h))
	i.XSurfaceSet(win.Id)
	i.XDraw()
	i.XPaint(win.Id)

	// Show the window
	win.Map()

	// Set the title
	ewmh.WmNameSet(x.conn, win.Id, title)

	ewmh.WmWindowTypeSet(xu, win, atomNames)

	// Set the event callback
	if ev == nil {
		ev = func(_ events.Event) {}
	}

	return &x11Window{
		x:  x,
		w:  win,
		ev: ev,
	}
}

func (xw *x11Window) Title() string {
	t, _ := ewmh.WmNameGet(xw.x.conn, xw.w.Id)
	return t
}

func (xw *x11Window) SetTitle(title string) {
	ewmh.WmNameSet(xw.x.conn, xw.w.Id, title)
}

func (xw *x11Window) Size() (int, int) {
	return xw.w.Geom.Width(), xw.w.Geom.Height()
}

func (xw *x11Window) SetSize(w int, h int) {
	xw.w.Resize(w, h)
}

func (xw *x11Window) Position() (int, int) {
	return xw.w.Geom.X(), xw.w.Geom.Y()
}

func (xw *x11Window) SetPosition(x int, y int) {
	xw.w.Move(x, y)
}

func (xw *x11Window) Render(s drivers.Surface, x, y int) {
	xs, ok := s.(*x11Surface)
	if !ok {
		// Not a X11 surface!!
		return
	}

	xs.i.CreatePixmap()
	xs.i.XExpPaint(xw.w.Id, x, y)
}

func (xw *x11Window) Close() {
	xw.w.Destroy()
}

func (x *x11) CreateSurface(w, h int) drivers.Surface {
	// Create the image
	i := xgraphics.New(x.conn, image.Rect(0, 0, w, h))
	i.CreatePixmap()

	return &x11Surface{
		x: x,
		i: i,
	}
}

func (xs *x11Surface) Set(x, y int, c color.Color) {
	xs.i.Set(x, y, c)
}

func (xs *x11Surface) ColorModel() color.Model {
	return xs.i.ColorModel()
}

func (xs *x11Surface) Bounds() image.Rectangle {
	return xs.i.Bounds()
}

func (xs *x11Surface) At(x, y int) color.Color {
	return xs.i.At(x, y)
}

func (x *x11) Release() error {
	xevent.Quit(x.conn)
	return nil
}
