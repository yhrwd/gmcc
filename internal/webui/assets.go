package webui

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"path"
	"strings"
	"time"
)

//go:embed all:dist
var embeddedFiles embed.FS

type AssetFile struct {
	Name        string
	Content     []byte
	Size        int64
	ModTime     time.Time
	ContentType string
}

type UIAssets interface {
	HasIndex() bool
	LookupAsset(requestPath string) (*AssetFile, error)
	OpenIndex() (*AssetFile, error)
	IsAssetLikePath(requestPath string) bool
}

type fsAssets struct {
	root fs.FS
}

func NewEmbeddedAssets() UIAssets {
	return NewFSAssets(embeddedFiles)
}

func NewFSAssets(root fs.FS) UIAssets {
	return &fsAssets{root: root}
}

func (a *fsAssets) HasIndex() bool {
	file, err := a.readFile("/index.html")
	return err == nil && file != nil
}

func (a *fsAssets) OpenIndex() (*AssetFile, error) {
	return a.readFile("/index.html")
}

func (a *fsAssets) LookupAsset(requestPath string) (*AssetFile, error) {
	return a.readFile(requestPath)
}

func (a *fsAssets) IsAssetLikePath(requestPath string) bool {
	clean := path.Clean("/" + strings.TrimSpace(requestPath))
	if strings.HasPrefix(clean, "/assets/") {
		return true
	}
	if strings.HasPrefix(clean, "/.") {
		return true
	}

	base := path.Base(clean)
	return strings.Contains(base, ".")
}

func (a *fsAssets) readFile(requestPath string) (*AssetFile, error) {
	cleanPath, err := normalizeRequestPath(requestPath)
	if err != nil {
		return nil, err
	}
	if cleanPath == ".keep" {
		return nil, nil
	}

	assetPath := path.Join("dist", cleanPath)
	data, err := fs.ReadFile(a.root, assetPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read asset %q: %w", cleanPath, err)
	}

	var modTime time.Time
	var size int64 = int64(len(data))
	if info, statErr := fs.Stat(a.root, assetPath); statErr == nil {
		modTime = info.ModTime()
		size = info.Size()
	}

	return &AssetFile{
		Name:        cleanPath,
		Content:     data,
		Size:        size,
		ModTime:     modTime,
		ContentType: detectContentType(cleanPath),
	}, nil
}

func normalizeRequestPath(requestPath string) (string, error) {
	trimmed := strings.TrimSpace(requestPath)
	if trimmed == "" || trimmed == "/" {
		return "index.html", nil
	}
	if !strings.HasPrefix(trimmed, "/") {
		trimmed = "/" + trimmed
	}

	for _, segment := range strings.Split(trimmed, "/") {
		if segment == ".." {
			return "", fmt.Errorf("invalid asset path %q", requestPath)
		}
	}

	clean := path.Clean(trimmed)
	clean = strings.TrimPrefix(clean, "/")
	if clean == "" || clean == "." {
		return "index.html", nil
	}

	return clean, nil
}

func detectContentType(name string) string {
	switch path.Ext(name) {
	case ".html":
		return "text/html; charset=utf-8"
	case ".js":
		return "text/javascript; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".svg":
		return "image/svg+xml"
	case ".json":
		return "application/json"
	case ".ico":
		return "image/x-icon"
	}

	contentType := mime.TypeByExtension(path.Ext(name))
	if contentType == "" {
		return "application/octet-stream"
	}
	return contentType
}
