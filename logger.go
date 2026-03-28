package kasa

import (
	"log"
)

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
