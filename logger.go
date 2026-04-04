package kasa

import (
	"log"
	// "log/slog"
	// "sync/atomic"
)

/*
var klogger atomic.Pointer[kasalogger]

func init() {
    var defaultLog kasalogger = log.Default()
    klogger.Store(&defaultLog)
}

func SetLogger(l kasalogger) {
    klogger.Store(&l)
}

// Internal helper to get the logger safely
func getLogger() kasalogger {
    return *klogger.Load()
} */

type NoopLogger struct{}

func (n NoopLogger) Println(...any)        {}
func (n NoopLogger) Printf(string, ...any) {}

// by default, use the standard logger, can be overwritten using kasa.SetLogger(l)
var klogger kasalogger = log.Default()

// Any log interface that has Println and Printf will do
type kasalogger interface {
	Println(...any)
	Printf(string, ...any)
}

// SetLogger allows applications to register their own logging interface
func SetLogger(l kasalogger) {
	klogger = l
}
