// Dependency: golang.org/x/crypto/acme/autocert
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

const (
	domain = "<your-domain.com"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `<html><body>Hello, world!</body></html>`)
}

func makeServerFromMux(mux *http.ServeMux) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func makeHTTPServer() *http.Server {
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleIndex)
	return makeServerFromMux(mux)

}

func makeHTTPToHTTPSRedirectServer() *http.Server {
	handleRedirect := func(w http.ResponseWriter, req *http.Request) {
		newURI := "https://" + req.Host + req.URL.String()
		http.Redirect(w, req, newURI, http.StatusFound)
	}
	mux := &http.ServeMux{}
	mux.HandleFunc("/", handleRedirect)
	return makeServerFromMux(mux)
}

func main() {
	redirect := true
	flag.BoolVar(&redirect, "redirect", true, "if true, redirect http to https")
	flag.Parse()

	// also start http.
	if redirect {
		httpSrv := makeHTTPToHTTPSRedirectServer()
		httpSrv.Addr = ":80"
		fmt.Printf("Starting HTTP server on :80\n")
		go func() {
			err := httpSrv.ListenAndServe()
			if err != nil {
				log.Fatalf("httpSrv.ListenAndServe() failed with %s", err)
			}
		}()
	}

	var httpsSrv *http.Server
	dataDir := "."
	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(domain),
		Cache:      autocert.DirCache(dataDir),
	}

	httpsSrv = makeHTTPServer()
	httpsSrv.Addr = ":443"
	httpsSrv.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

	fmt.Printf("Starting HTTPS server on %s\n", httpsSrv.Addr)
	err := httpsSrv.ListenAndServeTLS("", "")
	if err != nil {
		log.Fatalf("httpsSrv.ListendAndServeTLS() failed with %s", err)
	}
}
