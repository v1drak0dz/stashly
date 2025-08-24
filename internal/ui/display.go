package ui

import (
	git "github.com/go-git/go-git/v5"
)

func FormatStatus(filename string, status git.StatusCode) string {
	switch status {
	case git.Untracked:
		return "new"
	case git.Modified:
		return "modified"
	case git.Deleted:
		return "deleted"
	default:
		return filename
	}
}
