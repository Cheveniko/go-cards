package deploy

import (
	"github.com/charmbracelet/huh"
)

type Deploy struct {
	Deploy bool `json:"deploy"`
}

func CreateForm() (*huh.Form, *Deploy) {
	var deploy = Deploy{}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Desplegar las nuevas páginas a CloudFare?").
				Value(&deploy.Deploy).
				Affirmative("Sí").
				Negative("No"),
		),
	)

	return form, &deploy
}
