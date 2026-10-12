//go:build windows

package main

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	webview2 "github.com/jchv/go-webview2"
)

//go:embed ui.html
var uiHTML string

//go:embed assets/bg.webp
var bgWebp []byte

//go:embed assets/logo.webp
var logoWebp []byte

// Page size in CSS pixels; the window is scaled by the monitor's DPI so nothing gets cut off at 125%/150%.
const pageW, pageH = 960, 620

func pageHTML() string {
	d := func(b []byte) string { return "data:image/webp;base64," + base64.StdEncoding.EncodeToString(b) }
	return strings.NewReplacer("{{BG}}", d(bgWebp), "{{LOGO}}", d(logoWebp)).Replace(uiHTML)
}

type State struct {
	LauncherVersion string   `json:"launcherVersion"`
	Tps             string   `json:"tps"`
	TpsOk           bool     `json:"tpsOk"`
	SdkOk           bool     `json:"sdkOk"`
	Bl2             string   `json:"bl2"`
	Installed       string   `json:"installed"`
	Latest          string   `json:"latest"`
	Notes           []string `json:"notes"`
	Checked         bool     `json:"checked"`
	LauncherUpdate  string   `json:"launcherUpdate"`
	Busy            bool     `json:"busy"`
	Krieg           string   `json:"krieg"`
}

type App struct {
	w    webview2.WebView
	mu   sync.Mutex
	cfg  Config
	bl2  string
	feed *Feed
	busy bool
}

func (a *App) hwnd() uintptr { return uintptr(a.w.Window()) }

func (a *App) state() State {
	a.mu.Lock()
	defer a.mu.Unlock()
	s := State{LauncherVersion: LauncherVersion, Tps: a.cfg.TpsPath, TpsOk: IsTps(a.cfg.TpsPath), Bl2: a.bl2, Busy: a.busy}
	if s.TpsOk {
		s.SdkOk = HasModManager(s.Tps)
		s.Installed = InstalledVersion(s.Tps)
		s.Krieg = KriegStatus(s.Tps)
	}
	if a.feed != nil {
		s.Checked, s.Latest, s.Notes = true, a.feed.Version, a.feed.Notes
		if a.feed.LauncherFile != "" && Newer(a.feed.LauncherVersion, LauncherVersion) {
			s.LauncherUpdate = a.feed.LauncherVersion
		}
	}
	return s
}

// emit sends an event to the page from any goroutine.
func (a *App) emit(kind string, data interface{}) {
	b, _ := json.Marshal(map[string]interface{}{"kind": kind, "data": data})
	a.w.Dispatch(func() { a.w.Eval("window.onEvent && window.onEvent(" + string(b) + ")") })
}
func (a *App) log(format string, args ...interface{}) { a.emit("log", fmt.Sprintf(format, args...)) }
func (a *App) refresh()                               { a.emit("state", a.state()) }

// run does slow work off the UI thread, one job at a time.
func (a *App) run(job func()) bool {
	a.mu.Lock()
	if a.busy {
		a.mu.Unlock()
		return false
	}
	a.busy = true
	a.mu.Unlock()
	a.refresh()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				a.log("Something went wrong: %v", r)
			}
			a.mu.Lock()
			a.busy = false
			a.mu.Unlock()
			a.refresh()
		}()
		job()
	}()
	return true
}

func (a *App) checkUpdates(quiet bool) {
	a.run(func() {
		a.log("Checking for updates...")
		f, err := FetchFeed()
		if err != nil {
			if quiet {
				a.log("Update check skipped - couldn't reach GitHub (%v).", err)
			} else {
				a.log("Couldn't reach GitHub: %v", err)
			}
			return
		}
		a.mu.Lock()
		a.feed = f
		a.mu.Unlock()
		inst := ""
		if IsTps(a.cfg.TpsPath) {
			inst = InstalledVersion(a.cfg.TpsPath)
		}
		switch {
		case inst == "":
			a.log("Latest version: v%s. Click Install to get it.", f.Version)
		case Newer(f.Version, inst):
			a.log("Update available: v%s (you have v%s).", f.Version, inst)
		default:
			a.log("You're up to date (v%s).", inst)
		}
		if f.LauncherFile != "" && Newer(f.LauncherVersion, LauncherVersion) {
			a.log("A new launcher (v%s) is available - click \"Update launcher\".", f.LauncherVersion)
		}
	})
}

func (a *App) guardGame() bool {
	if gameRunning() {
		messageBox(a.hwnd(), "BL2 Arsenal", "Close The Pre-Sequel first, then try again.", mbOK|mbIconWarn)
		return false
	}
	return true
}

func (a *App) installLatest() {
	if !IsTps(a.cfg.TpsPath) || !a.guardGame() {
		return
	}
	a.run(func() {
		f := a.feed
		if f == nil {
			var err error
			if f, err = FetchFeed(); err != nil {
				a.log("Couldn't reach GitHub: %v", err)
				return
			}
			a.mu.Lock()
			a.feed = f
			a.mu.Unlock()
		}
		tmp := filepath.Join(os.TempDir(), f.File)
		defer os.Remove(tmp)
		a.log("Downloading v%s (%.1f MB)...", f.Version, float64(f.Size)/1e6)
		last := -1
		err := Download(f.File, f.Sha256, tmp, func(p float64) {
			if pc := int(p * 100); pc/5 != last/5 {
				last = pc
				a.emit("progress", p)
			}
		})
		a.emit("progress", -1)
		if err != nil {
			a.log("Download failed: %v", err)
			return
		}
		a.install(tmp)
	})
}

func (a *App) install(file string) {
	a.log("Installing %s...", filepath.Base(file))
	v, err := InstallSdkmod(a.cfg.TpsPath, file)
	if err != nil {
		a.log("Install failed: %v", err)
		return
	}
	a.log("Installed v%s. Your settings, save backups and built packages were kept.", v)
	if !HasModManager(a.cfg.TpsPath) {
		a.log("Note: willow2-mod-manager isn't installed yet - the mod needs it to load.")
	}
}

func (a *App) installFile() {
	if !IsTps(a.cfg.TpsPath) || !a.guardGame() {
		return
	}
	f := openFileDialog(a.hwnd(), "Choose a BL2 Arsenal package", "BL2 Arsenal package (*.sdkmod)\x00*.sdkmod\x00")
	if f != "" {
		a.run(func() { a.install(f) })
	}
}

func (a *App) chooseGame() {
	f := openFileDialog(a.hwnd(), "Find BorderlandsPreSequel.exe (in Binaries\\Win32)", "BorderlandsPreSequel.exe\x00BorderlandsPreSequel.exe\x00")
	if f == "" {
		return
	}
	dir := filepath.Dir(filepath.Dir(filepath.Dir(f)))
	if !IsTps(dir) {
		a.log("That isn't the Pre-Sequel's Binaries\\Win32\\BorderlandsPreSequel.exe.")
		return
	}
	a.mu.Lock()
	a.cfg.TpsPath = dir
	a.mu.Unlock()
	SaveConfig(a.cfg)
	a.log("Game folder set: %s", dir)
	a.refresh()
}

func (a *App) uninstall() {
	if !IsTps(a.cfg.TpsPath) || !a.guardGame() {
		return
	}
	if messageBox(a.hwnd(), "Uninstall BL2 Arsenal",
		"Remove BL2 Arsenal?\n\nPorted items disappear from your saves while the mod is removed (they come back if you reinstall). Your settings folder is kept.",
		mbYesNo|mbIconWarn) != idYes {
		return
	}
	a.run(func() {
		if err := Uninstall(a.cfg.TpsPath); err != nil {
			a.log("Uninstall failed: %v", err)
			return
		}
		a.log("BL2 Arsenal removed.")
	})
}

func (a *App) updateLauncher() {
	f := a.feed
	if f == nil || f.LauncherFile == "" {
		return
	}
	a.run(func() {
		exe, err := os.Executable()
		if err != nil {
			a.log("Can't find the launcher's own file: %v", err)
			return
		}
		next := exe + ".new"
		a.log("Downloading launcher v%s...", f.LauncherVersion)
		if err := Download(f.LauncherFile, f.LauncherSha256, next, nil); err != nil {
			a.log("Launcher download failed: %v", err)
			return
		}
		old := exe + ".old"
		os.Remove(old)
		if err := os.Rename(exe, old); err != nil { // Windows allows renaming a running .exe
			os.Remove(next)
			a.log("Couldn't replace the launcher: %v", err)
			return
		}
		if err := os.Rename(next, exe); err != nil {
			os.Rename(old, exe)
			a.log("Couldn't replace the launcher: %v", err)
			return
		}
		a.log("Launcher updated - restarting...")
		if err := exec.Command(exe).Start(); err == nil {
			a.w.Dispatch(func() { a.w.Terminate() })
		}
	})
}

func (a *App) play() {
	a.log("Starting The Pre-Sequel through Steam...")
	shellOpen(fmt.Sprintf("steam://rungameid/%d", TpsAppID))
}

func main() {
	if exe, err := os.Executable(); err == nil {
		os.Remove(exe + ".old") // left over from a launcher self-update
	}
	a := &App{cfg: LoadConfig()}
	if !IsTps(a.cfg.TpsPath) {
		exeDir := ""
		if exe, err := os.Executable(); err == nil {
			exeDir = filepath.Dir(exe)
		}
		if p := FindTpsNear(exeDir); p != "" {
			a.cfg.TpsPath = p
		} else {
			a.cfg.TpsPath = FindGame("BorderlandsPreSequel", "BorderlandsPreSequel.exe")
		}
	}
	if IsTps(a.cfg.TpsPath) {
		SaveConfig(a.cfg)
	}
	a.bl2 = FindGame("Borderlands 2", "Borderlands2.exe")

	w := webview2.NewWithOptions(webview2.WebViewOptions{
		AutoFocus: true,
		DataPath:  filepath.Join(configDir(), "webview"),
		WindowOptions: webview2.WindowOptions{
			Title: "BL2 Arsenal Launcher", Width: pageW, Height: pageH, IconId: 1, Center: true,
		},
	})
	if w == nil {
		if messageBox(0, "BL2 Arsenal Launcher",
			"This launcher needs Microsoft Edge WebView2, which isn't installed on this PC.\n\nOpen the download page now?",
			mbYesNo|mbIconError) == idYes {
			shellOpen("https://go.microsoft.com/fwlink/p/?LinkId=2124703")
		}
		return
	}
	defer w.Destroy()
	a.w = w
	cw, ch := scaledSize(uintptr(w.Window()), pageW, pageH)
	w.SetSize(cw, ch, webview2.HintFixed)
	centerWindow(uintptr(w.Window()))

	must := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	must(w.Bind("getState", func() State { return a.state() }))
	must(w.Bind("checkUpdates", func() { a.checkUpdates(false) }))
	must(w.Bind("installLatest", func() { a.installLatest() }))
	must(w.Bind("installFile", func() { a.installFile() }))
	must(w.Bind("chooseGame", func() { a.chooseGame() }))
	must(w.Bind("openModFolder", func() {
		if IsTps(a.cfg.TpsPath) {
			shellOpen(filepath.Join(a.cfg.TpsPath, "sdk_mods"))
		}
	}))
	must(w.Bind("uninstall", func() { a.uninstall() }))
	must(w.Bind("updateLauncher", func() { a.updateLauncher() }))
	must(w.Bind("play", func() { a.play() }))
	must(w.Bind("openUrl", func(u string) {
		if strings.HasPrefix(u, "https://") {
			shellOpen(u)
		}
	}))
	must(w.Bind("ready", func() { a.checkUpdates(true) }))
	w.SetHtml(pageHTML())
	w.Run()
}
