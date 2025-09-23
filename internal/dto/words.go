package dto

type WordSearchRequest struct {
	Q string `form:"q"  json:"q"`
	BaseList
}
