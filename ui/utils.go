package ui

// variables for indicating which panel/components/container or whatever the hell you wanna call it that the user is currently landed or selected, so that they can do precious action related to the part of whatever the hell you wanna call it
var (
	localBranchComponent  = "B1"
	filesChangesComponent = "B2"
	fileDiffComponent     = "B3"
)

// this is for tab ( there will be 4 tab for now, initialzation tab(only accesible when user's repo was not git initialized yet), home tab, commit logs tab, about gitti tab )
var (
	initializationTab = "A"
	homeTab           = "B"
	commitLogsTab     = "C"
	aboutGittiTab     = "D"
)
