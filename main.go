package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bihe/bookmarks-go/internal/bookmarks"
	"github.com/wangii/emoji"
)

// graceful stop taken from https://gist.github.com/peterhellberg/38117e546c217960747aacf689af3dc2
func main() {
	srv, logger := setup()

	go func() {
		logger.Printf("%s Starting server ...", emoji.EmojiTagToUnicode(`:rocket:`))
		logger.Printf("%s Listening on '%s' %s\n", emoji.EmojiTagToUnicode(`:computer:`), srv.Addr, emoji.EmojiTagToUnicode(`:+1:`))

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	graceful(srv, logger, 5*time.Second)
}

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

func setup() (*http.Server, *log.Logger) {
	// either get the server host:port from the environment
	// or use sensible defaults
	host, port, basePath, configFile := serverDefaults()
	r := bookmarks.SetupRouter(basePath, configFile)
	addr := host + ":" + port
	srv := &http.Server{Addr: addr, Handler: r}
	return srv, log.New(os.Stdout, "", 0)
}

func graceful(hs *http.Server, logger *log.Logger, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	logger.Printf("\nShutdown with timeout: %s\n", timeout)
	if err := hs.Shutdown(ctx); err != nil {
		logger.Printf("Error: %v\n", err)
	} else {
		logger.Println("Server stopped")
	}
}
