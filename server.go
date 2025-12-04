package shankskit

import (
	"fmt"
	"log"
	"net"
	"net/http"
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
