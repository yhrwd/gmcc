package webui

import (
	"io/fs"
	"testing"
	"testing/fstest"
	"time"
)

func TestHasIndexReturnsFalseWhenOnlyKeepFileExists(t *testing.T) {
	assets := NewFSAssets(fstest.MapFS{
		"dist/.keep": {Data: []byte{}},
	})

	if assets.HasIndex() {
		t.Fatal("expected HasIndex to be false when only .keep exists")
	}
}

func TestHasIndexReturnsTrueWhenIndexExists(t *testing.T) {
	assets := NewFSAssets(fstest.MapFS{
		"dist/.keep":      {Data: []byte{}},
		"dist/index.html": {Data: []byte("<html></html>")},
	})

	if !assets.HasIndex() {
		t.Fatal("expected HasIndex to be true")
	}
}

func TestLookupAssetReturnsNilForMissingFile(t *testing.T) {
	assets := NewFSAssets(fstest.MapFS{
		"dist/.keep": {Data: []byte{}},
	})

	file, err := assets.LookupAsset("/assets/app.js")
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if file != nil {
		t.Fatal("expected nil file for missing asset")
	}
}

func TestLookupAssetReturnsAssetMetadata(t *testing.T) {
	modTime := time.Unix(1700000000, 0)
	assets := NewFSAssets(fstest.MapFS{
		"dist/assets/app.js": {
			Data:    []byte("console.log('ok')"),
			Mode:    fs.FileMode(0644),
			ModTime: modTime,
		},
	})

	file, err := assets.LookupAsset("/assets/app.js")
	if err != nil {
		t.Fatalf("lookup failed: %v", err)
	}
	if file == nil {
		t.Fatal("expected asset file")
	}
	if file.Name != "assets/app.js" {
		t.Fatalf("expected asset name %q, got %q", "assets/app.js", file.Name)
	}
	if file.ContentType != "text/javascript; charset=utf-8" {
		t.Fatalf("unexpected content type: %q", file.ContentType)
	}
	if string(file.Content) != "console.log('ok')" {
		t.Fatalf("unexpected content: %q", string(file.Content))
	}
	if !file.ModTime.Equal(modTime) {
		t.Fatalf("expected mod time %v, got %v", modTime, file.ModTime)
	}
	if file.Size != int64(len("console.log('ok')")) {
		t.Fatalf("unexpected size: %d", file.Size)
	}
}

func TestLookupAssetRejectsDirectoryTraversal(t *testing.T) {
	assets := NewFSAssets(fstest.MapFS{
		"dist/index.html": {Data: []byte("<html></html>")},
	})

	_, err := assets.LookupAsset("/../../secret.txt")
	if err == nil {
		t.Fatal("expected traversal error")
	}
}

func TestIsAssetLikePathClassifiesExpectedPaths(t *testing.T) {
	assets := NewFSAssets(fstest.MapFS{})

	tests := []struct {
		path string
		want bool
	}{
		{path: "/assets/app.js", want: true},
		{path: "/favicon.ico", want: true},
		{path: "/robots.txt", want: true},
		{path: "/instances/demo", want: false},
		{path: "/instances/v1.2", want: true},
		{path: "/.well-known/test", want: true},
	}

	for _, tt := range tests {
		if got := assets.IsAssetLikePath(tt.path); got != tt.want {
			t.Fatalf("path %q: want %v, got %v", tt.path, tt.want, got)
		}
	}
}
