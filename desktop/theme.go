package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type blinkTheme struct{}

func (t *blinkTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 15, G: 15, B: 26, A: 255}
	case theme.ColorNameButton:
		return color.NRGBA{R: 74, G: 111, B: 255, A: 255}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 74, G: 111, B: 255, A: 255}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 224, G: 224, B: 224, A: 255}
	case theme.ColorNameInputBackground:
		return color.NRGBA{R: 30, G: 30, B: 50, A: 255}
	case theme.ColorNameSeparator:
		return color.NRGBA{R: 40, G: 40, B: 60, A: 255}
	case theme.ColorNameOverlayBackground:
		return color.NRGBA{R: 20, G: 20, B: 35, A: 240}
	}
	return theme.DefaultTheme().Color(name, theme.VariantDark)
}

func (t *blinkTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *blinkTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *blinkTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNamePadding:
		return 6
	case theme.SizeNameText:
		return 14
	}
	return theme.DefaultTheme().Size(name)
}
