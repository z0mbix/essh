package menu

import (
	"errors"
	"fmt"
	"github.com/z0mbix/essh/internal/aws"
	"github.com/z0mbix/essh/internal/config"
	"strings"

	"github.com/manifoldco/promptui"
)


func GetInstance(sess *aws.Session) (*aws.Instance, error) {
	var instances []aws.Instance

	reservations, err := sess.GetReservations()
	if err != nil {
		return nil, err
	}

	if len(reservations) == 0 {
		return nil, errors.New("no instance found, add better logging here")
	}

	for rIdx := range reservations {
		for _, inst := range reservations[rIdx].Instances {
			i, err := aws.NewInstance(sess, inst, config.ConnectPublicIP)
			if err != nil {
				return nil, fmt.Errorf("could not get instance/session: %s", err)
			}
			instances = append(instances, *i)
		}
	}
	if len(instances) == 1 {
		return &instances[0], nil
	}

	return show(instances)
}
func show(instances []aws.Instance) (*aws.Instance, error) {
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
