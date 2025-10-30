package settings

import (
	"encoding/json"
	"gitti/api/git"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

var GITTICONFIGSETTINGS *GittiConfigSettings

type GittiConfigSettings struct {
	FileWatcherDebounceMS           int     `json:"file_watcher_debounce_milli_second"`
	GitFilesActiveRefreshDurationMS int     `json:"git_files_active_refresh_duration_milli_second"`
	GitFetchDurationMS              int     `json:"git_fetch_duration_milli_second"`
	GitInitDefaultBranch            string  `json:"git_init_default_branch"`
	LeftPanelWidthRatio             float64 `json:"left_panel_width_ratio"`
	RightPanelWidthRatio            float64 `json:"right_panel_width_ratio"`
	GitBranchComponentHeightRatio   float64 `json:"git_branch_component_height_ratio"`
	GitFilesComponentHeightRatio    float64 `json:"git_files_component_height_ratio"`
	LanguageCode                    string  `json:"language_code"`
}

var GittiDefaultConfigSettings = GittiConfigSettings{
	FileWatcherDebounceMS:           450,
	GitFilesActiveRefreshDurationMS: 2500,
	GitFetchDurationMS:              60000,
	GitInitDefaultBranch:            "master",
	LeftPanelWidthRatio:             0.3,
	RightPanelWidthRatio:            0.7,
	GitBranchComponentHeightRatio:   0.4,
	GitFilesComponentHeightRatio:    0.6,
	LanguageCode:                    "EN",
}

const AppName = "gitti"
const AppVersion = "v0.1.0"

// getConfigPath returns the config.json path (creates directories if needed)
//
// *Example on MacOs : /Users/<USER_NAME>/Library/Application Support/gitti/config.json
func getConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, AppName)
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "config.json"), nil
}

// InitOrReadConfig loads existing config, ensures schema correctness, or creates default.
func InitOrReadConfig() {
	GITTICONFIGSETTINGS = &GittiDefaultConfigSettings

	cfgPath, err := getConfigPath()
	if err != nil {
		return
	}

	// If config doesn't exist, create a default one
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		writeDefaultConfig(cfgPath)
		return
	}

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		writeDefaultConfig(cfgPath)
		return
	}

	var cfg GittiConfigSettings
	if err := json.Unmarshal(data, &cfg); err != nil {
		// Bad JSON → reset
		writeDefaultConfig(cfgPath)
		return
	}

	// Validate and fix missing or invalid fields
	changed := ensureConfigIntegrity(&cfg, &GittiDefaultConfigSettings)
	if changed {
		saveConfig(cfgPath, cfg)
	}

	GITTICONFIGSETTINGS = &cfg
}

// ensureConfigIntegrity checks every field against the default.
// If a field is zero or invalid (type mismatch), it assigns the default value.
func ensureConfigIntegrity(cfg *GittiConfigSettings, def *GittiConfigSettings) bool {
	cfgVal := reflect.ValueOf(cfg).Elem()
	defVal := reflect.ValueOf(def).Elem()
	changed := false

	for i := 0; i < cfgVal.NumField(); i++ {
		field := cfgVal.Field(i)
		defaultField := defVal.Field(i)

		switch field.Kind() {
		case reflect.String:
			if field.String() == "" {
				field.SetString(defaultField.String())
				changed = true
			}
		case reflect.Int, reflect.Int64:
			if field.Int() == 0 {
				field.SetInt(defaultField.Int())
				changed = true
			}
		case reflect.Float64:
			if field.Float() == 0 {
				field.SetFloat(defaultField.Float())
				changed = true
			}
		default:
			// for unsupported types, just reset if zero
			if reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface()) {
				field.Set(defaultField)
				changed = true
			}
		}
	}
	return changed
}

func writeDefaultConfig(cfgPath string) {
	saveConfig(cfgPath, GittiDefaultConfigSettings)
}

func saveConfig(cfgPath string, cfg GittiConfigSettings) {
	file, err := os.Create(cfgPath)
	if err != nil {
		return
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	_ = enc.Encode(cfg)
}

func UpdateLanguageCode(languageCode string) {
	GITTICONFIGSETTINGS.LanguageCode = strings.ToUpper(languageCode)
	cfgPath, err := getConfigPath()
	if err == nil {
		saveConfig(cfgPath, *GITTICONFIGSETTINGS)
	}
}

func UpdatedDefaultBranch(branchName string, applyToGit bool, cwd string) {
	GITTICONFIGSETTINGS.GitInitDefaultBranch = branchName
	cfgPath, err := getConfigPath()
	if err == nil {
		saveConfig(cfgPath, *GITTICONFIGSETTINGS)
		if applyToGit {
			git.SetGitInitDefaultBranch(branchName, cwd)
		}
	}
}
