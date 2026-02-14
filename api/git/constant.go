package git

// Git Push operation types
const (
	PUSH               = "PUSH"               // Standard git push
	FORCEPUSHSAFE      = "FORCEPUSHSAFE"      // Force push with lease (--force-with-lease)
	FORCEPUSHDANGEROUS = "FORCEPUSHDANGEROUS" // Standard force push (--force)
)

// Git Tag push operation types
const (
	TAGPUSH         = "TAGPUSH"         // Push a specific tag
	TAGPUSHALL      = "TAGPUSHALL"      // Push all tags
	TAGPUSHFORCE    = "TAGPUSHFORCE"    // Force push a specific tag
	TAGPUSHALLFORCE = "TAGPUSHALLFORCE" // Force push all tags
)

// Branch creation operation types
const (
	NEWBRANCH              = "NEWBRANCH"              // Create new branch
	NEWBRANCHANDSWITCH     = "NEWBRANCHANDSWITCH"     // Create and switch to new branch
	NEWBRANCHBASEDONREMOTE = "NEWBRANCHBASEDONREMOTE" // Create new branch based on a remote branch
)

// Branch switching operation types
const (
	SWITCHBRANCH            = "SWITCHBRANCH"            // Switch branch (stashes changes first)
	SWITCHBRANCHWITHCHANGES = "SWITCHBRANCHWITHCHANGES" // Switch branch and bring changes over
)

// Git Pull operation strategies
const (
	GITPULL       = "GITPULL"       // pull and continue based on user git pull configuration
	GITPULLREBASE = "GITPULLREBASE" // pull and rebase (--rebase --autostash)
	GITPULLMERGE  = "GITPULLMERGE"  // pull and merge (--no-rebase)
)

// Git Tag fetch operation types
const (
	TAGFETCH          = "TAGFETCH"          // Fetch new tags from remote (skips existing mismatched tags)
	TAGFETCHPRUNE     = "TAGFETCHPRUNE"     // Fetch tags and prune local tags that no longer exist on remote
	TAGFETCHOVERWRITE = "TAGFETCHOVERWRITE" // Fetch tags and overwrite local tags if they differ
	TAGFETCHMIRROR    = "TAGFETCHMIRROR"    // Comprehensive sync: fetch, prune, and overwrite to mirror remote tags
)

// Git Tag deletion operation types
const (
	TAGDELETELOCAL  = "TAGDELETELOCAL"  // Delete tag from the local repository
	TAGDELETEREMOTE = "TAGDELETEREMOTE" // Delete tag from the remote repository
)

// Git Stash operation types
const (
	STASHALL   = "STASHALL"   // Stash all changes (including untracked)
	STASHFILE  = "STASHFILE"  // Stash specific files
	APPLYSTASH = "APPLYSTASH" // Apply a stash
	DROPSTASH  = "DROPSTASH"  // Drop/delete a stash
	POPSTASH   = "POPSTASH"   // Pop a stash (apply and drop)
)

// Discard changes operation types
const (
	DISCARDWHOLE              = "DISCARDWHOLE"              // Discard all changes in a file
	DISCARDUNSTAGE            = "DISCARDUNSTAGE"            // Unstage changes in a file (git checkout -- <file>)
	DISCARDUNTRACKED          = "DISCARDUNTRACKED"          // Delete untracked file (git clean -f)
	DISCARDNEWLYADDEDORCOPIED = "DISCARDNEWLYADDEDORCOPIED" // Discard newly added or copied file (git rm -f)
	DISCARDANDREVERTRENAME    = "DISCARDANDREVERTRENAME"    // Discard and revert rename operation
)

// Conflict resolution operation types
const (
	RESETCONFLICT               = "RESETCONFLICT"               // Reset file with conflicts to merge state
	CONFLICTACCEPTOURSCHANGES   = "CONFLICTACCEPTOURSCHANGES"   // Accept our changes (--ours)
	CONFLICTACCEPTTHEIRSCHANGES = "CONFLICTACCEPTTHEIRSCHANGES" // Accept their changes (--theirs)
)

// Stream and update configurations
const (
	STREAMUPDATETHROTTLEMS = 150 // Frequency in milliseconds to throttle UI updates for streaming output
)

// Diff retrieval types
const (
	GETCOMBINEDDIFF = "GETCOMBINEDDIFF" // Get combined diff (staged + unstaged)
	GETSTAGEDDIFF   = "GETSTAGEDDIFF"   // Get staged diff only
	GETUNSTAGEDDIFF = "GETUNSTAGEDDIFF" // Get unstaged diff only
)

// Git Reset modes
const (
	RESETSOFT  = "RESETSOFT"  // Soft reset (--soft)
	RESETHARD  = "RESETHARD"  // Hard reset (--hard)
	RESETMIXED = "RESETMIXED" // Mixed reset (--mixed)
)

// Workspace staging operations
const (
	STAGE   = "STAGE"   // Stage changes (git add)
	UNSTAGE = "UNSTAGE" // Unstage changes (git reset)
)
