package handlers

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Server starts Foreman's HTTP server and sets up routing to individual handlers
type Server struct {
	ServerMutex *sync.Mutex
	WaitGroup   *sync.WaitGroup

	HTTPListenAddr string
	HTTPLogger     *log.Logger

	Status http.Handler

	httpServer *http.Server
}

// Start spins up the server
func (server *Server) Start() {
	router := server.Router()

	server.ServerMutex.Lock()
	server.httpServer = &http.Server{
		Addr:     server.HTTPListenAddr,
		Handler:  router,
		ErrorLog: server.HTTPLogger,
	}
	server.ServerMutex.Unlock()

	log.Println("Starting server")
	if err := server.httpServer.ListenAndServe(); err != nil {
		server.HTTPLogger.Println("Unable to listen and server", err.Error())
	}
}

// Stop stops the server
func (server *Server) Stop(ctx context.Context) {
	server.ServerMutex.Lock()
	defer server.ServerMutex.Unlock()

	err := server.httpServer.Shutdown(ctx)

	if err != nil {
		server.HTTPLogger.Print("Unable to shutdown. Error:", err.Error())
	}

	log.Println("Stopping server")
	server.WaitGroup.Done()
}

// Router defines routes to individual handlers
func (server *Server) Router() *mux.Router {
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./cmd/foreman/static")))
	r.Handle("/status", server.Status).Methods("GET")

	return r
}
