package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/character-flat/Blink/desktop/internal/config"
	"github.com/character-flat/Blink/desktop/internal/ui"
)

func main() {
	a := app.NewWithID("com.cflat.blink-desktop")
	a.Settings().SetTheme(&blinkTheme{})

	cfg := config.Load()
	mw := ui.NewMainWindow(a, cfg)

	// System tray
	if desk, ok := a.(desktop.App); ok {
		menu := fyne.NewMenu("Blink",
			fyne.NewMenuItem("Show", func() {
				mw.Window().Show()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Start Timer", func() {
				mw.Timer().Start()
			}),
			fyne.NewMenuItem("Stop Timer", func() {
				mw.Timer().Stop()
			}),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Quit", func() {
				mw.Timer().Stop()
				a.Quit()
			}),
		)
		desk.SetSystemTrayMenu(menu)
	}

	if cfg.StartMinimized {
		// Window stays hidden
	} else {
		mw.Show()
	}

	a.Run()
}
