package log

import (
	"common"
	"os"
	"time"

	mdl "common/model"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	glogconf      LogConf
	glogConfigTag = "log_property"
	gViper        *viper.Viper
	L             = mdl.L
)

type LogConf struct {
	Level         string      `mapstructure:"level" json:"level" yaml:"level"`                                  // 级别
	Format        string      `mapstructure:"format" json:"format" yaml:"format"`                               // 输出
	Prefix        string      `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                               // 日志前缀
	Director      string      `mapstructure:"director" json:"director"  yaml:"director"`                        // 日志文件夹
	ShowLine      bool        `mapstructure:"show-line" json:"show-line" yaml:"show-line"`                      // 显示行
	EncodeLevel   string      `mapstructure:"encode-level" json:"encode-level" yaml:"encode-level"`             // 编码级
	StacktraceKey string      `mapstructure:"stacktrace-key" json:"stacktrace-key" yaml:"stacktrace-key"`       // 栈名
	LogInConsole  bool        `mapstructure:"log-in-console" json:"log-in-console" yaml:"log-in-console"`       // 输出控制台
	CallerEnable  bool        `mapstructure:"caller-enable" json:"caller-enable" yaml:"caller-enable"`          // 显示调用函数
	Rotated       RotatedConf `mapstructure:"rotated-property" json:"rotated-property" yaml:"rotated-property"` // 日志循环覆盖
}

type RotatedConf struct {
	// Filename is the file to write logs to.  Backup log files will be retained
	// in the same directory.  It uses <processname>-lumberjack.log in
	// os.TempDir() if empty.
	Filename string `mapstructure:"filename" json:"filename" yaml:"filename"`

	// MaxSize is the maximum size in megabytes of the log file before it gets
	// rotated. It defaults to 100 megabytes.
	MaxSize int `mapstructure:"maxsize" json:"maxsize" yaml:"maxsize"`

	// MaxAge is the maximum number of days to retain old log files based on the
	// timestamp encoded in their filename.  Note that a day is defined as 24
	// hours and may not exactly correspond to calendar days due to daylight
	// savings, leap seconds, etc. The default is not to remove old log files
	// based on age.
	MaxAge int `mapstructure:"maxage" json:"maxage" yaml:"maxage"`

	// MaxBackups is the maximum number of old log files to retain.  The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int `mapstructure:"maxbackups" json:"maxbackups" yaml:"maxbackups"`

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time.  The default is to use UTC
	// time.
	LocalTime bool `mapstructure:"localtime" json:"localtime" yaml:"localtime"`

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool `mapstructure:"compress" json:"compress" yaml:"compress"`
}

func InitLogConfig() {
	if err := gViper.UnmarshalKey(glogConfigTag, &glogconf); err != nil {
		L.Sugar().Fatalf("Read Config to struct err:%v\n", err)
	} else {
		L.Sugar().Infof("log config init : %+v\n", glogconf)
	}
}

func createWriteSyncer(filename string) zapcore.WriteSyncer {
	if len(filename) == 0 {
		filename = common.PathJoin(glogconf.Director, glogconf.Rotated.Filename)
		//filename = glogconf.Director + "//" + glogconf.Rotated.Filename
	}

	ljLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    glogconf.Rotated.MaxSize,
		MaxBackups: glogconf.Rotated.MaxBackups,
		MaxAge:     glogconf.Rotated.MaxAge,
		LocalTime:  glogconf.Rotated.LocalTime,
		Compress:   glogconf.Rotated.Compress,
	}

	if glogconf.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(ljLogger))
	} else {
		return zapcore.AddSync(ljLogger)
	}
}

func InitLogger() *zap.Logger {
	glogconf.Director = common.DealWithExecutingCurrentFilePath(glogconf.Director)
	// 判断是否有Director文件夹
	if ok, err := common.PathExists(glogconf.Director); !ok {
		L.Sugar().Infof("Path Not Found: %v", err)
		L.Sugar().Infof("Create directory: %+v", glogconf.Director)
		err = os.Mkdir(glogconf.Director, os.ModePerm)
		if err != nil {
			L.Sugar().Fatalf("Create directory err: %#v", err)
		}
	}

	w := createWriteSyncer("")

	var level zapcore.Level
	switch glogconf.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var core zapcore.Core

	if glogconf.Format == "json" {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			w,
			level,
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			w,
			level,
		)
	}
	logger := zap.New(core)
	if glogconf.CallerEnable {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}

func Zap() (logger *zap.Logger) {

	// 调试级别
	debugPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.DebugLevel
	})
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.InfoLevel
	})
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev == zap.WarnLevel
	})
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
		return lev >= zap.ErrorLevel
	})

	cores := [...]zapcore.Core{

		getEncoderCore(common.PathJoin(glogconf.Director, "debug.log"), debugPriority),
		getEncoderCore(common.PathJoin(glogconf.Director, "info.log"), infoPriority),
		getEncoderCore(common.PathJoin(glogconf.Director, "warn.log"), warnPriority),
		getEncoderCore(common.PathJoin(glogconf.Director, "error.log"), errorPriority),
	}
	logger = zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
	logger = logger.WithOptions(zap.AddCaller())

	return logger
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  glogconf.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	switch glogconf.EncodeLevel {
	case "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder() zapcore.Encoder {
	if glogconf.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(fileName string, level zapcore.LevelEnabler) (core zapcore.Core) {
	writer := createWriteSyncer(fileName) // 使用file-rotatelogs进行日志分割
	return zapcore.NewCore(getEncoder(), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(glogconf.Prefix + "2006/01/02 - 15:04:05.000"))
	//enc.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
}
