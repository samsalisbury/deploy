package logging

import (
	"time"
)

var (
	seqnum  int64 = 0
	seqchan       = make(chan interface{}, 1)
)

func seq() int64 {
	seqchan <- true
	seqnum++
	n := seqnum
	<-seqchan
	return n
}

func timestamp() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}
