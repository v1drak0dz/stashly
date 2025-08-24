package gitx

import (
	"fmt"
	"os/exec"
	"strings"

	git "github.com/go-git/go-git/v5"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func OpenRepo(path string) (*git.Repository, error) {
	return git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
}


func GetStatus(repo *git.Repository) (git.Status, error) {
    out, err := exec.Command("git", "status", "--porcelain").Output()
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(out), "\n")
    statusMap := make(git.Status)

    for _, l := range lines {
        if len(l) < 3 { // linhas curtas não valem
            continue
        }
        x := l[0] // staging
        y := l[1] // worktree
        file := strings.TrimSpace(l[3:])

        fs := &git.FileStatus{}

        // staging
        switch x {
        case 'M':
            fs.Staging = git.Modified
        case 'A':
            fs.Staging = git.Added
        case 'D':
            fs.Staging = git.Deleted
        }

        // worktree
        switch y {
        case 'M':
            fs.Worktree = git.Modified
        case 'A':
            fs.Worktree = git.Added
        case 'D':
            fs.Worktree = git.Deleted
        case '?':
            fs.Worktree = git.Untracked
        }

        statusMap[file] = fs
    }

    return statusMap, nil
}


func StageFiles(repo *git.Repository, files []string) error {
	wt, err := repo.Worktree()
	if err != nil {
		return err
	}
	for _, f := range files {
		if _, err := wt.Add(f); err != nil {
			return fmt.Errorf("error staging %s: %w", f, err)
		}
	}
	return nil
}

func Commit(repo *git.Repository, msg string) (string, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return "", err
	}
	hash, err := wt.Commit(msg, &git.CommitOptions{})
	if err != nil {
		return "", err
	}
	return hash.String(), nil
}

func Push(repo *git.Repository, auth gitssh.AuthMethod) error {
	return repo.Push(&git.PushOptions{Auth: auth})
}
