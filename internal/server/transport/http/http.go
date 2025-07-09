package httptransport

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gilwong00/file-streamer/internal/pkg/config"
	"github.com/gilwong00/file-streamer/internal/pkg/fileutils"
)

type httpServer struct {
	ctx              context.Context
	port             int
	uploadFolderName string
}

const (
	MinCompressionSize = 8 * 1024 // 8kb
)

func NewHttpServer(ctx context.Context, config *config.Config) *httpServer {
	return &httpServer{
		ctx:              ctx,
		port:             config.HTTPServerPort,
		uploadFolderName: config.FileDirectoryName,
	}
}

func (s *httpServer) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("HEAD /file/{fileName}", s.headHandler)
	mux.HandleFunc("GET /file/{fileName}", s.getHandler)

	server := http.Server{
		Addr:         fmt.Sprintf(":%v", s.port),
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	// start the server
	go func() {
		log.Printf("Starting server on port %v\n", s.port)
		if err := server.ListenAndServe(); err != nil {
			fmt.Printf("Error starting server: %s", err.Error())
			os.Exit(1)
		}
	}()
	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGINT)
	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}

func (s *httpServer) headHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue("filename")
	if err := fileutils.ValidateFileName(fileName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, stat, file, err := fileutils.OpenFile(fileName, s.uploadFolderName)
	if err != nil {
		fileutils.HandleFileOpenError(w, err)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	w.WriteHeader(http.StatusOK)
}

func (s *httpServer) getHandler(w http.ResponseWriter, r *http.Request) {}
