package gitx

import (
	"fmt"
	"os/exec"

	"stashly/internal/logger"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetCommits(repo *git.Repository, n int) ([]string, error) {
	logger.PrintLog(fmt.Sprintf("Obtendo últimos %d commits", n))
	ref, err := repo.Head()
	if err != nil {
		logger.Error("Erro ao obter HEAD: " + err.Error())
		return nil, err
	}

	iter, err := repo.Log(&git.LogOptions{From: ref.Hash()})
	if err != nil {
		logger.Error("Erro ao obter log do repositório: " + err.Error())
		return nil, err
	}

	var commits []string
	err = iter.ForEach(func(c *object.Commit) error {
		commits = append(commits, fmt.Sprintf("%s - %s", c.Hash.String()[:7], c.Message))
		if len(commits) >= n {
			return object.ErrUnsupportedObject
		}
		return nil
	})
	if err != nil && err != object.ErrUnsupportedObject {
		logger.Error("Erro ao iterar commits: " + err.Error())
		return nil, err
	}

	logger.PrintLog(fmt.Sprintf("Commits obtidos: %d", len(commits)))
	return commits, nil
}

func Commit(msg string) (string, error) {
	logger.PrintLog("Realizando commit: " + msg)
	cmd := exec.Command("git", "commit", "-m", msg)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Erro ao realizar commit: " + string(out))
		return "", fmt.Errorf("%s: %s", string(out), err)
	}

	result := string(out)
	logger.PrintLog("Commit realizado com sucesso: " + result)
	return result, nil
}
