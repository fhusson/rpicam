package main

import (
	"fmt"
	"github.com/whaly/rpicam"
	"net/http"
)

func WebcamHandler(w http.ResponseWriter, r *http.Request, pm *rpicam.Manager) {
	// maxWidth := 2592
	// maxHeight := 1944
	args := ""

	width := r.FormValue("width")
	height := r.FormValue("height")
	if width != "" && height != "" {
		args = "-w " + width + " -h " + height
	}
	result := pm.NewShot(args)

	if result.Err != nil {
		http.Error(w, result.Err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image")
	w.Write(result.Data)
}

func main() {
	fmt.Printf("RPiCam Webclient\n")

	pm := rpicam.NewManager()
	go pm.Serve()

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			WebcamHandler(w, r, pm)
		})
	http.ListenAndServe(":9001", nil)
}
