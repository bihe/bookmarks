package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"
	"time"

	"github.com/bihe/bookmarks-go/bookmarks"
	"github.com/bihe/bookmarks-go/bookmarks/conf"
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

func configFromEnv() (host, port, basePath, configFile string) {

	port = getOrDefault("APPLICATION_SERVER_PORT", "3000")
	host = getOrDefault("APPLICATION_SEVER_HOST", "localhost")
	basePath = getOrDefault("APPLICATION_BASE_PATH", "./")
	configFile = getOrDefault("APPLICATION_CONFIG_FILE", "application.json")
	return
}

func configFromFile(appBasePath, configFileName string) conf.Configuration {
	dir, err := filepath.Abs(appBasePath)
	if err != nil {
		panic("Could not get the application basepath!")
	}
	configFile, err := os.Open(path.Join(dir, configFileName))
	if err != nil {
		panic(fmt.Sprintf("Specified config file '%s' missing!", configFileName))
	}
	defer configFile.Close()

	config, err := conf.Settings(configFile)
	if err != nil {
		panic("No config values available to start the server. Missing config.json file!")
	}
	return *config
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
	host, port, basePath, configFile := configFromEnv()
	conf := configFromFile(basePath, configFile)
	s := bookmarks.SetupAPI(conf)
	addr := host + ":" + port
	srv := &http.Server{Addr: addr, Handler: s}
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
