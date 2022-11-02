package domain

type Row struct {
	Address string  `json:"a"`
	Rooms   uint    `json:"r"`
	Segment string  `json:"s"`
	Floors  uint    `json:"f"`
	Walls   string  `json:"w"`
	CFloor  uint    `json:"cf"`
	Total   float64 `json:"t"`
	Kitchen float64 `json:"k"`
	Balcony string  `json:"b"`
	Metro   string  `json:"m"`
	State   string  `json:"st"`
	Cost    float64 `json:"p,omitempty"`
}

type Table struct {
	ID   int    `json:"id"`
	Path string `json:"path"`
	Rows []*Row `json:"table"`
}
