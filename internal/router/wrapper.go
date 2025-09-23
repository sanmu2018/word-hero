package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanmu2018/word-hero/pkg/pke"
)

type RequestProcessor func(c *gin.Context) (interface{}, error)

func wrapper(handler RequestProcessor) func(c *gin.Context) {
	return func(c *gin.Context) {

		data, err := handler(c)
		resp := pke.APIResponse{Data: data}

		statusCode := http.StatusOK
		if err != nil { //process for error
			statusCode = http.StatusBadRequest
			var h *pke.APIResponse
			if errors.As(err, &h) {
				resp.Code = h.ErrorNo()
				resp.Msg = h.Error()
			} else { //如果不是规范的错误，则统一返回 "系统内部错误: << .msg >>"
				resp.Code = pke.CodeSystemError
				resp.Msg = pke.GetErrorMessage(pke.CodeSystemError)
			}
		}
		c.JSON(statusCode, resp)
	}
}
