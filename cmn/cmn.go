package cmn

import "regexp"

func NameMatchedPattern(name string) (match bool) {
	match, _ = regexp.MatchString("^[a-zA-Z0-9_.-]*$", name)
	return match
}
