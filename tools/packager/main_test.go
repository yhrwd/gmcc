package main

import (
	"os"
	"path/filepath"
	"testing"
)

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
