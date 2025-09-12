package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// InitLogger 初始化日志配置
func InitLogger() {
	// 设置时间格式为 2025-09-02 07:30:29 格式
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	
	// 设置日志级别
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	
	// 启用调用者信息（文件名和行号）
	log.Logger = log.With().Timestamp().Caller().Logger()
	
	// 设置字段顺序：level, timestamp, caller, message
	zerolog.TimestampFieldName = "time"
	zerolog.LevelFieldName = "level"
	zerolog.MessageFieldName = "message"
	zerolog.CallerFieldName = "caller"
	
	// 统一使用JSON格式输出
	log.Logger = log.Output(os.Stdout)
}

// GetLogger 获取日志实例
func GetLogger() zerolog.Logger {
	return log.Logger
}