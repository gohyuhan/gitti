package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

var GITTICONFIGSETTINGS *GittiConfigSettings

type GittiConfigSettings struct {
	FileWatcherDebounceMS           int     `json:"file_watcher_debounce_milli_second"`
	GitFilesActiveRefreshDurationMS int     `json:"git_files_active_refresh_duration_milli_second"`
	GitFetchDurationMS              int     `json:"git_fetch_duration_milli_second"`
	GitInitDefaultBranch            string  `json:"git_init_default_branch"`
	LeftPanelWidthRatio             float64 `json:"left_panel_width_ratio"`            // panel width ratio need to be add up to 1.00 in total
	RightPanelWidthRatio            float64 `json:"right_panel_width_ratio"`           // panel width ratio need to be add up to 1.00 in total
	GitBranchComponentHeightRatio   float64 `json:"git_branch_component_height_ratio"` // component height ratio need to be add up to 1.00 in total
	GitFilesComponentHeightRatio    float64 `json:"git_files_component_height_ratio"`  // component height ratio need to be add up to 1.00 in total
}

var GittiDefaultConfigSettings = GittiConfigSettings{
	FileWatcherDebounceMS:           450,   // 0.45 seconds
	GitFilesActiveRefreshDurationMS: 2500,  // 2.5 seconds
	GitFetchDurationMS:              60000, // 60 seconds
	GitInitDefaultBranch:            "master",
	LeftPanelWidthRatio:             0.3,
	RightPanelWidthRatio:            0.7,
	GitBranchComponentHeightRatio:   0.4,
	GitFilesComponentHeightRatio:    0.6,
}

// on MacOS will be ""/Users/<USER_NAME>/Library/Application Support/gitti/config.json"
func getConfigPath(appName string) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, appName)
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "config.json"), nil
}

// InitOrReadConfig loads existing config or creates default one
func InitOrReadConfig(appName string) {
	// default the setting config to default
	GITTICONFIGSETTINGS = &GittiDefaultConfigSettings
	cfgPath, err := getConfigPath(appName)
	if err != nil {
		return
	}

	// If config does not exist, create it
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		file, err := os.Create(cfgPath)
		if err != nil {
			return
		}
		defer file.Close()

		enc := json.NewEncoder(file)
		enc.SetIndent("", "  ") // pretty print
		if err := enc.Encode(GittiDefaultConfigSettings); err != nil {
			return
		}
		return
	}

	// Otherwise, read existing config
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return
	}

	var cfg GittiConfigSettings
	if err := json.Unmarshal(data, &cfg); err != nil {
		return
	}

	// set to the user's config settings if loaded successfully
	GITTICONFIGSETTINGS = &cfg
	return
}
