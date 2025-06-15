package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

const (
	TagBin          = "/opt/homebrew/bin/tag"
	DefaultTagColor = "Yellow"
	ApplicationsDir = "/Applications"
)

type CaskInfo struct {
	Casks []struct {
		Name      []string `json:"name"`
		Artifacts []struct {
			App []string `json:"app"`
		} `json:"artifacts"`
	} `json:"casks"`
}

func getCaskApps() ([]string, error) {
	cmd := exec.Command("brew", "list", "--cask")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	caskNames := strings.Fields(string(out))
	var appSet sync.Map
	var wg sync.WaitGroup
	errChan := make(chan error, len(caskNames))

	// Process each cask concurrently
	for _, cask := range caskNames {
		wg.Add(1)
		go func(caskName string) {
			defer wg.Done()
			infoCmd := exec.Command("brew", "info", "--json=v2", "--cask", caskName)
			infoOut, err := infoCmd.Output()
			if err != nil {
				errChan <- fmt.Errorf("failed to get info for cask %s: %v", caskName, err)
				return
			}

			var caskInfo CaskInfo
			if err := json.Unmarshal(infoOut, &caskInfo); err != nil {
				errChan <- fmt.Errorf("failed to parse JSON for cask %s: %v", caskName, err)
				return
			}

			if len(caskInfo.Casks) > 0 {
				for _, artifact := range caskInfo.Casks[0].Artifacts {
					for _, app := range artifact.App {
						if app != "" {
							appSet.Store(app, struct{}{})
						}
					}
				}
			}
		}(cask)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		slog.Error("Error processing cask", "error", err)
	}

	// Convert sync.Map to slice
	var apps []string
	appSet.Range(func(key, value interface{}) bool {
		apps = append(apps, key.(string))
		return true
	})

	return apps, nil
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fmt.Fprintf(os.Stderr, "  -tag-color string\n")
	fmt.Fprintf(os.Stderr, "        Color to use for tagging Homebrew cask apps (default \"%s\")\n", DefaultTagColor)
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s                    # Tag apps with default color (Yellow)\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -tag-color=Blue    # Tag apps with Blue color\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s -tag-color=Red     # Tag apps with Red color\n", os.Args[0])
	os.Exit(2)
}

func main() {
	// Configure slog to use JSON format
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Parse command line flags
	tagColor := flag.String("tag-color", DefaultTagColor, "Color to use for tagging Homebrew cask apps")
	flag.Usage = printUsage
	flag.Parse()

	slog.Info("Starting app tagging", "color", *tagColor)

	// Get cask apps
	caskApps, err := getCaskApps()
	if err != nil {
		slog.Error("Failed to get cask apps", "error", err)
		os.Exit(1)
	}

	slog.Info("Found cask apps", "count", len(caskApps))
	for _, app := range caskApps {
		slog.Info("Cask app", "name", app)
	}

	// Process all apps in /Applications
	apps, err := filepath.Glob(filepath.Join(ApplicationsDir, "*.app"))
	if err != nil {
		slog.Error("Failed to list /Applications", "error", err)
		os.Exit(1)
	}

	for _, appDir := range apps {
		appName := filepath.Base(appDir)
		for _, caskApp := range caskApps {
			if appName == caskApp {
				if isTagged(appDir, *tagColor) {
					slog.Info("Skipping app - already tagged", "app", appDir, "color", *tagColor)
					continue
				}
				slog.Info("Tagging app", "app", appDir, "color", *tagColor)
				if err := tagApp(appDir, *tagColor); err != nil {
					slog.Error("Failed to tag app", "app", appDir, "error", err)
				}
			}
		}
	}
}

func isTagged(appPath, color string) bool {
	cmd := exec.Command(TagBin, "--list", appPath)
	out, err := cmd.Output()
	if err != nil {
		slog.Error("Failed to check tags", "app", appPath, "error", err)
		return false
	}
	return strings.Contains(strings.ToLower(string(out)), strings.ToLower(color))
}

func tagApp(appPath, color string) error {
	cmd := exec.Command(TagBin, "--add", color, appPath)
	return cmd.Run()
}
