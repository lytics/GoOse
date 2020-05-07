package goose

//some word statistics
type wordStats struct {
	//total number of stopwords or good words that we can calculate
	stopWordCount int
	//total number of words on a node
	wordCount int
}

func (w *wordStats) getStopWordCount() int {
	return w.stopWordCount
}

func (w *wordStats) setStopWordCount(stopWordCount int) {
	w.stopWordCount = stopWordCount
}

func (w *wordStats) getWordCount() int {
	return w.wordCount
}

func (w *wordStats) setWordCount(wordCount int) {
	w.wordCount = wordCount
}
