package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	randomfs "github.com/TheEntropyCollective/randomfs-core"
	"github.com/gorilla/mux"
)

type Server struct {
	rfs     *randomfs.RandomFS
	router  *mux.Router
	port    string
	dataDir string
}

type StoreResponse struct {
	URL  string `json:"url"`
	Hash string `json:"hash"`
}

type StatsResponse struct {
	FilesStored     int64 `json:"files_stored"`
	BlocksGenerated int64 `json:"blocks_generated"`
	TotalSize       int64 `json:"total_size"`
	CacheHits       int64 `json:"cache_hits"`
	CacheMisses     int64 `json:"cache_misses"`
}

func NewServer(port, dataDir, ipfsAPI string) (*Server, error) {
	// Initialize RandomFS
	var rfs *randomfs.RandomFS
	var err error

	if ipfsAPI == "" {
		rfs, err = randomfs.NewRandomFSWithoutIPFS(dataDir, 100*1024*1024) // 100MB cache
	} else {
		rfs, err = randomfs.NewRandomFS(ipfsAPI, dataDir, 100*1024*1024) // 100MB cache
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initialize RandomFS: %v", err)
	}

	router := mux.NewRouter()
	server := &Server{
		rfs:     rfs,
		router:  router,
		port:    port,
		dataDir: dataDir,
	}

	server.setupRoutes()
	return server, nil
}

func (s *Server) setupRoutes() {
	// API routes
	s.router.HandleFunc("/api/v1/store", s.handleStore).Methods("POST")
	s.router.HandleFunc("/api/v1/stats", s.handleStats).Methods("GET")
	s.router.HandleFunc("/api/v1/retrieve/{hash}", s.handleRetrieve).Methods("GET")

	// File access routes
	s.router.HandleFunc("/rd/{encodedURL}", s.handleFileAccess).Methods("GET")

	// Static file serving - serve web files
	s.router.PathPrefix("/").HandlerFunc(s.handleStatic)
}

func (s *Server) handleStore(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file provided", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read file data
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	// Determine content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Store file using RandomFS
	randomURL, err := s.rfs.StoreFile(header.Filename, data, contentType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to store file: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare response
	response := StoreResponse{
		URL:  randomURL.String(),
		Hash: randomURL.RepHash,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	stats := s.rfs.GetStats()
	response := StatsResponse{
		FilesStored:     stats.FilesStored,
		BlocksGenerated: stats.BlocksGenerated,
		TotalSize:       stats.TotalSize,
		CacheHits:       stats.CacheHits,
		CacheMisses:     stats.CacheMisses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleRetrieve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]

	if hash == "" {
		http.Error(w, "No hash provided", http.StatusBadRequest)
		return
	}

	// Retrieve file from RandomFS
	data, rep, err := s.rfs.RetrieveFile(hash)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve file: %v", err), http.StatusInternalServerError)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", rep.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", rep.FileName))
	w.Header().Set("Content-Length", strconv.FormatInt(rep.FileSize, 10))

	// Write file data
	w.Write(data)
}

func (s *Server) handleFileAccess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	encodedURL := vars["encodedURL"]

	if encodedURL == "" {
		http.Error(w, "No URL provided", http.StatusBadRequest)
		return
	}

	// Decode the base64 encoded URL
	decodedBytes, err := io.ReadAll(strings.NewReader(encodedURL))
	if err != nil {
		http.Error(w, "Invalid URL encoding", http.StatusBadRequest)
		return
	}

	// Parse the RandomURL
	randomURL, err := randomfs.ParseRandomURL(string(decodedBytes))
	if err != nil {
		http.Error(w, "Invalid rd:// URL", http.StatusBadRequest)
		return
	}

	// Retrieve file from RandomFS
	data, rep, err := s.rfs.RetrieveFile(randomURL.RepHash)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve file: %v", err), http.StatusInternalServerError)
		return
	}

	// Set appropriate headers
	w.Header().Set("Content-Type", rep.ContentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", rep.FileName))
	w.Header().Set("Content-Length", strconv.FormatInt(rep.FileSize, 10))

	// Write file data
	w.Write(data)
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	webDir := filepath.Join(filepath.Dir(os.Args[0]), "web")
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		webDir = "web"
	}

	var filePath string
	if r.URL.Path == "/" {
		filePath = filepath.Join(webDir, "index.html")
	} else {
		filePath = filepath.Join(webDir, r.URL.Path)
	}

	log.Printf("[DEBUG] handleStatic: URL.Path=%q, resolved filePath=%q", r.URL.Path, filePath)

	absWebDir, _ := filepath.Abs(webDir)
	absFilePath, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(absFilePath, absWebDir) {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}

	http.ServeFile(w, r, filePath)
}

func (s *Server) Start() error {
	log.Printf("RandomFS Web Server starting on port %s", s.port)
	log.Printf("Data directory: %s", s.dataDir)
	log.Printf("Web interface available at: http://localhost:%s", s.port)

	return http.ListenAndServe(":"+s.port, s.router)
}

func main() {
	// Default configuration
	port := "8080"
	dataDir := "./data"
	ipfsAPI := "http://localhost:5001"

	// Check for environment variables
	if envPort := os.Getenv("RANDOMFS_PORT"); envPort != "" {
		port = envPort
	}
	if envDataDir := os.Getenv("RANDOMFS_DATA_DIR"); envDataDir != "" {
		dataDir = envDataDir
	}
	if envIPFSAPI := os.Getenv("RANDOMFS_IPFS_API"); envIPFSAPI != "" {
		ipfsAPI = envIPFSAPI
	} else {
		// If RANDOMFS_IPFS_API is not set, use empty string to disable IPFS
		ipfsAPI = ""
	}

	// Create server
	server, err := NewServer(port, dataDir, ipfsAPI)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start server
	if err := server.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
