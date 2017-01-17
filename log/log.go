package log

import (
	"log"
	"os"
)

var Info = log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
var Debug = log.New(os.Stdout, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile)
