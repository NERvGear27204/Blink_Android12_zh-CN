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
	"fyne.io/fyne/v2/widget"

	"github.com/character-flat/Blink/desktop/internal/config"
)

var overlayWindow fyne.Window

func ShowOverlay(app fyne.App, cfg *config.Settings) {
	if overlayWindow != nil {
		return
	}

	w := app.NewWindow("Blink - Rest")
	w.SetFullScreen(true)
	w.SetPadded(false)
	overlayWindow = w

	startTime := time.Now()
	restDuration := cfg.RestDuration

	// UI elements
	eyeTxt := canvas.NewText("👁", color.White)
	eyeTxt.TextSize = 56
	eyeTxt.Alignment = fyne.TextAlignCenter

	titleTxt := canvas.NewText("Look Away", color.White)
	titleTxt.TextSize = 32
	titleTxt.TextStyle = fyne.TextStyle{Bold: true}
	titleTxt.Alignment = fyne.TextAlignCenter

	descTxt := canvas.NewText("20 feet away for 20 seconds", colorSubtext)
	descTxt.TextSize = 14
	descTxt.Alignment = fyne.TextAlignCenter

	countTxt := canvas.NewText(fmt.Sprintf("%d", cfg.RestSeconds()), colorAccent)
	countTxt.TextSize = 48
	countTxt.TextStyle = fyne.TextStyle{Bold: true}
	countTxt.Alignment = fyne.TextAlignCenter

	hintTxt := canvas.NewText("Press ESC to dismiss", color.NRGBA{R: 80, G: 80, B: 80, A: 255})
	hintTxt.TextSize = 11
	hintTxt.Alignment = fyne.TextAlignCenter

	// Progress ring
	ring := NewOverlayRing()

	centerContent := container.NewVBox(
		eyeTxt,
		titleTxt,
		descTxt,
		layout.NewSpacer(),
		countTxt,
	)

	ringStack := container.NewStack(
		ring,
		container.NewCenter(centerContent),
	)

	content := container.NewVBox(
		layout.NewSpacer(),
		container.NewCenter(ringStack),
		layout.NewSpacer(),
		container.NewCenter(hintTxt),
		widget.NewSeparator(),
	)

	bg := canvas.NewRectangle(color.NRGBA{R: 10, G: 10, B: 18, A: 240})
	w.SetContent(container.NewStack(bg, container.NewPadded(content)))

	// ESC to dismiss
	w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyEscape {
			DismissOverlay()
		}
	})

	// Animation ticker
	go func() {
		ticker := time.NewTicker(50 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			if overlayWindow == nil {
				return
			}
			elapsed := time.Since(startTime)
			remaining := restDuration - elapsed
			if remaining < 0 {
				remaining = 0
			}
			sec := int(remaining.Seconds()) + 1
			progress := float64(remaining) / float64(restDuration)

			countTxt.Text = fmt.Sprintf("%d", sec)
			countTxt.Refresh()
			ring.Progress = float32(progress)
			ring.Refresh()

			if remaining <= 0 {
				DismissOverlay()
				return
			}
		}
	}()

	w.Show()
}

func DismissOverlay() {
	if overlayWindow != nil {
		overlayWindow.Close()
		overlayWindow = nil
	}
}

// --- Overlay Ring ---

type OverlayRing struct {
	widget.BaseWidget
	Progress float32
}

func NewOverlayRing() *OverlayRing {
	r := &OverlayRing{Progress: 1.0}
	r.ExtendBaseWidget(r)
	return r
}

func (r *OverlayRing) MinSize() fyne.Size {
	return fyne.NewSize(300, 300)
}

func (r *OverlayRing) CreateRenderer() fyne.WidgetRenderer {
	return &overlayRingRenderer{ring: r}
}

type overlayRingRenderer struct {
	ring *OverlayRing
}

func (rr *overlayRingRenderer) Layout(size fyne.Size)  {}
func (rr *overlayRingRenderer) MinSize() fyne.Size     { return fyne.NewSize(300, 300) }
func (rr *overlayRingRenderer) Refresh()               {}
func (rr *overlayRingRenderer) Destroy()               {}
func (rr *overlayRingRenderer) Objects() []fyne.CanvasObject {
	size := rr.ring.Size()
	if size.Width < 10 || size.Height < 10 {
		size = fyne.NewSize(300, 300)
	}

	cx := float64(size.Width) / 2
	cy := float64(size.Height) / 2
	radius := math.Min(cx, cy) - 20

	objects := make([]fyne.CanvasObject, 0, 64)
	totalSegments := 60
	gapAngle := 2.0
	segmentAngle := (360.0 / float64(totalSegments)) - gapAngle
	filledSegments := int(float64(totalSegments) * float64(rr.ring.Progress))

	for i := 0; i < totalSegments; i++ {
		angle := -90.0 + float64(i)*(segmentAngle+gapAngle) + segmentAngle/2
		rad := angle * math.Pi / 180

		innerR := radius - 6
		outerR := radius + 4

		x1 := float32(cx + innerR*math.Cos(rad))
		y1 := float32(cy + innerR*math.Sin(rad))
		x2 := float32(cx + outerR*math.Cos(rad))
		y2 := float32(cy + outerR*math.Sin(rad))

		var segColor color.Color
		if i < filledSegments {
			segColor = colorAccent
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
