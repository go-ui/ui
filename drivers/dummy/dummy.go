// Package dummy contains the dummy driver. It implements the driver interface but does nothing. Can be used for testing and development.
package dummy

import (
	"github.com/go-ui/ui/drivers"
	"github.com/go-ui/ui/events"
)

func init() {
	drivers.Set("dummy", func() drivers.Driver { return dummy{} })
}

type dummy struct{}

func (dummy) Release() error {
	return nil
}

type dummyWindow struct{}

func (dummy) CreateWindow(_ string, _ int, _ int, _ func(events.Event)) drivers.Window {
	return dummyWindow{}
}

func (dummyWindow) Title() string            { return "" }
func (dummyWindow) SetTitle(_ string)        {}
func (dummyWindow) Size() (int, int)         { return 0, 0 }
func (dummyWindow) SetSize(_ int, _ int)     {}
func (dummyWindow) Position() (int, int)     { return 0, 0 }
func (dummyWindow) SetPosition(_ int, _ int) {}
func (dummyWindow) Close()                   {}
