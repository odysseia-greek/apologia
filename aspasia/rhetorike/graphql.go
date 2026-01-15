package rhetorike

type gqlReq struct {
	Query     string         `json:"query"`
	Variables map[string]any `json:"variables,omitempty"`
}

type alexandrosFuzzyResp struct {
	Data struct {
		Fuzzy struct {
			Results []struct {
				Headword     string `json:"headword"`
				QuickGlosses []struct {
					Language string `json:"language"`
					Gloss    string `json:"gloss"`
				} `json:"quickGlosses"`
			} `json:"results"`
		} `json:"fuzzy"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

const fuzzyMiniQuery = `
query($input: SearchQueryInput!) {
  fuzzy(input: $input) {
    results {
      headword
      quickGlosses { language gloss }
    }
  }
}`
