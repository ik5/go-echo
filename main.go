package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// DefaultPort is the default port to use if once is not specified by the SERVER_PORT environment variable
const DefaultPort = "8080"

// PeekReadCloser return a copy of a reader without loosing the original content
func PeekReadCloser(stream *io.ReadCloser, cpy *[]byte) {
	*cpy, _ = ioutil.ReadAll(*stream)
	*stream = ioutil.NopCloser(bytes.NewReader(*cpy))
}

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
		body := make([]byte, request.ContentLength+1)
		PeekReadCloser(&request.Body, &body)
		log.Printf("Payload: %s\n", body)

		contentType := strings.Split(
			strings.ToLower(request.Header.Get("content-type")),
			";")
		switch contentType[0] {
		case "application/x-www-form-urlencoded":
			err := request.ParseForm()
			log.Printf("Form: %+v |err: %s\n", request.Form, err)
			log.Printf("PostForm: %+v\n", request.PostForm)
		case "multipart/form-data":
			err := request.ParseMultipartForm(http.DefaultMaxHeaderBytes)
			log.Printf("Form: %+v |err: %s\n", *request.MultipartForm, err)
		}
	} else {
		log.Println("No Payload")
	}
	log.Println("---------------------------------------------------------------\n")

	request.Write(writer)
}

func doServer() *http.Server {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", EchoHandler)

	var handler http.Handler = mux

	srv := &http.Server{
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second, // introduced in Go 1.8
		Handler:      handler,
		Addr:         fmt.Sprintf(":%s", getServerPort()),
	}

	return srv
}

func main() {
	log.Println("Starting server, listening on port " + getServerPort())

	doServer().ListenAndServe()
}
