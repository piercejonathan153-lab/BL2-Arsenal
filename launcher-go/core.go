package main

// Platform-independent launcher logic: versions, the update feed, downloads and installing the mod.
// Everything here is covered by core_test.go (runs on any OS).

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	LauncherVersion = "2.2.0"
	Owner           = "piercejonathan153-lab"
	Repo            = "BL2-Arsenal"
	ModFolder       = "bl2_arsenal"
	TpsAppID        = 261640
)

var FeedBase = "https://raw.githubusercontent.com/" + Owner + "/" + Repo + "/main"

// ---------------------------------------------------------------- versions

var verRe = regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)

// ParseVersion pulls "1.2.3" out of any string ("v1.2.3", "bl2_arsenal_v1.2.3.sdkmod"); ok=false if there is none.
func ParseVersion(s string) (v [3]int, ok bool) {
	m := verRe.FindStringSubmatch(s)
	if m == nil {
		return v, false
	}
	for i := 0; i < 3; i++ {
		v[i], _ = strconv.Atoi(m[i+1])
	}
	return v, true
}

// Newer reports whether version string a is newer than b (unparseable b counts as older).
func Newer(a, b string) bool {
	va, oka := ParseVersion(a)
	vb, okb := ParseVersion(b)
	if !oka {
		return false
	}
	if !okb {
		return true
	}
	for i := 0; i < 3; i++ {
		if va[i] != vb[i] {
			return va[i] > vb[i]
		}
	}
	return false
}

// ---------------------------------------------------------------- game folder

func IsTps(dir string) bool {
	if dir == "" {
		return false
	}
	st, err := os.Stat(filepath.Join(dir, "Binaries", "Win32", "BorderlandsPreSequel.exe"))
	return err == nil && !st.IsDir()
}

func HasModManager(tps string) bool {
	_, e1 := os.Stat(filepath.Join(tps, "Binaries", "Win32", "Plugins", "pyunrealsdk.dll"))
	_, e2 := os.Stat(filepath.Join(tps, "sdk_mods", "mods_base.sdkmod"))
	_, e3 := os.Stat(filepath.Join(tps, "sdk_mods", "mods_base"))
	return e1 == nil && (e2 == nil || e3 == nil)
}

var tomlVerRe = regexp.MustCompile(`(?m)^\s*version\s*=\s*"([^"]+)"`)

// InstalledVersion returns the installed mod version, or "" if the mod isn't installed.
func InstalledVersion(tps string) string {
	if b, err := os.ReadFile(filepath.Join(tps, "sdk_mods", ModFolder, "pyproject.toml")); err == nil {
		if m := tomlVerRe.FindSubmatch(b); m != nil {
			return string(m[1])
		}
		return "unknown"
	}
	if loose := looseSdkmods(tps); len(loose) > 0 {
		if v, ok := ParseVersion(filepath.Base(loose[0])); ok {
			return fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2])
		}
		return "unknown"
	}
	return ""
}

func looseSdkmods(tps string) []string {
	m, _ := filepath.Glob(filepath.Join(tps, "sdk_mods", ModFolder+"*.sdkmod"))
	return m
}

// KriegStatus says how the Krieg part of BL2 Arsenal can run: "standalone" (the separate KriegTPS mod is
// installed and used instead), "ready" (his character files from Vault Launcher are in DLC\Lilac) or "missing".
func KriegStatus(tps string) string {
	if st, err := os.Stat(filepath.Join(tps, "sdk_mods", "KriegTPS")); err == nil && st.IsDir() {
		return "standalone"
	}
	if m, _ := filepath.Glob(filepath.Join(tps, "sdk_mods", "KriegTPS*.sdkmod")); len(m) > 0 {
		return "standalone"
	}
	if st, err := os.Stat(filepath.Join(tps, "DLC", "Lilac", "Compat", "Content")); err == nil && st.IsDir() {
		return "ready"
	}
	return "missing"
}

// FindTpsNear walks up from dir looking for the game folder (for a launcher kept inside the TPS folder).
func FindTpsNear(dir string) string {
	for i := 0; i < 4 && dir != ""; i++ {
		if IsTps(dir) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// ---------------------------------------------------------------- config

type Config struct {
	TpsPath string `json:"TpsPath"`
}

func configDir() string {
	d, err := os.UserConfigDir()
	if err != nil {
		d = os.TempDir()
	}
	return filepath.Join(d, "BL2Arsenal")
}

func LoadConfig() Config {
	var c Config
	if b, err := os.ReadFile(filepath.Join(configDir(), "launcher.json")); err == nil {
		_ = json.Unmarshal(b, &c)
	}
	return c
}

func SaveConfig(c Config) {
	_ = os.MkdirAll(configDir(), 0o755)
	b, _ := json.MarshalIndent(c, "", "  ")
	_ = os.WriteFile(filepath.Join(configDir(), "launcher.json"), b, 0o644)
}

// ---------------------------------------------------------------- update feed

type Feed struct {
	Version         string   `json:"version"`
	File            string   `json:"file"`
	Size            int64    `json:"size"`
	Sha256          string   `json:"sha256"`
	Notes           []string `json:"notes"`
	LauncherVersion string   `json:"launcher_version"`
	LauncherFile    string   `json:"launcher_file"`
	LauncherSha256  string   `json:"launcher_sha256"`
}

var httpClient = &http.Client{Timeout: 60 * time.Second}

func get(url string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "BL2ArsenalLauncher/"+LauncherVersion)
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("server answered %s", resp.Status)
	}
	return resp, nil
}

func FetchFeed() (*Feed, error) {
	resp, err := get(fmt.Sprintf("%s/latest.json?t=%d", FeedBase, time.Now().Unix()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var f Feed
	if err := json.NewDecoder(resp.Body).Decode(&f); err != nil {
		return nil, fmt.Errorf("bad update info: %w", err)
	}
	if f.Version == "" || f.File == "" {
		return nil, errors.New("update info is incomplete")
	}
	return &f, nil
}

// Download saves releases/<name> to dest, checking the SHA-256 when one is given. progress gets 0..1.
func Download(name, wantSha, dest string, progress func(float64)) error {
	resp, err := get(FeedBase + "/releases/" + name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	h := sha256.New()
	total := resp.ContentLength
	var done int64
	buf := make([]byte, 64*1024)
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := out.Write(buf[:n]); werr != nil {
				out.Close()
				return werr
			}
			h.Write(buf[:n])
			done += int64(n)
			if progress != nil && total > 0 {
				progress(float64(done) / float64(total))
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			out.Close()
			return rerr
		}
	}
	if err := out.Close(); err != nil {
		return err
	}
	if wantSha != "" && !strings.EqualFold(hex.EncodeToString(h.Sum(nil)), wantSha) {
		os.Remove(dest)
		return errors.New("the download was damaged (checksum mismatch); nothing was changed")
	}
	return nil
}

// ---------------------------------------------------------------- install / uninstall

// InstallSdkmod unpacks a BL2 Arsenal .sdkmod into <tps>/sdk_mods/bl2_arsenal, replacing the old copy only after
// the new one is fully extracted. Settings, save backups and built packages (sdk_mods/settings) are never touched.
func InstallSdkmod(tps, file string) (string, error) {
	sdk := filepath.Join(tps, "sdk_mods")
	if st, err := os.Stat(sdk); err != nil || !st.IsDir() {
		return "", errors.New("there is no sdk_mods folder - install willow2-mod-manager first")
	}
	zr, err := zip.OpenReader(file)
	if err != nil {
		return "", errors.New("that file is not a valid .sdkmod package")
	}
	defer zr.Close()
	hasInit := false
	for _, f := range zr.File {
		if f.Name == ModFolder+"/__init__.py" {
			hasInit = true
		}
	}
	if !hasInit {
		return "", errors.New("that file is not a BL2 Arsenal package")
	}
	staging := filepath.Join(sdk, "."+ModFolder+".new")
	os.RemoveAll(staging)
	defer os.RemoveAll(staging)
	root := ModFolder + "/"
	for _, f := range zr.File {
		if !strings.HasPrefix(f.Name, root) || strings.Contains(f.Name, "__pycache__") {
			continue
		}
		rel := strings.TrimPrefix(f.Name, root)
		if rel == "" {
			continue
		}
		dst := filepath.Join(staging, filepath.FromSlash(rel))
		if !strings.HasPrefix(dst, staging+string(os.PathSeparator)) {
			return "", errors.New("the package contains an unsafe path")
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(dst, 0o755); err != nil {
				return "", err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return "", err
		}
		if err := extractOne(f, dst); err != nil {
			return "", err
		}
	}
	target := filepath.Join(sdk, ModFolder)
	old := filepath.Join(sdk, "."+ModFolder+".old")
	os.RemoveAll(old)
	if _, err := os.Stat(target); err == nil {
		if err := os.Rename(target, old); err != nil {
			return "", fmt.Errorf("couldn't replace the old version (is a file from it open?): %w", err)
		}
	}
	if err := os.Rename(staging, target); err != nil {
		os.Rename(old, target) // put the old version back
		return "", err
	}
	os.RemoveAll(old)
	for _, l := range looseSdkmods(tps) { // a loose .sdkmod copy would load the mod twice
		os.Remove(l)
	}
	return InstalledVersion(tps), nil
}

func extractOne(f *zip.File, dst string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, rc); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

func Uninstall(tps string) error {
	if err := os.RemoveAll(filepath.Join(tps, "sdk_mods", ModFolder)); err != nil {
		return err
	}
	for _, l := range looseSdkmods(tps) {
		os.Remove(l)
	}
	return nil
}
