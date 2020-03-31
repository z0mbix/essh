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
		Active:   "» {{ .NameTag | yellow }} {{ .ID | green }} ({{ .ConnectIP | red }})",
		Inactive: "  {{ .NameTag }} {{ .ID | cyan }} ({{ .ConnectIP }})",
		Selected: "{{ .NameTag | green }} {{ .ID | red }}",
		Details: `
--------- Instances ----------
{{ .NameTag | yellow }} {{ .ID | green }} ({{ .ConnectIP | red }})`,
	}

	prompt := promptui.Select{
		Label:     "Select an Instance",
		Items:     instances,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		return nil, fmt.Errorf("failed to do menu things, err:%s", err)
	}

	return &instances[i], nil
}
