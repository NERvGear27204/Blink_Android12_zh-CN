//go:build ignore

package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// Generates the Blink app icon: white eye on #2196F3 blue background
// Replicates the Android ic_eye vector drawable at 256x256
func main() {
	const size = 256
	img := image.NewNRGBA(image.Rect(0, 0, size, size))

	blue := color.NRGBA{R: 33, G: 150, B: 243, A: 255}
	white := color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	// Fill rounded-rect blue background
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if inRoundedRect(x, y, size, size, 40) {
				img.SetNRGBA(x, y, blue)
			}
		}
	}

	cx, cy := float64(size)/2, float64(size)/2
	scale := float64(size) / 24.0 // Android icon is 24x24 viewport

	// Draw the eye (matching ic_eye.xml pathData from Android)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Convert to Android viewport coords (0-24)
			vx := float64(x) / scale
			vy := float64(y) / scale

			fx := float64(x) - cx
			fy := float64(y) - cy

			// Outer eye shape: almond
			if inEyeShape(vx, vy) {
				img.SetNRGBA(x, y, white)
			}

			// Outer iris ring (radius 5 in viewport units = ~53px)
			irisR := 5.0 * scale
			dist := math.Sqrt(fx*fx + fy*fy)
			if dist <= irisR && dist >= irisR-3*scale/5 {
				img.SetNRGBA(x, y, white)
			}

			// Inner iris fill (radius 3 in viewport units)
			innerR := 3.0 * scale
			if dist <= innerR {
				img.SetNRGBA(x, y, white)
			}

			// Pupil hole (radius ~1.5 viewport units)
			pupilR := 1.4 * scale
			if dist <= pupilR {
				if inRoundedRect(x, y, size, size, 40) {
					img.SetNRGBA(x, y, blue)
				}
			}
		}
	}

	f, _ := os.Create("assets/icon.png")
	defer f.Close()
	png.Encode(f, img)
}

func inEyeShape(vx, vy float64) bool {
	// The Android eye path goes from x=1 to x=23, centered at 12,12
	// It's a visibility icon: two arcs forming an almond
	cx, cy := 12.0, 12.0
	dx := vx - cx
	dy := vy - cy

	// Eye almond shape using parametric form
	// halfWidth=11 (from center to edge), halfHeight varies
	halfW := 10.5
	if math.Abs(dx) > halfW {
		return false
	}

	// Eye opening height based on x position (elliptical taper)
	t := dx / halfW
	maxH := 4.2 * math.Sqrt(1-t*t) // elliptical envelope

	if math.Abs(dy) <= maxH && math.Abs(dy) >= maxH-0.8 {
		return true
	}
	return false
}

func inRoundedRect(x, y, w, h, radius int) bool {
	// Check corners
	corners := [][2]int{
		{radius, radius},
		{w - radius - 1, radius},
		{radius, h - radius - 1},
		{w - radius - 1, h - radius - 1},
	}

	if x < radius && y < radius {
		return dist(x, y, corners[0][0], corners[0][1]) <= float64(radius)
	}
	if x >= w-radius && y < radius {
		return dist(x, y, corners[1][0], corners[1][1]) <= float64(radius)
	}
	if x < radius && y >= h-radius {
		return dist(x, y, corners[2][0], corners[2][1]) <= float64(radius)
	}
	if x >= w-radius && y >= h-radius {
		return dist(x, y, corners[3][0], corners[3][1]) <= float64(radius)
	}
	return true
}

func dist(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Sqrt(dx*dx + dy*dy)
}
