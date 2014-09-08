// +build debuglogging

package logging

import (
	"errors"
	"fmt"
	"log"
)

func debug(args ...interface{}) {
	log.Println(append([]interface{}{"debuglogging:"}, args...)...)
}

func Error(args ...interface{}) error {
	debug(args...)
	return errors.New(fmt.Sprint(args...))
}
