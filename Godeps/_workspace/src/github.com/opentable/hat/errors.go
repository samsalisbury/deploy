package hat

import (
	"fmt"
	"strings"
)

type hatError struct {
	Message string
}

type HTTPError interface {
	error
	StatusCode() int
	Err() error
}

type httpError struct {
	statusCode int
	hatError
}

func (he httpError) StatusCode() int {
	return he.statusCode
}

func (he httpError) Err() error {
	return he.hatError
}

func (h hatError) Error() string {
	return h.Message
}

func Error(args ...interface{}) hatError {
	message := msgS(args...)
	return hatError{message}
}

func HttpError(statusCode int, args ...interface{}) HTTPError {
	return httpError{statusCode, Error(args...)}
}

// Error creates an error message prefixed by the node's entity pointer type name.
func (n *Node) Error(args ...interface{}) hatError {
	args = prepend(args, n.EntityPtrType)
	return Error(args...)
}

// MethodError creates an error message prefixed by the node's entity type name
// and a method name (you must supply the method name).
func (n *Node) MethodError(name string, args ...interface{}) hatError {
	args = prepend(args, msgNS(n.EntityPtrType, ".", name))
	return Error(args...)
}

func (n *Node) wrongNumIn(name string, o *Op, numIn int) hatError {
	var want string
	if o.MaxIn() == o.MinIn() {
		want = msgS(o.MinIn(), "parameters")
	} else {
		want = msgS(msgNS(o.MinIn(), "-", o.MaxIn()), "parameters, inclusive")
	}
	return n.MethodError(name, "got", numIn, "inputs; want", want)
}

func prepend(a []interface{}, b interface{}) []interface{} {
	return append([]interface{}{b}, a...)
}

func debug(args ...interface{}) {
	println(Error(args...).Error())
}

// msgS concatenates string representations of the args with spaces between them.
func msgS(args ...interface{}) string {
	return msg(" ", args...)
}

// msgS concatenates string representations of the args with no spaces between them.
func msgNS(args ...interface{}) string {
	return msg("", args...)
}

func msg(sep string, args ...interface{}) string {
	message := []string{}
	for _, a := range args {
		message = append(message, fmt.Sprint(a))
	}
	return strings.Join(message, sep)
}

func quot(args ...interface{}) string {
	return "'" + msgS(args...) + "'"
}
