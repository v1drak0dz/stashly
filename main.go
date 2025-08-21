package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	git "github.com/go-git/go-git/v5"
)

func main() {

  if len(os.Args) < 2 || os.Args[1] != "--review" {
    fmt.Println("Usage: stashly --review")
    os.Exit(1)
  }

  repo, err := git.PlainOpen(".")
  if err != nil {
    log.Fatal(err)
  }

  worktree, err := repo.Worktree()
  if err != nil {
    log.Fatal(err)
  }

  status, err := worktree.Status()
  if err != nil {
    log.Fatal(err)
  }

  var files []string
  for f := range status {
    files = append(files, f)
  }

  if len(files) == 0 {
    fmt.Println("Nothing to review")
    return
  }

  var selected []string
  prompt := &survey.MultiSelect{
    Message: "Select files to stage",
    Options: files,
  }
  survey.AskOne(prompt, &selected)

  for _, f := range selected {
    _, err := worktree.Add(f)
    if err != nil {
      log.Println("Error on adding:", f, err)
    }
  }

  var message string
  survey.AskOne(&survey.Input{Message: "Commit message:"}, &message)

  _,err = worktree.Commit(message, &git.CommitOptions{})
  if err != nil {
    log.Fatal(err)
  }

  var push bool
  survey.AskOne(&survey.Confirm{Message: "Do you want to push?"}, &push)
  if push {
    err = repo.Push(&git.PushOptions{})
    if err != nil {
      log.Fatal(err)
    }

    fmt.Println("Pushed successfully")
  } else {
    fmt.Println("Commit made successfully, but skipped pushing")
  }

}
