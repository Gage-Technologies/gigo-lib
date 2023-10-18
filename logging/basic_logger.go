package logging

import (
	"fmt"
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
	tailscaleLogger "tailscale.com/types/logger"
)

// BasicLoggerOptions
//
//	Options to configure a new BasicLogger.
//
//	Properties
//	  - FilePath (string): path to the log file
//	  - Level (string) (default - DEBUG): the log level
//	  - MaxSize (int) (default - 15MB): the maximum size of the log file in MB
//	  - MaxAge (int) (default - 7): the maximum number of days to retain old log files
//	  - MaxFiles (int) (default - 10): the maximum number of log files
type BasicLoggerOptions struct {
	FilePath string
	LogLevel LogLevel
	MaxSize  int
	MaxAge   int
	MaxFiles int
}

// BasicLogger
//
//	Generic logger that writes to the specified output file.
//	BasicLogger is generally used for tests and local logging
//	requirements.
type BasicLogger struct {
	writer *lumberjack.Logger
	debug  *log.Logger
	info   *log.Logger
	warn   *log.Logger
	error  *log.Logger
}

// DefaultBasicLoggerOptions
//
//	Default options for the BasicLogger.
var DefaultBasicLoggerOptions = BasicLoggerOptions{
	FilePath: "",
	LogLevel: DEBUG,
	MaxSize:  15,
	MaxAge:   7,
	MaxFiles: 10,
}

// NewDefaultBasicLoggerOptions
//
//	Creates a new DefaultBasicLoggerOptions setting only the FilePath attribute
func NewDefaultBasicLoggerOptions(path string) BasicLoggerOptions {
	opts := DefaultBasicLoggerOptions
	opts.FilePath = path
	return opts
}

// CreateBasicLogger
//
//	 Creates a new logging object with subsequent loggers
//	 Args:
//	   opts (BasicLoggerOptions): options for configuring a basic logger
//	 Returns:
//		  (BasicLogger): newly initialized logger
func CreateBasicLogger(opts BasicLoggerOptions) (Logger, error) {
	// create lumberjack log writer to manage log overflows
	logWriter := &lumberjack.Logger{
		Filename:   opts.FilePath,
		MaxSize:    opts.MaxSize,
		MaxAge:     opts.MaxAge,
		MaxBackups: opts.MaxFiles,
		// use local because we will explicitly set the TZ in all binaries
		LocalTime: false,
		Compress:  true,
	}

	// create basic logger
	logger := &BasicLogger{
		writer: logWriter,
	}

	// create level specific loggers
	if opts.LogLevel <= DEBUG {
		logger.debug = log.New(logWriter, "D: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	}
	if opts.LogLevel <= INFO {
		logger.info = log.New(logWriter, "I: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	}
	if opts.LogLevel <= WARN {
		logger.warn = log.New(logWriter, "W: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	}
	if opts.LogLevel <= ERROR {
		logger.error = log.New(logWriter, "E: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	}

	// return log object
	return logger, nil
}

func (l *BasicLogger) WithName(name string) Logger {
	return l
}

func (l *BasicLogger) Flush() {}

func (l *BasicLogger) Close() error {
	err := l.writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %v", err)
	}
	return nil
}

func (l *BasicLogger) Debug(args ...interface{}) {
	if l.debug == nil {
		return
	}
	l.debug.Print(args...)
}
func (l *BasicLogger) Debugf(msg string, args ...interface{}) {
	if l.debug == nil {
		return
	}
	l.debug.Printf(msg, args...)
}

func (l *BasicLogger) Info(args ...interface{}) {
	if l.info == nil {
		return
	}
	l.info.Print(args...)
}
func (l *BasicLogger) Infof(msg string, args ...interface{}) {
	if l.info == nil {
		return
	}
	l.info.Printf(msg, args...)
}

func (l *BasicLogger) Warn(args ...interface{}) {
	if l.warn == nil {
		return
	}
	l.warn.Print(args...)
}
func (l *BasicLogger) Warnf(msg string, args ...interface{}) {
	if l.warn == nil {
		return
	}
	l.warn.Printf(msg, args...)
}

func (l *BasicLogger) Error(args ...interface{}) {
	if l.error == nil {
		return
	}
	l.error.Print(args...)
}
func (l *BasicLogger) Errorf(msg string, args ...interface{}) {
	if l.error == nil {
		return
	}
	l.error.Printf(msg, args...)
}

func (l *BasicLogger) DebugLogger() InternalLogger {
	return l.debug
}
func (l *BasicLogger) InfoLogger() InternalLogger {
	return l.info
}
func (l *BasicLogger) WarnLogger() InternalLogger {
	return l.warn
}
func (l *BasicLogger) ErrorLogger() InternalLogger {
	return l.error
}

func (l *BasicLogger) TailscaleDebugLogger(_ string) tailscaleLogger.Logf {
	// return no-op if this level has been filtered
	if l.debug == nil {
		return func(format string, args ...any) {}
	}
	return l.debug.Printf
}
func (l *BasicLogger) TailscaleInfoLogger(_ string) tailscaleLogger.Logf {
	// return no-op if this level has been filtered
	if l.info == nil {
		return func(format string, args ...any) {}
	}
	return l.info.Printf
}
func (l *BasicLogger) TailscaleWarnLogger(_ string) tailscaleLogger.Logf {
	// return no-op if this level has been filtered
	if l.warn == nil {
		return func(format string, args ...any) {}
	}
	return l.warn.Printf
}
func (l *BasicLogger) TailscaleErrorLogger(_ string) tailscaleLogger.Logf {
	// return no-op if this level has been filtered
	if l.error == nil {
		return func(format string, args ...any) {}
	}
	return l.error.Printf
}

// LogDebugExternalAPI Formats the log message for an external_api API debug
// Args:
//
//	message      - string, main message of the log
//	endpoint     - string, endpoint that was executed in the call
//	method       - string, method that was called when error was thrown
//	methodType   - string, HTTP method type ex: [ GET, POST, PATCH, DELETE ]
//	ip           - string, IP address of the function caller
//	statusCode   - int, status code returned by HTTP server
//	err          - error, error that was thrown during method execution
func (l *BasicLogger) LogDebugExternalAPI(message string, endpoint string, method string, methodType string, reqId interface{}, ip string, username string, userId string, statusCode int, err error) {
	if l.debug == nil {
		return
	}
	l.Debugf("%s\n    Endpoint: %s\n    Method: %s\n    Method Type: %s\n    Request ID: %v\n    IP: %s\n    User Name: %s\n    User ID: %s\n    Status Code: %d\n    Error: %v\n", message, endpoint, method, methodType, reqId, ip, username, userId, statusCode, err)
}

// LogInfoExternalAPI Formats the log message for an external_api API info
// Args:
//
//	message      - string, main message of the log
//	endpoint     - string, endpoint that was executed in the call
//	method       - string, method that was called when error was thrown
//	methodType   - string, HTTP method type ex: [ GET, POST, PATCH, DELETE ]
//	ip           - string, IP address of the function caller
//	statusCode   - int, status code returned by HTTP server
//	err          - error, error that was thrown during method execution
func (l *BasicLogger) LogInfoExternalAPI(message string, endpoint string, method string, methodType string, reqId interface{}, ip string, username string, userId string, statusCode int, err error) {
	if l.info == nil {
		return
	}
	l.Infof("%s\n    Endpoint: %s\n    Method: %s\n    Method Type: %s\n    Request ID: %v\n    IP: %s\n    User Name: %s\n    User ID: %s\n    Status Code: %d\n    Error: %v\n", message, endpoint, method, methodType, reqId, ip, username, userId, statusCode, err)
}

// LogWarnExternalAPI Formats the log message for an external_api API warning
// Args:
//
//	message      - string, main message of the log
//	endpoint     - string, endpoint that was executed in the call
//	method       - string, method that was called when error was thrown
//	methodType   - string, HTTP method type ex: [ GET, POST, PATCH, DELETE ]
//	ip           - string, IP address of the function caller
//	statusCode   - int, status code returned by HTTP server
//	err          - error, error that was thrown during method execution
func (l *BasicLogger) LogWarnExternalAPI(message string, endpoint string, method string, methodType string, reqId interface{}, ip string, username string, userId string, statusCode int, err error) {
	if l.warn == nil {
		return
	}
	l.Warnf("%s\n    Endpoint: %s\n    Method: %s\n    Method Type: %s\n    Request ID: %v\n    IP: %s\n    User Name: %s\n    User ID: %s\n    Status Code: %d\n    Error: %v\n", message, endpoint, method, methodType, reqId, ip, username, userId, statusCode, err)
}

// LogErrorExternalAPI Formats the log message for an external_api API error
// Args:
//
//	message      - string, main message of the log
//	endpoint     - string, endpoint that was executed in the call
//	method       - string, method that was called when error was thrown
//	methodType   - string, HTTP method type ex: [ GET, POST, PATCH, DELETE ]
//	ip           - string, IP address of the function caller
//	statusCode   - int, status code returned by HTTP server
//	err          - error, error that was thrown during method execution
func (l *BasicLogger) LogErrorExternalAPI(message string, endpoint string, method string, methodType string, reqId interface{}, ip string, username string, userId string, statusCode int, err error) {
	if l.error == nil {
		return
	}
	l.Errorf("%s\n    Endpoint: %s\n    Method: %s\n    Method Type: %s\n    Request ID: %v\n    IP: %s\n    User Name: %s\n    User ID: %s\n    Status Code: %d\n    Error: %v\n", message, endpoint, method, methodType, reqId, ip, username, userId, statusCode, err)
}
