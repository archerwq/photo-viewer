package pvlog

import (
	"log"
	"os"
)

var PVLog = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
