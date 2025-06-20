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

func StartApp(settings AppSettings) (*astilectron.Window, *astilectron.Astilectron, *http.Server) { // Create a new Astilectron instance
	appName := settings.AppName
	port := settings.Port
	transparent := settings.Transparent
	alwaysontop := settings.AlwaysOnTop
	fullscreen := settings.Fullscreen
	visibleUI := settings.VisibleUI
	routes := settings.Routes
	width := settings.Width
	height := settings.Height

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
	if err := a.Start(); err != nil {
		log.Fatal(err)
	}

	// Create a new window (frameless window, no UI elements)
	w, err := a.NewWindow(url, &astilectron.WindowOptions{
		Fullscreen:  &fullscreen,
		Width:       &width,
		Height:      &height,
		Transparent: &transparent,
		Frame:       &visibleUI,
		AlwaysOnTop: &alwaysontop,
	})

	if err != nil {
		log.Fatal(err)
	}

	// Show the window
	if err := w.Create(); err != nil {
		log.Fatal(err)
	}

	return w, a, server
}

func HandleShutDown(a *astilectron.Astilectron, server *http.Server) {
	defer a.Close()
	// Wait for the window to be closed
	a.Wait() // Blocks here until the window is closed

	// Gracefully shut down the server after window is closed
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Error shutting down server:", err)
	}

	fmt.Println("Server stopped")

}

type AppSettings struct {
	Port        string
	Routes      map[string]http.HandlerFunc
	AppName     string
	Width       int
	Height      int
	Fullscreen  bool
	VisibleUI   bool
	Transparent bool
	AlwaysOnTop bool
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

