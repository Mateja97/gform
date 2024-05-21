# Google Forms Image Uploader (gform)

This project creates a Google Form with multiple-choice questions for images uploaded from a specified folder. The form is created using the Google Forms API and images are uploaded to Google Drive and made publicly accessible.

## Features

- Create a Google Form with a title.
- Upload images from a specified folder to Google Drive.
- Make uploaded images publicly accessible.
- Add multiple-choice questions for each image in the Google Form.
- Share the form with a specified Google account.
- Print both the edit URL and public URL of the created form.

## Prerequisites

- Go (version 1.16 or later)
- Google Cloud project with the following APIs enabled:
    - Google Forms API
    - Google Drive API
- Service account credentials JSON file

## Setup

1. **Enable APIs in Google Cloud Console:**

    - Go to the [Google Cloud Console](https://console.cloud.google.com/).
    - Select your project or create a new project.
    - Enable the Google Forms API and Google Drive API for your project.

2. **Create a Service Account:**

    - Navigate to **APIs & Services** > **Credentials**.
    - Click on **Create credentials** and select **Service account**.
    - Fill in the service account details and create the service account.
    - Go to the service account and create a key in JSON format.
    - Download the JSON key file and save it as `credentials.json` in your project directory.

3. **Clone the Repository:**

   ```bash
   git clone https://github.com/your-username/gform.git
   cd gform

4. **Usage:**
    - Run the Program:
    ```bash
    go run main.go -credentials=credentials.json -folder=images -email=your_email@gmail.com
   ```
    Replace your_email@gmail.com with your personal Google account email.

5. **Example**
    - The program will print the edit URL and public URL of the created form:
   ```bash
    Form created with ID: 1dQw4w9WgXcQ
    Edit the form at: https://docs.google.com/forms/d/1dQw4w9WgXcQ/edit
    Public form URL: https://docs.google.com/forms/d/e/1dQw4w9WgXcQ/viewform