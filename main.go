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

	"github.com/bihe/bookmarks/api"
	"github.com/bihe/bookmarks/core"
	"github.com/wangii/emoji"
)

var (
	// Version exports the application version
	Version string
	// Build provides information about the application build
	Build string
)

// graceful stop taken from https://gist.github.com/peterhellberg/38117e546c217960747aacf689af3dc2
func main() {
	srv := setupServer()
	go func() {
		log.Printf("%s Starting server ...", emoji.EmojiTagToUnicode(`:rocket:`))
		log.Printf("%s Listening on '%s'", emoji.EmojiTagToUnicode(`:computer:`), srv.Addr)
		log.Printf("%s Version: '%s-%s'\n", emoji.EmojiTagToUnicode(`:bookmark:`), Version, Build)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	graceful(srv, 5*time.Second)
}

func setupServer() *http.Server {
	host, port, basePath, configFile := configFromEnv()
	conf := configFromFile(basePath, configFile)
	core.InitLogger(conf.Log)
	h := api.SetupAPI(conf, Version, Build)
	addr := host + ":" + port
	srv := &http.Server{Addr: addr, Handler: h}
	return srv
}

func graceful(hs *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	log.Printf("\nShutdown with timeout: %s\n", timeout)
	if err := hs.Shutdown(ctx); err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Println("Server stopped")
	}
}

func configFromEnv() (host, port, basePath, configFile string) {
	port = getOrDefault("APPLICATION_SERVER_PORT", "3000")
	host = getOrDefault("APPLICATION_SEVER_HOST", "localhost")
	basePath = getOrDefault("APPLICATION_BASE_PATH", "./")
	configFile = getOrDefault("APPLICATION_CONFIG_FILE", "application.json")
	return
}

func configFromFile(appBasePath, configFileName string) core.Configuration {
	dir, err := filepath.Abs(appBasePath)
	if err != nil {
		panic(fmt.Sprintf("Could not get the application basepath: %v", err))
	}
	f, err := os.Open(path.Join(dir, configFileName))
	if err != nil {
		panic(fmt.Sprintf("Could not open specific config file '%s': %v", configFileName, err))
	}
	defer f.Close()

	config, err := core.Settings(f)
	if err != nil {
		panic(fmt.Sprintf("Could not get server config values from file '%s': %v", configFileName, err))
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
