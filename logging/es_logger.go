package logging

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/elastic/go-elasticsearch/v7"
	elastic_logrus "github.com/gage-technologies/go-logrus-elasticsearch"
	"github.com/sirupsen/logrus"
	tailscaleLogger "tailscale.com/types/logger"
)

type ESLogger struct {
	debug     *logrus.Logger
	info      *logrus.Logger
	warn      *logrus.Logger
	error     *logrus.Logger
	debugHook *elastic_logrus.ElasticSearchHook
	infoHook  *elastic_logrus.ElasticSearchHook
	warnHook  *elastic_logrus.ElasticSearchHook
	errorHook *elastic_logrus.ElasticSearchHook
	id        string
	indexName string
}

// CreateESLogger Creates a new logging object with subsequent loggers
// Args:
//
//	logFile   - string, path to log output file
//
// Returns:
//
//	out       - Logger, freshly created logging object with subsequent loggers
func CreateESLogger(esNodes []string, esUser string, esPassword string, esIndexName string, id string) (Logger, error) {
	// retrieve hostname
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	client, err := elasticsearch.NewClient(elasticsearch.Config{
		// connecting to port
		Addresses: esNodes,
		// creating transport
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Minute * 2,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion: tls.VersionTLS11,
			},
		},
		// set authentication
		Username: esUser,
		Password: esPassword,
	})
	if err != nil {
		return nil, err
	}

	debugHook, err := elastic_logrus.NewElasticHook(client, hostname, logrus.DebugLevel, func() string {
		return esIndexName
	}, time.Second*5)
	if err != nil {
		return nil, err
	}

	debugLog := logrus.New()
	debugLog.Hooks.Add(debugHook)
	debugLog.SetFormatter(&logrus.JSONFormatter{})
	debugLog.SetLevel(logrus.DebugLevel)
	debugLog.SetOutput(ioutil.Discard)

	infoHook, err := elastic_logrus.NewElasticHook(client, hostname, logrus.InfoLevel, func() string {
		return esIndexName
	}, time.Second*5)
	if err != nil {
		return nil, err
	}

	infoLog := logrus.New()
	infoLog.Hooks.Add(infoHook)
	infoLog.SetFormatter(&logrus.JSONFormatter{})
	infoLog.SetLevel(logrus.InfoLevel)
	infoLog.SetOutput(ioutil.Discard)

	warnHook, err := elastic_logrus.NewElasticHook(client, hostname, logrus.WarnLevel, func() string {
		return esIndexName
	}, time.Second*5)
	if err != nil {
		return nil, err
	}

	warningLog := logrus.New()
	warningLog.Hooks.Add(warnHook)
	warningLog.SetFormatter(&logrus.JSONFormatter{})
	warningLog.SetLevel(logrus.WarnLevel)
	warningLog.SetOutput(ioutil.Discard)

	errorHook, err := elastic_logrus.NewElasticHook(client, hostname, logrus.ErrorLevel, func() string {
		return esIndexName
	}, time.Second*5)
	if err != nil {
		return nil, err
	}

	errorLog := logrus.New()
	errorLog.Hooks.Add(errorHook)
	errorLog.SetFormatter(&logrus.JSONFormatter{})
	errorLog.SetLevel(logrus.ErrorLevel)
	errorLog.SetOutput(ioutil.Discard)

	// create and return log object
	return &ESLogger{
		debug:     debugLog,
		info:      infoLog,
		warn:      warningLog,
		error:     errorLog,
		debugHook: debugHook,
		infoHook:  infoHook,
		warnHook:  warnHook,
		errorHook: errorHook,
		id:        id,
		indexName: esIndexName,
	}, nil
}

func (l *ESLogger) WithName(name string) Logger {
	return &ESLogger{
		debug:     l.debug,
		info:      l.info,
		warn:      l.warn,
		error:     l.error,
		debugHook: l.debugHook,
		infoHook:  l.infoHook,
		warnHook:  l.warnHook,
		errorHook: l.errorHook,
		id:        l.id,
		indexName: name,
	}
}

func (l *ESLogger) Flush() {
	l.infoHook.Flush()
	l.debugHook.Flush()
	l.warnHook.Flush()
	l.errorHook.Flush()
}

func (l *ESLogger) Debug(args ...interface{}) {
	l.debug.WithFields(logrus.Fields{
		"worker_id":   l.id,
		"logger_name": l.indexName,
	}).Debug(args...)
}
func (l *ESLogger) Debugf(msg string, args ...interface{}) {
	l.debug.WithFields(logrus.Fields{
		"worker_id":   l.id,
		"logger_name": l.indexName,
	}).Debugf(msg, args...)
}

func (l *ESLogger) Info(args ...interface{}) {
	l.info.WithFields(logrus.Fields{
		"worker_id":   l.id,
		"logger_name": l.indexName,
	}).Info(args...)
}
func (l *ESLogger) Infof(msg string, args ...interface{}) {
	l.info.WithFields(logrus.Fields{
		"worker_id":   l.id,
		"logger_name": l.indexName,
	}).Infof(msg, args...)
}

func (l *ESLogger) Warn(args ...interface{}) {
	l.warn.WithFields(logrus.Fields{
		"worker_id":   l.id,
		"logger_name": l.indexName,
	}).Warn(args...)
}
func (l *ESLogger) Warnf(msg string, args ...interface{}) {
	l.warn.WithFields(logrus.Fields{
		"worker_id":   l.id,
		"logger_name": l.indexName,
	}).Warnf(msg, args...)
}

func (l *ESLogger) Error(args ...interface{}) {
	l.error.WithFields(logrus.Fields{
		"worker_id":   l.id,
		"logger_name": l.indexName,
	}).Error(args...)
}
func (l *ESLogger) Errorf(msg string, args ...interface{}) {
	l.error.WithFields(logrus.Fields{
		"worker_id":   l.id,
		"logger_name": l.indexName,
	}).Errorf(msg, args...)
}

func (l *ESLogger) DebugLogger() InternalLogger { return l.debug }
func (l *ESLogger) InfoLogger() InternalLogger  { return l.info }
func (l *ESLogger) WarnLogger() InternalLogger  { return l.warn }
func (l *ESLogger) ErrorLogger() InternalLogger { return l.error }

func (l *ESLogger) TailscaleDebugLogger(subName string) tailscaleLogger.Logf {
	if subName != "" {
		subName = "-" + subName
	}
	return l.debug.WithFields(logrus.Fields{
		"worker_id":   l.id,
		"logger_name": l.indexName + subName,
	}).Debugf
}
func (l *ESLogger) TailscaleInfoLogger(subName string) tailscaleLogger.Logf {
	if subName != "" {
		subName = "-" + subName
	}
	return l.info.WithFields(logrus.Fields{
		"worker_id":   l.id,
		"logger_name": l.indexName + subName,
	}).Infof
}
func (l *ESLogger) TailscaleWarnLogger(subName string) tailscaleLogger.Logf {
	if subName != "" {
		subName = "-" + subName
	}
	return l.debug.WithFields(logrus.Fields{
		"worker_id":   l.id,
		"logger_name": l.indexName + subName,
	}).Warnf
}
func (l *ESLogger) TailscaleErrorLogger(subName string) tailscaleLogger.Logf {
	if subName != "" {
		subName = "-" + subName
	}
	return l.debug.WithFields(logrus.Fields{
		"worker_id":   l.id,
		"logger_name": l.indexName + subName,
	}).Errorf
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
func (l *ESLogger) LogDebugExternalAPI(message string, endpoint string, method string, methodType string, reqId interface{}, ip string, username string, userId string, statusCode int, err error) {
	startTime := time.Now()
	if reqId != nil {
		// parse snowflake into start time
		startTime = time.UnixMilli(snowflake.ParseInt64(reqId.(int64)).Time())
	}

	l.Debugf("%s\n    Exec Time: %v\n    Endpoint: %s\n    Method: %s\n    Method Type: %s\n    Request ID: %v\n    IP: %s\n    User Name: %s\n    User ID: %s\n    Status Code: %d\n    Error: %v\n", message, time.Since(startTime), endpoint, method, methodType, reqId, ip, username, userId, statusCode, err)
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
func (l *ESLogger) LogInfoExternalAPI(message string, endpoint string, method string, methodType string, reqId interface{}, ip string, username string, userId string, statusCode int, err error) {
	startTime := time.Now()
	if reqId != nil {
		// parse snowflake into start time
		startTime = time.UnixMilli(snowflake.ParseInt64(reqId.(int64)).Time())
	}

	l.Infof("%s\n    Exec Time: %v\n    Endpoint: %s\n    Method: %s\n    Method Type: %s\n    Request ID: %v\n    IP: %s\n    User Name: %s\n    User ID: %s\n    Status Code: %d\n    Error: %v\n", message, time.Since(startTime), endpoint, method, methodType, reqId, ip, username, userId, statusCode, err)
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
func (l *ESLogger) LogWarnExternalAPI(message string, endpoint string, method string, methodType string, reqId interface{}, ip string, username string, userId string, statusCode int, err error) {
	startTime := time.Now()
	if reqId != nil {
		// parse snowflake into start time
		startTime = time.UnixMilli(snowflake.ParseInt64(reqId.(int64)).Time())
	}

	l.Warnf("%s\n    Exec Time: %v\n    Endpoint: %s\n    Method: %s\n    Method Type: %s\n    Request ID: %v\n    IP: %s\n    User Name: %s\n    User ID: %s\n    Status Code: %d\n    Error: %v\n", message, time.Since(startTime), endpoint, method, methodType, reqId, ip, username, userId, statusCode, err)
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
func (l *ESLogger) LogErrorExternalAPI(message string, endpoint string, method string, methodType string, reqId interface{}, ip string, username string, userId string, statusCode int, err error) {
	startTime := time.Now()
	if reqId != nil {
		// parse snowflake into start time
		startTime = time.UnixMilli(snowflake.ParseInt64(reqId.(int64)).Time())
	}

	l.Errorf("%s\n    Exec Time: %v\n    Endpoint: %s\n    Method: %s\n    Method Type: %s\n    Request ID: %v\n    IP: %s\n    User Name: %s\n    User ID: %s\n    Status Code: %d\n    Error: %v\n", message, time.Since(startTime), endpoint, method, methodType, reqId, ip, username, userId, statusCode, err)
}
