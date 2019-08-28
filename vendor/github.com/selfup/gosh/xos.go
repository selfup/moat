package gosh

import "runtime"

// Slash returns cross platform (xos) specific slashes for file paths
func Slash() string {
	var separator string

	if runtime.GOOS == "windows" {
		separator = "\\"
	} else {
		separator = "/"
	}

	return separator
}
