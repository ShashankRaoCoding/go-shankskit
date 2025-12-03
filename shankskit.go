package shankskit

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"text/template"
	"io/fs"
	"github.com/asticode/go-astilectron"
)

func StartApp(settings AppSettings) (*astilectron.Window, *astilectron.Astilectron, *http.Server) { // Create a new Astilectron instance
	appName := settings.AppName
	transparent := settings.Transparent
	alwaysontop := settings.AlwaysOnTop
	fullscreen := settings.Fullscreen
	visibleUI := settings.VisibleUI
	routes := settings.Routes
	width := settings.Width
	height := settings.Height

	ln, _ := net.Listen("tcp", ":0")
	port := fmt.Sprintf("%v", ln.Addr().(*net.TCPAddr).Port)
	mux := http.NewServeMux()
	for url, handler := range routes {
		mux.HandleFunc(url, handler)
	}
	server := &http.Server{
		Handler: mux,
	}

	go func() {
		if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server:", err)
		}
	}()

	url := "http://localhost:" + port

	fmt.Println("Server running on", url)

	logger := log.New(os.Stderr, "", log.LstdFlags)
	a, err := astilectron.New(logger, astilectron.Options{
		AppName:        appName,
		SingleInstance: false,
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
	Routes      map[string]http.HandlerFunc
	AppName     string
	Width       int
	Height      int
	Fullscreen  bool
	VisibleUI   bool
	Transparent bool
	AlwaysOnTop bool
}

func RespondWithEmbed(w http.ResponseWriter, filePath string, fs fs.FS, data interface{}) {
	tmpl, err := template.ParseFS(fs, filePath)
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

func StartServer(routes map[string]http.HandlerFunc) string {

	ln, _ := net.Listen("tcp", ":0")
	port := fmt.Sprintf("%v", ln.Addr().(*net.TCPAddr).Port)
	mux := http.NewServeMux()
	for url, handler := range routes {
		mux.HandleFunc(url, handler)
	}
	server := &http.Server{
		Handler: mux,
	}

	go func() {
		if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
			fmt.Println("Error starting server:", err)
		}
	}()

	log.Printf("Access the server at localhost:%s", port)

	return port

}
