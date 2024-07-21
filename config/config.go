package config

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"io"
	"log"
	"os"
	"time"
)

// Config 结构体
type Config struct {
	Redis struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		DB       int    `json:"db"` // 添加数据库选择字段
	} `json:"redis"`

	MySQL struct {
		DSN string `json:"dsn"`
	} `json:"mysql"`

	MongoDB struct {
		Addr string `json:"addr"`
		DB   string `json:"db"`
	} `json:"mongodb"`

	OrderStatus struct {
		// 成功状态
		SuccessStatus int `json:"success_status"`
		// 失败状态
		FailStatus int `json:"fail_status"`
		// 本地状态
		LocalStatus int `json:"local_status"`
		// 提交到三方的状态
		CommitStatus int `json:"commit_status"`
	} `json:"order_status"`
}

// LoadConfig 从配置文件中加载配置
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开配置文件: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, fmt.Errorf("无法解析配置文件: %v", err)
	}

	return config, nil
}

// CustomLogger 结构体用于自定义日志输出
type CustomLogger struct {
	writer io.Writer
}

// Write 实现 io.Writer 接口，用于将日志消息写入指定的 writer
func (c *CustomLogger) Write(p []byte) (n int, err error) {
	return c.writer.Write(p)
}

func GetLog() {
	// 获取当前日期，并格式化为字符串
	currentDate := time.Now().Format("2006-01-02")
	// 设置总日志文件
	totalLogFile := &lumberjack.Logger{
		Filename:   "logs/total.log",
		MaxSize:    500,  // 每个日志文件的最大大小（以MB为单位）
		MaxBackups: 3,    // 保留的旧日志文件的最大数量
		MaxAge:     28,   // 保留旧文件的最大天数
		Compress:   true, // 是否压缩/归档旧文件
	}

	// 设置每日日志文件
	dailyLogFile := &lumberjack.Logger{
		Filename:   "logs/daily_" + currentDate + ".log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     1,
		Compress:   true,
	}

	// 创建一个多路复用器，将日志输出到多个目的地
	multiWriter := io.MultiWriter(os.Stdout, totalLogFile, dailyLogFile)

	// 创建自定义的日志记录器
	customLogger := &CustomLogger{writer: multiWriter}

	// 配置 Gin 使用自定义日志记录器
	gin.DefaultWriter = customLogger
	gin.DefaultErrorWriter = customLogger
	// 设置 Gin 日志级别为 debug
	gin.SetMode(gin.DebugMode)

	// 配置日志格式
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
