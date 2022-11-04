package domain

type Row struct {
	ID        int     `json:"id,omitempty"`
	Address   string  `json:"a"`
	Rooms     string  `json:"r"`
	Segment   string  `json:"s"`
	Floors    uint    `json:"f"`
	Walls     string  `json:"w"`
	CFloor    uint    `json:"cf"`
	Total     float64 `json:"t"`
	Kitchen   float64 `json:"k"`
	Balcony   string  `json:"b"`
	Metro     float64 `json:"m"`
	State     string  `json:"st"`
	Cost      int     `json:"p,omitempty"`
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"lng,omitempty"`
}

type Table struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
	Rows []*Row `json:"table"`
}
