package goose

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"gopkg.in/fatih/set.v0"
)

var (
	initErr           error
	stopwordsFromFile = StopWords{}
	punctuationRegex  = regexp.MustCompile(`[^\p{Ll}\p{Lu}\p{Lt}\p{Lo}\p{Nd}\p{Pc}\s]`)
)

const path = "resources/stopwords"

func init() {
	stopwordsFromFile, initErr = NewStopwords()
	if initErr != nil {
		panic(initErr.Error())
	}
}

// StopWords implements a simple language detector
type StopWords struct {
	cachedStopWords map[string]*set.Set
}

// NewStopwords returns an instance of a stop words detector
// new stopword lists can be added to the "resources/stopwords" directory as .txt with the filename
// prefixed as "ISOLangCode_language_stopwords.txt"
// ie. en_english_stopwords.txt
func NewStopwords() (StopWords, error) {
	cachedStopWords := make(map[string]*set.Set)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return StopWords{}, fmt.Errorf("could not load goose stopwords from files: %v", err)
	}
	for _, file := range files {
		name := file.Name()
		nlist := strings.Split(name, "_")
		if nlist[len(nlist)-1] != "stopwords.txt" {
			continue
		}
		lang := nlist[0]
		lines := ReadLinesOfFile(path + "/" + file.Name())
		cachedStopWords[lang] = set.New()
		for _, line := range lines {
			if strings.HasPrefix(line, "#") {
				continue
			}
			line = strings.TrimSpace(line)
			cachedStopWords[lang].Add(line)
		}
	}
	return StopWords{
		cachedStopWords: cachedStopWords,
	}, nil
}

func (stop *StopWords) GetStopWordsByLanguage(lang string) *set.Set {
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
	stopWords := set.New()
	text = strings.ToLower(text)
	items := strings.Split(text, " ")
	stops := stop.cachedStopWords[lang]
	count := 0
	if stops != nil {
		for _, item := range items {
			if stops.Has(item) {
				stopWords.Add(item)
				count++
			}
		}
	}

	ws.stopWordCount = stopWords.Size()
	ws.wordCount = len(items)
	ws.stopWords = stopWords

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

// ReadLinesOfFile returns the lines from a file as a slice of strings
func ReadLinesOfFile(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Println(err.Error())
	}
	lines := strings.Split(string(content), "\n")
	return lines
}
