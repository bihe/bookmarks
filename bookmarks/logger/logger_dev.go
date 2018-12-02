// +build !prod

package logger

import (
	"log"
	"os"

	"github.com/bihe/bookmarks-go/bookmarks/conf"
)

// InitLogger performs a setup for the logging mechanism
func InitLogger(conf conf.LogConfig) {
	log.SetPrefix(LogPrefix(conf))
	log.SetOutput(os.Stdout)
}
