package image

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// GoogleSearchResponse is used to parse the JSON response from the Google Custom Search API
type GoogleSearchResponse struct {
	Items []struct {
		Link string `json:"link"`
	} `json:"items"`
}

/*
SearchSmartphoneImage searches for a smartphone image based on the brand and model.
This function sends a request to the Google Custom Search API and retrieves the first image URL found.
- brand: The smartphone brand name (e.g., "Samsung", "Apple")
- model: The smartphone model (e.g., "Galaxy S24", "iPhone 15")
- Returns the image URL if found or an error if the request fails.
*/
func SearchSmartphoneImage(brand, model string) (string, error) {
	APIKey := os.Getenv("GOOGLE_API_KEY")
	CX := os.Getenv("SEARCH_ENGINE_ID")
	CUSTOM_URL := os.Getenv("CUSTOM_URL")

	if APIKey == "" || CX == "" {
		return "", errors.New("API key and search engine ID not found")
	}

	query := fmt.Sprintf("%s %s smartphone", brand, model)
	params := url.Values{}
	params.Add("q", query)
	params.Add("key", APIKey)
	params.Add("cx", CX)
	params.Add("searchType", "image")
	params.Add("num", "1") // Retrieve only one image

	// Construct the request URL
	requestURL := fmt.Sprintf("%s?%s", CUSTOM_URL, params.Encode())

	// Send HTTP GET request
	resp, err := http.Get(requestURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the JSON response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the JSON response
	var result GoogleSearchResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	// Ensure at least one image result is available
	if len(result.Items) > 0 {
		return result.Items[0].Link, nil
	}

	return "", errors.New("no image found")
}
