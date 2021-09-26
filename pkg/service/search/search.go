package search

type Searcher interface {
	Search(locale string, text string) (SearchResult, error)
}

type SearchResult struct {
	Source string
	URL    string
	Text   string
}
