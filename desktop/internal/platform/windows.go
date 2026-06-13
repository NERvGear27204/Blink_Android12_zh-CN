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

var (
	kernel32       = syscall.NewLazyDLL("kernel32.dll")
	createMutex    = kernel32.NewProc("CreateMutexW")
	releaseMutex   = kernel32.NewProc("ReleaseMutex")
	mutexHandle    uintptr
)

const errorAlreadyExists = 183

// AcquireSingleInstance tries to grab a named mutex.
// Returns true if this is the only instance, false if another is running.
func AcquireSingleInstance() bool {
	name, _ := syscall.UTF16PtrFromString("Global\\BlinkDesktopMutex")
	handle, _, err := createMutex.Call(0, 0, uintptr(unsafe.Pointer(name)))
	if handle == 0 {
		return false
	}
	// ERROR_ALREADY_EXISTS means another instance holds it
	if errno, ok := err.(syscall.Errno); ok && errno == errorAlreadyExists {
		syscall.CloseHandle(syscall.Handle(handle))
		return false
	}
	mutexHandle = handle
	return true
}

// ReleaseSingleInstance releases the mutex (called on exit).
func ReleaseSingleInstance() {
	if mutexHandle != 0 {
		releaseMutex.Call(mutexHandle)
		syscall.CloseHandle(syscall.Handle(mutexHandle))
	}
}

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

// PlaySound plays a Windows system sound without showing a terminal.
func PlaySound() {
	go func() {
		cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command",
			"[System.Media.SystemSounds]::Exclamation.Play()")
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		cmd.Run()
	}()
}

// SetWindowOpacity sets the opacity of a window by its title using Win32 API.
// opacity should be between 0.0 (transparent) and 1.0 (opaque).
func SetWindowOpacity(title string, opacity float64) {
	user32 := syscall.NewLazyDLL("user32.dll")
	findWindow := user32.NewProc("FindWindowW")
	getWindowLong := user32.NewProc("GetWindowLongW")
	setWindowLong := user32.NewProc("SetWindowLongW")
	setLayeredAttrs := user32.NewProc("SetLayeredWindowAttributes")

	titlePtr, _ := syscall.UTF16PtrFromString(title)
	hwnd, _, _ := findWindow.Call(0, uintptr(unsafe.Pointer(titlePtr)))
	if hwnd == 0 {
		return
	}

	// GWL_EXSTYLE = -20 → 0xFFFFFFEC as uint32
	var nExStyle int32 = -20
	const wsExLayered = 0x00080000
	const lwaAlpha = 0x02

	// Add WS_EX_LAYERED style
	style, _, _ := getWindowLong.Call(hwnd, uintptr(nExStyle))
	setWindowLong.Call(hwnd, uintptr(nExStyle), style|wsExLayered)

	// Set alpha (0-255)
	alpha := byte(opacity * 255)
	if alpha < 76 { // minimum ~30%
		alpha = 76
	}
	setLayeredAttrs.Call(hwnd, 0, uintptr(alpha), lwaAlpha)
}

// ForceForegroundWindow brings the overlay window to focus using Windows API.
func ForceForegroundWindow() {
	user32 := syscall.NewLazyDLL("user32.dll")
	setForeground := user32.NewProc("SetForegroundWindow")
	findWindow := user32.NewProc("FindWindowW")

	title, _ := syscall.UTF16PtrFromString("Blink - Rest")
	hwnd, _, _ := findWindow.Call(0, uintptr(unsafe.Pointer(title)))
	if hwnd != 0 {
		setForeground.Call(hwnd)
	}
}
