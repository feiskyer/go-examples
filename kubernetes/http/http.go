package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	//var body []byte
	var response *http.Response
	var request *http.Request

	url := "https://127.0.0.1:10250/containerLogs/default/nginx/nginx?follow=true"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	//request, err := client.Dj(url)
	request, err := http.NewRequest("GET", url, nil)

	resp, err := client.Get(url)
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	return

	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}
}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}

// Wrap wraps an io.Writer into a writer that flushes after every write if
// the writer implements the Flusher interface.
func Wrap(w io.Writer) io.Writer {
	fw := &flushWriter{
		writer: w,
	}
	if flusher, ok := w.(http.Flusher); ok {
		fw.flusher = flusher
	}
	return fw
}

// flushWriter provides wrapper for responseWriter with HTTP streaming capabilities
type flushWriter struct {
	flusher http.Flusher
	writer  io.Writer
}

// Write is a FlushWriter implementation of the io.Writer that sends any buffered
// data to the client.
func (fw *flushWriter) Write(p []byte) (n int, err error) {
	n, err = fw.writer.Write(p)
	if err != nil {
		return
	}
	if fw.flusher != nil {
		fw.flusher.Flush()
	}
	return
}
