package logging

import (
	"base_structure/src/config"
	"fmt"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var (
	zapSingleLogger *zap.SugaredLogger
	atomicLevel     zap.AtomicLevel
	zapLogLevelMap  = map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
	}
)

type zapLogger struct {
	cfg    *config.Config
	logger *zap.SugaredLogger
}

func newZapLogger(cfg *config.Config) *zapLogger {
	l := &zapLogger{cfg: cfg}
	l.Init()
	return l
}

func prepareLogInfo(extra map[ExtraKey]interface{}, cat Category, sub SubCategory) []interface{} {
	if extra == nil {
		extra = make(map[ExtraKey]interface{})
	}
	extra["Category"] = cat
	extra["SubCategory"] = sub
	return logParamsToZapParams(extra)
}

func (l *zapLogger) Init() {
	once.Do(func() {
		atomicLevel = zap.NewAtomicLevelAt(l.getLogLevel())
		fileName := fmt.Sprintf("%s%s.%s", l.cfg.Logger.FilePath, time.Now().Format("2006-01-02"), "log")
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(prodEncoderCfg()),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   fileName,
				MaxSize:    10, // MB
				MaxBackups: 30,
				MaxAge:     14, // days
				Compress:   true,
				LocalTime:  true,
			}),
			atomicLevel,
		)
		consoleCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(prodEncoderCfg()),
			zapcore.AddSync(os.Stdout),
			atomicLevel,
		)
		core := zapcore.NewTee(fileCore, consoleCore)
		zl := zap.New(core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		).Sugar()
		_ = godotenv.Load()
		zapSingleLogger = zl.With(string(AppName), os.Getenv("APP_NAME"), string(LoggerName), "ZapLog")
	})
	l.logger = zapSingleLogger
}

func prodEncoderCfg() zapcore.EncoderConfig {
	enc := zap.NewProductionEncoderConfig()
	enc.EncodeTime = zapcore.ISO8601TimeEncoder
	return enc
}

func (l *zapLogger) getLogLevel() zapcore.Level {
	if lvl, ok := zapLogLevelMap[l.cfg.Logger.Level]; ok {
		return lvl
	}
	return zapcore.DebugLevel
}

func (l *zapLogger) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.Debugw(msg, prepareLogInfo(extra, cat, sub)...)
}
func (l *zapLogger) Debugf(template string, args ...interface{}) { l.logger.Debugf(template, args...) }

func (l *zapLogger) Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.Infow(msg, prepareLogInfo(extra, cat, sub)...)
}
func (l *zapLogger) Infof(template string, args ...interface{}) { l.logger.Infof(template, args...) }

func (l *zapLogger) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.Warnw(msg, prepareLogInfo(extra, cat, sub)...)
}
func (l *zapLogger) Warnf(template string, args ...interface{}) { l.logger.Warnf(template, args...) }

func (l *zapLogger) Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.Errorw(msg, prepareLogInfo(extra, cat, sub)...)
}
func (l *zapLogger) Errorf(template string, args ...interface{}) { l.logger.Errorf(template, args...) }

func (l *zapLogger) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	l.logger.Fatalw(msg, prepareLogInfo(extra, cat, sub)...)
}
func (l *zapLogger) Fatalf(template string, args ...interface{}) { l.logger.Fatalf(template, args...) }

func (l *zapLogger) Sync() error { return l.logger.Sync() }
