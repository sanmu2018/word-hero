package pke

type AnyResponse[R any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data R      `json:"data"`
}
type CommResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
type BaseListResp struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}
type APIResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg,omitempty"`
}

func (ar *APIResponse) Error() string {
	return ar.Msg
}

func (ar *APIResponse) ErrorNo() int {
	return ar.Code
}

func NewApiError(code int) *APIResponse {
	return &APIResponse{Code: code, Msg: GetErrorMessage(code)}
}
