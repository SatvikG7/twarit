package router

import (
	"net/http"

	"github.com/satvikg7/twarit/internal/handler"
)

func Route() {
	// serve static files from the src directory
	srcDir := http.FileServer(http.Dir("./src/"))
	http.Handle("/", http.StripPrefix("/", srcDir))
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "src/send.html")
	})

	http.HandleFunc("/explore", handler.ExplorePageHandler)
	http.HandleFunc("/explore/", handler.ExplorePageHandler)

	http.Handle("/favicon.ico", http.StripPrefix("/", srcDir))
	http.Handle("/logo.png", http.StripPrefix("/", srcDir))
	http.HandleFunc("/files/", handler.ExploreHandler)

	// Handle file uploads with progress tracking
	http.HandleFunc("/upload", handler.UploadHandler)
	http.HandleFunc("/progress", handler.ProgressHandler)
	http.HandleFunc("/ws-progress", handler.WSProgressHandler)
}
