package main

import (
	"log"
	"time"

	"github.com/go-ui/ui"
	_ "github.com/go-ui/ui/drivers/x11"
)

func main() {
	// Initialize a new UI instance
	gui, err := ui.New("x11")
	if err != nil {
		log.Fatalf("GUI error: %s\n", err)
	}

	// Initialize the main window and add a key event handler
	mainWindow := gui.NewWindow(ui.WithTitle("Test window"), ui.WithSize(800, 600))
	_ = mainWindow

	time.Sleep(10 * time.Second)

	gui.Release()
}
