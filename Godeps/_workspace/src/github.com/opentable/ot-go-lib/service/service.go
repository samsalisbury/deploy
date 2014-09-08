// +build linux darwin

// Package service provides a standard well-behaved HTTP service base, indended for use on
// Mesos, but potentially usable anywhere. It includes logging, discovery, and graceful
// startup/shutdown.
//
// NOTE: this package is only available for Linux and Mac (darwin). Windows support is not
// planned but would be trivial to implement.
package service

import (
	"github.com/opentable/ot-go-lib/disco"
	"github.com/opentable/ot-go-lib/env"
	"github.com/opentable/ot-go-lib/logging"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	addr = env.RequireListenAddr()
	url  = env.RequireAppURL()
)

type Service interface {
	// Start starts the service. It will shut down gracefully on SIGINT or SIGTERM.
	Start()
	// Stop shuts down the service gracefully.
	Stop()
	// Healthy gets the health of the service.
	Healthy() bool
	// SetHealthy sets the service's health.
	SetHealthy(bool)
	// StartupLog returns the generated startup log.
	StartupLog() logging.StartupLog
	// Discovery returns the generated discovery client.
	Discovery() *disco.Client
}

var rwTimeout = 3 * time.Second

type httpService struct {
	ServiceType string
	startupLog  logging.StartupLog
	disco       *disco.Client
	server      *http.Server
	done        chan os.Signal
	healthy     bool
	wg          sync.WaitGroup
}

// RequireHTTPServiceFromEnv is very similar to NewHTTPServiceFromEnv, except that
// it does not return an error but rather records a fatal log message which should
// cause your app to exit immediately on error. This is useful in a var block.
func RequireHTTPServiceFromEnv(serviceType string, handler http.HandlerFunc) Service {
	s, err := NewHTTPServiceFromEnv(serviceType, handler)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	return s
}

// NewHTTPServiceFromEnv creates a new HTTP service using standard OT env vars.
// Incoming requests for GET /health will be handled by the internal health
// handler. Everything else will be passed to the handlerFunc you pass as the last arg.
func NewHTTPServiceFromEnv(serviceType string, handler http.HandlerFunc) (Service, error) {
	//env.AssertServiceType(serviceType)
	logConfig := logging.StandardConfig(serviceType)
	startupLog := logConfig.StartupLog(0)
	disco, err := disco.NewClientFromEnv(serviceType, "No comment.")
	if err != nil {
		return nil, err
	}
	s := &httpService{
		ServiceType: serviceType,
		startupLog:  startupLog,
		disco:       disco,
		healthy:     true,
		wg:          sync.WaitGroup{},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.healthHandler)
	//mux.HandleFunc("/disco", s.discoHandler)
	mux.HandleFunc("/", handler)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  rwTimeout,
		WriteTimeout: rwTimeout,
		// Track open connections (we don't care about idle connections atm)
		ConnState: func(_ net.Conn, state http.ConnState) {
			switch state {
			case http.StateNew:
				s.wg.Add(1)
			case http.StateClosed, http.StateHijacked, http.StateIdle:
				s.wg.Done()
			}
		},
	}
	s.server.SetKeepAlivesEnabled(false)
	return s, nil
}

var discoSleep = 5 * time.Second

// Start performs a graceful startup of the service, by doing these things in this order:
//
// 1: start logging;
// 2: start web server;
// 3: announce to discovery.
//
// Once started, the server listens for SIGINT and SIGTERM, and initiates graceful shut
// down. This can be manually initiated from outside by calling Stop. The graceful shut
// down does these things in this order:
//
// 1. Unannounce to discovery;
// 2. bleed connections (or timeout after 5s);
// 3. stop logging.
func (s *httpService) Start() {
	s.done = make(chan os.Signal)
	go func() {
		s.startupLog.Info("Listening on " + addr)
		s.startupLog.Fatal(s.server.ListenAndServe())
	}()
	go func() {
		s.startupLog.Info("Announicing at " + url)
		s.disco.AnnounceEvery500ms()
	}()
	signal.Notify(s.done, os.Interrupt, syscall.SIGTERM)
	<-s.done
	s.disco.Unannounce()
	s.startupLog.Info("Shutting down... Waiting " + discoSleep.String() + " for discovery clients to catch up.")
	time.Sleep(discoSleep)
	s.SetHealthy(false)
	connectionsBled := make(chan struct{})
	go func() {
		s.wg.Wait()
		connectionsBled <- struct{}{}
	}()
	select {
	case <-time.After(discoSleep):
		s.startupLog.Error(nil, "Service stopped after "+discoSleep.String()+", some connections were still open.")
	case <-connectionsBled:
	}
	s.startupLog.Info("Shut down successful; killing startup log.")
	s.startupLog.Stop()
}

// Stop invokes the graceful shutdown. See documntation for Start for details.
func (s *httpService) Stop() {
	s.done <- os.Interrupt
}

func (s *httpService) StartupLog() logging.StartupLog {
	return s.startupLog
}

func (s *httpService) Discovery() *disco.Client {
	return s.disco
}

func (s *httpService) Healthy() bool {
	return s.healthy
}

func (s *httpService) SetHealthy(to bool) {
	s.healthy = to
}

func (s *httpService) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	if s.healthy {
		w.WriteHeader(200)
		w.Write([]byte("healthy"))
	} else {
		w.WriteHeader(503)
		w.Write([]byte("unhealthy"))
	}
}

// func (s *httpService) discoHandler(w http.ResponseWriter, r *http.Request) {
// 	state := s.disco.State()
// 	b, err := json.Marshal(state)
// 	if err != nil {
// 		w.WriteHeader(500)
// 		w.Write([]byte(err.Error()))
// 	}
// 	w.WriteHeader(200)
// 	w.Write(b)
// }
