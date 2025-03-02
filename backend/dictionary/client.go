package dictionary

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiURL = "https://api.dictionaryapi.dev/api/v2/entries/en/%s"

func GetWordDefinition(word string) (string, error) {
	url := fmt.Sprintf(apiURL, word)

	response, err := http.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		if err == nil {
			err = fmt.Errorf("API returned status code: %d", response.StatusCode)
		}
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var definitionResponse DefinitionResponse
	if err := json.Unmarshal(body, &definitionResponse); err != nil {
		return "", err
	}

	if len(definitionResponse) > 0 && len(definitionResponse[0].Meanings) > 0 && len(definitionResponse[0].Meanings[0].Definitions) > 0 {
		return definitionResponse[0].Meanings[0].Definitions[0].Definition, nil
	}

	return "", fmt.Errorf("Definition not found for word: %s", word)
}
