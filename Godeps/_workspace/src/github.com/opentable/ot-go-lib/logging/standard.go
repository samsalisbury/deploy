package logging

import "github.com/opentable/ot-go-lib/env"

type Std struct {
	Config *LogConfig
	Redis  Target
	Stderr Target
}

// StandardConfig creates a special kind of LogConfig (Std), based on environment
// variables. From an Std, you can create Logs easily.
// If any of the variables are missing or malformed,
// this function will cause your program to exit immediately. The required
// variables are:
//
// OT_LOGGING_REDIS_HOST, OT_LOGGING_REDIS_PORT, OT_LOGGING_REDIS_LIST,
// OT_LOGGING_REDIS_TIMEOUT_MS, HOSTNAME, PORT0
//
// The parameter should be the name of your app, without its version.
func StandardConfig(serviceType string) *Std {
	return &Std{
		NewLogConfig(
			serviceType,
			env.RequireAPP_HOST(),
			env.RequirePORT0()),
		NewRedisTarget(
			env.RequireString("OT_LOGGING_REDIS_HOST"),
			env.RequireInt("OT_LOGGING_REDIS_PORT"),
			env.RequireString("OT_LOGGING_REDIS_LIST"),
			env.RequireInt("OT_LOGGING_REDIS_TIMEOUT_MS")),
		NewStderrTarget(),
	}
}

// Log returns a Log with a Redis target.
func (s *Std) Log(name string, version int) Log {
	return s.Config.NewLog(name, version, s.Redis)
}

// StartupLog returns a Log named "startup" with a Redis target and a stderr target.
func (s *Std) StartupLog(version int) StartupLog {
	return s.Config.NewStartupLog(version, s.Redis, s.Stderr)
}

// RequestLog returns a RequestLog with a Redis target.
func (s *Std) RequestLog(version int) RequestLog {
	return s.Config.NewRequestLog(version, s.Redis)
}
