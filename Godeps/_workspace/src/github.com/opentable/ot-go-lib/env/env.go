// Package env provides helper methods for reading environment variables.
//
// These methods add a small amount of convention and some convenience methods
// for OpenTable Go applications. The Require(String|Int) methods are highly
// recommended, as they require less code to use, and ensure your program exits
// when its required configuration is missing.
//
// Here is an example, showing how to require some environment variables, and
// immediately exit with a log message if any of those variables are missing:
//
//    package main
//
//    import "github.com/opentable/ot-go-lib/env"
//
//    var (
//        redisHost = env.RequireString("OT_LOGGING_REDIS_HOST")
//        redisList = env.RequireString("OT_LOGGING_REDIS_LIST")
//        redisPort = env.RequireInt("OT_LOGGING_REDIS_PORT")
//    )
//
//    func main() {
//        println("Redis host:", redisHost)
//        println("Redis list:", redisList)
//        println("Redis port:", redisPort)
//    }
//
package env

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

// RequireString is a convenience wrapper around String that calls
// log.Fatal (exiting the program) if the env var is not set.
//
// RequireString is designed to be used in program initialisation to
// exit early if a requires variable is not set.
func RequireString(name string) string {
	s, err := String(name)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

// String gets the named environment variable. It returns an error
// if the variable is not set, or is set to an empty string.
func String(name string) (string, error) {
	v, ok := s(name)
	if !ok {
		return "", notset(name)
	}
	return v, nil
}

// StringOrDefault gets the named environment variable, or, if
// that variable is not set, returns the default string provided.
//
// Use of StringOrDefault is probably a code smell. This function
// may be removed in future.
func StringOrDefault(name string, def string) string {
	v, ok := s(name)
	if !ok {
		return def
	}
	return v
}

// RequireInt is a convenience wrapper around Int that calls
// log.Fatal (exiting the calling program) if the env var is not set
// or is not parseable as a decimal integer on this system.
//
// RequireInt is designed to be used in program initialisation to
// exit early if a requires variable is not set.
func RequireInt(name string) int {
	i, err := Int(name)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

// Int gets the names environment variable as an int. It returns an
// error if the variable is not set, or is not parseable as a decimal
// integer on this system.
func Int(name string) (int, error) {
	s, err := String(name)
	if err != nil {
		return 0, err
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("env var %v not parseable as int; was: '%v' (error: %s)", name, s, err))
	}
	return i, nil
}

// RequirePORT0 is equivalent to calling RequireInt("PORT0").
// PORT0 is a variable typically set by Mesos, it is used
// to tell your application which port to listen on. It is
// recommended to use this variable on local environments as
// well.
func RequirePORT0() int {
	return RequireInt("PORT0")
}

func RequireAPP_HOST() string {
	return RequireString("APP_HOST")
}

func RequireAppURL() string {
	return "http://" + RequireAPP_HOST() + RequireListenAddr()
}

func RequireListenAddr() string {
	return ":" + strconv.Itoa(RequirePORT0())
}

// AssertServiceType ensures that OT_SERVICE_TYPE has been set to a specific
// string. This is important for some of the standard ot-go-lib libraries.
// NOTE: Service Type is OpenTable parlance for "Service name without version",
// e.g. Discovery Service v1.0.0's serviceType would be 'discovery'.
func AssertServiceType(is string) {
	if is == "" {
		log.Fatal("AssertServiceType requires a non-empty string.")
	}
	t := RequireString("OT_SERVICE_TYPE")
	if t != is {
		log.Fatal("Got OT_SERVICE_TYPE == "+t+"; want: ", is)
	}
}

func s(name string) (v string, ok bool) {
	v = os.Getenv(name)
	ok = v != ""
	return
}

func notset(name string) error {
	msg := fmt.Sprintf("env var %s not set", name)
	return errors.New(msg)
}
