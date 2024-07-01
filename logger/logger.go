package logger

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
	once        sync.Once
	logConfig   *Config
)

// LogLevel defines the severity of a logger internal.
type LogLevel int8

const (
	// DebugLevel logs debug messages.
	DebugLevel LogLevel = iota - 1
	// InfoLevel logs informational messages.
	InfoLevel
	// WarnLevel logs warning messages.
	WarnLevel
	// ErrorLevel logs error messages.
	ErrorLevel
	FatalLevel
)

// Config holds logger configuration options.
type Config struct {
	LogPath      string
	MaxSize      int
	MaxBackups   int
	MaxAge       int
	LogLevel     LogLevel
	Output       string
	Mode         string
	KafkaBrokers string
	KafkaTopic   string
}

// init initializes the logger when the package is imported.
func init() {
	once.Do(func() {
		var err error
		sugarLogger, logger, err = NewLogger()
		if err != nil {
			fmt.Printf("Could not initialize logger: %v\n", err)
			os.Exit(1)
		}
	})
}

// NewLogger initializes the logger and returns a sugared logger.
func NewLogger() (*zap.SugaredLogger, *zap.Logger, error) {
	var (
		writeSyncer zapcore.WriteSyncer
		encoder     zapcore.Encoder
		core        zapcore.Core
	)

	// Load configuration from environment variables or use default values
	logConfig = defaultLoggerConfig()

	encoder = getEncoder(logConfig.Mode)

	switch logConfig.Output {
	case "stdout":
		writeSyncer = zapcore.AddSync(os.Stdout)
	case "file":
		writeSyncer = getFileLogWriter(logConfig)
	case "kafka":
		kafkaSyncer, err := getKafkaLogWriter(logConfig)
		if err != nil {
			return nil, nil, err
		}
		writeSyncer = kafkaSyncer
	default:
		writeSyncer = zapcore.AddSync(os.Stdout)
	}

	switch logConfig.LogLevel {
	case DebugLevel:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)
	case InfoLevel:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	case WarnLevel:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.WarnLevel)
	case ErrorLevel:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.ErrorLevel)
	case FatalLevel:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.FatalLevel)
	default:
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	}

	logger = zap.New(core, zap.AddCaller())
	sugaredLogger := logger.Sugar()
	sugaredLogger.Infof("log output is %s", logConfig.Output)
	sugaredLogger.Infof("log mode is %s", strings.ToTitle(logConfig.Mode))
	return sugaredLogger, logger, nil
}

// Logger returns the initialized zap logger.
func Logger() *zap.Logger {
	return logger
}

// SugaredLogger returns the initialized sugared logger.
func SugaredLogger() *zap.SugaredLogger {
	return sugarLogger
}

// getEncoder returns the appropriate encoder based on the mode.
func getEncoder(mode string) zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig

	if mode == "prod" {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = customTimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		return zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 大写编码器
	encoderConfig.EncodeTime = customTimeEncoder            // ISO8601 UTC 时间格式
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // 相对路径编码器
	encoderConfig.EncodeName = zapcore.FullNameEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

// customTimeEncoder formats the time as 2024-06-08 00:51:55.
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

// defaultLoggerConfig loads configuration from environment variables or uses default values.
func defaultLoggerConfig() *Config {
	return &Config{
		LogPath:      getEnv("LOG_PATH", "./logs"),
		MaxSize:      getIntEnv("LOG_MAX_SIZE", 100),
		MaxBackups:   getIntEnv("LOG_MAX_BACKUPS", 30),
		MaxAge:       getIntEnv("LOG_MAX_AGE", 7),
		LogLevel:     getLogLevelEnv("LOG_LEVEL", InfoLevel),
		Output:       getEnv("LOG_OUTPUT", "stdout"),
		Mode:         getEnv("LOG_MODE", "prod"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", ""),
		KafkaTopic:   getEnv("KAFKA_TOPIC", ""),
	}
}

// getEnv retrieves the value of the environment variable named by the key or returns the default value if the variable is not present.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// getLogLevelEnv retrieves the value of the environment variable named by the key, converts it to LogLevel, or returns the default value if the variable is not present or invalid.
func getLogLevelEnv(key string, logLevel LogLevel) LogLevel {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return logLevel
	}
	switch strings.ToLower(valueStr) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return logLevel
	}
}

// getIntEnv retrieves the value of the environment variable named by the key, converts it to an integer, or returns the default value if the variable is not present or invalid.
func getIntEnv(key string, intValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists {
		return intValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return intValue
	}
	return value
}

// SetLogPath sets the log path.
func SetLogPath(logPath string) {
	if logConfig == nil {
		logConfig = &Config{}
	}
	logConfig.LogPath = logPath
}

// SetMaxSize sets the maximum size of a log file.
func SetMaxSize(maxSize int) {
	if logConfig == nil {
		logConfig = &Config{}
	}
	logConfig.MaxSize = maxSize
}

// SetMaxBackups sets the maximum number of backup log files.
func SetMaxBackups(maxBackups int) {
	if logConfig == nil {
		logConfig = &Config{}
	}
	logConfig.MaxBackups = maxBackups
}

// SetMaxAge sets the maximum age of a log file.
func SetMaxAge(maxAge int) {
	if logConfig == nil {
		logConfig = &Config{}
	}
	logConfig.MaxAge = maxAge
}

// SetLogLevel sets the log level.
func SetLogLevel(logLevel LogLevel) {
	if logConfig == nil {
		logConfig = &Config{}
	}
	logConfig.LogLevel = logLevel
}

// SetOutput sets the log output.
func SetOutput(output string) {
	if logConfig == nil {
		logConfig = &Config{}
	}
	logConfig.Output = output
}

// SetMode sets the log mode.
func SetMode(mode string) {
	if logConfig == nil {
		logConfig = &Config{}
	}
	logConfig.Mode = mode
}

// SetKafkaBrokers sets the Kafka brokers.
func SetKafkaBrokers(kafkaBrokers string) {
	if logConfig == nil {
		logConfig = &Config{}
	}
	logConfig.KafkaBrokers = kafkaBrokers
}

// SetKafkaTopic sets the Kafka topic.
func SetKafkaTopic(kafkaTopic string) {
	if logConfig == nil {
		logConfig = &Config{}
	}
	logConfig.KafkaTopic = kafkaTopic
}
