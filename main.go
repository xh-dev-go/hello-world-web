package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"net/http"
	"os"
)

type Headers map[string]interface{}
type ResponseBody struct {
	Host    string  `yaml:"host"`
	Url     string  `yaml:"url"`
	Ip      string  `yaml:"ip"`
	Referer string  `yaml:"referer"`
	Headers Headers `yaml:"headers"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		var response ResponseBody
		var headers = make(Headers)
		for k, v := range request.Header {
			headers[k] = v
		}
		response.Host = request.Host
		response.Url = request.RequestURI
		response.Ip = request.RemoteAddr
		response.Referer = request.Referer()
		response.Headers = headers
		b, err := yaml.Marshal(&response)
		if err != nil {
			panic(err)
		}
		_, err = writer.Write(b)
		if err != nil {
			panic(err)
		}
	})
	portStr := fmt.Sprintf(":%v", port)
	err := http.ListenAndServe(portStr, nil)
	if err != nil {
		fmt.Println(err)
	}

}
