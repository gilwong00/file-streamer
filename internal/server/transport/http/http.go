package httptransport

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gilwong00/file-streamer/internal/pkg/config"
	"github.com/gilwong00/file-streamer/internal/pkg/fileutils"
	"github.com/gilwong00/file-streamer/internal/pkg/storage"
	"github.com/minio/minio-go/v7"
)

type httpServer struct {
	ctx              context.Context
	port             int
	uploadFolderName string
	storageClient    storage.Client
	bucketName       string
}

const (
	MinCompressionSize = 8 * 1024 // 8kb
)

func NewHttpServer(
	ctx context.Context,
	config *config.Config,
	storageClient storage.Client,
	bucketName string,
) *httpServer {
	return &httpServer{
		ctx:              ctx,
		port:             config.HTTPServerPort,
		uploadFolderName: config.FileDirectoryName,
		storageClient:    storageClient,
		bucketName:       bucketName,
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
	// fileName := r.PathValue("filename")
	// if err := fileutils.ValidateFileName(fileName); err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// _, stat, file, err := fileutils.OpenFile(fileName, s.uploadFolderName)
	// if err != nil {
	// 	fileutils.HandleFileOpenError(w, err)
	// 	return
	// }
	// defer file.Close()
	// w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	// w.WriteHeader(http.StatusOK)

	// using minio instead of writing to disk
	fileName := r.PathValue("filename")
	if fileName == "" {
		http.Error(w, "missing fileName", http.StatusBadRequest)
		return
	}
	if err := fileutils.ValidateFileName(fileName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	info, err := s.storageClient.GetObjectInfo(r.Context(), s.bucketName, fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size))
	w.Header().Set("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusOK)
}

func (s *httpServer) getHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.PathValue("filename")
	if fileName == "" {
		http.Error(w, "missing fileName", http.StatusBadRequest)
		return
	}
	if err := fileutils.ValidateFileName(fileName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	obj, err := s.storageClient.GetObject(r.Context(), s.bucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	defer obj.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err := io.Copy(w, obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
