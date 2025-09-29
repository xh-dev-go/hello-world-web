package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/xh-dev-go/hello-world-web/interfaces"
	"gopkg.in/yaml.v2"
)

func LaunchServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		response := interfaces.ResponseBody{
			Host:    request.Host,                       // Host header
			URL:     request.RequestURI,                 // Request URI
			Ip:      request.RemoteAddr,                 // IP address of the client
			Referer: request.Referer(),                  // Referer header
			Headers: interfaces.Headers(request.Header), // All request headers
		}

		format := request.URL.Query().Get("format")
		if format == "" {
			format = "yaml" // Default to yaml
		}

		var data []byte
		var err error

		switch format {
		case "json":
			writer.Header().Set("Content-Type", "application/json")
			data, err = json.MarshalIndent(&response, "", "  ")
		case "yaml":
			writer.Header().Set("Content-Type", "application/x-yaml")
			data, err = yaml.Marshal(&response)
		default:
			http.Error(writer, fmt.Sprintf("error: unsupported format '%s'", format), http.StatusBadRequest)
			return
		}

		if err != nil {
			log.Printf("Error marshaling response: %v", err)
			http.Error(writer, "Error creating response", http.StatusInternalServerError)
			return
		}

		_, err = writer.Write(data)
		if err != nil {
			log.Printf("Error writing response: %v", err)
		}
	})
	portStr := fmt.Sprintf(":%v", port)
	log.Printf("Starting server on port %s", port)
	err := http.ListenAndServe(portStr, nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}

}
