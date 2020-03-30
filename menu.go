package main

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

func showMenu(instances []AwsInstance) (*AwsInstance, error) {

	searcher := func(i string, index int) bool {
		sInst := instances[index]
		name := sInst.NameTag
		input := i

		return strings.Contains(name, input) || strings.Contains(sInst.ID, input)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F5A5  {{ .NameTag }}  {{ .ID | cyan }} ({{ .CoonectIP | red }})",
		Inactive: " {{ .NameTag }} {{ .ID | cyan }} ({{ .CoonectIP | red }})",
		Selected: "{{ .NameTag }} {{ .ID | red | cyan }}",
		Details: `
--------- Instances ----------
{{ "ID:" | faint }}	{{ .ID }}
{{ "CoonectIP IP:" | faint }}	{{ .CoonectIP }}`,
	}

	prompt := promptui.Select{
		Label:     "Select an Instance",
		Items:     instances,
		Templates: templates,
		Size:      4,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, fmt.Errorf("failed to do menu things, err:%s", err)
	}

	return &instances[i], nil

}
