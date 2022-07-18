package logger

import (
	"bytes"
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	buf           bytes.Buffer
	InfoLogger    = log.New(os.Stdout, color.GreenString("INFO: "), log.Ltime|log.Lshortfile)
	WarningLogger = log.New(os.Stdout, color.YellowString("WARNING : "), log.Ltime)
	ErrorLogger   = log.New(os.Stderr, color.RedString("ERROR: "), log.Ltime)
)
