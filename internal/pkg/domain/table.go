package domain

type Row struct {
	Address string `json:"addr"`
}

type Table struct {
	Path string `json:"path"`
	Rows []*Row `json:"table"`
}
