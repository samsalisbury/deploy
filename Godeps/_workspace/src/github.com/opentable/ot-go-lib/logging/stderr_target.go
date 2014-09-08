package logging

import (
	"fmt"
	"log"
)

type stderrTarget struct {
}

func NewStderrTarget() Target {
	return &stderrTarget{}
}

func (d *stderrTarget) WriteLog(e interface{}) error {
	entry, ok := e.(Entry)
	if !ok {
		return writeOther(e)
	}
	return writeEntry(entry)
}

func writeEntry(e Entry) error {
	s := e.Standard()
	r := e.Request()
	rs := ""
	if r != nil {
		rs = fmt.Sprint("(Serving request: )", r.Method, r.Url, " ")
	}
	log.Println(s.SequenceNumber, s.Severity, s.LogMessage, rs)
	return nil
}

func writeOther(o interface{}) error {
	data, err := json(o)
	if err != nil {
		return err
	}
	log.Println(string(data))
	return nil
}
