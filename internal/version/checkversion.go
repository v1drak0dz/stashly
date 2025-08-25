package version

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Release struct {
	TagName string `json:"tag_name"`
}

func CheckNewVersion(currentVersion string, repo string) (bool, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("erro ao buscar release: %d", resp.StatusCode)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, "", err
	}

	if release.TagName != currentVersion {
		return true, release.TagName, nil
	}
	return false, release.TagName, nil
}
