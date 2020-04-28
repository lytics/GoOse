package goose

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/lytics/GoOse/resources/stopwords"
	"gopkg.in/fatih/set.v0"
)

var punctuationRegex = regexp.MustCompile(`[^\p{Ll}\p{Lu}\p{Lt}\p{Lo}\p{Nd}\p{Pc}\s]`)

// StopWords implements a simple language detector
type StopWords struct {
	cachedStopWords map[string]*set.Set
}

// NewStopwords returns an instance of a stop words detector
func NewStopwords() StopWords {
	cachedStopWords := make(map[string]*set.Set)
	for lang, stopwords := range sw {
		lines := strings.Split(stopwords, "\n")
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
	}
}

/*
func NewStopwords(path string) StopWords {
	cachedStopWords := make(map[string]*set.Set)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err.Error())
	}
	for _, file := range files {
		name := strings.Replace(file.Name(), ".txt", "", -1)
		name = strings.Replace(name, "stopwords-", "", -1)
		name = strings.ToLower(name)

		stops := set.New()
		lines := ReadLinesOfFile(path + "/" + file.Name())
		for _, line := range lines {
			line = strings.Trim(line, " ")
			stops.Add(line)
		}
		cachedStopWords[name] = stops
	}

	return StopWords{
		cachedStopWords: cachedStopWords,
	}
}
*/

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

// SimpleLanguageDetector returns the language code for the text, based on its stop words
func (stop StopWords) SimpleLanguageDetector(text string) string {
	max := 0
	currentLang := "en"

	for k := range sw {
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

var sw = map[string]string{
	"ar": stopwords.ArabicStopwords,
	"bg": stopwords.BulgarianStopwords,
	"bn": stopwords.BengaliStopwords,
	"ca": stopwords.CatalanStopwords,
	"cs": stopwords.CzechStopwords,
	"da": stopwords.DanishStopwords,
	"de": stopwords.GermanStopwords,
	"el": stopwords.GreekStopwords,
	"en": stopwords.EnglishStopwords,
	"es": stopwords.SpanishStopwords,
	"eu": stopwords.BasqueStopwords,
	"fa": stopwords.PersianFarsiStopwords,
	"fi": stopwords.FinnishStopwords,
	"fr": stopwords.FrenchStopwords,
	"ga": stopwords.IrishStopwords,
	"gl": stopwords.GalicianStopwords,
	"he": stopwords.HebrewStopwords,
	"hi": stopwords.HindiStopwords,
	"hr": stopwords.CroatianStopwords,
	"hu": stopwords.HungarianStopwords,
	"hy": stopwords.ArmenianStopwords,
	"id": stopwords.IndonesianStopwords,
	"it": stopwords.ItalianStopwords,
	"ja": stopwords.JapaneseStopwords,
	"ko": stopwords.KoreanStopwords,
	"ku": stopwords.KurdishStopwords,
	"lt": stopwords.LithuanianStopwords,
	"lv": stopwords.LatvianStopwords,
	"mr": stopwords.MarathiStopwords,
	"nb": stopwords.NorwegianBokmalStopwords,
	"nl": stopwords.DutchStopwords,
	"no": stopwords.NorwegianStopwords,
	"pl": stopwords.PolishStopwords,
	"pt": stopwords.PortugueseStopwords,
	"ro": stopwords.RomanianStopwords,
	"ru": stopwords.RussianStopwords,
	"sk": stopwords.SlovakStopwords,
	"sr": stopwords.SerbianCyrillicStopwords,
	"sv": stopwords.SwedishStopwords,
	"th": stopwords.ThaiStopwords,
	"tr": stopwords.TurkishStopwords,
	"uk": stopwords.UkrainianStopwords,
	"ur": stopwords.UrduStopwords,
	"vi": stopwords.VietnameseStopwords,
	"zh": stopwords.ChineseStopwords,
}
