package image

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func SearchSmartphoneImage(brand, model string) (string, error) {
	// Define the local directory where images are stored
	localDir := "data/Smartphone_Scrapping"

	// Sanitize the brand and model names to match the file naming convention
	fileName := fmt.Sprintf("%s_%s.jpg", sanitizeFileName(brand), sanitizeFileName(model))

	// Construct the full path to the image file
	imagePath := filepath.Join(localDir, fileName)

	// Check if the file exists
	if _, err := os.Stat(imagePath); err != nil {
		if os.IsNotExist(err) {
			log.Printf("No image found in local directory for: %s %s\n", brand, model)
			return "", nil
		}
		log.Println("Error checking file:", err)
		return "", err
	}

	// Return the backend URL for the image
	backendURL := fmt.Sprintf("/static/%s", fileName)
	log.Println("Image found! URL:", backendURL)
	return backendURL, nil
}

func sanitizeFileName(name string) string {
	return strings.ReplaceAll(strings.TrimSpace(name), " ", "_")
}
