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

	http.HandleFunc("/explore", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "src/explore.html")
	})
	
	http.Handle("/favicon.ico", http.StripPrefix("/", srcDir))
	http.Handle("/logo.png", http.StripPrefix("/", srcDir))
	http.HandleFunc("/files/", handler.ExploreHandler)

	// Handle file uploads
	http.HandleFunc("/upload", handler.UploadHandler)
}
