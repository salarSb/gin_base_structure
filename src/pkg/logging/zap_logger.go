package logging

import (
	"base_structure/src/config"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var zapLogLevelMap = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}
var zapSingleLogger *zap.SugaredLogger

type zapLogger struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func newZapLogger(cfg *config.Config) *zapLogger {
	logger := &zapLogger{cfg: cfg}
	logger.Init()
	return logger
}

func prepareLogInfo(extra map[ExtraKey]interface{}, cat Category, sub SubCategory) []interface{} {
	if extra == nil {
		extra = make(map[ExtraKey]interface{})
	}
	extra["Category"] = cat
	extra["SubCategory"] = sub
	params := logParamsToZapParams(extra)
	return params
}

func (l *zapLogger) getLogLevel() zapcore.Level {
	level, exists := zapLogLevelMap[l.cfg.Logger.Level]
	if !exists {
		return zapcore.DebugLevel
	}
	return level
}

func (l *zapLogger) Init() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			panic("unable to read .env file")
		}
		appName := os.Getenv("APP_NAME")
		fileName := fmt.Sprintf("%s%s-%s.%s", l.cfg.Logger.FilePath, time.Now().Format("2006-01-02"), uuid.New(), "log")
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    1,
			MaxAge:     5,
			LocalTime:  true,
			MaxBackups: 10,
			Compress:   true,
		})
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), w, l.getLogLevel())
		logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
		zapSingleLogger = logger.With(string(AppName), appName, string(LoggerName), "ZapLog")
	})
	l.logger = zapSingleLogger
}

func (l *zapLogger) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(extra, cat, sub)
	l.logger.Debugw(msg, params...)
}

func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args)
}

func (l *zapLogger) Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(extra, cat, sub)
	l.logger.Infow(msg, params...)
}

func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args)
}

func (l *zapLogger) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(extra, cat, sub)
	l.logger.Warnw(msg, params...)
}

func (l *zapLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args)
}

func (l *zapLogger) Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(extra, cat, sub)
	l.logger.Errorw(msg, params...)
}

func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args)
}

func (l *zapLogger) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(extra, cat, sub)
	l.logger.Fatalw(msg, params...)
}

func (l *zapLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args)
}
