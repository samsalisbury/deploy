package logging

import (
	"fmt"
)

// LogConfig contains a basic set of log configuration.
// It acts as a factory to create Logs and RequestLogs.
type LogConfig struct {
	// ServiceType is the name of your service. E.g. the Deployment Service has ServiceType = "deployment".
	// It should be be an alphanumeric string, matching the regular expression "[A-Za-z0-9]+"
	ServiceType string
	// AppHost should be set to the host name of the machine where this instance of the app is accessible.
	AppHost string
	// AppPort should be set to the port number this instance of the app is available on.
	// Currently, only a single port is supported.
	AppPort int
	// hostString is AppHost:AppPort
	hostString string
}

// NewLogConfig returns an initialised LogConfig. Please see
// the documentation for type LogConfig for the meanings of
// the fields.
func NewLogConfig(serviceType string, appHost string, appPort int) *LogConfig {
	hostString := fmt.Sprintf("%v:%v", appHost, appPort)
	return &LogConfig{serviceType, appHost, appPort, hostString}
}
