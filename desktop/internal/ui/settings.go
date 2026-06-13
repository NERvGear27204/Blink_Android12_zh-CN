package ui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/character-flat/Blink/desktop/internal/config"
	"github.com/character-flat/Blink/desktop/internal/platform"
)

func ShowSettings(app fyne.App, cfg *config.Settings, onSave func()) {
	w := app.NewWindow("Blink - Settings")
	w.Resize(fyne.NewSize(420, 520))
	w.SetFixedSize(true)

	// Auto-save helper — saves config and notifies caller on every change
	save := func() {
		_ = cfg.Save()
		_ = platform.SetAutoStartBoot(cfg.AutoStartBoot)
		if onSave != nil {
			onSave()
		}
	}

	// Work duration
	workLabel := widget.NewLabel(fmt.Sprintf("Work Duration: %d min", cfg.WorkMinutes()))
	workSlider := widget.NewSlider(1, 60)
	workSlider.Value = float64(cfg.WorkMinutes())
	workSlider.Step = 1
	workSlider.OnChanged = func(v float64) {
		workLabel.SetText(fmt.Sprintf("Work Duration: %d min", int(v)))
		cfg.WorkDuration = time.Duration(int(v)) * time.Minute
		save()
	}

	// Rest duration
	restLabel := widget.NewLabel(fmt.Sprintf("Rest Duration: %d sec", cfg.RestSeconds()))
	restSlider := widget.NewSlider(5, 60)
	restSlider.Value = float64(cfg.RestSeconds())
	restSlider.Step = 1
	restSlider.OnChanged = func(v float64) {
		restLabel.SetText(fmt.Sprintf("Rest Duration: %d sec", int(v)))
		cfg.RestDuration = time.Duration(int(v)) * time.Second
		save()
	}

	// Alert mode
	alertModeSelect := widget.NewSelect([]string{"both", "overlay", "notification"}, func(v string) {
		cfg.AlertMode = v
		save()
	})
	alertModeSelect.SetSelected(cfg.AlertMode)

	// Overlay style
	styleSelect := widget.NewSelect([]string{"expressive", "calm", "minimal"}, func(v string) {
		cfg.OverlayStyle = v
		save()
	})
	styleSelect.SetSelected(cfg.OverlayStyle)

	// Opacity
	opacityLabel := widget.NewLabel(fmt.Sprintf("Overlay Opacity: %d%%", int(cfg.Opacity*100)))
	opacitySlider := widget.NewSlider(30, 100)
	opacitySlider.Value = cfg.Opacity * 100
	opacitySlider.Step = 5
	opacitySlider.OnChanged = func(v float64) {
		opacityLabel.SetText(fmt.Sprintf("Overlay Opacity: %d%%", int(v)))
		cfg.Opacity = v / 100.0
		save()
	}

	// Options
	soundCheck := widget.NewCheck("Play sound on rest", func(v bool) {
		cfg.SoundEnabled = v
		save()
	})
	soundCheck.Checked = cfg.SoundEnabled

	autoStartCheck := widget.NewCheck("Start timer on app launch", func(v bool) {
		cfg.AutoStart = v
		save()
	})
	autoStartCheck.Checked = cfg.AutoStart

	autoBootCheck := widget.NewCheck("Start Blink on Windows boot", func(v bool) {
		cfg.AutoStartBoot = v
		save()
	})
	autoBootCheck.Checked = cfg.AutoStartBoot

	minimizedCheck := widget.NewCheck("Start minimized to tray", func(v bool) {
		cfg.StartMinimized = v
		save()
	})
	minimizedCheck.Checked = cfg.StartMinimized

	content := container.NewVBox(
		widget.NewLabel("⏱  Timers"),
		widget.NewSeparator(),
		workLabel,
		workSlider,
		restLabel,
		restSlider,
		layout.NewSpacer(),
		widget.NewLabel("🔔  Alert Mode"),
		widget.NewSeparator(),
		alertModeSelect,
		layout.NewSpacer(),
		widget.NewLabel("🎨  Overlay"),
		widget.NewSeparator(),
		widget.NewLabel("Style"),
		styleSelect,
		opacityLabel,
		opacitySlider,
		layout.NewSpacer(),
		widget.NewLabel("⚡  Options"),
		widget.NewSeparator(),
		soundCheck,
		autoStartCheck,
		autoBootCheck,
		minimizedCheck,
	)

	scrollable := container.NewVScroll(container.NewPadded(content))
	w.SetContent(scrollable)
	w.Show()
}
