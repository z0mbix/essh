package main

import (
	"strings"

	"github.com/manifoldco/promptui"
)

func showMenu(instances []AwsInstance) (*AwsInstance, error) {
	searcher := func(i string, index int) bool {
		sInst := instances[index]
		name := sInst.NameTag
		input := i
		return strings.Contains(name, input) || strings.Contains(sInst.ID, input) || strings.Contains(sInst.ConnectIP, input)
	}

	templates := &promptui.SelectTemplates{
		Label:    `{{ . }}`,
		Active:   `{{ "»" | magenta }} {{ .NameTag | yellow }} {{ .ID | green }} ({{ .ConnectIP | red }})`,
		Inactive: `  {{ .NameTag }} {{ .ID | cyan }} ({{ .ConnectIP }})`,
		Selected: `{{ .NameTag | green }} {{ .ID | red }}`,
	}

	prompt := promptui.Select{
		Label:     "Select an instance:",
		Items:     instances,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, err
	}

	return &instances[i], nil
}
