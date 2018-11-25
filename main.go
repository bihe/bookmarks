package main

import (
	"net/http"
	"os"

	"github.com/bihe/bookmarks-go/internal/bookmarks"
)

func serverDefaults() (host, port, basePath, configFile string) {

	port = getOrDefault("APPLICATION_SERVER_PORT", "3000")
	host = getOrDefault("APPLICATION_SEVER_HOST", "localhost")
	basePath = getOrDefault("APPLICATION_BASE_PATH", "./")
	basePath = getOrDefault("APPLICATION_CONFIG_FILE", "application.json")
	return
}

func getOrDefault(env, def string) string {
	value := os.Getenv(env)
	if value == "" {
		value = def
	}
	return value
}

func main() {
	// either get the server host:port from the environment
	// or use sensible defaults
	host, port, basePath, configFile := serverDefaults()
	router := bookmarks.SetupRouter(basePath, configFile)
	http.ListenAndServe(host+":"+port, router)
}
