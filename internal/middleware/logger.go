package middleware

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sanmu2018/word-hero/log"
)

func LoggerHandler(c *gin.Context) {
	// Start timer
	path := c.Request.URL.Path
	start := time.Now()
	c.Next()
	stop := time.Since(start)
	//latency := int(math.Ceil(float64(stop.Nanoseconds()) / 1000.0))
	statusCode := c.Writer.Status()
	clientIP := c.ClientIP()
	clientUserAgent := c.Request.UserAgent()
	referer := c.Request.Referer()
	dataLength := c.Writer.Size()
	if dataLength < 0 {
		dataLength = 0
	}
	log.Info().Str("method", c.Request.Method).
		Int("statusCode", statusCode).
		Str("path", path).
		Str("client_ip", clientIP).
		Str("cost", fmt.Sprintf("%v", stop)).
		Str("referer", referer).
		Int("data_length", dataLength).
		Str("user_agent", clientUserAgent).
		Send()
	//log.Info().Fields(map[string]interface{}{
	//	"method":     c.Request.Method,
	//	"path":       path,
	//	"statusCode": statusCode,
	//	"cost":       fmt.Sprintf("(%dμs)", latency),
	//	"clientIP":   clientIP,
	//	"referer":    referer,
	//	"dataLength": dataLength,
	//	"userAgent":  clientUserAgent,
	//}).Send()
	if len(c.Errors) > 0 {
		log.Error(errors.New(c.Errors.ByType(gin.ErrorTypePrivate).String())).Send()
	} else {
		if statusCode > 499 {
			log.Error(errors.New("statusCode error")).Msgf("status_code:%d", statusCode)
		} else if statusCode > 399 {
			log.Warn().Int("status_code", statusCode).Send()
		} else if statusCode == 200 {
		} else {
			log.Info().Int("status_code", statusCode).Send()
		}
	}
}
