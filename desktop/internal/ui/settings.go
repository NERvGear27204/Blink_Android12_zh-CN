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
	w.Resize(fyne.NewSize(420, 560))
	w.SetFixedSize(true)

	// Work duration
	workLabel := widget.NewLabel(fmt.Sprintf("Work Duration: %d min", cfg.WorkMinutes()))
	workSlider := widget.NewSlider(1, 60)
	workSlider.Value = float64(cfg.WorkMinutes())
	workSlider.Step = 1
	workSlider.OnChanged = func(v float64) {
		workLabel.SetText(fmt.Sprintf("Work Duration: %d min", int(v)))
	}

	// Rest duration
	restLabel := widget.NewLabel(fmt.Sprintf("Rest Duration: %d sec", cfg.RestSeconds()))
	restSlider := widget.NewSlider(5, 60)
	restSlider.Value = float64(cfg.RestSeconds())
	restSlider.Step = 1
	restSlider.OnChanged = func(v float64) {
		restLabel.SetText(fmt.Sprintf("Rest Duration: %d sec", int(v)))
	}

	// Alert mode
	alertModeSelect := widget.NewSelect([]string{"both", "overlay", "notification"}, nil)
	alertModeSelect.SetSelected(cfg.AlertMode)

	// Overlay style
	styleSelect := widget.NewSelect([]string{"expressive", "calm", "minimal"}, nil)
	styleSelect.SetSelected(cfg.OverlayStyle)

	// Opacity
	opacityLabel := widget.NewLabel(fmt.Sprintf("Overlay Opacity: %d%%", int(cfg.Opacity*100)))
	opacitySlider := widget.NewSlider(30, 100)
	opacitySlider.Value = cfg.Opacity * 100
	opacitySlider.Step = 5
	opacitySlider.OnChanged = func(v float64) {
		opacityLabel.SetText(fmt.Sprintf("Overlay Opacity: %d%%", int(v)))
	}

	// Options
	soundCheck := widget.NewCheck("Play sound on rest", nil)
	soundCheck.Checked = cfg.SoundEnabled

	autoStartCheck := widget.NewCheck("Start timer on app launch", nil)
	autoStartCheck.Checked = cfg.AutoStart

	autoBootCheck := widget.NewCheck("Start Blink on Windows boot", nil)
	autoBootCheck.Checked = cfg.AutoStartBoot

	minimizedCheck := widget.NewCheck("Start minimized to tray", nil)
	minimizedCheck.Checked = cfg.StartMinimized

	// Save
	saveBtn := widget.NewButton("Save", func() {
		cfg.WorkDuration = time.Duration(int(workSlider.Value)) * time.Minute
		cfg.RestDuration = time.Duration(int(restSlider.Value)) * time.Second
		cfg.AlertMode = alertModeSelect.Selected
		cfg.OverlayStyle = styleSelect.Selected
		cfg.Opacity = opacitySlider.Value / 100.0
		cfg.SoundEnabled = soundCheck.Checked
		cfg.AutoStart = autoStartCheck.Checked
		cfg.AutoStartBoot = autoBootCheck.Checked
		cfg.StartMinimized = minimizedCheck.Checked
		_ = cfg.Save()

		// Register/unregister startup
		_ = platform.SetAutoStartBoot(cfg.AutoStartBoot)

		if onSave != nil {
			onSave()
		}
		w.Close()
	})
	saveBtn.Importance = widget.HighImportance

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
		layout.NewSpacer(),
		container.NewCenter(saveBtn),
	)

	w.SetContent(container.NewPadded(content))
	w.Show()
}
