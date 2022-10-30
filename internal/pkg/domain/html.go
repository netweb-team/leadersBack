package domain

type SearchPage struct {
	Links []string `json:"adlinks"`
	Count int      `json:"count"`
}

type AdPage struct {
	Latitude    string `json:"lat"`
	Longitude   string `json:"lng"`
	Address     string `json:"addr"`
	Price       string `json:"price"`
	Floor       string `json:"current_floor"`
	TotalArea   string `json:"total_area"`
	KitchenArea string `json:"kitchen_area"`
	Balcony     string `json:"balcony"`
	Renovation  string `json:"renovation,omitempty"`
	Year        string `json:"year,omitempty"`
	Metro       string `json:"metro"`
}

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
