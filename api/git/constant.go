package git

const (
	PUSH               = "PUSH"
	FORCEPUSHSAFE      = "FORCEPUSHSAFE"
	FORCEPUSHDANGEROUS = "FORCEPUSHDANGEROUS"
)

const (
	NEWBRANCH          = "NEWBRANCH"
	NEWBRANCHANDSWITCH = "NEWBRANCHANDSWITCH"
)

const (
	SWITCHBRANCH            = "SWITCHBRANCH"
	SWITCHBRANCHWITHCHANGES = "SWITCHBRANCHWITHCHANGES"
)

const (
	GITPULL       = "GITPULL"       // this pull and continue based on user git pull onfiguration
	GITPULLREBASE = "GITPULLREBASE" // this pull and rebase
	GITPULLMERGE  = "GITPULLMERGE"  // this pull and merge with local
)

const (
	STASHALL        = "STASHALL"
	STASHINDIVIDUAL = "STASHINDIVIDUAL"
)
