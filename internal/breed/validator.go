package breed

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CatBreedResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func ValidateCatBreed(breed string) (bool, error) {
	breed = strings.ReplaceAll(breed, " ", "%20")

	url := fmt.Sprintf("https://api.thecatapi.com/v1/breeds/search?q=%s&attach_image=0", breed)

	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("breed: failed to make request to TheCatAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("breed: received non-OK response: %d", resp.StatusCode)
	}

	var breeds []CatBreedResponse
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		return false, fmt.Errorf("breed: failed to decode JSON response: %w", err)
	}

	if len(breeds) == 0 {
		return false, nil
	}

	return true, nil
}
