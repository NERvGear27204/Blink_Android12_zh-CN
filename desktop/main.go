package main

import (
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"

	"github.com/character-flat/Blink/desktop/assets"
	"github.com/character-flat/Blink/desktop/internal/config"
	"github.com/character-flat/Blink/desktop/internal/platform"
	"github.com/character-flat/Blink/desktop/internal/ui"
)

func main() {
	// Prevent multiple instances
	if !platform.AcquireSingleInstance() {
		// Another instance is already running — exit silently
		return
	}
	defer platform.ReleaseSingleInstance()

	a := app.NewWithID("com.cflat.blink-desktop")
	a.Settings().SetTheme(&blinkTheme{})

	// Set app icon (matches Android launcher icon)
	iconRes := fyne.NewStaticResource("icon.png", assets.IconPNG)
	a.SetIcon(iconRes)

	cfg := config.Load()

	// Check if launched with --minimized (e.g. from Windows startup)
	startMinimized := cfg.StartMinimized
	for _, arg := range os.Args[1:] {
		if strings.TrimLeft(arg, "-") == "minimized" {
			startMinimized = true
			break
		}
	}

	mw := ui.NewMainWindow(a, cfg)

	// System tray
	if desk, ok := a.(desktop.App); ok {
		desk.SetSystemTrayIcon(iconRes)
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

	if startMinimized {
		// Window stays hidden, runs in tray
	} else {
		mw.Show()
	}

	a.Run()
}
