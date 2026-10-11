//go:build windows

package main

import (
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	user32            = windows.NewLazySystemDLL("user32.dll")
	comdlg32          = windows.NewLazySystemDLL("comdlg32.dll")
	procMessageBoxW   = user32.NewProc("MessageBoxW")
	procGetOpenFileW  = comdlg32.NewProc("GetOpenFileNameW")
	procGetDpiForWnd  = user32.NewProc("GetDpiForWindow")
	procGetWindowRect = user32.NewProc("GetWindowRect")
	procSetWindowPos  = user32.NewProc("SetWindowPos")
	procSysParamInfo  = user32.NewProc("SystemParametersInfoW")
)

type rect struct{ Left, Top, Right, Bottom int32 }

func workArea() rect {
	var r rect
	procSysParamInfo.Call(0x0030 /* SPI_GETWORKAREA */, 0, uintptr(unsafe.Pointer(&r)), 0)
	return r
}

// scaledSize turns a size in CSS pixels into window pixels for the window's monitor DPI, shrunk to fit the screen.
func scaledSize(hwnd uintptr, w, h int) (int, int) {
	dpi := uintptr(96)
	if procGetDpiForWnd.Find() == nil {
		if d, _, _ := procGetDpiForWnd.Call(hwnd); d > 0 {
			dpi = d
		}
	}
	scale := float64(dpi) / 96
	wa := workArea()
	if aw, ah := float64(wa.Right-wa.Left)*0.95, float64(wa.Bottom-wa.Top-40)*0.95; aw > 0 && ah > 0 {
		if f := aw / (float64(w) * scale); f < 1 {
			scale *= f
		}
		if f := ah / (float64(h) * scale); f < 1 {
			scale *= f
		}
	}
	return int(float64(w)*scale + .5), int(float64(h)*scale + .5)
}

func centerWindow(hwnd uintptr) {
	var r rect
	procGetWindowRect.Call(hwnd, uintptr(unsafe.Pointer(&r)))
	wa := workArea()
	x := wa.Left + (wa.Right-wa.Left-(r.Right-r.Left))/2
	y := wa.Top + (wa.Bottom-wa.Top-(r.Bottom-r.Top))/2
	if y < wa.Top {
		y = wa.Top
	}
	procSetWindowPos.Call(hwnd, 0, uintptr(x), uintptr(y), 0, 0, 0x0001|0x0004|0x0010) // NOSIZE|NOZORDER|NOACTIVATE
}

func steamRoots() []string {
	var roots []string
	for _, k := range []struct {
		root registry.Key
		path string
		val  string
	}{
		{registry.CURRENT_USER, `Software\Valve\Steam`, "SteamPath"},
		{registry.LOCAL_MACHINE, `SOFTWARE\WOW6432Node\Valve\Steam`, "InstallPath"},
		{registry.LOCAL_MACHINE, `SOFTWARE\Valve\Steam`, "InstallPath"},
	} {
		key, err := registry.OpenKey(k.root, k.path, registry.QUERY_VALUE)
		if err != nil {
			continue
		}
		if v, _, err := key.GetStringValue(k.val); err == nil && v != "" {
			roots = append(roots, v)
		}
		key.Close()
	}
	return roots
}

func gameRunning() bool {
	snap, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return false
	}
	defer windows.CloseHandle(snap)
	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))
	for err = windows.Process32First(snap, &pe); err == nil; err = windows.Process32Next(snap, &pe) {
		if strings.EqualFold(windows.UTF16ToString(pe.ExeFile[:]), "BorderlandsPreSequel.exe") {
			return true
		}
	}
	return false
}

func shellOpen(target string) {
	verb, _ := syscall.UTF16PtrFromString("open")
	t, _ := syscall.UTF16PtrFromString(target)
	_ = windows.ShellExecute(0, verb, t, nil, nil, windows.SW_SHOWNORMAL)
}

const (
	mbOK        = 0x0
	mbYesNo     = 0x4
	mbIconWarn  = 0x30
	mbIconInfo  = 0x40
	mbIconError = 0x10
	idYes       = 6
)

func messageBox(owner uintptr, title, text string, flags uint32) int {
	t, _ := syscall.UTF16PtrFromString(text)
	c, _ := syscall.UTF16PtrFromString(title)
	r, _, _ := procMessageBoxW.Call(owner, uintptr(unsafe.Pointer(t)), uintptr(unsafe.Pointer(c)), uintptr(flags))
	return int(r)
}

type openFileNameW struct {
	lStructSize       uint32
	hwndOwner         uintptr
	hInstance         uintptr
	lpstrFilter       *uint16
	lpstrCustomFilter *uint16
	nMaxCustFilter    uint32
	nFilterIndex      uint32
	lpstrFile         *uint16
	nMaxFile          uint32
	lpstrFileTitle    *uint16
	nMaxFileTitle     uint32
	lpstrInitialDir   *uint16
	lpstrTitle        *uint16
	flags             uint32
	nFileOffset       uint16
	nFileExtension    uint16
	lpstrDefExt       *uint16
	lCustData         uintptr
	lpfnHook          uintptr
	lpTemplateName    *uint16
	pvReserved        uintptr
	dwReserved        uint32
	flagsEx           uint32
}

// openFileDialog shows the standard Windows "Open" dialog; filter is "Name\x00*.ext\x00". Returns "" if cancelled.
func openFileDialog(owner uintptr, title, filter string) string {
	buf := make([]uint16, 1024)
	f := utf16z(filter + "\x00")
	t, _ := syscall.UTF16PtrFromString(title)
	ofn := openFileNameW{
		hwndOwner:   owner,
		lpstrFilter: &f[0],
		lpstrFile:   &buf[0],
		nMaxFile:    uint32(len(buf)),
		lpstrTitle:  t,
		flags:       0x00001000 | 0x00000800 | 0x00080000 | 0x00000008, // FILEMUSTEXIST|PATHMUSTEXIST|EXPLORER|NOCHANGEDIR
	}
	ofn.lStructSize = uint32(unsafe.Sizeof(ofn))
	r, _, _ := procGetOpenFileW.Call(uintptr(unsafe.Pointer(&ofn)))
	if r == 0 {
		return ""
	}
	return windows.UTF16ToString(buf)
}

func utf16z(s string) []uint16 {
	u := make([]uint16, 0, len(s)+1)
	for _, r := range s {
		if r == 0 {
			u = append(u, 0)
			continue
		}
		e, _ := syscall.UTF16FromString(string(r))
		u = append(u, e[:len(e)-1]...)
	}
	return append(u, 0)
}
