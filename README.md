<p align="center">
  <img src="https://img.shields.io/badge/Android-14%2B-brightgreen?logo=android&logoColor=white" />
  <img src="https://img.shields.io/badge/Windows-10%2B-0078D4?logo=windows&logoColor=white" />
  <img src="https://img.shields.io/badge/Kotlin-2.1-7F52FF?logo=kotlin&logoColor=white" />
  <img src="https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white" />
  <img src="https://img.shields.io/github/v/release/character-flat/Blink" />
</p>

# Blink

**A reliable 20-20-20 eye-care timer for Android and Windows.**

Blink reminds you to rest your eyes every 20 minutes: look at something 20 feet away for 20 seconds. The Android app is designed to survive Android battery restrictions with exact alarms and a foreground service. The Windows desktop app runs quietly from the system tray and presents the same timer experience in a single executable.

## Platform isolation

The two applications are intentionally independent:

- **Android:** [app/](app) is the only Gradle module. Its Kotlin, Compose, manifest, permissions, and APK builds are unaffected by desktop work.
- **Windows desktop:** [desktop/](desktop) is a standalone Go module built with Fyne. Its settings, timer, system-tray integration, and executable do not depend on the Android app.

They share the Blink name and the 20-20-20 concept, but have no shared runtime, configuration, build output, or platform-specific code.

---

## Features

### Android

| Area | Feature | Description |
|---|---|---|
| Reliability | Exact alarms | Uses `AlarmManager.setExactAndAllowWhileIdle` to fire during deep sleep and Doze. |
| Reliability | Foreground service | Persistent `specialUse` foreground service with a live countdown notification. |
| Overlay | True overlay | `TYPE_APPLICATION_OVERLAY` appears above other apps without switching context. |
| Overlay | Three styles | Expressive, Calm, and Minimal rest overlays. |
| Overlay | Overlay controls | Adjustable opacity, optional dismiss button, and responsive landscape support. |
| Customization | Alert modes | Overlay only, notification only, or both. |
| Customization | Adjustable timers | Work duration from 1–60 minutes and rest duration from 10–60 seconds. |
| Customization | Dynamic Material 3 color | Uses wallpaper-derived colors through `dynamicColorScheme()`. |
| Power | Battery controls | Optional Shizuku integration, unrestricted-battery guidance, `START_STICKY`, and a watchdog alarm. |
| Lifecycle | Device handling | Smart screen-lock detection, boot receiver, direct boot, and optional startup after reboot. |
| Power | OEM guidance | Battery-setting instructions for Xiaomi, Samsung, OPPO, Vivo, and other Android variants. |

### Windows desktop

| Area | Feature | Description |
|---|---|---|
| Timer | 20-20-20 timer | Configurable work duration from 1–60 minutes and rest duration from 5–60 seconds. |
| Overlay | Fullscreen rest screen | Fullscreen countdown overlay with Expressive, Calm, and Minimal styles. |
| System integration | System tray | Runs in the notification area, with controls to show the window, start, stop, and quit. |
| System integration | Single instance | Prevents duplicate timer instances. |
| System integration | Windows startup | Optional launch at sign-in and start-minimized support. |
| Notifications | Alert modes | Overlay, native notification, or both, with optional system sound. |
| Presentation | Timer ring | Segmented progress ring and dark theme designed for comfortable viewing. |
| Distribution | Single executable | No installer or separate runtime required for the released Windows binary. |

---

## Download

[Latest release](https://github.com/character-flat/Blink/releases/latest)

### Android

| Build | Description |
|---|---|
| `blink-v1.1.0-release.apk` | Signed, R8-minified release build, approximately 1.8 MB. |
| `blink-v1.1.0-debug.apk` | Debug build with logging, approximately 10 MB. |

### Windows desktop

| Build | Description |
|---|---|
| `blink-desktop-windows-amd64.exe` | Windows 10+ 64-bit single executable, approximately 24 MB. |

See [desktop/README.md](desktop/README.md) for desktop-specific details and [CHANGELOG.md](CHANGELOG.md) for release notes.

---

## Requirements

### Android

- Android 14+ (API 34)
- Display-over-other-apps permission for the overlay
- Notification permission
- Exact alarm permission

### Windows desktop

- Windows 10 or later, 64-bit

---

## Architecture

### Android

```
MainActivity (Jetpack Compose / Material 3)
  - Timer ring, start/stop controls, and settings

EyeCareService (Foreground Service)
  - AlarmManager, screen-state receiver, countdown, and notification

RestOverlayService (TYPE_APPLICATION_OVERLAY)
  - ComposeView, three overlay styles, raw progress animation, and vibration

Supporting components
  - ShizukuUtils, BootReceiver, AlarmReceiver, and PrefsManager
```

### Windows desktop

```
main.go
  - App lifecycle, command-line startup mode, and system tray

internal/timer
  - Goroutine-based timer with cancellation-safe lifecycle transitions

internal/ui
  - Main window, settings window, fullscreen rest overlay, and custom timer rings

internal/config and internal/platform
  - Per-user JSON settings and Windows-specific startup, sound, opacity, and focus support
```

---

## Key design decisions

### Android

- **Exact scheduling:** `AlarmManager.setExactAndAllowWhileIdle` is used instead of WorkManager because the timer needs predictable 20-minute intervals.
- **Overlay design:** `TYPE_APPLICATION_OVERLAY` preserves the current app context; launching an activity would force an app switch.
- **Service-hosted Compose:** The overlay service explicitly provides `LifecycleOwner`, `SavedStateRegistryOwner`, and `ViewModelStoreOwner` for its `ComposeView`.
- **Smooth countdown:** A direct 16 ms update loop avoids jitter from animating state that is already updated continuously.
- **Small releases:** R8 removes unused Material icons from the release APK.

### Windows desktop

- **Standalone module:** The Go module in [desktop/](desktop) is intentionally outside the Android Gradle project.
- **Cancellation-safe timer:** A run identifier prevents a stopped timer worker from entering a later rest cycle.
- **UI-thread updates:** Background timer and overlay work is marshalled to Fyne's main runtime before changing UI state.
- **Native integration:** Windows-specific functionality is limited to the desktop platform package for startup, sound, focus, opacity, and the single-instance guard.

---

## Build from source

### Android

```bash
# Debug APK
./gradlew assembleDebug

# Release APK (requires signing key)
./gradlew assembleRelease
```

### Windows desktop

Requirements: Go 1.21+ and a Windows C compiler supported by Fyne, such as MinGW-w64.

```bash
cd desktop

# Debug build with a console window
go build -o blink.exe .

# Release build without a console window
go build -ldflags="-s -w -H windowsgui" -o blink.exe .
```

### Technology

| Platform | Stack |
|---|---|
| Android | Kotlin 2.1, Jetpack Compose, Material 3, Gradle 8.9, Android Gradle Plugin 8.7.3, Shizuku, and Coroutines |
| Windows desktop | Go, Fyne, Windows APIs, and MinGW-w64 for source builds |

---

## Disclaimer

Blink is a general wellness reminder, not medical advice or a substitute for professional eye-care guidance. It does not diagnose, treat, cure, or prevent any condition. If you have eye discomfort, vision changes, or other health concerns, consult a qualified healthcare professional.

Android and desktop reliability can still be affected by operating-system settings, permissions, battery policies, notification settings, and device or desktop environment behavior. Review the applicable permissions and platform requirements before relying on Blink for scheduled breaks.

---

## License

GNU General Public License v3.
