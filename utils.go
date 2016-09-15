package main

import (
	"fmt"
	"time"

	"github.com/aybabtme/rgbterm"
)

// shortcut for "no" colour
func nc(s string) string {
	return rgbterm.String(s, 255, 255, 255, 0, 0, 0)
}

// HumanDuration returns a human-readable approximation of a duration
// This function is taken from the Docker project, and slightly modified
// to cap units at days.
// (eg. "About a minute", "4 hours ago", etc.)
// (c) 2013 Docker, inc. and the Docker authors (http://docker.io)
func HumanDuration(d time.Duration) string {
	if seconds := int(d.Seconds()); seconds < 1 {
		return "Less than a second"
	} else if seconds < 60 {
		return fmt.Sprintf("%d seconds", seconds)
	} else if minutes := int(d.Minutes()); minutes == 1 {
		return "About a minute"
	} else if minutes < 60 {
		return fmt.Sprintf("%d minutes", minutes)
	} else if hours := int(d.Hours()); hours == 1 {
		return "About an hour"
	} else if hours < 48 {
		return fmt.Sprintf("%d hours", hours)
	}
	return fmt.Sprintf("%d days", int(d.Hours()/24))
}
