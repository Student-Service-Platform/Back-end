package utils

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger

func init() {
	// 创建日志文件夹
	logDir := "/log"
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}

	// 创建日志文件
	logFile := filepath.Join(logDir, "log.log")
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("无法创建日志文件")
	}

	// 初始化 zerolog
	logger = zerolog.New(file).With().Timestamp().Logger()
}

// LogError 记录错误信息和发生时间
func LogError(err error) {
	if err != nil {
		logger.Error().Time("time", time.Now()).Err(err).Msg("发生错误")
	}
}

// GetLogger 返回初始化的 logger
func GetLogger() zerolog.Logger {
	return logger
}

// 上述工具用法：utils.LogError(err)，就可以了
