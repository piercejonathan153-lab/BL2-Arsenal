package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var vdfPathRe = regexp.MustCompile(`"path"\s+"([^"]+)"`)

// LibrariesFromVdf lists the Steam library folders named in a libraryfolders.vdf.
func LibrariesFromVdf(vdf string) []string {
	var out []string
	for _, m := range vdfPathRe.FindAllStringSubmatch(vdf, -1) {
		out = append(out, strings.ReplaceAll(m[1], `\\`, `\`))
	}
	return out
}

// SteamLibraries returns every Steam library folder (roots come from the registry on Windows).
func SteamLibraries() []string {
	seen := map[string]bool{}
	var libs []string
	add := func(p string) {
		p = filepath.Clean(strings.ReplaceAll(p, "/", string(os.PathSeparator)))
		if p != "." && !seen[strings.ToLower(p)] {
			seen[strings.ToLower(p)] = true
			libs = append(libs, p)
		}
	}
	for _, root := range steamRoots() {
		add(root)
		if b, err := os.ReadFile(filepath.Join(root, "steamapps", "libraryfolders.vdf")); err == nil {
			for _, l := range LibrariesFromVdf(string(b)) {
				add(l)
			}
		}
	}
	return libs
}

// FindGame looks for steamapps/common/<folder> containing Binaries/Win32/<exe> in every library.
func FindGame(folder, exe string) string {
	for _, lib := range SteamLibraries() {
		p := filepath.Join(lib, "steamapps", "common", folder)
		if _, err := os.Stat(filepath.Join(p, "Binaries", "Win32", exe)); err == nil {
			return p
		}
	}
	return ""
}
