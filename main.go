package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Cheveniko/go-cards/cmd/deploy"
	"github.com/Cheveniko/go-cards/cmd/form"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	storage_go "github.com/supabase-community/storage-go"
	"github.com/supabase-community/supabase-go"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func toBase64(path string) string {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	base64Encoding := base64.StdEncoding.EncodeToString(bytes)

	return base64Encoding
}

const (
	imageBucket = "images"
	vcfBucket   = "vcf"
)

// My first Go program, don't judge 🤓
func main() {
	logo := ` 
  ____          ____              _     
 / ___| ___    / ___|__ _ _ __ __| |___ 
| |  _ / _ \  | |   / _' | '__/ _' / __|
| |_| | (_) | | |__| (_| | | | (_| \__ \
 \____|\___/   \____\__,_|_|  \__,_|___/
`
	logoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)

	fmt.Printf("%s\n", logoStyle.Render(logo))

	// Create a new form
	form, cardInfo := form.CreateForm()

	err := form.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	whatsappLink := "https://wa.me/" + cardInfo.Phone
	cardInfo.Whatsapp = whatsappLink

	// Initialize Supabase client
	SUPABASE_URL, _ := os.LookupEnv("API_URL")
	SUPABASE_KEY, _ := os.LookupEnv("API_KEY")

	client, err := supabase.NewClient(SUPABASE_URL, SUPABASE_KEY, nil)
	if err != nil {
		fmt.Println("Cannot initalize supabase client", err)
		os.Exit(1)
	}

	_ = spinner.New().Title("Creando tarjeta...").Action(func() {
		var base64Image string = ""

		// Edit image name and upload image to storage
		if cardInfo.ImageUrl != "" {

			image, _ := os.Open(cardInfo.ImageUrl)
			defer image.Close()

			newFilename := cardInfo.Slug + filepath.Ext(cardInfo.ImageUrl)

			mimeType := "image/jpeg"

			if filepath.Ext(cardInfo.ImageUrl) == ".png" {
				mimeType = "image/png"
			}

			_, err = client.Storage.UploadFile(imageBucket, newFilename, image, storage_go.FileOptions{ContentType: &mimeType})
			if err != nil {
				fmt.Println("Error uploading image", err)
				os.Exit(1)
			}

			imageUrl := client.Storage.GetPublicUrl(imageBucket, newFilename).SignedURL

			base64Image = toBase64(cardInfo.ImageUrl)

			cardInfo.ImageUrl = imageUrl
		}

		// Create vcf file
		vcf := `BEGIN:VCARD
VERSION:3.0
FN:` + cardInfo.FirstName + ` ` + cardInfo.LastName + `
N:` + cardInfo.LastName + `;` + cardInfo.FirstName + `;
TEL;TYPE=CELL:+` + cardInfo.Phone + `
EMAIL:` + cardInfo.Email + `
TITLE:` + cardInfo.Profession + `
URL:` + cardInfo.Website + `
NOTE:` + "Contacto agregado a través de Smart Cards" + `
PHOTO;ENCODING=b;TYPE=JPEG:` + base64Image + `
END:VCARD
`
		reader := strings.NewReader(vcf)

		// Upload vcf file to storage
		vcfType := "text/vcard"
		vcfFilename := cardInfo.Slug + ".vcf"
		_, err = client.Storage.UploadFile(vcfBucket, vcfFilename, reader, storage_go.FileOptions{ContentType: &vcfType})
		if err != nil {
			fmt.Println("Error uploading vcf card", err)
			os.Exit(1)
		}
		vcfUrl := client.Storage.GetPublicUrl(vcfBucket, vcfFilename).SignedURL

		cardInfo.VcfUrl = vcfUrl

		// Insert card info into database
		_, _, err = client.From("cards").Insert(cardInfo, false, "", "", "exact").Execute()
		if err != nil {
			fmt.Println("Error inserting data", err)
			os.Exit(1)
		}

	}).Run()

	fmt.Println("Tajeta creada con éxito!")

	// Form to ask the user if they want to deploy the app
	deployForm, deploy := deploy.CreateForm()

	err = deployForm.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// If the user doesn't want to deploy, exit the program
	if !deploy.Deploy {
		fmt.Println("Gracias por usar la aplicación!")
		os.Exit(0)
	}

	// Spinner to show the user that the app is being deployed
	_ = spinner.New().Title("Ejecutando build y deploy...").Action(func() {

		cmd := exec.Command("sh", "-c", "cd ~/Developer/Github/Astro/smart-cards && pnpm build && wrangler pages deploy dist --project-name=smart-cards")

		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println("Error deploying", err)
			os.Exit(1)
		}

		fmt.Println(string(stdout))

	}).Run()

	fmt.Println("\nAplicación desplegada con éxito!")

}
