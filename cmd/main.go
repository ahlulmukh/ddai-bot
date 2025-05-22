package main

import (
	"bufio"
	"ddai-bot/internal/menu"
	"ddai-bot/internal/updater"
	"ddai-bot/internal/utils"
	"fmt"
	"os"
	"runtime"
	"strings"
)

const version = "1.0.0"

func main() {
	if updated := checkForUpdates(); updated {
		return
	}

	app := menu.NewMenuHandler()
	app.ShowMainMenu(version)
}

func checkForUpdates() bool {
	update, err := updater.CheckUpdate(version)
	if err != nil {
		utils.LogMessage(0, 0, "Update check failed: "+err.Error(), "warning")
		return false
	}

	if update != nil {
		utils.LogMessage(0, 0,
			fmt.Sprintf("Update v%s available! Current: v%s", update.Version, version),
			"info")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Update now? (y/n): ")
		input, _ := reader.ReadString('\n')

		if strings.TrimSpace(input) == "y" {
			performUpdate(update)
			return true
		}
	}
	return false
}

func performUpdate(update *updater.UpdateInfo) {
	var url, checksum string

	switch runtime.GOOS {
	case "windows":
		url = update.Windows.URL
		checksum = update.Windows.Checksum
	case "linux":
		if runtime.GOARCH == "amd64" {
			url = update.Linux.Amd64.URL
			checksum = update.Linux.Amd64.Checksum
		} else if runtime.GOARCH == "arm64" {
			url = update.Linux.Arm64.URL
			checksum = update.Linux.Arm64.Checksum
		} else {
			utils.LogMessage(0, 0, "Unsupported architecture for Linux", "error")
			return
		}
	default:
		utils.LogMessage(0, 0, "Unsupported OS for auto-update", "error")
		return
	}

	utils.LogMessage(0, 0, "Downloading update...", "process")
	tmpFile, err := updater.DownloadUpdate(url)
	if err != nil {
		utils.LogMessage(0, 0, "Download failed: "+err.Error(), "error")
		return
	}

	if !updater.VerifyChecksum(tmpFile, checksum) {
		utils.LogMessage(0, 0, "Checksum verification failed!", "error")
		os.Remove(tmpFile)
		return
	}

	utils.LogMessage(0, 0, "Applying update...", "process")
	if err := updater.ApplyUpdate(tmpFile); err != nil {
		utils.LogMessage(0, 0, "Update failed: "+err.Error(), "error")
		return
	}

	utils.LogMessage(0, 0, "Update successful! Restarting...", "success")
	os.Exit(0)
}
