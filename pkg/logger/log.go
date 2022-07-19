package logger

import (
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	InfoLogger    = log.New(os.Stdout, color.GreenString("INFO: "), log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, color.YellowString("WARNING : "), log.Ltime|log.Lshortfile)
	ErrorLogger   = log.New(os.Stderr, color.RedString("ERROR: "), log.Ltime|log.Lshortfile)
)
