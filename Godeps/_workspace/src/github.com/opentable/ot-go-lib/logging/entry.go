package logging

import (
	"net/http"
	"time"
)

// Entry encapsulates the data logged by a single log entry.
type Entry interface {
	Standard() *StandardFields
	Optional() *OptionalFields
	Request() *RequestFields
	Value() interface{}
}

type entryBase struct {
	*StandardFields
	*OptionalFields
	*RequestFields
}

func (e *entryBase) Standard() *StandardFields {
	return e.StandardFields
}

func (e *entryBase) Optional() *OptionalFields {
	return e.OptionalFields
}

func (e *entryBase) Request() *RequestFields {
	return e.RequestFields
}

func (e *entryBase) Value() interface{} {
	return e
}

func newEntry(l *logBase, s Severity, m interface{}, r *http.Request) Entry {
	var o *OptionalFields
	if r != nil {
		o = newOptionalFields(r)
	}
	return &entryBase{newStandardFields(l, s, m), o, nil}
}

func newRequestEntry(l *logBase, r *http.Request, status int, responseSize int64, duration time.Duration) Entry {
	var s Severity
	if (status - 500) >= 0 {
		s = ERROR
	} else {
		s = INFO
	}
	return &entryBase{
		newStandardFields(l, s, nil),
		newOptionalFields(r),
		newRequestFields(l, r, status, responseSize, duration),
	}
}
