package logging

import tailscaleLogger "tailscale.com/types/logger"

type InternalLogger interface {
	Print(...interface{})
	Println(...interface{})
	Printf(string, ...interface{})

	Fatal(...interface{})
	Fatalln(...interface{})
	Fatalf(string, ...interface{})
}

type Logger interface {
	Flush()

	Debug(...interface{})
	Debugf(msg string, args ...interface{})

	Info(...interface{})
	Infof(msg string, args ...interface{})

	Warn(...interface{})
	Warnf(msg string, args ...interface{})

	Error(...interface{})
	Errorf(msg string, args ...interface{})

	WithName(name string) Logger

	DebugLogger() InternalLogger
	InfoLogger() InternalLogger
	WarnLogger() InternalLogger
	ErrorLogger() InternalLogger

	TailscaleDebugLogger(subName string) tailscaleLogger.Logf
	TailscaleInfoLogger(subName string) tailscaleLogger.Logf
	TailscaleWarnLogger(subName string) tailscaleLogger.Logf
	TailscaleErrorLogger(subName string) tailscaleLogger.Logf

	LogDebugExternalAPI(message string, endpoint string, method string, methodType string, reqId interface{}, ip string, username string, userId string, statusCode int, err error)
	LogInfoExternalAPI(message string, endpoint string, method string, methodType string, reqId interface{}, ip string, username string, userId string, statusCode int, err error)
	LogWarnExternalAPI(message string, endpoint string, method string, methodType string, reqId interface{}, ip string, username string, userId string, statusCode int, err error)
	LogErrorExternalAPI(message string, endpoint string, method string, methodType string, reqId interface{}, ip string, username string, userId string, statusCode int, err error)
}
