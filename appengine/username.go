package main

import (
	"regexp"
	"strings"
)

var usernameRe = regexp.MustCompile(`^[a-z0-9]*$`)

func isAllowedUsername(s string) bool {
	if len(s) <= 1 {
		return false
	}
	if len(s) >= 31 {
		return false
	}
	if !usernameRe.MatchString(s) {
		return false
	}
	return strings.Index(s, "scrobble") == -1 && s != "username"
}
