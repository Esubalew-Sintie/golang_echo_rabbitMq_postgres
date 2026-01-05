package logger

import (
	"context"
	"fmt"
	"log"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...interface{})
	Error(ctx context.Context, msg string, args ...interface{})
	Warn(ctx context.Context, msg string, args ...interface{})
	Panic(ctx context.Context, msg string, args ...interface{})
	Fatal(ctx context.Context, msg string, args ...interface{})
}

type SimpleLogger struct {
	prefix string
}

func NewSimpleLogger(prefix string) Logger {
	return &SimpleLogger{prefix: prefix}
}

func (l *SimpleLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	log.Printf("[INFO] %s: %s", l.prefix, msg)
}

func (l *SimpleLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	log.Printf("[ERROR] %s: %s", l.prefix, msg)
}

func (l *SimpleLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	log.Printf("[WARN] %s: %s", l.prefix, msg)
}

func (l *SimpleLogger) Panic(ctx context.Context, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	log.Panicf("[PANIC] %s: %s", l.prefix, msg)
}

func (l *SimpleLogger) Fatal(ctx context.Context, msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	log.Fatalf("[FATAL] %s: %s", l.prefix, msg)
}

func InitLogger() Logger {
	return NewSimpleLogger("payment-gateway")
}
