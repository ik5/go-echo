package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// DefaultPort is the default port to use if once is not specified by the SERVER_PORT environment variable
const DefaultPort = "8080"

func getServerPort() string {
	port := os.Getenv("SERVER_PORT")
	if port != "" {
		return port
	}

	return DefaultPort
}

// EchoHandler echos back the request as a response
func EchoHandler(writer http.ResponseWriter, request *http.Request) {
	log.Println("---------------------------------------------------------------")
	log.Printf("request: %+v\n", request)
	log.Printf("RemoteAddress: %s. Path: %+v\n", request.RemoteAddr, request.URL)
	log.Printf("Headers: %+v\n", request.Header)
	if request.ContentLength > 0 {
		defer request.Body.Close()
		body, err := ioutil.ReadAll(request.Body)
		log.Printf("Payload: %s | err: %s\n", body, err)
	} else {
		log.Println("No Payload")
	}
	log.Println("---------------------------------------------------------------")

	request.Write(writer)
}

func main() {

	log.Println("starting server, listening on port " + getServerPort())

	http.HandleFunc("/", EchoHandler)
	http.ListenAndServe(":"+getServerPort(), nil)
}
