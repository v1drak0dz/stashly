package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"stashly/internal/gitx"
	"stashly/internal/logger"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func getFileStatusColor(status gitx.FileStatus) string {
	switch status.Status {
	case "modified":
		return tview.TranslateANSI("[ ] " + status.Path + " [yellow](modified)[-]")
	case "new":
		return tview.TranslateANSI("[ ] " + status.Path + " [green](new)[-]")
	case "deleted":
		return tview.TranslateANSI("[ ] " + status.Path + " [red](deleted)[-]")
	default:
		return status.Path
	}
}

func RunView(files map[string]*gitx.FileStatus, commits []string, branches map[string]bool) {
	app := tview.NewApplication()
	selectedFiles := map[string]bool{}

	// Listas à esquerda
	filesList := tview.NewList().ShowSecondaryText(false).SetSelectedFocusOnly(true).SetSelectedBackgroundColor(tcell.ColorDefault).SetSelectedTextColor(tcell.ColorPurple)
	commitsList := tview.NewList().ShowSecondaryText(false).SetSelectedFocusOnly(true).SetSelectedBackgroundColor(tcell.ColorDefault).SetSelectedTextColor(tcell.ColorPurple)
	branchesList := tview.NewList().ShowSecondaryText(false).SetSelectedFocusOnly(true).SetSelectedBackgroundColor(tcell.ColorDefault).SetSelectedTextColor(tcell.ColorPurple)

	currentFocus := 0
	focusables := []tview.Primitive{filesList, commitsList, branchesList}

	var _files []string
	for f, status := range files {
		if status.Status != "untracked" && status.Status != "ignored" {
			_files = append(_files, f)
			filesList.AddItem(getFileStatusColor(*status), "", 0, nil)
			selectedFiles[f] = false
		}
	}
	for _, c := range commits {
		commitsList.AddItem(c, "", 0, nil)
	}
	for b, current := range branches {
		if current {
			branchesList.AddItem("[yellow]● "+b+"[-]", "", 0, nil)
		} else {
			branchesList.AddItem(b, "", 0, nil)
		}
	}

	filesList.SetBorder(true).SetTitle("Files").SetTitleAlign(tview.AlignLeft)
	commitsList.SetBorder(true).SetTitle("Commits").SetTitleAlign(tview.AlignLeft)
	branchesList.SetBorder(true).SetTitle("Branches").SetTitleAlign(tview.AlignLeft)

	// Painel do diff à direita
	diffText := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(false)

	// Display o arquivo inicial
	out, _ := exec.Command("git", "diff", "--color=always", _files[0]).CombinedOutput()
	content := string(out)
	if content == "" {
		content = "(sem alterações para mostrar)"
	}
	diffText.SetText(tview.TranslateANSI(content)).SetBorder(true).SetTitle("Diff").SetTitleAlign(tview.AlignLeft)

	// Atualiza diff ao mudar arquivo selecionado
	filesList.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		out, _ := exec.Command("git", "diff", "--color=always", _files[index]).CombinedOutput()
		content := string(out)
		if content == "" {
			content = "(sem alterações para mostrar)"
		}
		diffText.SetText(tview.TranslateANSI(content)).SetBorder(true)
	})

	// Layout vertical à esquerda
	leftFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(filesList, 0, 1, true).
		AddItem(commitsList, 0, 1, false).
		AddItem(branchesList, 0, 1, false)

	helperText := []string{
		"[::b]a:[::-] Stage file | [::b]c:[::-] Commit | [::b]p:[::-] Push | [::b](ctrl+c) q:[::-] Quit",
		"[::b]q:[::-] Quit",
		"[::b]c:[::-] Checkout | [::b]n:[::-] New | [::b]p:[::-] Pull | [::b](ctrl+c) q:[::-] Quit",
	}

	commandsText := tview.NewTextView().
		SetDynamicColors(true).
		SetTextColor(tcell.ColorWhite).
		SetWrap(false).
		SetTextAlign(tview.AlignLeft).
		SetText(helperText[currentFocus])

	// Layout horizontal principal
	mainFlex := tview.NewFlex().
		AddItem(leftFlex, 40, 1, true).
		AddItem(diffText, 0, 2, false)

	mainVertical := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(mainFlex, 0, 1, true).
		AddItem(commandsText, 1, 1, false)

	pages := tview.NewPages()
	pages.AddPage("main", mainVertical, true, true)

	// Atalho para sair
	app.SetRoot(pages, true).SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch {
		case event.Rune() == 'q':
			app.Stop()
			return nil

		case event.Key() == tcell.KeyTab:
			currentFocus = (currentFocus + 1) % len(focusables)
			app.SetFocus(focusables[currentFocus])
			commandsText.SetText(helperText[currentFocus])
			return nil

		case event.Key() == tcell.KeyBacktab: // Shift+Tab
			currentFocus = (currentFocus - 1 + len(focusables)) % len(focusables)
			app.SetFocus(focusables[currentFocus])
			commandsText.SetText(helperText[currentFocus])
			return nil

		case event.Key() == tcell.KeyRune && event.Rune() == 'a':
			currentItem := filesList.GetCurrentItem()
			mainText, secondaryText := filesList.GetItemText(currentItem)

			selectedFiles[mainText] = !selectedFiles[mainText]
			if strings.HasPrefix(mainText, "[ ]") {
				filesList.SetItemText(currentItem, strings.Replace(mainText, "[ ]", "[green][✓][-]", 1), secondaryText)
			} else if strings.HasPrefix(mainText, "[✓]") {
				filesList.SetItemText(currentItem, strings.Replace(mainText, "[green][✓][-]", "[ ]", 1), secondaryText)
			}

		case event.Key() == tcell.KeyRune && event.Rune() == 'c':
			switch currentFocus {
			case 0:
				var staged []string
				for f, selected := range selectedFiles {
					if selected {
						if err := gitx.StageFiles(f); err != nil {
							err = gitx.StageFiles(f)
							if err != nil {
								logger.Error(fmt.Sprintf("Error staging file %s: %s", f, err))
							}
							staged = append(staged, f)
						}
					}
				}

				if len(staged) == 0 {
					return nil
				}

				var input *tview.InputField
				input = tview.NewInputField().
					SetLabel("Commit message: ").
					SetFieldWidth(0).
					SetDoneFunc(func(key tcell.Key) {
						msg := input.GetText()
						if msg != "" {
							_, err := gitx.Commit(msg)
							if err != nil {
								logger.Error(fmt.Sprintf("Error committing: %s", err))
							}
						}
						pages.RemovePage("commitModal")
					})

				modalFlex := tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(input, 3, 1, true)

				pages.AddPage("commitModal", modalFlex, true, true)
				app.SetFocus(input)
			case 2: // branch list
				idx := branchesList.GetCurrentItem()
				branchName, _ := branchesList.GetItemText(idx)
				branchName = strings.TrimSpace(tview.Escape(branchName))

				if err := gitx.CheckoutBranch(branchName); err != nil {
					logger.Error(fmt.Sprintf("Error checking out branch %s: %s", branchName, err))
				}

				// Atualiza cores da lista
				for i := 0; i < branchesList.GetItemCount(); i++ {
					name, _ := branchesList.GetItemText(i)
					if name == branchName {
						branchesList.SetItemText(i, fmt.Sprintf("[yellow]● %s[-]", name), "")
					} else {
						branchesList.SetItemText(i, name, "")
					}
				}
			}

			return nil

		case event.Key() == tcell.KeyRune && event.Rune() == 'n':
			if currentFocus == 2 { // branch list
				// Input para nova branch
				var input *tview.InputField
				input = tview.NewInputField().
					SetLabel("New branch name: ").
					SetFieldWidth(30).
					SetDoneFunc(func(key tcell.Key) {
						branchName := input.GetText()
						if branchName != "" {
							if err := gitx.NewBranch(branchName); err != nil {
								logger.Error(fmt.Sprintf("Error creating branch %s: %s", branchName, err))
							} else {
								// Atualiza lista de branches
								branchesList.AddItem(fmt.Sprintf("[yellow]%s[-]", branchName), "", 0, nil)
							}
						}
						pages.RemovePage("newBranchModal")
						app.SetFocus(branchesList)
					})

				modalFlex := tview.NewFlex().
					SetDirection(tview.FlexRow).
					AddItem(input, 3, 1, true)

				pages.AddPage("newBranchModal", modalFlex, true, true)
				app.SetFocus(input)
			}

		case event.Key() == tcell.KeyRune && event.Rune() == 'p':
			curBranch, _ := gitx.GetCurrentBranch()

			switch currentFocus {
			case 0:
				err := gitx.PushBranch(curBranch)
				if err != nil {
					logger.Error(fmt.Sprintf("Error pushing branch: %s", err))
				}
			case 2:
				err := gitx.PullBranch(curBranch)
				if err != nil {
					logger.Error(fmt.Sprintf("Error pulling branch: %s", err))
				}
			}
		}
		return event
	})
	app.SetFocus(focusables[currentFocus])

	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
