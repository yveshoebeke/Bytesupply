package googleapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	// SEARCHAPIURL ...
	SEARCHAPIURL = "https://www.googleapis.com/customsearch/v1?key={KEY}&cx=017576662512468239146:omuauf_lfve&q={SUBJECT}"
	// SEARCHLIMIT ...
	SEARCHLIMIT = "10"
)

var (
	searchAPIKey = os.Getenv("BS_GOOGLE_SEARCH_API_KEY")
)

//APIResult search result(s)
type APIResult struct {
	Kind string `json:"kind"`
	URL  struct {
		Type     string
		Template string
	}
	Queries struct {
		Request []struct {
			Title          string
			TotalResults   string
			SearchTerms    string
			count          int
			StartIndex     int
			InputEncoding  string
			OutputEncoding string
			Safe           string
			Cx             string
		}
		NextPage []struct {
			Title          string
			TotalResults   string
			SearchTerms    string
			count          int
			StartIndex     int
			InputEncoding  string
			OutputEncoding string
			Safe           string
			Cx             string
		}
	}
	Context struct {
		Title  string
		Facets []struct {
			Anchor      string
			Label       string
			LabelWithOp string
		}
	}
	SearchInformation struct {
		SearchTime            float32
		FormattedSearchTime   float32
		TotalResults          string
		FormattedTotalResults string
	}
	Items []struct {
		Kind             string
		Title            string
		HTMLTitle        string
		Link             string
		DisplayLink      string
		Snippet          string
		HTMLSnippet      string
		CacheID          string
		FormattedURL     string
		HTMLFormattedURL string
		Mime             string
		FileFormat       string
	}
}

// GetSearchResults ...
func GetSearchResults(searchKey string) (APIResult, error) {
	var googleapi APIResult
	searchURL := SEARCHAPIURL
	searchURL = strings.Replace(searchURL, "{KEY}", searchAPIKey, 1)
	searchURL = strings.Replace(searchURL, "{SUBJECT}", searchKey, 1)
	searchURL = strings.Replace(searchURL, "{NUM}", SEARCHLIMIT, 1)

	resp, err := http.Get(searchURL)
	if err != nil {
		return googleapi, err
	}
	//app.log.Println("GET", searchURL)

	body, ioerr := ioutil.ReadAll(resp.Body)
	if ioerr != nil {
		return googleapi, err
	}
	defer resp.Body.Close()

	umerr := json.Unmarshal(body, &googleapi)
	if umerr != nil {
		fmt.Println("Unmarchal error:", umerr)
		return googleapi, err
	}

	return googleapi, nil
}
