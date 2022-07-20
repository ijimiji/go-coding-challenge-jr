package logger

import (
	"log"
	"os"

	"github.com/fatih/color"
)

// define loggers the following way to get a colorful prefix
var (
	Info  = log.New(os.Stdout, color.GreenString("INFO: "), log.Ltime|log.Lshortfile)
	Warn  = log.New(os.Stdout, color.YellowString("WARNING : "), log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, color.RedString("ERROR: "), log.Ltime|log.Lshortfile)
)
