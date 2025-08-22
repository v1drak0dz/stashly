package ui

import (
	"github.com/fatih/color"
	git "github.com/go-git/go-git/v5"
)

func FormatStatus(filename string, status git.StatusCode) string {
	switch status {
	case git.Untracked:
		return color.GreenString(filename + " (new)")
	case git.Modified:
		return color.YellowString(filename + " (modified)")
	case git.Deleted:
		return color.RedString(filename + " (deleted)")
	default:
		return filename
	}
}
