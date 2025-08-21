package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	git "github.com/go-git/go-git/v5"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] != "--review" {
		fmt.Println("Usage: stashly --review")
		os.Exit(1)
	}

	// Abre o repositório atual
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}

	// Configura autenticação SSH
	auth, err := getGitAuth()
	if err != nil {
		log.Fatal(err)
	}

	// Obtém worktree e status
	worktree, err := repo.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	status, err := worktree.Status()
	if err != nil {
		log.Fatal(err)
	}

	// Prepara lista colorida de arquivos
	var files []string
	fileMap := make(map[string]string) // mapeia display -> nome real
	for f, s := range status {
		display := f
		if s.Worktree == git.Untracked {
			display = color.GreenString(f + " (new)")
		} else if s.Worktree == git.Modified {
			display = color.YellowString(f + " (modified)")
		} else if s.Worktree == git.Deleted {
			display = color.RedString(f + " (deleted)")
		}
		files = append(files, display)
		fileMap[display] = f
	}

	if len(files) == 0 {
		fmt.Println("Nothing to review")
		return
	}

	// Seleção de arquivos
	var selected []string
	prompt := &survey.MultiSelect{
		Message: "Select files to stage",
		Options: files,
	}
	survey.AskOne(prompt, &selected)

	for _, f := range selected {
		_, err := worktree.Add(fileMap[f])
		if err != nil {
			log.Println("Error on adding:", f, err)
		}
	}

	// Confirmação visual de staging
	fmt.Println("\nStaged files:")
	for _, f := range selected {
		fmt.Println("  ", f)
	}

	// Mensagem de commit
	var message string
	survey.AskOne(&survey.Input{Message: "Commit message:"}, &message)
	if message == "" {
		fmt.Println("Commit message cannot be empty")
		return
	}

	_, err = worktree.Commit(message, &git.CommitOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Commit made successfully")

	// Push opcional
	var push bool
	survey.AskOne(&survey.Confirm{Message: "Do you want to push?"}, &push)
	if push {
		err = repo.Push(&git.PushOptions{Auth: auth})
		if err != nil {
			if err == git.NoErrAlreadyUpToDate {
				fmt.Println("Already up-to-date, nothing to push")
			} else {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Pushed successfully")
		}
	} else {
		fmt.Println("Skipped pushing")
	}
}

// getGitAuth retorna um gitssh.AuthMethod usando SSH agent ou chave privada
func getGitAuth() (gitssh.AuthMethod, error) {
	// Tenta usar SSH agent
	sshAgentSock := os.Getenv("SSH_AUTH_SOCK")
	if sshAgentSock != "" {
		conn, err := net.Dial("unix", sshAgentSock)
		if err == nil {
			return &gitssh.PublicKeysCallback{
				User:     "git",
				Callback: agent.NewClient(conn).Signers,
			}, nil
		}
	}

	// Fallback para chave privada
	home, _ := os.UserHomeDir()
	keyPath := filepath.Join(home, ".ssh", "id_ed25519_vrs10")
	authKey, err := gitssh.NewPublicKeysFromFile("git", keyPath, "")
	if err != nil {
		return nil, err
	}
	return authKey, nil
}
