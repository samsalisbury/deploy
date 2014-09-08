package logging

// Target represents a log target (somewhere where Entries should be written).
// You can extend this package by implementing this yourself.
//
// WriteLog() should write the log provided somewhere.
// This interface is primarily used internally.
// Whenever it is called by ot-go-lib/logging, it is passed a struct of the form:
//
//   {*StandardFields, *OptionalFields, *RequestFields}
//
// Which you may use to compose logs in whatever format you like. Note: Any of the fields may be nil.
type Target interface {
	WriteLog(interface{}) error
}
