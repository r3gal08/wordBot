package dictionary

type WordData []struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic"`
	Phonetics []struct {
		Text  string `json:"text"`
		Audio string `json:"audio,omitempty"`
	} `json:"phonetics"`
	Origin   string `json:"origin"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string `json:"definition"`
			Example    string `json:"example"`
			Synonyms   []any  `json:"synonyms"`
			Antonyms   []any  `json:"antonyms"`
		} `json:"definitions"`
	} `json:"meanings"`
}

// Can add in additional fields to the struct as needed
type WordResponse struct {
	Word         string `json:"word,omitempty"`
	Definition   string `json:"definition,omitempty"`
	PartOfSpeech string `json:"partofspeech,omitempty"`
	ConfidenceRating int `json:"confidencerating,omitempty"`
}