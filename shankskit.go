package shankskit

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/asticode/go-astilectron"
)

func StartApp(port string, routes map[string]http.HandlerFunc, appName string, transparent bool, alwaysontop bool) { // Create a new Astilectron instance

	for url, handlerfunc := range routes {
		http.HandleFunc(url, handlerfunc)
	}

	server := &http.Server{
		Addr: ":" + port,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server:", err)
		}
	}()

	url := "http://localhost:" + port

	fmt.Println("Server running on", url)

	logger := log.New(os.Stderr, "", log.LstdFlags)
	a, err := astilectron.New(logger, astilectron.Options{
		AppName:        appName,
		SingleInstance: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()
	if err := a.Start(); err != nil {
		log.Fatal(err)
	}

	frameless := false
	width := 600
	height := 800
	fullscreen := true
	// Create a new window (frameless window, no UI elements)
	w, err := a.NewWindow(url, &astilectron.WindowOptions{
		Fullscreen: &fullscreen,
		Width:      &width,
		Height:     &height,
		Transparent: &transparent, 
		Frame:      &frameless,
		AlwaysOnTop:       &alwaysontop,
	})

	if err != nil {
		log.Fatal(err)
	}

	// Show the window
	if err := w.Create(); err != nil {
		log.Fatal(err)
	}

	// Wait for the window to be closed
	a.Wait() // Blocks here until the window is closed

	// Gracefully shut down the server after window is closed
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Error shutting down server:", err)
	}

	fmt.Println("Server stopped")
}

func Respond(w http.ResponseWriter, filePath string, data interface{}) {
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func ListSubPaths(dir string) []string {
	var output []string

	// Open the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return output
	}

	// Iterate through the files and check if they are directories or files
	for _, file := range files {
		if file.IsDir() {
			subDir := dir + string(os.PathSeparator) + file.Name()
			output = append(output, ListSubPaths(subDir)...)
		} else {
			output = append(output, dir+string(os.PathSeparator)+file.Name())
		}
	}

	return output
}
