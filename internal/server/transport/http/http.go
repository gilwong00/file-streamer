package httptransport

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gilwong00/file-streamer/internal/pkg/config"
	"github.com/gilwong00/file-streamer/internal/pkg/fileutils"
	"github.com/gilwong00/file-streamer/internal/pkg/storage"
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
	if err := fileutils.ValidateFileName(fileName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	info, err := s.storageClient.GetObjectInfo(r.Context(), s.bucketName, fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Parse Range header (supports bytes=start-end, bytes=start-, bytes=-suffix)
	start, end, err := parseRange(r.Header.Get("Range"), info.Size)
	if err != nil {
		// If parsing fails or header is missing, default to sending the entire file
		start = 0           // Start at the first byte
		end = info.Size - 1 // End at the last byte (zero-based index)
	}
	obj, err := s.storageClient.GetObject(
		r.Context(),
		s.bucketName,
		fileName, storage.GetObjectOptions{
			Start: start,
			End:   end,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer obj.Close()
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, info.Size))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusPartialContent)
	// Optional gzip compression
	if shouldCompress(r, fileName, end-start+1) {
		gz := gzip.NewWriter(w)
		defer gz.Close()
		io.Copy(gz, obj)
	} else {
		io.Copy(w, obj)
	}
}

// parseRange parses a simple "Range: bytes=start-end" header.
// Returns start and end (inclusive) byte offsets.
func parseRange(s string, size int64) (int64, int64, error) {
	if s == "" {
		return 0, size - 1, nil
	}
	if !strings.HasPrefix(s, "bytes=") {
		return 0, size - 1, fmt.Errorf("invalid range unit")
	}
	parts := strings.Split(strings.TrimPrefix(s, "bytes="), "-")
	if len(parts) != 2 {
		return 0, size - 1, fmt.Errorf("invalid range format")
	}
	var start, end int64
	var err error
	if parts[0] == "" {
		// suffix range: "-500"
		suffix, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return 0, size - 1, err
		}
		if suffix > size {
			suffix = size
		}
		start = size - suffix
		end = size - 1
	} else {
		start, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, size - 1, err
		}
		if parts[1] != "" {
			end, err = strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return 0, size - 1, err
			}
		} else {
			end = size - 1
		}
	}
	if start < 0 || end >= size || start > end {
		return 0, size - 1, fmt.Errorf("invalid range values")
	}
	return start, end, nil
}

func shouldCompress(r *http.Request, fileName string, size int64) bool {
	// Only compress if client supports gzip
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		return false
	}
	// Skip compression for already-compressed file types
	compressedExts := []string{".gz", ".zip", ".rar", ".7z", ".mp4", ".mp3", ".png", ".jpg", ".jpeg"}
	for _, ext := range compressedExts {
		if strings.HasSuffix(strings.ToLower(fileName), ext) {
			return false
		}
	}
	// Optionally: skip small files (compression overhead not worth it)
	if size < 1024 {
		return false
	}
	return true
}
