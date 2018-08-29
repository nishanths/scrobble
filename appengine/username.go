package main

import (
	"regexp"
)

var usernameRe = regexp.MustCompile(`^[a-z0-9]*$`)

const (
	reasonLengthShort = iota
	reasonLengthLong
	reasonDisallowedChar
	reasonAlreadyTaken
)

func isAllowedUsername(s string) (bool, int) {
	if len(s) <= 1 {
		return false, reasonLengthShort
	}
	if len(s) >= 31 {
		return false, reasonLengthLong
	}
	if !usernameRe.MatchString(s) {
		return false, reasonDisallowedChar
	}
	return true, 0
}
