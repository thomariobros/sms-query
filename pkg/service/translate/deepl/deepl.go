package deepl

import (
	"encoding/json"

	"sms-query/pkg/http"
)

const (
	deepLRootURL = "https://www.deepl.com/jsonrpc"
)

type DeepLTranslateService struct {
	rootURL string
}

func NewDeepLTranslateService() *DeepLTranslateService {
	return &DeepLTranslateService{
		rootURL: deepLRootURL,
	}
}

func NewDeepLTranslateServiceRootURL(rootURL string) *DeepLTranslateService {
	return &DeepLTranslateService{
		rootURL: rootURL,
	}
}

// Translate call deepL translate api to translate text
func (s *DeepLTranslateService) Translate(source string, target string, text string) (string, error) {
	data := `{
    "jsonrpc": "2.0",
    "method": "LMT_handle_jobs",
    "params": {
      "jobs": [
        {
          "kind": "default",
          "raw_en_sentence": "` + text + `"
        }
      ],
      "lang": {
        "source_lang": "` + source + `",
        "target_lang": "` + target + `"
      }
    }
  }`
	resp, err := http.Post(s.rootURL, "application/json", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var response response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}
	var result []string
	for _, beam := range response.Result.Translations[0].Beams {
		result = append(result, beam.Text)
	}

	return result[0], nil
}

type response struct {
	Result result `json:"result"`
}

type result struct {
	Translations []translation `json:"translations"`
}

type translation struct {
	Beams []beam `json:"beams"`
}

type beam struct {
	Text string `json:"postprocessed_sentence"`
}
