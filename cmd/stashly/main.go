package main

import (
	"stashly/internal/gitx"
	"stashly/internal/logger"
	"stashly/internal/ui"
	"stashly/internal/version"
)

const appVersion = "1.4.0"
const owner = "v1drak0dz"
const repoName = "stashly"

func main() {
	// 0. Inicializa o logger (se precisar de algum setup extra, pode chamar aqui)
	logger.PrintLog("Aplicativo iniciado")

	// 1. Checa nova versão
	newAvailable, latest, err := version.CheckNewVersion(appVersion, owner+"/"+repoName)
	if err != nil {
		logger.Error("Erro ao checar nova versão: " + err.Error())
	} else if newAvailable {
		logger.PrintLog("Nova versão disponível: " + latest)
	}

	// 2. Abre o repositório Git
	repo, err := gitx.OpenRepo(".")
	if err != nil {
		logger.Error("Erro ao abrir repositório: " + err.Error())
		return
	}
	logger.PrintLog("Repositório aberto com sucesso")

	// 3. Pega status do repositório
	status, err := gitx.GetStatus()
	if err != nil {
		logger.Error("Erro ao pegar status do repositório: " + err.Error())
		return
	}
	logger.PrintLog("Status do repositório obtido com sucesso")

	// 4. Pega últimos 20 commits
	commits, err := gitx.GetCommits(repo, 20)
	if err != nil {
		logger.Error("Erro ao pegar commits: " + err.Error())
		return
	}
	logger.PrintLog("Commits carregados com sucesso")

	// 5. Pega todas as branches locais
	branches, err := gitx.GetBranches(repo)
	if err != nil {
		logger.Error("Erro ao pegar branches: " + err.Error())
		return
	}
	logger.PrintLog("Branches carregadas com sucesso")

	// 6. Chama a TUI
	logger.PrintLog("Iniciando interface TUI")
	ui.RunView(status, commits, branches)
	logger.PrintLog("Aplicativo finalizado")
}
