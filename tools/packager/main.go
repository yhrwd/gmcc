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
	"sort"
	"strings"
)

const frontendBuildScript = "build"

type buildTarget struct {
	GOOS   string
	GOARCH string
}

var defaultBuildTargets = []buildTarget{
	{GOOS: "linux", GOARCH: "amd64"},
	{GOOS: "linux", GOARCH: "arm64"},
	{GOOS: "windows", GOARCH: "amd64"},
	{GOOS: "windows", GOARCH: "arm64"},
	{GOOS: "darwin", GOARCH: "amd64"},
	{GOOS: "darwin", GOARCH: "arm64"},
}

var rootWhitelist = map[string]struct{}{
	"index.html":           {},
	"favicon.ico":          {},
	"favicon.svg":          {},
	"manifest.webmanifest": {},
	"robots.txt":           {},
}

func main() {
	frontendDist := flag.String("frontend-dist", filepath.Join("frontend", "dist"), "frontend dist input")
	frontendDir := flag.String("frontend-dir", "frontend", "frontend project directory")
	embedDist := flag.String("embed-dist", filepath.Join("internal", "webui", "dist"), "embed dist output")
	output := flag.String("output", "", "single binary output path (legacy mode)")
	outputDir := flag.String("output-dir", "build", "binary output directory")
	targetsFlag := flag.String("targets", "", "comma-separated os/arch targets, such as linux/amd64,darwin/arm64")
	version := flag.String("version", "", "version injected into main.Version")
	flag.Parse()

	targets, err := resolveBuildTargets(*targetsFlag, *output)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve build targets: %v\n", err)
		os.Exit(1)
	}

	frontendBuilt, err := buildFrontendIfPresent(*frontendDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build frontend: %v\n", err)
		os.Exit(1)
	}
	if frontendBuilt {
		fmt.Printf("frontend build complete: %s\n", *frontendDir)
	} else {
		fmt.Println("frontend project unavailable; using existing dist if present")
	}

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

	if *output != "" {
		if err := buildBinary(*output, targets[0], *version); err != nil {
			fmt.Fprintf(os.Stderr, "build binary: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("binary written to %s\n", *output)
		return
	}

	outputs, err := buildBinaries(*outputDir, targets, *version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "build binaries: %v\n", err)
		os.Exit(1)
	}

	for _, outputPath := range outputs {
		fmt.Printf("binary written to %s\n", outputPath)
	}
}

func defaultOutputPath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join("build", "gmcc.exe")
	}

	return filepath.Join("build", "gmcc")
}

func defaultTargetsFlagValue() string {
	parts := make([]string, 0, len(defaultBuildTargets))
	for _, target := range defaultBuildTargets {
		parts = append(parts, target.String())
	}
	return strings.Join(parts, ",")
}

func resolveBuildTargets(targetsFlag, output string) ([]buildTarget, error) {
	if output != "" {
		if strings.TrimSpace(targetsFlag) == "" {
			return []buildTarget{{GOOS: runtime.GOOS, GOARCH: runtime.GOARCH}}, nil
		}

		targets, err := parseBuildTargets(targetsFlag)
		if err != nil {
			return nil, err
		}
		if len(targets) != 1 {
			return nil, fmt.Errorf("-output requires exactly one target, got %d", len(targets))
		}
		return targets, nil
	}

	if strings.TrimSpace(targetsFlag) == "" {
		return append([]buildTarget(nil), defaultBuildTargets...), nil
	}

	return parseBuildTargets(targetsFlag)
}

func parseBuildTargets(value string) ([]buildTarget, error) {
	parts := strings.Split(value, ",")
	seen := make(map[string]struct{}, len(parts))
	targets := make([]buildTarget, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		segments := strings.Split(part, "/")
		if len(segments) != 2 {
			return nil, fmt.Errorf("invalid target %q, want os/arch", part)
		}

		target := buildTarget{
			GOOS:   strings.TrimSpace(segments[0]),
			GOARCH: strings.TrimSpace(segments[1]),
		}
		if target.GOOS == "" || target.GOARCH == "" {
			return nil, fmt.Errorf("invalid target %q, want os/arch", part)
		}

		key := target.String()
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		targets = append(targets, target)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no build targets configured")
	}

	return targets, nil
}

func (t buildTarget) String() string {
	return t.GOOS + "/" + t.GOARCH
}

func (t buildTarget) outputName() string {
	name := fmt.Sprintf("gmcc-%s-%s", t.GOOS, t.GOARCH)
	if t.GOOS == "windows" {
		return name + ".exe"
	}
	return name
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

func buildFrontendIfPresent(frontendDir string) (bool, error) {
	packageJSONPath := filepath.Join(frontendDir, "package.json")
	if _, err := os.Stat(packageJSONPath); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	if err := runCommand(frontendDir, npmCommandArgs(frontendBuildScript)...); err != nil {
		return false, err
	}

	return true, nil
}

func npmCommandArgs(script string) []string {
	if runtime.GOOS == "windows" {
		return []string{"npm.cmd", "run", script}
	}

	return []string{"npm", "run", script}
}

func runCommand(dir string, args ...string) error {
	return runCommandWithEnv(dir, nil, args...)
}

func runCommandWithEnv(dir string, env []string, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}

	return cmd.Run()
}

func buildBinary(output string, target buildTarget, version string) error {
	if err := os.MkdirAll(filepath.Dir(output), 0o755); err != nil {
		return err
	}

	args := []string{"go", "build"}
	if version != "" {
		args = append(args, "-ldflags", fmt.Sprintf("-s -w -X main.Version=%s", version))
	}
	args = append(args, "-o", output, "./cmd/gmcc")

	env := []string{
		"GOOS=" + target.GOOS,
		"GOARCH=" + target.GOARCH,
	}

	return runCommandWithEnv("", env, args...)
}

func buildBinaries(outputDir string, targets []buildTarget, version string) ([]string, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	outputs := make([]string, 0, len(targets))
	for _, target := range targets {
		outputPath := filepath.Join(outputDir, target.outputName())
		if err := buildBinary(outputPath, target, version); err != nil {
			return nil, fmt.Errorf("build %s: %w", target.String(), err)
		}
		outputs = append(outputs, outputPath)
	}

	sort.Strings(outputs)
	return outputs, nil
}
