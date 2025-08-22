package gitx

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func GetAuth() (gitssh.AuthMethod, error) {
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

	// fallback: chave padrão id_rsa
	home, _ := os.UserHomeDir()
	entries, err := os.ReadDir(filepath.Join(home, ".ssh"))
	if err != nil || len(entries) == 0 {
		return nil, err
	}

	var authfiles []string
	for _, entry := range entries {
		authfiles = append(authfiles, entry.Name())
	}

	var keyFile string
	prompt := &survey.Select{
		Message: "Choose a file with key to auth:",
		Options: authfiles,
	}
	err = survey.AskOne(prompt, &keyFile)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil, err
	}

	keyPath := filepath.Join(home, ".ssh", keyFile)
	return gitssh.NewPublicKeysFromFile("git", keyPath, "")
}
