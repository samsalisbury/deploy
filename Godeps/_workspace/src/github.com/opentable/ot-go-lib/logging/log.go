package logging

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

// Log is a general-purpose logger.
//
// Each of the main methods (Fatal(), Error(), Warn(), Info(), Trace()) takes an *http.Request
// as its first argument. You should provide the current *http.Request wherever possible, as it
// may contain some OpenTable-specific headers that will enrich the logs. When there is no
// associated request, you should set this to nil.
type Log interface {
	// Fatal logs a message with severity FATAL, and calls os.Exit(1)
	// The *http.Requst parameter is documented above.
	Fatal(req *http.Request, args ...interface{}) Entry
	// Error logs a message with severity ERROR, and returns the log entry written.
	// The *http.Requst parameter is documented above.
	Error(req *http.Request, args ...interface{}) Entry
	// Warn logs a message with severity WARN, and returns the log entry written.
	// The *http.Requst parameter is documented above.
	Warn(req *http.Request, args ...interface{}) Entry
	// Info logs a message with severity INFO, and returns the log entry written.
	// The *http.Requst parameter is documented above.
	Info(req *http.Request, args ...interface{}) Entry
	// Trace logs a message with severity TRACE, and returns the log entry written.
	// The *http.Requst parameter is documented above.
	Trace(req *http.Request, args ...interface{}) Entry
	stopper
}

// RequestLog is a special-purpose logger, designed for writing request logs.
type RequestLog interface {
	// LogRequest writes a special type of Entry that includes an instance of RequestFields.
	LogRequest(request *http.Request, statusCode int, responseSize int64, duration time.Duration) Entry
	stopper
}

// StartupLog is similar to Log, except that its methods do not expect an *http.Request parameter
// as it is assumed that on app startup, there will be no http request in context.
type StartupLog interface {
	// Fatal logs a message with severity FATAL, and calls 	os.Exit(1)
	// The *http.Requst parameter is documented above.
	Fatal(args ...interface{}) Entry
	// Error logs a message with severity ERROR, and returns the log entry written.
	// The *http.Requst parameter is documented above.
	Error(args ...interface{}) Entry
	// Warn logs a message with severity WARN, and returns the log entry written.
	// The *http.Requst parameter is documented above.
	Warn(args ...interface{}) Entry
	// Info logs a message with severity INFO, and returns the log entry written.
	// The *http.Requst parameter is documented above.
	Info(args ...interface{}) Entry
	// Trace logs a message with severity TRACE, and returns the log entry written.
	// The *http.Requst parameter is documented above.
	Trace(args ...interface{}) Entry
	stopper
}

type stopper interface {
	Stop()
}

// NewLog creates a new Log based on a LogConfig. Please see the documentation for types LogConfig and Log.
func (c *LogConfig) NewLog(name string, version int, targets ...Target) Log {
	return c.newLogBase(name, version, targets)
}

// NewStartupLog creates a new StartupLog based on a LogConfig. Please see the documentation for types LogConfig and Log.
func (c *LogConfig) NewStartupLog(version int, targets ...Target) StartupLog {
	return &startupLogBase{c.newLogBase("startup", version, targets)}
}

// NewRequestLog creates a new RequestLog based on a LogConfig. Please see the documentation for types LogConfig and RequestLog.
func (c *LogConfig) NewRequestLog(version int, targets ...Target) RequestLog {
	return c.newLogBase("request", version, targets)
}

type logBase struct {
	Config              *LogConfig
	LogName             string
	FormatVersion       int
	formatVersionString string
	Type                string
	Targets             []Target
	wg                  sync.WaitGroup
	stopped             bool
}

func (l *logBase) log(s Severity, r *http.Request, a ...interface{}) Entry {
	if l.stopped {
		panic(fmt.Sprint(l.LogName, ".log called after it was stopped"))
	}
	var m interface{}
	if len(a) > 1 {
		m = fmt.Sprint(a...)
	} else {
		m = a[0]
	}
	e := newEntry(l, s, m, r)
	l.broadcast(e)
	return e
}

func (l *logBase) Fatal(r *http.Request, a ...interface{}) Entry {
	l.log(FATAL, r, a...)
	l.Stop() // Allow logs to complete before exiting!
	os.Exit(1)
	return nil
}

func (l *logBase) Error(r *http.Request, a ...interface{}) Entry {
	return l.log(ERROR, r, a...)
}

func (l *logBase) Warn(r *http.Request, a ...interface{}) Entry {
	return l.log(WARN, r, a...)
}

func (l *logBase) Info(r *http.Request, a ...interface{}) Entry {
	return l.log(INFO, r, a...)
}

func (l *logBase) Trace(r *http.Request, a ...interface{}) Entry {
	return l.log(TRACE, r, a...)
}

func (l *logBase) LogRequest(r *http.Request, status int, responseSize int64, duration time.Duration) Entry {
	e := newRequestEntry(l, r, status, responseSize, duration)
	l.broadcast(e)
	return e
}

func (l *logBase) broadcast(e Entry) {
	for _, t := range l.Targets {
		l.wg.Add(1)
		go func(t Target) {
			defer l.wg.Done()
			t.WriteLog(e.Value())
		}(t)
	}
}

func (l *logBase) Stop() {
	l.stopped = true
	// TODO: Prevent incoming messages from keeping this alive too long
	l.wg.Wait()
}

func (c *LogConfig) newLogBase(name string, version int, targets []Target) *logBase {
	versionString := fmt.Sprintf("v%v", version)
	return &logBase{
		Config:              c,
		LogName:             name,
		FormatVersion:       version,
		formatVersionString: versionString,
		Type:                fmt.Sprintf("%v-%v-%v", c.ServiceType, name, versionString),
		Targets:             targets,
	}
}

type startupLogBase struct {
	base *logBase
}

func (l *startupLogBase) Fatal(a ...interface{}) Entry {
	return l.base.Fatal(nil, a...)
}

func (l *startupLogBase) Error(a ...interface{}) Entry {
	return l.base.Error(nil, a...)
}

func (l *startupLogBase) Warn(a ...interface{}) Entry {
	return l.base.Warn(nil, a...)
}

func (l *startupLogBase) Info(a ...interface{}) Entry {
	return l.base.Info(nil, a...)
}

func (l *startupLogBase) Trace(a ...interface{}) Entry {
	return l.base.Trace(nil, a...)
}

func (l *startupLogBase) Stop() {
	l.base.Stop()
}
