package logger

import (
	"log"
	"os"
)

var L = log.New(os.Stdout, "[chat] ", log.LstdFlags|log.Lshortfile)
