package api

import "regexp"

var baseNameValidChars = regexp.MustCompile(`^[A-Za-z0-9 _-]{3,24}$`)

// isValidBaseName validates a name against the BaseName rules (used for usernames & group names):
// must be 3-24 characters, only alphanumeric + spaces/underscores/hyphens,
// and must contain at least one non-space character.
func isValidBaseName(name string) bool {
	if !baseNameValidChars.MatchString(name) {
		return false
	}
	for _, c := range name {
		if c != ' ' {
			return true
		}
	}
	return false
}
