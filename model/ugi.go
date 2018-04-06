package model

type UGI struct {
	Sha1      string `json:"sha1"`
	Tags      string `json:"tags"`
	Story     string `json:"story"`
	UpdatedOn int64  `json:"updatedOn"`
}
