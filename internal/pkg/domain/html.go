package domain

type SearchPage struct {
	Links []string `json:"adlink"`
	Count int      `json:"count"`
}

type AdPage struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}
