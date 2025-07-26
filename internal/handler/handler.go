package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Progress tracking structures
type UploadProgress struct {
	TotalBytes     int64   `json:"totalBytes"`
	UploadedBytes  int64   `json:"uploadedBytes"`
	Percentage     float64 `json:"percentage"`
	Speed          int64   `json:"speed"`         // bytes per second
	RemainingTime  int64   `json:"remainingTime"` // seconds
	CurrentFile    string  `json:"currentFile"`
	CompletedFiles int     `json:"completedFiles"`
	TotalFiles     int     `json:"totalFiles"`
	Status         string  `json:"status"` // "uploading", "completed", "error"
}

var (
	uploadSessions = make(map[string]*UploadProgress)
	sessionMutex   = sync.RWMutex{}
	upgrader       = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for development
		},
	}
)

// ProgressReader wraps an io.Reader to track progress
type ProgressReader struct {
	reader     io.Reader
	total      int64
	read       int64
	progress   *UploadProgress
	sessionID  string
	startTime  time.Time
	lastUpdate time.Time
	baseBytes  int64 // bytes already processed from previous files
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.read += int64(n)

	now := time.Now()
	if now.Sub(pr.lastUpdate) > 100*time.Millisecond { // Update every 100ms
		sessionMutex.Lock()
		totalProcessed := pr.baseBytes + pr.read
		pr.progress.UploadedBytes = totalProcessed
		if pr.progress.TotalBytes > 0 {
			pr.progress.Percentage = float64(totalProcessed) / float64(pr.progress.TotalBytes) * 100
		}

		elapsed := now.Sub(pr.startTime).Seconds()
		if elapsed > 0 {
			pr.progress.Speed = int64(float64(totalProcessed) / elapsed)
			if pr.progress.Speed > 0 {
				remaining := (pr.progress.TotalBytes - totalProcessed) / pr.progress.Speed
				pr.progress.RemainingTime = remaining
			}
		}
		sessionMutex.Unlock()
		pr.lastUpdate = now
	}

	return n, err
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Generate session ID for progress tracking
	sessionID := fmt.Sprintf("%d", time.Now().UnixNano())

	// Enable CORS for progress endpoint
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Session-ID")
	w.Header().Set("X-Session-ID", sessionID)

	// Get content length for progress tracking
	contentLength := r.ContentLength
	if contentLength <= 0 {
		http.Error(w, "Content-Length required", http.StatusBadRequest)
		return
	}

	// Initialize progress tracking
	sessionMutex.Lock()
	uploadSessions[sessionID] = &UploadProgress{
		TotalBytes:     contentLength,
		UploadedBytes:  0,
		Percentage:     0,
		Speed:          0,
		RemainingTime:  0,
		CurrentFile:    "",
		CompletedFiles: 0,
		TotalFiles:     0,
		Status:         "uploading",
	}
	progress := uploadSessions[sessionID]
	sessionMutex.Unlock()

	startTime := time.Now()

	// 32 MB is the default used by FormFile()
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		sessionMutex.Lock()
		progress.Status = "error"
		sessionMutex.Unlock()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get a reference to the fileHeaders.
	files := r.MultipartForm.File["file"]

	sessionMutex.Lock()
	progress.TotalFiles = len(files)
	sessionMutex.Unlock()

	// Calculate total file size for more accurate progress
	var totalFileSize int64
	for _, fileHeader := range files {
		totalFileSize += fileHeader.Size
	}

	sessionMutex.Lock()
	progress.TotalBytes = totalFileSize
	sessionMutex.Unlock()

	var processedBytes int64 = 0

	for i, fileHeader := range files {
		sessionMutex.Lock()
		progress.CurrentFile = fileHeader.Filename
		sessionMutex.Unlock()

		// Open the file
		file, err := fileHeader.Open()
		if err != nil {
			sessionMutex.Lock()
			progress.Status = "error"
			sessionMutex.Unlock()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		err = os.MkdirAll("./Received", os.ModePerm)
		if err != nil {
			sessionMutex.Lock()
			progress.Status = "error"
			sessionMutex.Unlock()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		f, err := os.Create(fmt.Sprintf("./Received/%s", fileHeader.Filename))
		if err != nil {
			sessionMutex.Lock()
			progress.Status = "error"
			sessionMutex.Unlock()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()

		// Create a progress reader for this file
		progressReader := &ProgressReader{
			reader:     file,
			total:      fileHeader.Size,
			read:       0,
			progress:   progress,
			sessionID:  sessionID,
			startTime:  startTime,
			lastUpdate: time.Now(),
			baseBytes:  processedBytes,
		}

		_, err = io.Copy(f, progressReader)
		if err != nil {
			sessionMutex.Lock()
			progress.Status = "error"
			sessionMutex.Unlock()
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		processedBytes += fileHeader.Size
		sessionMutex.Lock()
		progress.CompletedFiles = i + 1
		progress.UploadedBytes = processedBytes
		progress.Percentage = float64(processedBytes) / float64(totalFileSize) * 100
		sessionMutex.Unlock()
	}

	// Mark upload as completed
	sessionMutex.Lock()
	progress.Status = "completed"
	progress.Percentage = 100
	progress.UploadedBytes = totalFileSize
	sessionMutex.Unlock()

	// Return success response with session ID
	response := map[string]interface{}{
		"message":   "Upload successful",
		"sessionId": sessionID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ProgressHandler provides upload progress for a session
func ProgressHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Content-Type", "application/json")

	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		http.Error(w, "Session ID required", http.StatusBadRequest)
		return
	}

	sessionMutex.RLock()
	progress, exists := uploadSessions[sessionID]
	sessionMutex.RUnlock()

	if !exists {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(progress)
}

// WebSocket handler for real-time progress updates
func WSProgressHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	sessionID := r.URL.Query().Get("sessionId")
	if sessionID == "" {
		conn.WriteJSON(map[string]string{"error": "Session ID required"})
		return
	}

	// Send progress updates every 100ms
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		sessionMutex.RLock()
		progress, exists := uploadSessions[sessionID]
		sessionMutex.RUnlock()

		if !exists {
			conn.WriteJSON(map[string]string{"error": "Session not found"})
			return
		}

		if err := conn.WriteJSON(progress); err != nil {
			return
		}

		// Stop sending updates when upload is completed or errored
		if progress.Status == "completed" || progress.Status == "error" {
			return
		}
	}
}

// struct to store file info
type fileInfo struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

func ExplorePageHandler(w http.ResponseWriter, r *http.Request) {
	// Simply serve the explore.html file
	// The JavaScript will handle extracting the path from the URL
	http.ServeFile(w, r, "src/explore.html")
}

func ExploreHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/files"):]

	// Convert URL path to Windows path
	if path == "" || path == "/" {
		path = "C:/"
	} else {
		// Remove leading slash and ensure proper Windows path format
		if path[0] == '/' {
			path = path[1:]
		}
		if len(path) >= 2 && path[1] == ':' {
			// Already a Windows path like C:/something
		} else {
			// Prepend C:/ if not a full Windows path
			path = "C:/" + path
		}
	}

	// Convert forward slashes to backslashes for Windows
	path = strings.ReplaceAll(path, "/", "\\")

	files, err := os.Stat(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if files.IsDir() {
		// get all files in directory
		files, err := os.ReadDir(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// store file info
		var fileInfos []fileInfo
		for _, file := range files {
			fileInfos = append(fileInfos, fileInfo{
				Name:  file.Name(),
				IsDir: file.IsDir(),
			})
		}

		// send file info as json
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(fileInfos)
		return
	}

	// serve file
	http.ServeFile(w, r, path)
}
