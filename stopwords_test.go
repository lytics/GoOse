package goose

import (
	"testing"
)

func TestStopwordList(t *testing.T) {
	// making new stopwords from resources/stopwords directory
	newSW, err := NewStopwords()
	if err != nil {
		t.Errorf("Error when getting new stopwords")
	}
	got := len(newSW.cachedStopWords)
	if got < 10 {
		t.Errorf("Number of stopword lists from file = %d; want > 10", got)
	}
	for _, s := range newSW.cachedStopWords {
		got := s.Size()
		if got < 1 {
			t.Errorf("Size of language set = %d; want > 1", got)
		}
	}

	// making new stopwords from variable
	defaultSW := GetDefaultStopwords()
	got = len(defaultSW.cachedStopWords)
	if got < 10 {
		t.Errorf("Number of stopword lists from file = %d; want > 10", got)
	}
	for _, s := range defaultSW.cachedStopWords {
		got := s.Size()
		if got < 1 {
			t.Errorf("Size of language set = %d; want > 1", got)
		}
	}

	// getting stopword set by language
	stops := newSW.GetStopWordsByLanguage("en")
	if stops == nil {
		t.Errorf("Could not get language set")
	}

	// getting word stats
	text := "Why is this such a long thought. Will it ever stop or will it continue?"
	swStats := newSW.stopWordsCount("en", text)
	got = swStats.stopWordCount
	if got != 8 {
		t.Errorf("Stopword count in text = %d; want 8", got)
	}
	got = swStats.wordCount
	if got != 15 {
		t.Errorf("Word count in text = %d; want 15", got)
	}

	// language detection
	text = "Mais le week-end, il est assez différent. Pendant le week-end, nous ne sommes pas très occupés comme les autres jours. Le samedi matin, mon père qui est très sportif fait de la natation, et ma mère fait la cuisine parce que chaque samedi, mes parents invitent ma tante à dîner avec nous."
	gotLang := newSW.SimpleLanguageDetector(text)
	if gotLang != "fr" {
		t.Errorf("Detected language = %s; want fr", gotLang)
	}
}
