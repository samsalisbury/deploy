// +build !debuglogging

package logging

import (
	"errors"
	"fmt"
)

func debug(...interface{}) {}

func Error(args ...interface{}) error {
	return errors.New(fmt.Sprint(args...))
}
