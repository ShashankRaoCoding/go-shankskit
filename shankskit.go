package shankskit

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/asticode/go-astilectron"
	"syscall"
)

// StartApp launches the app and returns the window handle (HWND)
func StartApp(port string, routes map[string]http.HandlerFunc, appName string, transparent bool, alwaysontop bool) (uintptr, error) {
	// Set up routes
	for url, handlerfunc := range routes {
		http.HandleFunc(url, handlerfunc)
	}

	server := &http.Server{Addr: ":" + port}
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
		return 0, err
	}
	defer a.Close()
	if err := a.Start(); err != nil {
		return 0, err
	}

	frameless := false
	width := 600
	height := 800
	fullscreen := true

	w, err := a.NewWindow(url, &astilectron.WindowOptions{
		Fullscreen:  &fullscreen,
		Width:       &width,
		Height:      &height,
		Transparent: &transparent,
		Frame:       &frameless,
		AlwaysOnTop: &alwaysontop,
	})
	if err != nil {
		return 0, err
	}

	// Get HWND from the window handle
	windowID := w.Handle()
	hwnd := getHWNDFromWindowID(windowID)
	if hwnd == 0 {
		return 0, fmt.Errorf("failed to retrieve HWND")
	}

	// Show the window
	if err := w.Create(); err != nil {
		return 0, err
	}

	// Wait for the window to be closed
	go func() {
		a.Wait()
		server.Shutdown(context.Background())
		fmt.Println("Server stopped")
	}()

	return hwnd, nil
}

// getHWNDFromWindowID converts an Astilectron window ID to a real Windows HWND
func getHWNDFromWindowID(windowID int) uintptr {
	// Convert window ID to string (Astilectron uses window names)
	windowName := fmt.Sprintf("Astilectron %d", windowID)

	// Find the HWND using Win32 API
	hwnd, _, _ := syscall.NewLazyDLL("user32.dll").NewProc("FindWindowW").Call(0, syscall.StringToUTF16Ptr(windowName))
	return hwnd
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
