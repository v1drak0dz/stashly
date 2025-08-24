package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	"encoding/json"
	"net/http"
	"time"

	"stashly/internal/gitx"
	"stashly/internal/ui"

	git "github.com/go-git/go-git/v5"
)

const version = "1.3.0"
const repoOwner = "ryuvi"
const repoName = "stashly"

type Release struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

func checkLatestRelease() {
	client := http.Client{Timeout: 3 * time.Second}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return
	}

	if release.TagName != "v"+version {
		fmt.Printf("A new version is available: %s (you are on v%s)\n", release.HTMLURL, version)
		fmt.Printf("Download: %s\n\n", release.HTMLURL)
	}
}

func main() {
	include := flag.String("include", "", "Filter files to include (comma separated)")
	includer := flag.String("includer", "", "regex pattern to include files (comma separated)")
	exclude := flag.String("exclude", "", "Filter files to exclude (comma separated)")
	excluder := flag.String("excluder", "", "regex pattern to exclude files (comma separated)")
	lines := flag.Int("lines", 10, "Number of items do display")
	flag.Parse()

	checkLatestRelease()

	// abre repo
	repo, err := gitx.OpenRepo(".")
	if err != nil {
		log.Fatal(err)
	}

	// autenticação
	auth, err := gitx.GetAuth()
	if err != nil {
		log.Fatal(err)
	}

	// status
	status, err := gitx.GetStatus(repo)
	if err != nil {
		log.Fatal(err)
	}

	if len(status) == 0 {
		fmt.Println("Nothing to review")
		return
	}

	var includeRegex, excludeRegex *regexp.Regexp

	if *includer != "" {
		includeRegex, err = regexp.Compile(*includer)
		if err != nil {
			log.Fatal("Invalid includer regex:", err)
		}
	}

	if *excluder != "" {
		excludeRegex, err = regexp.Compile(*excluder)
		if err != nil {
			log.Fatal("Invalid excluder regex:", err)
		}
	}

	// montar lista formatada
	files := []string{}
	fileMap := map[string]string{}
	for f, s := range status {
		if passesFilter(f, *include, *exclude, includeRegex, excludeRegex) {
			display := ui.FormatStatus(f, s.Worktree)
			files = append(files, f)
			fileMap[f] = display
		}
	}

	if len(files) == 0 {
		fmt.Println("Nothing to review")
		return
	}

	// selecionar arquivos
	selected, err := ui.AskMultiSelectColored("Select files to stage:", files, fileMap)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	var realFiles []string
	for _, f := range selected {
		realFiles = append(realFiles, fileMap[f])
	}

	if err := gitx.StageFiles(repo, realFiles); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Staged files:")
	for _, f := range selected {
		fmt.Println("  ", f)
	}

	// commit
	msg, err := ui.AskInput("Commit message:")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if msg == "" {
		fmt.Println("Commit message cannot be empty")
		return
	}

	hash, err := gitx.Commit(repo, msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Commit made successfully:", hash)

	// push opcional
	push, err := ui.AskConfirm("Do you want to push?")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if push {
		err = gitx.Push(repo, auth)
		if err != nil {
			if err == git.NoErrAlreadyUpToDate {
				fmt.Println("Already up-to-date")
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

func passesFilter(file, includeSub, excludeSub string, includeRegex, excludeRegex *regexp.Regexp) bool {
	// include substring
	if includeSub != "" && !strings.Contains(file, includeSub) {
		return false
	}
	// exclude substring
	if excludeSub != "" && strings.Contains(file, excludeSub) {
		return false
	}
	// include regex
	if includeRegex != nil && !includeRegex.MatchString(file) {
		return false
	}
	// exclude regex
	if excludeRegex != nil && excludeRegex.MatchString(file) {
		return false
	}
	return true
}

