package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/ZachGill/Foreman/cmd/foreman/handlers"
)

func main() {
	var (
		waitGroup = &sync.WaitGroup{}

		httpServer = &handlers.Server{
			ServerMutex:    &sync.Mutex{},
			WaitGroup:      waitGroup,
			HTTPListenAddr: ":8080",
			HTTPLogger:     log.New(os.Stderr, "HTTPLogger: ", log.Lshortfile),

			Status: &handlers.Status{},
		}
	)

	waitGroup.Add(1)
	go httpServer.Start()
	go waitForSignal(make(chan os.Signal, 1), httpServer)
	waitGroup.Wait()
}

func waitForSignal(c <-chan os.Signal, server *handlers.Server) {
	<-c
	ctx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	server.Stop(ctx)
	cancelFunc()
}
