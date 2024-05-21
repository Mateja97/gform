package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

var (
	credentialsFile = flag.String("credentials", "credentials.json", "Path to the Google Cloud credentials JSON file")
	imageFolder     = flag.String("folder", "images", "Path to the folder containing images")
	personalEmail   = flag.String("email", "", "Your personal Google account email")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	// Load Google Cloud credentials
	creds, err := ioutil.ReadFile(*credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Use OAuth2 for authentication
	config, err := google.JWTConfigFromJSON(creds, drive.DriveFileScope, "https://www.googleapis.com/auth/forms.body")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := config.Client(ctx)

	// Initialize Google Forms API client
	formService, err := forms.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Forms service: %v", err)
	}

	// Initialize Google Drive API client
	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Drive service: %v", err)
	}

	// Create a new form with just the title
	form := &forms.Form{
		Info: &forms.Info{
			Title: "Image Selection Form",
		},
	}

	createdForm, err := formService.Forms.Create(form).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to create form: %v", err)
	}

	formID := createdForm.FormId
	editURL := fmt.Sprintf("https://docs.google.com/forms/d/%s/edit", formID)
	fmt.Printf("Form created with ID: %s\n", formID)

	// Share the form with your personal Google account
	err = shareFormWithEmail(ctx, driveService, formID, *personalEmail)
	if err != nil {
		log.Fatalf("Unable to share form: %v", err)
	}

	// Read images from the specified folder
	imageFiles, err := ioutil.ReadDir(*imageFolder)
	if err != nil {
		log.Fatalf("Unable to read image folder: %v", err)
	}

	// Prepare requests to add items to the form
	var requests []*forms.Request
	index := 0

	for _, file := range imageFiles {
		if !file.IsDir() {
			imagePath := filepath.Join(*imageFolder, file.Name())
			imageContent, err := os.Open(imagePath)
			if err != nil {
				log.Printf("Failed to read image %s: %v", imagePath, err)
				continue
			}
			defer imageContent.Close()

			// Upload image to Google Drive and get the URL
			imageURL, err := uploadAndMakeImagePublic(ctx, driveService, imageContent)
			if err != nil {
				log.Printf("Failed to upload image %s: %v", imagePath, err)
				continue
			}

			// Create a request to add an image and multiple-choice question
			request := &forms.Request{
				CreateItem: &forms.CreateItemRequest{
					Item: &forms.Item{
						Title: "Select an option for the image below",
						QuestionItem: &forms.QuestionItem{
							Image: &forms.Image{
								SourceUri: imageURL,
							},
							Question: &forms.Question{
								Required: true,
								ChoiceQuestion: &forms.ChoiceQuestion{
									Options: []*forms.Option{
										{Value: "Option 1"},
										{Value: "Option 2"},
									},
									Type: "RADIO",
								},
							},
						},
					},
					Location: &forms.Location{
						Index:           int64(index),
						ForceSendFields: []string{"Index"},
					},
				},
			}
			requests = append(requests, request)
			index++
		}
	}

	// Perform a batch update to add items to the form
	_, err = formService.Forms.BatchUpdate(formID, &forms.BatchUpdateFormRequest{
		Requests: requests,
	}).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Unable to update form: %v", err)
	}

	fmt.Printf("Form created and can be edited at: %s\n", editURL)
	publicURL := createdForm.ResponderUri
	fmt.Printf("Public form URL: %s\n", publicURL)
}

func uploadAndMakeImagePublic(ctx context.Context, driveService *drive.Service, imageFile *os.File) (string, error) {
	file, err := driveService.Files.Create(&drive.File{
		Name: filepath.Base(imageFile.Name()),
	}).Media(imageFile).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %v", err)
	}

	// Make the uploaded file public
	_, err = driveService.Permissions.Create(file.Id, &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to make image public: %v", err)
	}

	// Return the public URL of the image
	imageURL := fmt.Sprintf("https://drive.google.com/uc?id=%s", file.Id)
	return imageURL, nil
}

func shareFormWithEmail(ctx context.Context, driveService *drive.Service, fileId, email string) error {
	permission := &drive.Permission{
		Type:         "user",
		Role:         "writer",
		EmailAddress: email,
	}
	_, err := driveService.Permissions.Create(fileId, permission).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to share form with email: %v", err)
	}
	return nil
}
