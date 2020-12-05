// Package cmd implement common functions
package cmn

import "regexp"

// NameMatchedPattern checking the name using a regex template, returns true if it matched
func NameMatchedPattern(name string) (match bool) {
	match, _ = regexp.MatchString("^[a-zA-Z0-9_.-]*$", name)
	return match
}
