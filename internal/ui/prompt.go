package ui

import "github.com/AlecAivazis/survey/v2"

func AskSelect(msg string, options []string) (string, error) {
	var selected string
	err := survey.AskOne(&survey.Select{
		Message: msg,
		Options: options,
		PageSize: 20,
	}, &selected)
	return selected, err
}

func AskMultiSelect(msg string, options []string) ([]string, error) {
	var selected []string
	err := survey.AskOne(&survey.MultiSelect{
		Message: msg,
		Options: options,
		PageSize: 20,
	}, &selected)
	return selected, err
}

func AskInput(msg string) (string, error) {
	var result string
	err := survey.AskOne(&survey.Input{Message: msg}, &result)
	return result, err
}

func AskConfirm(msg string) (bool, error) {
	var result bool
	err := survey.AskOne(&survey.Confirm{Message: msg}, &result)
	return result, err
}
