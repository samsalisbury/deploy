package logging

// Severity. One of FATAL, ERROR, WARN, INFO, TRACE
type Severity int

const (
	none  = iota
	FATAL = iota
	ERROR = iota
	WARN  = iota
	INFO  = iota
	TRACE = iota
)

// String returns the string representation of a Severity.
func (s Severity) String() string {
	switch s {
	default:
		return "NONE"
	case FATAL:
		return "FATAL"
	case ERROR:
		return "ERROR"
	case WARN:
		return "WARN"
	case INFO:
		return "INFO"
	case TRACE:
		return "TRACE"
	}
}
