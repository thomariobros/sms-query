package qwant

import (
	"encoding/json"
	"net/url"
	"strings"

	"main/pkg/http"
	"main/pkg/service/search"
)

const (
	qwantRootURL = "https://api.qwant.com/egp/search/web"
)

type QwantSearchService struct {
	rootURL string
}

func NewQwantSearchService() *QwantSearchService {
	return &QwantSearchService{
		rootURL: qwantRootURL,
	}
}

func NewQwantSearchServiceRootURL(rootURL string) *QwantSearchService {
	return &QwantSearchService{
		rootURL: rootURL,
	}
}

// Search Qwant search
func (s *QwantSearchService) Search(locale string, text string) (*search.SearchResult, error) {
	params := "q=" + url.QueryEscape(text) +
		"&locale=" + strings.ToLower(locale) +
		"&count=1"
	resp, err := http.Get(s.rootURL + "?" + params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	if len(response.Data.Result.Items) > 0 {
		return &search.SearchResult{
			Source: response.Data.Result.Items[0].URL, // url is better
			URL:    response.Data.Result.Items[0].URL,
			Text:   response.Data.Result.Items[0].Desc,
		}, nil
	}
	return nil, nil
}

type response struct {
	Data data `json:"data"`
}

type data struct {
	Result result `json:"result"`
}

type result struct {
	Items []item `json:"items"`
}

type item struct {
	URL  string `json:"url"`
	Desc string `json:"desc"`
}
