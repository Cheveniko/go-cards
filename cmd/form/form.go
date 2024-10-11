package form

import (
	"github.com/charmbracelet/huh"
)

type CardInfo struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Profession    string `json:"profession"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	Location      string `json:"location"`
	GoogleMapsUrl string `json:"google_maps_url"`
	ImagePath     string `json:"image_path"`
	Slug          string `json:"slug"`
	Website       string `json:"website"`
	Linkedin      string `json:"linkedin"`
	Twitter       string `json:"twitter"`
	Instagram     string `json:"instagram"`
	Github        string `json:"github"`
	Facebook      string `json:"facebook"`
	Whatsapp      string `json:"whatsapp"`
	VcfUrl        string `json:"vcf_url"`
}

func CreateForm() (*huh.Form, *CardInfo) {

	var cardInfo = CardInfo{}

	form := huh.NewForm(

		// Gather card info
		huh.NewGroup(
			huh.NewInput().
				Value(&cardInfo.FirstName).
				Title("Bievenid@ a la aplicación de creación de tarjetas de presentación").
				Description("Por favor ingresa la siguiente información:\n").
				Placeholder("Nombre"),
			huh.NewInput().
				Value(&cardInfo.LastName).
				Placeholder("Apellido"),
			huh.NewInput().
				Value(&cardInfo.Profession).
				Placeholder("Profesión"),
			huh.NewInput().
				Value(&cardInfo.Email).
				Placeholder("Mail"),
			huh.NewInput().
				Value(&cardInfo.Phone).
				Placeholder("Teléfono (fomato 593...)"),
			huh.NewInput().
				Value(&cardInfo.Location).
				Placeholder("Dirección"),
			huh.NewInput().
				Value(&cardInfo.GoogleMapsUrl).
				Placeholder("Link google maps"),
			huh.NewInput().
				Value(&cardInfo.Website).
				Placeholder("Web"),
			huh.NewInput().
				Value(&cardInfo.Linkedin).
				Placeholder("Linkedin"),
			huh.NewInput().
				Value(&cardInfo.Twitter).
				Placeholder("Twitter (X)"),
			huh.NewInput().
				Value(&cardInfo.Instagram).
				Placeholder("Instagram"),
			huh.NewInput().
				Value(&cardInfo.Github).
				Placeholder("Github"),
			huh.NewInput().
				Value(&cardInfo.Facebook).
				Placeholder("Facebook"),
			huh.NewInput().
				Value(&cardInfo.ImagePath).
				Placeholder("Ruta de la foto de perfil:"),
			huh.NewInput().
				Value(&cardInfo.Slug).
				Placeholder("Slug"),
		),
	)

	return form, &cardInfo
}
