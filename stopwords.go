package goose

import (
	"regexp"
	"strings"

	"github.com/willf/bloom"
)

var (
	punctuationRegex = regexp.MustCompile(`[^\p{Ll}\p{Lu}\p{Lt}\p{Lo}\p{Nd}\p{Pc}\s]`)
)

// StopWords implements a simple language detector
type StopWords struct {
	cachedStopWords map[string]*bloom.BloomFilter
}

// NewStopwords returns an instance of a stop words detector
// new stopword lists can be added to the "resources/stopwords" directory as .txt with the filename
// prefixed as "ISOLangCode_language_stopwords.txt"
// ie. en_english_stopwords.txt
func NewStopwords() StopWords {
	return StopWords{
		cachedStopWords: filters,
	}
}

func (stop *StopWords) GetStopWordsByLanguage(lang string) *bloom.BloomFilter {
	s, ok := stop.cachedStopWords[lang]
	if ok {
		return s
	}
	return nil
}

func (stop *StopWords) removePunctuation(text string) string {
	return punctuationRegex.ReplaceAllString(text, "")
}

func (stop *StopWords) stopWordsCount(lang string, text string) wordStats {
	if text == "" {
		return wordStats{}
	}
	ws := wordStats{}
	text = strings.ToLower(text)
	items := strings.Split(text, " ")
	stops := stop.cachedStopWords[lang]
	count := 0
	if stops != nil {
		for _, item := range items {
			if stops.Test([]byte(item)){
				count++
			}
		}
	}

	ws.stopWordCount = count
	ws.wordCount = len(items)
	return ws
}

// SimpleLanguageDetector returns the language code for the text, based on its stop words; defaults to "en"
func (stop StopWords) SimpleLanguageDetector(text string) string {
	max := 0
	currentLang := "en"

	for k := range stop.cachedStopWords {
		ws := stop.stopWordsCount(k, text)
		if ws.stopWordCount > max {
			max = ws.stopWordCount
			currentLang = k
		}
	}

	return currentLang
}
