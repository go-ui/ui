package main

import (
	"log"
	"time"

	"github.com/go-ui/ui"
	_ "github.com/go-ui/ui/drivers/dummy"
	"github.com/go-ui/ui/events"
)

func main() {
	// Initialize a new UI instance
	gui, err := ui.New("dummy")
	if err != nil {
		log.Fatalf("GUI error: %s\n", err)
	}

	// Initialize the main window and add a key event handler
	mainWindow := gui.NewWindow(ui.WithTitle("Main window"), ui.WithSize(800, 600))
	mainWindow.OnKeyDown = func(e events.KeyEvent) {
		log.Println("Key event")
	}

	time.Sleep(10 * time.Second)

	gui.Release()
}
