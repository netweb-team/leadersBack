package domain

import "time"

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
	AvgCost   float64 `json:"ap,omitempty"`
	Latitude  float64 `json:"lat,omitempty"`
	Longitude float64 `json:"lng,omitempty"`
	Good      int     `json:"good,omitempty"`
}

type Table struct {
	ID   int               `json:"id"`
	Path string            `json:"path"`
	Rows []*Row            `json:"table,omitempty"`
	Flat int               `json:"count,omitempty"`
	Time time.Time         `json:"time,omitempty"`
	PA   []*PatternAnalogs `json:"pa,omitempty"`
}

type CorrectCoefs struct {
	Sale    float64 `json:"sale"`
	Floor   float64 `json:"floor"`
	Total   float64 `json:"total"`
	Kitchen float64 `json:"kitchen"`
	Balcony float64 `json:"balcony"`
	Metro   float64 `json:"metro"`
	State   float64 `json:"state"`
}
