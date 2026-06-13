package ui

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/character-flat/Blink/desktop/internal/config"
	"github.com/character-flat/Blink/desktop/internal/timer"
)

var (
	colorBg       = color.NRGBA{R: 15, G: 15, B: 26, A: 255}
	colorAccent   = color.NRGBA{R: 74, G: 111, B: 255, A: 255}
	colorWarning  = color.NRGBA{R: 255, G: 170, B: 68, A: 255}
	colorRest     = color.NRGBA{R: 255, G: 85, B: 85, A: 255}
	colorInactive = color.NRGBA{R: 30, G: 30, B: 45, A: 255}
	colorText     = color.NRGBA{R: 224, G: 224, B: 224, A: 255}
	colorSubtext  = color.NRGBA{R: 136, G: 136, B: 153, A: 255}
)

type MainWindow struct {
	window fyne.Window
	cfg    *config.Settings
	tmr    *timer.Timer
	app    fyne.App

	// UI elements
	ring      *TimerRing
	timeTxt   *canvas.Text
	secTxt    *canvas.Text
	subtitle  *canvas.Text
	statusTxt *canvas.Text
	startBtn  *widget.Button
}

func NewMainWindow(app fyne.App, cfg *config.Settings) *MainWindow {
	w := app.NewWindow("Blink")
	w.Resize(fyne.NewSize(400, 560))
	w.SetFixedSize(true)

	mw := &MainWindow{
		window: w,
		cfg:    cfg,
		app:    app,
	}

	mw.tmr = timer.New(cfg.WorkDuration, cfg.RestDuration, timer.Callbacks{
		OnTick:        mw.onTick,
		OnRestStart:   mw.onRestStart,
		OnRestEnd:     mw.onRestEnd,
		OnStateChange: mw.onStateChange,
	})

	mw.buildUI()

	w.SetCloseIntercept(func() {
		w.Hide()
	})

	if cfg.AutoStart {
		go func() {
			time.Sleep(300 * time.Millisecond)
			mw.tmr.Start()
		}()
	}

	return mw
}

func (mw *MainWindow) Window() fyne.Window {
	return mw.window
}

func (mw *MainWindow) Timer() *timer.Timer {
	return mw.tmr
}

func (mw *MainWindow) Show() {
	mw.window.Show()
}

func (mw *MainWindow) buildUI() {
	// Title
	title := canvas.NewText("Blink", colorText)
	title.TextSize = 32
	title.TextStyle = fyne.TextStyle{Bold: true}
	title.Alignment = fyne.TextAlignCenter

	// Subtitle
	mw.subtitle = canvas.NewText("Start to protect your eyes", colorSubtext)
	mw.subtitle.TextSize = 13
	mw.subtitle.Alignment = fyne.TextAlignCenter

	// Timer ring
	mw.ring = NewTimerRing()
	mw.ring.Progress = 1.0
	mw.ring.State = timer.Idle

	// Time display
	mw.timeTxt = canvas.NewText(fmt.Sprintf("%02d", mw.cfg.WorkMinutes()), colorText)
	mw.timeTxt.TextSize = 60
	mw.timeTxt.TextStyle = fyne.TextStyle{Bold: true}
	mw.timeTxt.Alignment = fyne.TextAlignCenter

	mw.secTxt = canvas.NewText("00", colorSubtext)
	mw.secTxt.TextSize = 24
	mw.secTxt.Alignment = fyne.TextAlignCenter

	mw.statusTxt = canvas.NewText("", colorSubtext)
	mw.statusTxt.TextSize = 10
	mw.statusTxt.Alignment = fyne.TextAlignCenter

	// Buttons
	mw.startBtn = widget.NewButtonWithIcon("Start", theme.MediaPlayIcon(), mw.toggleTimer)
	mw.startBtn.Importance = widget.HighImportance

	settingsBtn := widget.NewButtonWithIcon("", theme.SettingsIcon(), mw.openSettings)

	buttons := container.NewHBox(
		layout.NewSpacer(),
		mw.startBtn,
		settingsBtn,
		layout.NewSpacer(),
	)

	// Stack time text over ring
	timeStack := container.NewStack(
		mw.ring,
		container.NewCenter(
			container.NewVBox(
				mw.timeTxt,
				mw.secTxt,
				mw.statusTxt,
			),
		),
	)

	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(title),
		container.NewCenter(mw.subtitle),
		widget.NewSeparator(),
		layout.NewSpacer(),
		container.NewCenter(timeStack),
		layout.NewSpacer(),
		buttons,
		layout.NewSpacer(),
	)

	mw.window.SetContent(content)
}

func (mw *MainWindow) toggleTimer() {
	if mw.tmr.State() != timer.Idle {
		mw.tmr.Stop()
	} else {
		mw.tmr.Configure(mw.cfg.WorkDuration, mw.cfg.RestDuration)
		mw.tmr.Start()
	}
}

func (mw *MainWindow) onTick(remaining time.Duration, state timer.State) {
	if state == timer.Resting {
		sec := int(remaining.Seconds()) + 1
		mw.timeTxt.Text = fmt.Sprintf("%d", sec)
		mw.timeTxt.Color = colorRest
		mw.secTxt.Text = "sec"
		mw.statusTxt.Text = "LOOK AWAY"
		mw.statusTxt.Color = colorRest
		progress := float64(remaining) / float64(mw.cfg.RestDuration)
		mw.ring.Progress = float32(progress)
		mw.ring.State = timer.Resting
	} else {
		totalSec := int(remaining.Seconds())
		min := totalSec / 60
		sec := totalSec % 60
		mw.timeTxt.Text = fmt.Sprintf("%02d", min)
		mw.timeTxt.Color = colorText
		mw.secTxt.Text = fmt.Sprintf("%02d", sec)
		mw.statusTxt.Text = ""
		progress := float64(remaining) / float64(mw.cfg.WorkDuration)
		mw.ring.Progress = float32(progress)
		mw.ring.State = timer.Running
	}
	mw.ring.Refresh()
	mw.timeTxt.Refresh()
	mw.secTxt.Refresh()
	mw.statusTxt.Refresh()
}

func (mw *MainWindow) onStateChange(s timer.State) {
	switch s {
	case timer.Running:
		mw.startBtn.SetText("Stop")
		mw.startBtn.SetIcon(theme.MediaStopIcon())
		mw.subtitle.Text = "Protecting your vision"
		mw.subtitle.Refresh()
	case timer.Resting:
		mw.subtitle.Text = "Time to rest your eyes"
		mw.subtitle.Refresh()
	case timer.Idle:
		mw.startBtn.SetText("Start")
		mw.startBtn.SetIcon(theme.MediaPlayIcon())
		mw.subtitle.Text = "Start to protect your eyes"
		mw.subtitle.Refresh()
		mw.statusTxt.Text = "PAUSED"
		mw.statusTxt.Color = colorSubtext
		mw.statusTxt.Refresh()
		mw.timeTxt.Text = fmt.Sprintf("%02d", mw.cfg.WorkMinutes())
		mw.timeTxt.Color = colorText
		mw.timeTxt.Refresh()
		mw.secTxt.Text = "00"
		mw.secTxt.Refresh()
		mw.ring.Progress = 1.0
		mw.ring.State = timer.Idle
		mw.ring.Refresh()
	}
}

func (mw *MainWindow) onRestStart() {
	ShowOverlay(mw.app, mw.cfg)
}

func (mw *MainWindow) onRestEnd() {
	DismissOverlay()
}

func (mw *MainWindow) openSettings() {
	ShowSettings(mw.app, mw.cfg, func() {
		mw.tmr.Configure(mw.cfg.WorkDuration, mw.cfg.RestDuration)
		if mw.tmr.State() == timer.Idle {
			mw.timeTxt.Text = fmt.Sprintf("%02d", mw.cfg.WorkMinutes())
			mw.timeTxt.Refresh()
		}
	})
}

// --- Timer Ring Widget ---

type TimerRing struct {
	widget.BaseWidget
	Progress float32 // 0.0 to 1.0
	State    timer.State
}

func NewTimerRing() *TimerRing {
	r := &TimerRing{Progress: 1.0, State: timer.Idle}
	r.ExtendBaseWidget(r)
	return r
}

func (r *TimerRing) MinSize() fyne.Size {
	return fyne.NewSize(260, 260)
}

func (r *TimerRing) CreateRenderer() fyne.WidgetRenderer {
	return &timerRingRenderer{ring: r}
}

type timerRingRenderer struct {
	ring *TimerRing
}

func (rr *timerRingRenderer) Layout(size fyne.Size)  {}
func (rr *timerRingRenderer) MinSize() fyne.Size     { return fyne.NewSize(260, 260) }
func (rr *timerRingRenderer) Refresh()               {}
func (rr *timerRingRenderer) Destroy()               {}
func (rr *timerRingRenderer) Objects() []fyne.CanvasObject {
	size := rr.ring.Size()
	if size.Width < 10 || size.Height < 10 {
		size = fyne.NewSize(260, 260)
	}

	cx := float64(size.Width) / 2
	cy := float64(size.Height) / 2
	radius := math.Min(cx, cy) - 15

	objects := make([]fyne.CanvasObject, 0, 64)

	totalSegments := 60
	gapAngle := 2.0
	segmentAngle := (360.0 / float64(totalSegments)) - gapAngle
	filledSegments := int(float64(totalSegments) * float64(rr.ring.Progress))

	var activeColor color.Color
	switch rr.ring.State {
	case timer.Resting:
		activeColor = colorRest
	case timer.Running:
		if rr.ring.Progress < 0.15 {
			activeColor = colorWarning
		} else {
			activeColor = colorAccent
		}
	default:
		activeColor = color.NRGBA{R: 50, G: 50, B: 80, A: 255}
	}

	for i := 0; i < totalSegments; i++ {
		angle := -90.0 + float64(i)*(segmentAngle+gapAngle) + segmentAngle/2
		rad := angle * math.Pi / 180

		innerR := radius - 8
		outerR := radius + 2

		x1 := float32(cx + innerR*math.Cos(rad))
		y1 := float32(cy + innerR*math.Sin(rad))
		x2 := float32(cx + outerR*math.Cos(rad))
		y2 := float32(cy + outerR*math.Sin(rad))

		var segColor color.Color
		if i < filledSegments {
			segColor = activeColor
		} else {
			segColor = colorInactive
		}

		line := canvas.NewLine(segColor)
		line.StrokeWidth = 4
		line.Position1 = fyne.NewPos(x1, y1)
		line.Position2 = fyne.NewPos(x2, y2)
		objects = append(objects, line)
	}

	return objects
}
