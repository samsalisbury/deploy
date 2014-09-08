package logging

import (
	"fmt"
	"github.com/fzzy/radix/redis"
	"time"
)

type redisTarget struct {
	Host    string
	Port    int
	List    string
	Timeout time.Duration
}

// NewRedisTarget creates a Target that writes logs to a Redis list as JSON.
func NewRedisTarget(host string, port int, list string, timeoutMS int) Target {
	return &redisTarget{host, port, list, time.Millisecond * time.Duration(timeoutMS)}
}

func (t *redisTarget) WriteLog(e interface{}) error {
	debug("CALLED: redisTarget.WriteLog", fmt.Sprintf("%+v", e))
	data, err := json(e)
	if err != nil {
		debug("ERR: json.Marshal failed. ", err)
		return err
	}
	debug("JSON = ", string(data))
	fullHost := fmt.Sprintf("%v:%v", t.Host, t.Port)
	c, err := redis.DialTimeout("tcp", fullHost, t.Timeout)
	if err != nil {
		debug("ERR: Dialling Redis failed. ", err)
		return err
	}
	defer c.Close()
	s := string(data)
	c.Cmd("RPUSH", t.List, s)
	debug("WROTE:", s)
	return nil
}
