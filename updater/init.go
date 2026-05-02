package updater

// ------------------------------------
//
//	Auto-check for updates and prompt user to install if a newer version is available
//
// ------------------------------------
func AutoUpdater() {
	if ShouldCheckForUpdate() {
		latestVersion, isNewer, err := CheckForUpdates()
		if err != nil {
			return
		}

		if isNewer {
			if PromptUserForUpdate(latestVersion) {
				Update()
			}
		}
	}
}
