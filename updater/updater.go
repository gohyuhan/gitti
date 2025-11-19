package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"gitti/constant"
	"gitti/i18n"
	"gitti/settings"

	"golang.org/x/mod/semver"
)

const GittiRepoURL = "https://api.github.com/repos/gohyuhan/gitti/releases/latest"

// CheckForUpdates checks if a new version is available
func CheckForUpdates() (string, bool, error) {
	// Use GitHub API to fetch the latest release information
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", GittiRepoURL, nil)
	if err != nil {
		return "", false, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", false, fmt.Errorf("failed to fetch latest release: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", false, fmt.Errorf("failed to read response body: %v", err)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return "", false, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	currentVersion := constant.APPVERSION
	latestVersion := release.TagName
	isNewer := compareVersions(currentVersion, latestVersion)
	// Update the last fetch time after a successful update check
	SaveUpdateInfo()
	return latestVersion, isNewer, nil
}

// compareVersions compares two version strings to determine if the latest is newer
//
//	func compareVersions(current, latest string) (bool, error) {
//		// Normalize: remove leading 'v' if present
//		current = strings.TrimPrefix(current, "v")
//		latest = strings.TrimPrefix(latest, "v")
//
//		// Parse using the gold-standard semver library
//		c, err := semver.NewVersion(current)
//		if err != nil {
//			return false, err // invalid current version
//		}
//		l, err := semver.NewVersion(latest)
//		if err != nil {
//			return false, err // invalid latest version
//		}
//
//		return l.GreaterThan(c), nil
//	}
func compareVersions(current, latest string) bool {
	// Add 'v' prefix if missing (common for GitHub tags)
	if !strings.HasPrefix(current, "v") {
		current = "v" + current
	}
	if !strings.HasPrefix(latest, "v") {
		latest = "v" + latest
	}

	// Validate both
	if !semver.IsValid(current) || !semver.IsValid(latest) {
		return false // or handle error as needed
	}

	return semver.Compare(latest, current) > 0 // true if latest > current
}

// ShouldCheckForUpdate determines if an update check is due based on last fetch time
func ShouldCheckForUpdate() bool {
	lastFetchTime := LoadLastFetchTime()

	sevenDaysAgo := time.Now().UTC().AddDate(0, 0, -7)
	return lastFetchTime.Before(sevenDaysAgo) || lastFetchTime.IsZero()
}

// LoadUpdateInfo reads the last fetch time from the settings file
func LoadLastFetchTime() time.Time {
	return settings.GITTICONFIGSETTINGS.LastUpdateCheckTime
}

// SaveUpdateInfo saves the current time as the last fetch time
func SaveUpdateInfo() {
	settings.UpdateLastFetchTime()
}

// PromptUserForUpdate prompts the user to download the latest version
func PromptUserForUpdate(latestVersion string) bool {
	fmt.Printf(i18n.LANGUAGEMAPPING.UpdaterDownloadPrompt, latestVersion)
	var response string
	fmt.Scanln(&response)
	return response == "y" || response == "Y"
}

// Update handles the download and replacement of the current TUI application with the latest version
func Update() {
	// Fetch the latest version information
	latestVersion, isNewer, err := CheckForUpdates()
	if err != nil {
		fmt.Printf(i18n.LANGUAGEMAPPING.UpdaterFailToCheckForUpdate, err)
		os.Exit(1)
	}

	if !isNewer {
		fmt.Printf(i18n.LANGUAGEMAPPING.UpdaterAlreadyLatest, constant.APPVERSION)
		os.Exit(0)
	}

	fmt.Printf(i18n.LANGUAGEMAPPING.UpdaterDownloading, latestVersion)
	// Determine the correct binary URL based on OS and architecture
	osName := runtime.GOOS
	arch := runtime.GOARCH
	binaryURL := getBinaryURL(osName, arch, latestVersion)
	if binaryURL == "" {
		fmt.Printf(i18n.LANGUAGEMAPPING.UpdaterUnSupportedOS, osName, arch)
		os.Exit(1)
	}

	// Download the binary
	tempFile, err := downloadBinary(binaryURL)
	if err != nil {
		fmt.Printf(i18n.LANGUAGEMAPPING.UpdaterDownloadFail, err)
		os.Exit(1)
	}
	defer os.Remove(tempFile) // Clean up temporary file after use

	// Replace the current binary with the downloaded one
	err = replaceBinary(tempFile)
	if err != nil {
		fmt.Printf(i18n.LANGUAGEMAPPING.UpdaterBinaryReplaceFail, err)
		os.Exit(1)
	}

	fmt.Printf(i18n.LANGUAGEMAPPING.UpdaterDownloadSuccess, latestVersion)

	os.Exit(0)
}

// getBinaryURL constructs the URL for the binary based on OS, architecture, and version
func getBinaryURL(osName, arch, version string) string {
	// Clean version string if needed (remove 'v' prefix if present)
	version = strings.TrimPrefix(version, "v")

	// Map of OS and architecture to binary suffix
	binarySuffixes := map[string]map[string]string{
		"darwin": {
			"amd64": "gitti-%s-darwin-amd64.tar.gz",
			"arm64": "gitti-%s-darwin-arm64.tar.gz",
		},
		"linux": {
			"amd64": "gitti-%s-linux-amd64.tar.gz",
			"arm64": "gitti-%s-linux-arm64.tar.gz",
		},
		"windows": {
			"amd64": "gitti-%s-windows-amd64.zip",
			"arm64": "gitti-%s-windows-arm64.zip",
		},
	}

	if osMap, ok := binarySuffixes[osName]; ok {
		if suffix, ok := osMap[arch]; ok {
			fileName := fmt.Sprintf(suffix, version)
			return fmt.Sprintf("https://github.com/gohyuhan/gitti/releases/download/v%s/%s", version, fileName)
		}
	}
	return ""
}

// downloadBinary downloads the binary from the specified URL to a temporary file
func downloadBinary(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(i18n.LANGUAGEMAPPING.UpdaterDownloadUnexpectedStatusCode, resp.StatusCode)
	}

	tempFile, err := os.CreateTemp("", "gitti-update-*.tmp")
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

// replaceBinary replaces the current executable with the downloaded binary
func replaceBinary(tempFile string) error {
	// Get the path of the current executable
	execPath, err := os.Executable()
	if err != nil {
		return err
	}

	// Handle different OS behaviors
	if runtime.GOOS == "windows" {
		// On Windows, we can't replace a running executable, so rename the old one and move the new one
		backupPath := execPath + ".old"
		err = os.Rename(execPath, backupPath)
		if err != nil {
			return err
		}
		err = os.Rename(tempFile, execPath)
		if err != nil {
			// Try to restore the original if rename fails
			os.Rename(backupPath, execPath)
			return err
		}
		os.Remove(backupPath) // Clean up backup if successful
	} else {
		// On Unix-like systems, we can replace the executable directly
		err = os.Rename(tempFile, execPath)
		if err != nil {
			return err
		}
		// Set executable permissions
		err = os.Chmod(execPath, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}
