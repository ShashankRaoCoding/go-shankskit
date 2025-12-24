package shankskit

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"text/template"
	"io/fs"
)

func StartServer(routes map[string]http.HandlerFunc) (string, *http.Server) {

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

	return port, server 

}

func RespondWithEmbed(w http.ResponseWriter, filePath string, fs fs.FS, data interface{}) {
    tmpl, err := template.New(filePath.Base(filePath)).ParseFS(fs, filePath)
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


