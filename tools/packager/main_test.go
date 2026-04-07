package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestDefaultTargetsFlagValueIncludesAllDefaultTargets(t *testing.T) {
	got := defaultTargetsFlagValue()
	want := "linux/amd64,linux/arm64,windows/amd64,windows/arm64,darwin/amd64,darwin/arm64"
	if got != want {
		t.Fatalf("want targets %q, got %q", want, got)
	}
}

func TestParseBuildTargetsDeduplicatesAndTrims(t *testing.T) {
	targets, err := parseBuildTargets(" linux/amd64 , windows/arm64,linux/amd64 ")
	if err != nil {
		t.Fatal(err)
	}

	want := []buildTarget{
		{GOOS: "linux", GOARCH: "amd64"},
		{GOOS: "windows", GOARCH: "arm64"},
	}
	if !reflect.DeepEqual(targets, want) {
		t.Fatalf("want targets %#v, got %#v", want, targets)
	}
}

func TestParseBuildTargetsRejectsInvalidValue(t *testing.T) {
	if _, err := parseBuildTargets("linux-amd64"); err == nil {
		t.Fatal("expected invalid target error")
	}
}

func TestResolveBuildTargetsForLegacyOutputRequiresSingleTarget(t *testing.T) {
	_, err := resolveBuildTargets("linux/amd64,windows/amd64", filepath.Join("build", "gmcc.exe"))
	if err == nil {
		t.Fatal("expected single-target validation error")
	}
}

func TestResolveBuildTargetsDefaultsToCurrentTargetForLegacyOutput(t *testing.T) {
	targets, err := resolveBuildTargets("", filepath.Join("build", "gmcc.exe"))
	if err != nil {
		t.Fatal(err)
	}

	want := []buildTarget{{GOOS: runtime.GOOS, GOARCH: runtime.GOARCH}}
	if !reflect.DeepEqual(targets, want) {
		t.Fatalf("want targets %#v, got %#v", want, targets)
	}
}

func TestBuildTargetOutputNameUsesExeForWindows(t *testing.T) {
	if got, want := (buildTarget{GOOS: "windows", GOARCH: "amd64"}).outputName(), "gmcc-windows-amd64.exe"; got != want {
		t.Fatalf("want output name %q, got %q", want, got)
	}

	if got, want := (buildTarget{GOOS: "linux", GOARCH: "arm64"}).outputName(), "gmcc-linux-arm64"; got != want {
		t.Fatalf("want output name %q, got %q", want, got)
	}
}

func TestPrepareEmbedDirWithoutFrontendDistKeepsOnlyKeepFile(t *testing.T) {
	tmp := t.TempDir()
	embedDir := filepath.Join(tmp, "embed")
	if err := os.MkdirAll(embedDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(embedDir, ".keep"), []byte{}, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(embedDir, "old.txt"), []byte("old"), 0o644); err != nil {
		t.Fatal(err)
	}

	prepared, err := prepareEmbedDir(filepath.Join(tmp, "missing-frontend-dist"), embedDir)
	if err != nil {
		t.Fatal(err)
	}
	if prepared {
		t.Fatal("expected prepared to be false")
	}
	if _, err := os.Stat(filepath.Join(embedDir, ".keep")); err != nil {
		t.Fatalf("expected .keep to remain: %v", err)
	}
	if _, err := os.Stat(filepath.Join(embedDir, "old.txt")); !os.IsNotExist(err) {
		t.Fatalf("expected old artifact removed, got %v", err)
	}
}

func TestPrepareEmbedDirCopiesOnlyWhitelistedFiles(t *testing.T) {
	tmp := t.TempDir()
	frontendDist := filepath.Join("testdata", "frontend-dist")
	embedDir := filepath.Join(tmp, "embed")

	prepared, err := prepareEmbedDir(frontendDist, embedDir)
	if err != nil {
		t.Fatal(err)
	}
	if !prepared {
		t.Fatal("expected prepared to be true")
	}

	wantExists := []string{
		".keep",
		"index.html",
		filepath.Join("assets", "app.js"),
		filepath.Join("assets", "nested", "chunk.js"),
		"favicon.ico",
		"favicon.svg",
		"manifest.webmanifest",
		"robots.txt",
	}
	for _, rel := range wantExists {
		if _, err := os.Stat(filepath.Join(embedDir, rel)); err != nil {
			t.Fatalf("expected %s to exist: %v", rel, err)
		}
	}

	wantMissing := []string{
		filepath.Join("assets", "app.js.map"),
		filepath.Join("assets", ".temp.js"),
		filepath.Join("assets", "nested", ".chunk.js.tmp"),
		"extra.txt",
		".temp.html",
	}
	for _, rel := range wantMissing {
		if _, err := os.Stat(filepath.Join(embedDir, rel)); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be absent, got %v", rel, err)
		}
	}
}

func TestPrepareEmbedDirCopiesViteStyleFrontendAssets(t *testing.T) {
	tmp := t.TempDir()
	frontendDist := filepath.Join(tmp, "frontend-dist")
	embedDir := filepath.Join(tmp, "embed")

	files := map[string]string{
		"index.html": `<!doctype html>
<html lang="zh-CN">
  <head>
    <meta charset="UTF-8" />
    <script type="module" crossorigin src="/assets/index-m3WTdUe_.js"></script>
    <link rel="stylesheet" crossorigin href="/assets/index-BIKtDdUB.css">
  </head>
  <body>
    <div id="app"></div>
  </body>
</html>
`,
		filepath.Join("assets", "index-m3WTdUe_.js"):     "console.log('vite js');",
		filepath.Join("assets", "index-BIKtDdUB.css"):    "body{background:#020617;}",
		"manifest.webmanifest":                           "{}",
		"robots.txt":                                     "User-agent: *\nAllow: /\n",
		filepath.Join("assets", "index-m3WTdUe_.js.map"): "{}",
		filepath.Join("assets", ".hidden.js"):            "console.log('hidden');",
		"vite.svg":                                       `<svg xmlns="http://www.w3.org/2000/svg"></svg>`,
	}

	for rel, content := range files {
		path := filepath.Join(frontendDist, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	prepared, err := prepareEmbedDir(frontendDist, embedDir)
	if err != nil {
		t.Fatal(err)
	}
	if !prepared {
		t.Fatal("expected prepared to be true")
	}

	wantExists := []string{
		"index.html",
		filepath.Join("assets", "index-m3WTdUe_.js"),
		filepath.Join("assets", "index-BIKtDdUB.css"),
		"manifest.webmanifest",
		"robots.txt",
	}
	for _, rel := range wantExists {
		if _, err := os.Stat(filepath.Join(embedDir, rel)); err != nil {
			t.Fatalf("expected %s to exist: %v", rel, err)
		}
	}

	wantMissing := []string{
		filepath.Join("assets", "index-m3WTdUe_.js.map"),
		filepath.Join("assets", ".hidden.js"),
		"vite.svg",
	}
	for _, rel := range wantMissing {
		if _, err := os.Stat(filepath.Join(embedDir, rel)); !os.IsNotExist(err) {
			t.Fatalf("expected %s to be absent, got %v", rel, err)
		}
	}
}

func TestPrepareEmbedDirWithoutIndexFallsBackToAPIOnly(t *testing.T) {
	tmp := t.TempDir()
	frontendDist := filepath.Join(tmp, "frontend-dist")
	embedDir := filepath.Join(tmp, "embed")
	if err := os.MkdirAll(filepath.Join(frontendDist, "assets"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(frontendDist, "assets", "app.js"), []byte("console.log('ok')"), 0o644); err != nil {
		t.Fatal(err)
	}

	prepared, err := prepareEmbedDir(frontendDist, embedDir)
	if err != nil {
		t.Fatal(err)
	}
	if prepared {
		t.Fatal("expected prepared to be false without index.html")
	}
	if _, err := os.Stat(filepath.Join(embedDir, ".keep")); err != nil {
		t.Fatalf("expected .keep to remain: %v", err)
	}
	if _, err := os.Stat(filepath.Join(embedDir, "assets", "app.js")); !os.IsNotExist(err) {
		t.Fatalf("expected copied assets to be cleared when index.html is missing, got %v", err)
	}
}

func TestBuildFrontendIfPresentReturnsFalseWhenProjectMissing(t *testing.T) {
	tmp := t.TempDir()

	built, err := buildFrontendIfPresent(filepath.Join(tmp, "missing-frontend"))
	if err != nil {
		t.Fatal(err)
	}
	if built {
		t.Fatal("expected built to be false")
	}
}

func TestBuildFrontendIfPresentRunsNpmBuildScript(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("script execution test is POSIX-only")
	}

	tmp := t.TempDir()
	frontendDir := filepath.Join(tmp, "frontend")
	binDir := filepath.Join(tmp, "bin")
	markerPath := filepath.Join(tmp, "npm-build-marker.txt")

	if err := os.MkdirAll(frontendDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(frontendDir, "package.json"), []byte(`{"name":"test-frontend"}`), 0o644); err != nil {
		t.Fatal(err)
	}

	npmPath := filepath.Join(binDir, "npm")
	script := fmt.Sprintf("#!/bin/sh\nprintf '%%s' \"$@\" > %q\n", markerPath)
	if err := os.WriteFile(npmPath, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}

	originalPath := os.Getenv("PATH")
	t.Setenv("PATH", fmt.Sprintf("%s%c%s", binDir, os.PathListSeparator, originalPath))

	built, err := buildFrontendIfPresent(frontendDir)
	if err != nil {
		t.Fatal(err)
	}
	if !built {
		t.Fatal("expected built to be true")
	}

	marker, err := os.ReadFile(markerPath)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(marker), "runbuild"; got != want {
		t.Fatalf("expected npm args %q, got %q", want, got)
	}
}
