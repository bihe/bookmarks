package core

import (
	"fmt"
)

// LogPrefix is used to display a meaningful prefix for log-messages
func LogPrefix(config LogConfig) string {
	return fmt.Sprintf("[%s] ", config.Prefix)
}
