package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		hostname, err := os.Hostname()
		if err!=nil {
			fmt.Fprintf(writer, "hello world")
		}
		fmt.Fprintf(writer, "hello from %s", hostname)
	})
	err:=http.ListenAndServe(":8080", nil)
	if err!=nil {
		fmt.Println(err)
	}

}
