package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVersions(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{{"1.0.12", "1.0.9", true}, {"v1.0.12", "1.0.12", false}, {"1.0.9", "1.0.12", false},
		{"2.0.0", "1.99.99", true}, {"1.0.1", "unknown", true}, {"x", "1.0.0", false}}
	for _, c := range cases {
		if got := Newer(c.a, c.b); got != c.want {
			t.Errorf("Newer(%q,%q)=%v", c.a, c.b, got)
		}
	}
}

func TestVdf(t *testing.T) {
	vdf := `"libraryfolders" { "0" { "path" "C:\\Program Files (x86)\\Steam" } "1" { "path" "D:\\SteamLibrary" } }`
	l := LibrariesFromVdf(vdf)
	if len(l) != 2 || l[1] != `D:\SteamLibrary` {
		t.Fatal(l)
	}
}

func fakeTps(t *testing.T, ver string) string {
	d := t.TempDir()
	for _, f := range []string{"Binaries/Win32/BorderlandsPreSequel.exe", "Binaries/Win32/Plugins/pyunrealsdk.dll",
		"sdk_mods/mods_base.sdkmod", "sdk_mods/settings/bl2_arsenal/keep.txt", "sdk_mods/bl2_arsenal_v1.0.1.sdkmod"} {
		os.MkdirAll(filepath.Join(d, filepath.Dir(f)), 0o755)
		os.WriteFile(filepath.Join(d, f), []byte("x"), 0o644)
	}
	if ver != "" {
		os.MkdirAll(filepath.Join(d, "sdk_mods/bl2_arsenal"), 0o755)
		os.WriteFile(filepath.Join(d, "sdk_mods/bl2_arsenal/pyproject.toml"), []byte("[project]\nversion = \""+ver+"\"\n"), 0o644)
		os.WriteFile(filepath.Join(d, "sdk_mods/bl2_arsenal/stale.py"), []byte("old"), 0o644)
	}
	return d
}

func makeSdkmod(t *testing.T, ver string, extra map[string]string) string {
	p := filepath.Join(t.TempDir(), "bl2_arsenal_v"+ver+".sdkmod")
	f, _ := os.Create(p)
	z := zip.NewWriter(f)
	files := map[string]string{"bl2_arsenal/__init__.py": "x", "bl2_arsenal/pyproject.toml": "version = \"" + ver + "\"\n",
		"bl2_arsenal/b2p/builder.py": "y", "bl2_arsenal/__pycache__/junk.pyc": "z"}
	for k, v := range extra {
		files[k] = v
	}
	for k, v := range files {
		w, _ := z.Create(k)
		w.Write([]byte(v))
	}
	z.Close()
	f.Close()
	return p
}

func TestInstallUpgrade(t *testing.T) {
	tps := fakeTps(t, "1.0.9")
	if !IsTps(tps) || !HasModManager(tps) || InstalledVersion(tps) != "1.0.9" {
		t.Fatal("fake game not recognised")
	}
	v, err := InstallSdkmod(tps, makeSdkmod(t, "1.0.12", nil))
	if err != nil || v != "1.0.12" {
		t.Fatal(v, err)
	}
	mod := filepath.Join(tps, "sdk_mods/bl2_arsenal")
	if _, err := os.Stat(filepath.Join(mod, "stale.py")); err == nil {
		t.Error("old files left behind")
	}
	if _, err := os.Stat(filepath.Join(mod, "b2p/builder.py")); err != nil {
		t.Error("subfolder missing")
	}
	if _, err := os.Stat(filepath.Join(mod, "__pycache__")); err == nil {
		t.Error("pycache copied")
	}
	if _, err := os.Stat(filepath.Join(tps, "sdk_mods/settings/bl2_arsenal/keep.txt")); err != nil {
		t.Error("settings touched")
	}
	if len(looseSdkmods(tps)) != 0 {
		t.Error("loose sdkmod not removed")
	}
	ents, _ := os.ReadDir(filepath.Join(tps, "sdk_mods"))
	for _, e := range ents {
		if strings.HasPrefix(e.Name(), ".") {
			t.Error("staging folder left:", e.Name())
		}
	}
}

func TestInstallRejects(t *testing.T) {
	tps := fakeTps(t, "1.0.9")
	bad := filepath.Join(t.TempDir(), "bad.sdkmod")
	os.WriteFile(bad, []byte("not a zip"), 0o644)
	if _, err := InstallSdkmod(tps, bad); err == nil {
		t.Error("accepted a corrupt file")
	}
	other := makeSdkmod(t, "1.0.0", nil)
	// a zip without our __init__.py
	p := filepath.Join(t.TempDir(), "other.sdkmod")
	f, _ := os.Create(p)
	z := zip.NewWriter(f)
	w, _ := z.Create("someothermod/__init__.py")
	w.Write([]byte("x"))
	z.Close()
	f.Close()
	if _, err := InstallSdkmod(tps, p); err == nil {
		t.Error("accepted another mod")
	}
	evil := makeSdkmod(t, "1.0.13", map[string]string{"bl2_arsenal/../../evil.txt": "x"})
	if _, err := InstallSdkmod(tps, evil); err == nil {
		t.Error("accepted a path escaping the mod folder")
	}
	if InstalledVersion(tps) != "1.0.9" {
		t.Error("failed install changed the installed version:", InstalledVersion(tps))
	}
	_ = other
}

func TestUninstall(t *testing.T) {
	tps := fakeTps(t, "1.0.9")
	if err := Uninstall(tps); err != nil || InstalledVersion(tps) != "" {
		t.Fatal(err, InstalledVersion(tps))
	}
	if _, err := os.Stat(filepath.Join(tps, "sdk_mods/settings/bl2_arsenal/keep.txt")); err != nil {
		t.Error("settings removed")
	}
}

func TestFeedAndDownload(t *testing.T) {
	payload := []byte("sdkmod bytes")
	sum := sha256.Sum256(payload)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/latest.json"):
			w.Write([]byte(`{"version":"1.0.12","file":"m.sdkmod","size":12,"sha256":"` + hex.EncodeToString(sum[:]) + `","notes":["a","b"]}`))
		case r.URL.Path == "/releases/m.sdkmod":
			w.Write(payload)
		default:
			w.WriteHeader(404)
		}
	}))
	defer srv.Close()
	old := FeedBase
	FeedBase = srv.URL
	defer func() { FeedBase = old }()
	f, err := FetchFeed()
	if err != nil || f.Version != "1.0.12" || len(f.Notes) != 2 {
		t.Fatal(f, err)
	}
	dst := filepath.Join(t.TempDir(), "m.sdkmod")
	if err := Download(f.File, f.Sha256, dst, nil); err != nil {
		t.Fatal(err)
	}
	if err := Download(f.File, strings.Repeat("0", 64), dst, nil); err == nil {
		t.Error("bad checksum accepted")
	}
	if _, err := os.Stat(dst); err == nil {
		t.Error("damaged download kept")
	}
	if err := Download("missing.sdkmod", "", dst, nil); err == nil {
		t.Error("404 accepted")
	}
}

func TestFindTpsNear(t *testing.T) {
	tps := fakeTps(t, "")
	sub := filepath.Join(tps, "BL2 Arsenal Launcher")
	os.MkdirAll(sub, 0o755)
	if FindTpsNear(sub) != tps {
		t.Error("launcher inside the game folder didn't find it")
	}
}
