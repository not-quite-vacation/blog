package main

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"not-quite-vacation/blog"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

func main() {
	httpSrv := http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Connection", "close")
			url := "https://" + req.Host + req.URL.String()
			http.Redirect(w, req, url, http.StatusMovedPermanently)
		}),
	}
	go func() { log.Fatal(httpSrv.ListenAndServe()) }()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(blog.FS(false)))

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("notquitevacation.com", "www.notquitevacation.com"),
	}
	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		GetCertificate: m.GetCertificate,
	}

	srv := http.Server{
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  2 * time.Minute,
		TLSConfig:    tlsConfig,
	}

	lis, err := net.Listen("tcp", ":https")
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	err = srv.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
