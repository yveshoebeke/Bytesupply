package googleapi

import "os"

var (
	// SearchAPIURL ...
	SearchAPIURL = "https://www.googleapis.com/customsearch/v1?key={KEY}&cx=017576662512468239146:omuauf_lfve&q={SUBJECT}"
	// SearchAPIKey ...
	SearchAPIKey = os.Getenv("BS_GOOGLE_SEARCH_API_KEY")
	// SearchLimit ...
	SearchLimit = "10"
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
