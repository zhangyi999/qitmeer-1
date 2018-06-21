// Copyright (c) 2017-2018 The nox developers

package main

import (
	"os"
	"io"
	"github.com/mattn/go-colorable"
	"github.com/jrick/logrotate/rotator"
	"github.com/noxproject/nox/log"
	"github.com/noxproject/nox/log/term"
	"path/filepath"
	"fmt"
)

var (
	glogger *log.GlogHandler

	// logRotator is one of the logging outputs.  It should be closed on
	// application shutdown.
	logRotator *rotator.Rotator
)



func init() {

	// init a colorful logger if possible
	usecolor := term.IsTty(os.Stderr.Fd()) && os.Getenv("TERM") != "dumb"

	// output set to Stderr
	// it's easier to handle when run as a daemon through systemd or supervisord,
	// and Go runtime exceptions are printed to stderr as well.
	output := io.Writer(os.Stderr)
	if usecolor {
		output = colorable.NewColorableStderr()
	}
	glogger = log.NewGlogHandler(log.StreamHandler(output, log.TerminalFormat(usecolor)))

	// print log location (file:line) (useful for debug)
	// TODO config & command line flag
	log.PrintOrigins(false)

	// set log level to info
	// TODO config & comand line flag
	glogger.Verbosity(log.LvlInfo)

	log.Root().SetHandler(glogger)
}

// initLogRotator initializes the logging rotater to write logs to logFile and
// create roll files in the same directory.  It must be called before the
// package-global log rotater variables are used.
func initLogRotator(logFile string) {
	logDir, _ := filepath.Split(logFile)
	err := os.MkdirAll(logDir, 0700)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create log directory: %v\n", err)
		os.Exit(1)
	}
	r, err := rotator.New(logFile, 10*1024, false, 3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create file rotator: %v\n", err)
		os.Exit(1)
	}

	logRotator = r
}
