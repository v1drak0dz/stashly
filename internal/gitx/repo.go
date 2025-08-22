package gitx

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func OpenRepo(path string) (*git.Repository, error) {
	return git.PlainOpen(path)
}

func GetStatus(repo *git.Repository) (git.Status, error) {
	wt, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	return wt.Status()
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
