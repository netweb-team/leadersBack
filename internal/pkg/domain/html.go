package domain

type FlatObj struct {
	Content struct {
		Address string `json:"text"`
	} `json:"content"`
}

type Coordinates struct {
	Data struct {
		Points map[string]*FlatObj `json:"points"`
	} `json:"data"`
}

type MapResponse struct {
	Features []struct {
		Center []float64 `json:"center"`
	} `json:"features"`
}
