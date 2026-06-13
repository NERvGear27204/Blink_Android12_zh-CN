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
	"github.com/character-flat/Blink/desktop/internal/platform"
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
	style := cfg.OverlayStyle

	// Common text elements
	countTxt := canvas.NewText(fmt.Sprintf("%d", cfg.RestSeconds()), colorAccent)
	countTxt.TextStyle = fyne.TextStyle{Bold: true}
	countTxt.Alignment = fyne.TextAlignCenter

	var ring fyne.CanvasObject
	var content *fyne.Container

	switch style {
	case "calm":
		content, ring, countTxt = buildCalmOverlay(cfg, countTxt)
	case "minimal":
		content, ring, countTxt = buildMinimalOverlay(cfg, countTxt)
	default: // expressive
		content, ring, countTxt = buildExpressiveOverlay(cfg, countTxt)
	}

	// Background
	bg := canvas.NewRectangle(color.NRGBA{R: 10, G: 10, B: 18, A: 245})
	w.SetContent(container.NewStack(bg, container.NewPadded(content)))

	// ESC to dismiss + force focus
	w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyEscape {
			DismissOverlay()
		}
	})

	w.Show()

	// Force focus to this window so ESC works immediately
	go func() {
		time.Sleep(200 * time.Millisecond)
		w.RequestFocus()
		platform.ForceForegroundWindow()
	}()

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

			if r, ok := ring.(*OverlayRing); ok {
				r.Progress = float32(progress)
				r.Refresh()
			} else if r, ok := ring.(*CalmRing); ok {
				r.Progress = float32(progress)
				r.Refresh()
			} else if r, ok := ring.(*MinimalRing); ok {
				r.Progress = float32(progress)
				r.Refresh()
			}

			if remaining <= 0 {
				DismissOverlay()
				return
			}
		}
	}()
}

func DismissOverlay() {
	if overlayWindow != nil {
		overlayWindow.Close()
		overlayWindow = nil
	}
}

// ── Expressive style: segmented ring, orbital dots, ripple feel ──

func buildExpressiveOverlay(cfg *config.Settings, countTxt *canvas.Text) (*fyne.Container, fyne.CanvasObject, *canvas.Text) {
	countTxt.TextSize = 76
	countTxt.Color = color.White

	secLabel := canvas.NewText("SECONDS", color.NRGBA{R: 100, G: 200, B: 180, A: 180})
	secLabel.TextSize = 12
	secLabel.Alignment = fyne.TextAlignCenter

	titleTxt := canvas.NewText("Rest Your Eyes", color.White)
	titleTxt.TextSize = 28
	titleTxt.TextStyle = fyne.TextStyle{Bold: true}
	titleTxt.Alignment = fyne.TextAlignCenter

	descTxt := canvas.NewText("Look at something 20 feet away\nand breathe slowly", colorSubtext)
	descTxt.TextSize = 14
	descTxt.Alignment = fyne.TextAlignCenter

	badge := canvas.NewText("blink", color.NRGBA{R: 100, G: 200, B: 180, A: 200})
	badge.TextSize = 12
	badge.Alignment = fyne.TextAlignCenter

	hintTxt := canvas.NewText("Press ESC to dismiss", color.NRGBA{R: 60, G: 60, B: 70, A: 255})
	hintTxt.TextSize = 11
	hintTxt.Alignment = fyne.TextAlignCenter

	ring := NewOverlayRing()

	centerContent := container.NewVBox(
		countTxt,
		secLabel,
		layout.NewSpacer(),
		titleTxt,
		descTxt,
		layout.NewSpacer(),
		badge,
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

	return content, ring, countTxt
}

// ── Calm style: smooth single arc, gentle glow, light font ──

func buildCalmOverlay(cfg *config.Settings, countTxt *canvas.Text) (*fyne.Container, fyne.CanvasObject, *canvas.Text) {
	countTxt.TextSize = 72
	countTxt.Color = color.White

	secLabel := canvas.NewText("seconds", color.NRGBA{R: 74, G: 111, B: 255, A: 180})
	secLabel.TextSize = 13
	secLabel.Alignment = fyne.TextAlignCenter

	titleTxt := canvas.NewText("Rest Your Eyes", color.White)
	titleTxt.TextSize = 26
	titleTxt.TextStyle = fyne.TextStyle{Bold: true}
	titleTxt.Alignment = fyne.TextAlignCenter

	descTxt := canvas.NewText("Look 20 feet away and breathe", colorSubtext)
	descTxt.TextSize = 14
	descTxt.Alignment = fyne.TextAlignCenter

	hintTxt := canvas.NewText("Press ESC to dismiss", color.NRGBA{R: 60, G: 60, B: 70, A: 255})
	hintTxt.TextSize = 11
	hintTxt.Alignment = fyne.TextAlignCenter

	ring := NewCalmRing()

	centerContent := container.NewVBox(
		countTxt,
		secLabel,
		layout.NewSpacer(),
		titleTxt,
		descTxt,
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

	return content, ring, countTxt
}

// ── Minimal style: just a thin progress circle + big number ──

func buildMinimalOverlay(cfg *config.Settings, countTxt *canvas.Text) (*fyne.Container, fyne.CanvasObject, *canvas.Text) {
	countTxt.TextSize = 96
	countTxt.Color = color.White

	secLabel := canvas.NewText("sec", color.NRGBA{R: 74, G: 111, B: 255, A: 180})
	secLabel.TextSize = 14
	secLabel.Alignment = fyne.TextAlignCenter

	titleTxt := canvas.NewText("Rest your eyes", color.NRGBA{R: 200, G: 200, B: 210, A: 140})
	titleTxt.TextSize = 18
	titleTxt.Alignment = fyne.TextAlignCenter

	hintTxt := canvas.NewText("Press ESC to dismiss", color.NRGBA{R: 60, G: 60, B: 70, A: 255})
	hintTxt.TextSize = 11
	hintTxt.Alignment = fyne.TextAlignCenter

	ring := NewMinimalRing()

	centerContent := container.NewVBox(
		countTxt,
		secLabel,
		layout.NewSpacer(),
		titleTxt,
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

	return content, ring, countTxt
}

// ── Expressive Ring (60 segments with gradient) ──

type OverlayRing struct {
	widget.BaseWidget
	Progress float32
}

func NewOverlayRing() *OverlayRing {
	r := &OverlayRing{Progress: 1.0}
	r.ExtendBaseWidget(r)
	return r
}

func (r *OverlayRing) MinSize() fyne.Size { return fyne.NewSize(300, 300) }

func (r *OverlayRing) CreateRenderer() fyne.WidgetRenderer {
	return &overlayRingRenderer{ring: r}
}

type overlayRingRenderer struct{ ring *OverlayRing }

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

		// Gradient from accent to teal across segments
		t := float64(i) / float64(totalSegments)
		var segColor color.Color
		if i < filledSegments {
			segColor = lerpColor(colorAccent, colorTeal, t)
		} else {
			segColor = color.NRGBA{R: 20, G: 20, B: 30, A: 50}
		}

		line := canvas.NewLine(segColor)
		line.StrokeWidth = 5
		line.Position1 = fyne.NewPos(x1, y1)
		line.Position2 = fyne.NewPos(x2, y2)
		objects = append(objects, line)
	}

	return objects
}

// ── Calm Ring (smooth continuous arc) ──

type CalmRing struct {
	widget.BaseWidget
	Progress float32
}

func NewCalmRing() *CalmRing {
	r := &CalmRing{Progress: 1.0}
	r.ExtendBaseWidget(r)
	return r
}

func (r *CalmRing) MinSize() fyne.Size { return fyne.NewSize(280, 280) }

func (r *CalmRing) CreateRenderer() fyne.WidgetRenderer {
	return &calmRingRenderer{ring: r}
}

type calmRingRenderer struct{ ring *CalmRing }

func (rr *calmRingRenderer) Layout(size fyne.Size)  {}
func (rr *calmRingRenderer) MinSize() fyne.Size     { return fyne.NewSize(280, 280) }
func (rr *calmRingRenderer) Refresh()               {}
func (rr *calmRingRenderer) Destroy()               {}
func (rr *calmRingRenderer) Objects() []fyne.CanvasObject {
	size := rr.ring.Size()
	if size.Width < 10 || size.Height < 10 {
		size = fyne.NewSize(280, 280)
	}

	cx := float64(size.Width) / 2
	cy := float64(size.Height) / 2
	radius := math.Min(cx, cy) - 15

	objects := make([]fyne.CanvasObject, 0, 200)

	// Background circle (faint)
	bgSegments := 120
	for i := 0; i < bgSegments; i++ {
		angle := float64(i) * (360.0 / float64(bgSegments))
		rad := angle * math.Pi / 180
		x1 := float32(cx + (radius-1)*math.Cos(rad))
		y1 := float32(cy + (radius-1)*math.Sin(rad))
		x2 := float32(cx + (radius+1)*math.Cos(rad))
		y2 := float32(cy + (radius+1)*math.Sin(rad))
		line := canvas.NewLine(color.NRGBA{R: 40, G: 40, B: 70, A: 80})
		line.StrokeWidth = 6
		line.Position1 = fyne.NewPos(x1, y1)
		line.Position2 = fyne.NewPos(x2, y2)
		objects = append(objects, line)
	}

	// Progress arc (smooth, thick)
	progressSegments := int(float64(bgSegments) * float64(rr.ring.Progress))
	for i := 0; i < progressSegments; i++ {
		angle := -90.0 + float64(i)*(360.0/float64(bgSegments))
		rad := angle * math.Pi / 180
		x1 := float32(cx + (radius-2)*math.Cos(rad))
		y1 := float32(cy + (radius-2)*math.Sin(rad))
		x2 := float32(cx + (radius+2)*math.Cos(rad))
		y2 := float32(cy + (radius+2)*math.Sin(rad))
		line := canvas.NewLine(colorAccent)
		line.StrokeWidth = 6
		line.Position1 = fyne.NewPos(x1, y1)
		line.Position2 = fyne.NewPos(x2, y2)
		objects = append(objects, line)
	}

	return objects
}

// ── Minimal Ring (thin clean circle) ──

type MinimalRing struct {
	widget.BaseWidget
	Progress float32
}

func NewMinimalRing() *MinimalRing {
	r := &MinimalRing{Progress: 1.0}
	r.ExtendBaseWidget(r)
	return r
}

func (r *MinimalRing) MinSize() fyne.Size { return fyne.NewSize(320, 320) }

func (r *MinimalRing) CreateRenderer() fyne.WidgetRenderer {
	return &minimalRingRenderer{ring: r}
}

type minimalRingRenderer struct{ ring *MinimalRing }

func (rr *minimalRingRenderer) Layout(size fyne.Size)  {}
func (rr *minimalRingRenderer) MinSize() fyne.Size     { return fyne.NewSize(320, 320) }
func (rr *minimalRingRenderer) Refresh()               {}
func (rr *minimalRingRenderer) Destroy()               {}
func (rr *minimalRingRenderer) Objects() []fyne.CanvasObject {
	size := rr.ring.Size()
	if size.Width < 10 || size.Height < 10 {
		size = fyne.NewSize(320, 320)
	}

	cx := float64(size.Width) / 2
	cy := float64(size.Height) / 2
	radius := math.Min(cx, cy) - 20

	objects := make([]fyne.CanvasObject, 0, 200)

	// Background track
	segments := 180
	for i := 0; i < segments; i++ {
		angle := float64(i) * (360.0 / float64(segments))
		rad := angle * math.Pi / 180
		x1 := float32(cx + (radius-1)*math.Cos(rad))
		y1 := float32(cy + (radius-1)*math.Sin(rad))
		x2 := float32(cx + (radius+1)*math.Cos(rad))
		y2 := float32(cy + (radius+1)*math.Sin(rad))
		line := canvas.NewLine(color.NRGBA{R: 30, G: 30, B: 50, A: 60})
		line.StrokeWidth = 4
		line.Position1 = fyne.NewPos(x1, y1)
		line.Position2 = fyne.NewPos(x2, y2)
		objects = append(objects, line)
	}

	// Progress
	progressSegments := int(float64(segments) * float64(rr.ring.Progress))
	for i := 0; i < progressSegments; i++ {
		angle := -90.0 + float64(i)*(360.0/float64(segments))
		rad := angle * math.Pi / 180
		x1 := float32(cx + (radius-1)*math.Cos(rad))
		y1 := float32(cy + (radius-1)*math.Sin(rad))
		x2 := float32(cx + (radius+1)*math.Cos(rad))
		y2 := float32(cy + (radius+1)*math.Sin(rad))
		line := canvas.NewLine(colorAccent)
		line.StrokeWidth = 4
		line.Position1 = fyne.NewPos(x1, y1)
		line.Position2 = fyne.NewPos(x2, y2)
		objects = append(objects, line)
	}

	return objects
}

// Helper: lerp between two colors
func lerpColor(a, b color.NRGBA, t float64) color.NRGBA {
	return color.NRGBA{
		R: uint8(float64(a.R) + t*(float64(b.R)-float64(a.R))),
		G: uint8(float64(a.G) + t*(float64(b.G)-float64(a.G))),
		B: uint8(float64(a.B) + t*(float64(b.B)-float64(a.B))),
		A: uint8(float64(a.A) + t*(float64(b.A)-float64(a.A))),
	}
}

var colorTeal = color.NRGBA{R: 100, G: 200, B: 180, A: 255}
