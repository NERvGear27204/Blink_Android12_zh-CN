# Blink Desktop

**The 20-20-20 rule eye care timer — for Windows.**

A lightweight, single-binary desktop companion to the [Blink Android app](../README.md). Built in Go with the Fyne GUI toolkit.

---

## ✨ Features

| Feature | Description |
|---|---|
| **20-20-20 Rule** | Every 20 minutes, look 20 feet away for 20 seconds |
| **Fullscreen Overlay** | Beautiful animated rest screen with countdown ring |
| **System Tray** | Runs silently in the background, always accessible |
| **Segmented Timer Ring** | 60-segment animated progress ring (matches Android) |
| **Configurable Timers** | Work: 1–60 min, Rest: 5–60 sec |
| **3 Overlay Styles** | Expressive, Calm, Minimal |
| **Dark Theme** | Custom dark theme easy on the eyes |
| **Single Binary** | No installer, no dependencies — just run |
| **~24 MB** | Stripped release build |

---

## 📦 Download

**[→ Latest Release](https://github.com/character-flat/Blink/releases/latest)**

Download `blink-desktop-windows-amd64.exe`, place it anywhere, and run.

---

## 🔨 Build from Source

### Requirements
- Go 1.21+
- GCC (MinGW-w64 on Windows) — required for Fyne's OpenGL bindings
- On Windows: [w64devkit](https://github.com/skeeto/w64devkit) or TDM-GCC works great

### Build

```bash
cd desktop

# Debug build (with console)
go build -o blink.exe .

# Release build (stripped, no console window)
go build -ldflags="-s -w -H windowsgui" -o blink.exe .
```

---

## 🏗️ Architecture

```
desktop/
├── main.go              # Entry point, tray setup
├── theme.go             # Custom dark theme
└── internal/
    ├── config/          # JSON settings persistence (%APPDATA%/Blink/)
    ├── timer/           # Core timer logic (goroutine-based)
    └── ui/
        ├── window.go    # Main window + segmented timer ring widget
        ├── overlay.go   # Fullscreen rest overlay with animated ring
        └── settings.go  # Settings dialog
```

### Key Design Decisions

- **Fyne** — Cross-platform Go GUI with OpenGL rendering. Single binary, no runtime deps.
- **Goroutine timer** — Precise 50ms tick loop using `time.NewTicker`, no drift.
- **Custom widgets** — `TimerRing` and `OverlayRing` draw segmented arcs via Fyne canvas lines.
- **System tray** — App minimizes to tray on close; keeps running until explicitly quit.
- **Thread-safe state** — Timer uses `sync.Mutex` for safe cross-goroutine access.

---

## ⚙️ Settings

Settings are stored in `%APPDATA%/Blink/settings.json`:

```json
{
  "work_duration_ns": 1200000000000,
  "rest_duration_ns": 20000000000,
  "overlay_style": "expressive",
  "opacity": 0.9,
  "sound_enabled": true,
  "auto_start": false,
  "start_minimized": false
}
```

---

## 📄 License

GNU General Public License v3 — same as the Android app.
