package domain

type SearchPage struct {
	Links []string `json:"adlink"`
	Count int      `json:"count"`
}
