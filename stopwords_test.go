package goose

import (
	"testing"
)

func TestStopwordList(t *testing.T) {
	// getting stopword set by language
	stopWords := NewStopwords()
	stops := stopWords.GetStopWordsByLanguage("en")
	if stops == nil {
		t.Errorf("Could not get language set")
	}

	// getting word stats
	text := "Why is this such a long thought. Will it ever stop or will it continue?"
	swStats := stopWords.stopWordsCount("en", text)
	got := swStats.stopWordCount
	if got != 10 {
		t.Errorf("Stopword count in text = %d; want 8", got)
	}
	got = swStats.wordCount
	if got != 15 {
		t.Errorf("Word count in text = %d; want 15", got)
	}

	// language detection
	text = "Mais le week-end, il est assez différent. Pendant le week-end, nous ne sommes pas très occupés comme les autres jours. Le samedi matin, mon père qui est très sportif fait de la natation, et ma mère fait la cuisine parce que chaque samedi, mes parents invitent ma tante à dîner avec nous."
	gotLang := stopWords.SimpleLanguageDetector(text)
	if gotLang != "fr" {
		t.Errorf("Detected language = %s; want fr", gotLang)
	}
}
