package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	TagBin          = "/opt/homebrew/bin/tag"
	CaskTagColor    = "Yellow"
	MASTagColor     = "Blue"
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

func checkMasInstalled() error {
	cmd := exec.Command("brew", "list")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to check brew list: %v", err)
	}
	if !strings.Contains(string(out), "mas") {
		return fmt.Errorf("mas command not found. Please install it with: brew install mas")
	}
	return nil
}

func getMasApps() ([]string, error) {
	cmd := exec.Command("mas", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get mas list: %v", err)
	}

	var apps []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			// The app name is everything after the ID and before the version
			// Example: "1365531024  1Blocker               (6.1.3)"
			appName := strings.TrimSpace(fields[1])
			apps = append(apps, appName)
		}
	}
	return apps, nil
}

func main() {
	// Check if mas is installed
	if err := checkMasInstalled(); err != nil {
		log.Fatal(err)
	}

	// Get cask apps
	caskApps, err := getCaskApps()
	if err != nil {
		log.Fatalf("Failed to get cask apps: %v", err)
	}

	// Get MAS apps
	masApps, err := getMasApps()
	if err != nil {
		log.Fatalf("Failed to get MAS apps: %v", err)
	}

	fmt.Println("Cask apps found:")
	for _, app := range caskApps {
		fmt.Printf("  %s\n", app)
	}

	fmt.Println("\nMAS apps found:")
	for _, app := range masApps {
		fmt.Printf("  %s\n", app)
	}

	// Process all apps in /Applications
	apps, err := filepath.Glob(filepath.Join(ApplicationsDir, "*.app"))
	if err != nil {
		log.Fatalf("Failed to list /Applications: %v", err)
	}

	for _, appDir := range apps {
		appName := filepath.Base(appDir)
		appNameWithoutExt := strings.TrimSuffix(appName, ".app")

		// Check if it's a cask app
		for _, caskApp := range caskApps {
			if appName == caskApp {
				if isTagged(appDir, CaskTagColor) {
					fmt.Printf("Skipping %s - already has %s tag\n", appDir, CaskTagColor)
					continue
				}
				fmt.Printf("Tagging %s with color %s\n", appDir, CaskTagColor)
				if err := tagApp(appDir, CaskTagColor); err != nil {
					log.Printf("Failed to tag %s: %v", appDir, err)
				}
			}
		}

		// Check if it's a MAS app
		for _, masApp := range masApps {
			if appNameWithoutExt == masApp {
				if isTagged(appDir, MASTagColor) {
					fmt.Printf("Skipping %s - already has %s tag\n", appDir, MASTagColor)
					continue
				}
				fmt.Printf("Tagging %s with color %s\n", appDir, MASTagColor)
				if err := tagApp(appDir, MASTagColor); err != nil {
					log.Printf("Failed to tag %s: %v", appDir, err)
				}
			}
		}
	}
}

func getCaskApps() ([]string, error) {
	cmd := exec.Command("brew", "list", "--cask")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	caskNames := strings.Fields(string(out))
	appSet := make(map[string]struct{})

	for _, cask := range caskNames {
		infoCmd := exec.Command("brew", "info", "--json=v2", "--cask", cask)
		infoOut, err := infoCmd.Output()
		if err != nil {
			log.Printf("Failed to get info for cask %s: %v", cask, err)
			continue
		}

		var caskInfo CaskInfo
		if err := json.Unmarshal(infoOut, &caskInfo); err != nil {
			log.Printf("Failed to parse JSON for cask %s: %v", cask, err)
			continue
		}

		if len(caskInfo.Casks) > 0 {
			for _, artifact := range caskInfo.Casks[0].Artifacts {
				for _, app := range artifact.App {
					if app != "" {
						appSet[app] = struct{}{}
					}
				}
			}
		}
	}

	var apps []string
	for app := range appSet {
		apps = append(apps, app)
	}
	return apps, nil
}

func isTagged(appPath, color string) bool {
	cmd := exec.Command(TagBin, "--list", appPath)
	out, err := cmd.Output()
	if err != nil {
		log.Printf("Failed to check tags for %s: %v", appPath, err)
		return false
	}
	return strings.Contains(strings.ToLower(string(out)), strings.ToLower(color))
}

func tagApp(appPath, color string) error {
	cmd := exec.Command(TagBin, "--add", color, appPath)
	return cmd.Run()
}
