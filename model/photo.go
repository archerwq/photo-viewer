package model

type Photo struct {
	Sha1     string   `json:"sha1"`
	Path     string   `json:"path"`
	Time     string   `json:"time"`
	Width    int64    `json:"width"`
	Length   int64    `json:"length"`
	Location string   `json:"location"`
	Type     string   `json:"type"`
	Tags     []string `json:"tags"`
	Story    string   `json:"story"`
}
