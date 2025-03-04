package dictionary

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiURL = "https://api.dictionaryapi.dev/api/v2/entries/en/%s"

func GetWordData(word string) (WordData, error) {
	url := fmt.Sprintf(apiURL, word)

	response, err := http.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("API returned status code: %d", response.StatusCode)
		}
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Load word data into memory
	var wordData WordData
	if err := json.Unmarshal(body, &wordData); err != nil || len(wordData) == 0 {
		if err == nil {
			err = fmt.Errorf("word data not found for: %s", word)
		}
		return nil, err
	}

	return wordData, nil
}
