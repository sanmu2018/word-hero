package pke

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
