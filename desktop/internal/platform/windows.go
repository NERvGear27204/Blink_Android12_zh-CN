package platform

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

const startupRegKey = `Software\Microsoft\Windows\CurrentVersion\Run`
const appName = "Blink"

// SetAutoStartBoot registers/unregisters the app in Windows startup.
func SetAutoStartBoot(enabled bool) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, startupRegKey, registry.SET_VALUE|registry.QUERY_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()

	if enabled {
		exePath, err := os.Executable()
		if err != nil {
			return err
		}
		exePath, _ = filepath.Abs(exePath)
		return key.SetStringValue(appName, `"`+exePath+`" --minimized`)
	}
	_ = key.DeleteValue(appName)
	return nil
}

// IsAutoStartBoot checks if app is registered in startup.
func IsAutoStartBoot() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, startupRegKey, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()
	_, _, err = key.GetStringValue(appName)
	return err == nil
}

// PlaySound plays a Windows system sound.
func PlaySound() {
	// Play the Windows notification sound
	go exec.Command("powershell", "-NoProfile", "-Command",
		"[System.Media.SystemSounds]::Exclamation.Play()").Run()
}

// ShowNotification shows a Windows toast notification using PowerShell.
func ShowNotification(title, message string) {
	script := `
Add-Type -AssemblyName System.Windows.Forms
$notify = New-Object System.Windows.Forms.NotifyIcon
$notify.Icon = [System.Drawing.SystemIcons]::Information
$notify.BalloonTipTitle = '` + title + `'
$notify.BalloonTipText = '` + message + `'
$notify.Visible = $true
$notify.ShowBalloonTip(5000)
Start-Sleep -Milliseconds 5500
$notify.Dispose()
`
	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	go cmd.Run()
}

// ForceForegroundWindow brings the window to focus using Windows API.
func ForceForegroundWindow() {
	user32 := syscall.NewLazyDLL("user32.dll")
	getForeground := user32.NewProc("GetForegroundWindow")
	setForeground := user32.NewProc("SetForegroundWindow")
	findWindow := user32.NewProc("FindWindowW")

	// Find the Blink - Rest window
	title, _ := syscall.UTF16PtrFromString("Blink - Rest")
	hwnd, _, _ := findWindow.Call(0, uintptr(unsafe.Pointer(title)))
	if hwnd != 0 {
		setForeground.Call(hwnd)
	} else {
		// Fallback: set the current foreground window
		hwnd, _, _ = getForeground.Call()
		if hwnd != 0 {
			setForeground.Call(hwnd)
		}
	}
}
