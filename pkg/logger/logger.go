package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type Log struct {
	*zap.Logger
}

var initLog = new(InitStruct)

// InitStruct 日志初始化所需参数
type InitStruct struct {
	LogSavePath   string //日志保存路径
	LogFileExt    string //日志文件后缀名
	MaxSize       int    //单个日志文件的大小（M）
	MaxBackups    int    //最大备份数
	MaxAge        int    //最大备份天数
	Compress      bool   //是否压缩过期日志
	LowLevelFile  string //低级别文件名（通常用于记录信息级别和调试级别的日志）
	HighLevelFile string //高级别文件名（通常用于记录错误级别的日志）
	//分级日志可以帮助开发人员更好地查看和分析不同级别的日志信息，便于故障排查和系统优化
}

// NewLogger 初始化 logger
func NewLogger(x *InitStruct, level string) *Log {
	initLog = x
	//定义不同级别的日志处理逻辑
	//error 级别
	highPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})
	//info 和 debug 级别（debug 是最低级别）
	lowPriority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level < zapcore.ErrorLevel && level >= zapcore.DebugLevel
	})
	//多个日志文件
	var cores []zapcore.Core
	//获取低级别日志写入器
	lowFileWriteSyncer := getLogWriter(initLog.LogSavePath + initLog.LowLevelFile + initLog.LogFileExt)
	//获取高级别日志写入器
	highFileWriteSyncer := getLogWriter(initLog.LogSavePath + initLog.HighLevelFile + initLog.LogFileExt)
	//获取日志编码器
	encoder := getEncoder()
	//创建日志核心
	lowFileCore := zapcore.NewCore(encoder, lowFileWriteSyncer, lowPriority)
	highFileCore := zapcore.NewCore(encoder, highFileWriteSyncer, highPriority)
	//添加日志核心
	cores = append(cores, lowFileCore, highFileCore)
	if level == "debug" { //如果是 debug 级别，需要输出到终端
		consoleEncode := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		cores = append(cores, zapcore.NewCore(consoleEncode, zapcore.Lock(os.Stdout), zapcore.DebugLevel))
	}
	core := zapcore.NewTee(cores...)
	return &Log{zap.New(core, zap.AddCaller())} //增加函数调用信息
}

// getEncoder 获取日志编码器
func getEncoder() zapcore.Encoder {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		SkipLineEnding: false,
		LineEnding:     zapcore.DefaultLineEnding,   //结尾字符
		EncodeLevel:    zapcore.CapitalLevelEncoder, //将 Level 序列化为全部大写的字符串。例如 InfoLevel 被序列化为 INFO
		EncodeTime:     zapcore.ISO8601TimeEncoder,  //格式化时间戳
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	return zapcore.NewConsoleEncoder(encoderConfig) //普通的 Log Encoder
}

// getLogWriter 获取日志写入器
func getLogWriter(filename string) zapcore.WriteSyncer {
	//日志切割，滚动记录日志
	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    initLog.MaxSize,    //单个日志文件的大小（M）
		MaxAge:     initLog.MaxAge,     //最大备份天数
		MaxBackups: initLog.MaxBackups, //最大备份数量
		Compress:   initLog.Compress,   //是否压缩过期文件
	}
	return zapcore.AddSync(lumberjackLogger)
}
