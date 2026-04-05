package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var rootWhitelist = map[string]struct{}{
	"index.html":           {},
	"favicon.ico":          {},
	"favicon.svg":          {},
	"manifest.webmanifest": {},
	"robots.txt":           {},
}

func main() {
	frontendDist := flag.String("frontend-dist", filepath.Join("frontend", "dist"), "frontend dist input")
	embedDist := flag.String("embed-dist", filepath.Join("internal", "webui", "dist"), "embed dist output")
	output := flag.String("output", defaultOutputPath(), "binary output path")
	flag.Parse()

	prepared, err := prepareEmbedDir(*frontendDist, *embedDist)
	if err != nil {
		fmt.Fprintf(os.Stderr, "prepare embed dir: %v\n", err)
		os.Exit(1)
	}

	if prepared {
		fmt.Println("frontend assets prepared for embedding")
	} else {
		fmt.Println("frontend assets unavailable; building API-only binary")
	}

	if err := buildBinary(*output); err != nil {
		fmt.Fprintf(os.Stderr, "build binary: %v\n", err)
		os.Exit(1)
	}
}

func defaultOutputPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join("build", "gmcc.exe")
	}

	return filepath.Join("build", "gmcc")
}

func prepareEmbedDir(frontendDist, embedDir string) (bool, error) {
	if err := ensureCleanEmbedDir(embedDir); err != nil {
		return false, err
	}

	info, err := os.Stat(frontendDist)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if !info.IsDir() {
		return false, fmt.Errorf("frontend dist is not a directory: %s", frontendDist)
	}

	if err := copyWhitelistedFiles(frontendDist, embedDir); err != nil {
		return false, err
	}

	if _, err := os.Stat(filepath.Join(embedDir, "index.html")); err != nil {
		if os.IsNotExist(err) {
			if cleanErr := ensureCleanEmbedDir(embedDir); cleanErr != nil {
				return false, cleanErr
			}
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func ensureCleanEmbedDir(embedDir string) error {
	if err := os.MkdirAll(embedDir, 0o755); err != nil {
		return err
	}

	entries, err := os.ReadDir(embedDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Name() == ".keep" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(embedDir, entry.Name())); err != nil {
			return err
		}
	}

	keepPath := filepath.Join(embedDir, ".keep")
	if _, err := os.Stat(keepPath); err != nil {
		if os.IsNotExist(err) {
			return os.WriteFile(keepPath, []byte{}, 0o644)
		}
		return err
	}

	return nil
}

func copyWhitelistedFiles(frontendDist, embedDir string) error {
	return filepath.WalkDir(frontendDist, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(frontendDist, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)

		if rel == "." {
			return nil
		}

		if d.IsDir() {
			if shouldSkipDir(rel) {
				return filepath.SkipDir
			}
			return nil
		}

		if !isWhitelistedFile(rel) {
			return nil
		}

		return copyFile(path, filepath.Join(embedDir, filepath.FromSlash(rel)))
	})
}

func shouldSkipDir(rel string) bool {
	if strings.HasPrefix(rel, "assets/") {
		return hasHiddenSegment(strings.TrimPrefix(rel, "assets/"))
	}

	return strings.HasPrefix(filepath.Base(rel), ".")
}

func isWhitelistedFile(rel string) bool {
	if strings.Contains(rel, "/") {
		if !strings.HasPrefix(rel, "assets/") {
			return false
		}
		assetRel := strings.TrimPrefix(rel, "assets/")
		if assetRel == "" || hasHiddenSegment(assetRel) || strings.HasSuffix(assetRel, ".map") {
			return false
		}
		return true
	}

	_, ok := rootWhitelist[rel]
	return ok
}

func hasHiddenSegment(rel string) bool {
	for _, segment := range strings.Split(rel, "/") {
		if strings.HasPrefix(segment, ".") {
			return true
		}
	}
	return false
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return err
	}

	if err := out.Close(); err != nil {
		return err
	}

	return nil
}

func buildBinary(output string) error {
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return err
	}

	cmd := exec.Command("go", "build", "-o", output, "./cmd/gmcc")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
