package gitx

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"stashly/internal/logger"

	git "github.com/go-git/go-git/v5"
)

// Abre o repositório Git na pasta atual
func OpenRepo(path string) (*git.Repository, error) {
	logger.PrintLog(fmt.Sprintf("Abrindo repositório em: %s", path))
	repo, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		logger.Error("Erro ao abrir repositório: " + err.Error())
		return nil, err
	}
	logger.PrintLog("Repositório aberto com sucesso")
	return repo, nil
}

// Retorna o status do repositório como map[file]status
type FileStatus struct {
	Status string
	Path   string
}

func parseStatus(index, worktree string) string {
	code := index
	if code == " " {
		code = worktree
	}

	switch code {
	case "A":
		return "new"
	case "M":
		return "modified"
	case "D":
		return "deleted"
	case "R":
		return "renamed"
	case "C":
		return "copied"
	case "?":
		return "untracked"
	case "!":
		return "ignored"
	default:
		return "unknown"
	}
}

// GetStatusPorcelain lê o status do git via "git status --porcelain"
func GetStatus() (map[string]*FileStatus, error) {
	logger.PrintLog("Obtendo status do repositório")
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		logger.Error("Erro ao obter status: " + err.Error())
		return nil, err
	}

	files := make(map[string]*FileStatus)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 4 {
			continue
		}
		indexStatus := string(line[0])
		worktreeStatus := string(line[1])
		path := strings.TrimSpace(line[3:])

		files[path] = &FileStatus{
			Status: parseStatus(indexStatus, worktreeStatus),
			Path:   path,
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Erro ao ler status linha a linha: " + err.Error())
		return nil, err
	}

	logger.PrintLog(fmt.Sprintf("Status obtido: %d arquivos", len(files)))
	return files, nil
}

// Staging de arquivos
func StageFiles(file string) error {
	logger.PrintLog("Fazendo stage do arquivo: " + file)
	cmd := exec.Command("git", "add", file)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(fmt.Sprintf("Erro ao adicionar arquivo %s: %s", file, string(out)))
		return fmt.Errorf("%s: %s", string(out), err)
	}
	logger.PrintLog("Arquivo adicionado ao stage: " + file)
	return nil
}

// Pega diff colorido usando git diff
func GetDiff(file string) (string, error) {
	logger.PrintLog("Obtendo diff para o arquivo: " + file)
	out, err := exec.Command("git", "diff", "--color=always", file).CombinedOutput()
	if err != nil {
		logger.Error("Erro ao gerar diff: " + err.Error())
		return "", err
	}

	diff := string(out)
	if diff == "" {
		diff = "(sem alterações para mostrar)"
	}

	logger.PrintLog("Diff obtido para: " + file)
	return diff, nil
}
