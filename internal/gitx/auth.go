package gitx

import (
	"net"
	"os"
	"path/filepath"

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
	keyPath := filepath.Join(home, ".ssh", "id_rsa")
	return gitssh.NewPublicKeysFromFile("git", keyPath, "")
}
