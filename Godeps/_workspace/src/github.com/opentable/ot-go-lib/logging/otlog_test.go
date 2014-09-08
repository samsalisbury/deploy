package logging

import (
	"bytes"
	"testing"
	"time"
)

type TestTarget struct {
	lastEntry interface{}
}

func (target *TestTarget) WriteLog(e interface{}) error {
	target.lastEntry = e
	return nil
}

func Test_LogLifecycle(t *testing.T) {
	buf := new(bytes.Buffer)
	testtarget := &TestTarget{buf}
	target := Target(testtarget)
	ts := []Target{target}
	c := NewLogConfig("test-type", "localhost", 12345, ts)
	l := c.NewLog("application", 0)
	l.Info("hello", nil)
	l.Stop()
	e := testtarget.lastEntry.(entryS)
	if e.Host != "localhost:12345" {
		t.Error("e.Host was ", e.Host)
	}
}

func Test_RedisLogGetsCalled(y *testing.T) {
	targets := []Target{NewRedisTarget("host", 1, "list", time.Duration(1000)*time.Millisecond)}
	config := NewLogConfig("serviceType", "appHost", 12345, targets)
	log := config.NewLog("name", 0)
	log.Info("hello?", nil)
}
