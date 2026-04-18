package lsp

import (
	"log"
	"os"
)

// debugEnabled gates verbose internal logging so qhugo doesn't spam
// stdout (and potentially leak document contents) by default.
// Set QHUGO_DEBUG=1 to enable.
var debugEnabled = os.Getenv("QHUGO_DEBUG") != ""

func dlog(format string, args ...any) {
	if debugEnabled {
		log.Printf(format, args...)
	}
}
