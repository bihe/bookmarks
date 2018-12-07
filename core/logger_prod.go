// +build prod

package core

import (
	"io"
	"log"
	"os"

	"github.com/bihe/bookmarks/bookmarks/conf"
	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger performs a setup for the logging mechanism
func InitLogger(conf conf.LogConfig) {
	f := &lumberjack.Logger{
		Filename:   conf.Rolling.FilePath,
		MaxSize:    conf.Rolling.MaxSize,
		MaxBackups: conf.Rolling.MaxBackups,
		MaxAge:     conf.Rolling.MaxAge,
		Compress:   conf.Rolling.Compress,
	}
	mw := io.MultiWriter(os.Stdout, f)
	log.SetPrefix(LogPrefix(conf))
	log.SetOutput(mw)
}
