package updater

// InitUpdater initializes the updater and checks for updates if needed
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
