package gitx

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"

	"stashly/internal/logger"

	git "github.com/go-git/go-git/v5"
)

func GetBranches(repo *git.Repository) (map[string]bool, error) {
	logger.PrintLog("Obtendo branches locais")
	cmd := exec.Command("git", "branch")
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Erro ao obter branches: " + string(out))
		return nil, fmt.Errorf("%s: %s", string(out), err)
	}

	branches := make(map[string]bool)
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "* ") {
			branches[line[2:]] = true
		} else {
			branches[line] = false
		}
	}

	if err := scanner.Err(); err != nil {
		logger.Error("Erro ao escanear branches: " + err.Error())
		return nil, err
	}

	logger.PrintLog(fmt.Sprintf("Branches obtidas: %v", branches))
	return branches, nil
}

func CheckoutBranch(branch string) error {
	logger.PrintLog("Fazendo checkout na branch: " + branch)
	cmd := exec.Command("git", "checkout", branch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Erro ao fazer checkout: " + string(out))
		return fmt.Errorf("%s: %s", string(out), err)
	}
	logger.PrintLog("Checkout realizado com sucesso em: " + branch)
	return nil
}

func NewBranch(branch string) error {
	logger.PrintLog("Criando nova branch: " + branch)
	cmd := exec.Command("git", "checkout", "-b", branch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Erro ao criar nova branch: " + string(out))
		return fmt.Errorf("%s: %s", string(out), err)
	}
	logger.PrintLog("Nova branch criada: " + branch)
	return nil
}

func PushBranch(branch string) error {
	logger.PrintLog("Fazendo push da branch: " + branch)
	cmd := exec.Command("git", "push", "origin", branch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Erro ao dar push: " + string(out))
		return fmt.Errorf("%s: %s", string(out), err)
	}
	logger.PrintLog("Push realizado com sucesso na branch: " + branch)
	return nil
}

func PullBranch(branch string) error {
	logger.PrintLog("Fazendo pull da branch: " + branch)
	cmd := exec.Command("git", "pull", "origin", branch)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Erro ao dar pull: " + string(out))
		return fmt.Errorf("%s: %s", string(out), err)
	}
	logger.PrintLog("Pull realizado com sucesso na branch: " + branch)
	return nil
}

func GetCurrentBranch() (string, error) {
	logger.PrintLog("Obtendo branch atual")
	cmd := exec.Command("git", "branch")
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Erro ao obter branch atual: " + string(out))
		return "", fmt.Errorf("%s: %s", string(out), err)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "* ") {
			current := line[2:]
			logger.PrintLog("Branch atual: " + current)
			return current, nil
		}
	}
	logger.PrintLog("Nenhuma branch atual encontrada")
	return "", nil
}
