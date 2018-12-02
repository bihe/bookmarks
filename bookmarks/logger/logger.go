package logger

import (
	"fmt"

	"github.com/bihe/bookmarks-go/bookmarks/conf"
)

// LogPrefix is used to display a meaningful prefix for log-messages
func LogPrefix(config conf.LogConfig) string {
	return fmt.Sprintf("[%s] ", config.Prefix)
}
