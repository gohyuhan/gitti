package updater

// InitUpdater initializes the updater and checks for updates if needed
func InitUpdater(currentVersion string) error {
	if ShouldCheckForUpdate() {
		latestVersion, isNewer, err := CheckForUpdates()
		if err != nil {
			return err
		}

		if isNewer {
			if PromptUserForUpdate(latestVersion) {
				// Placeholder for download logic
				// Download and update logic will go here
			}
		}

		SaveUpdateInfo()
	}
	return nil
}
