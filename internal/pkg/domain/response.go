package domain

type Response struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body,omitempty"`
}

type PatternAnalogs struct {
	Pattern *Row            `json:"pattern"`
	Analogs []*Row          `json:"analogs"`
	Correct []*CorrectCoefs `json:"coefs"`
}
