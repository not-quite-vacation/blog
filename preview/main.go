package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/not-quite-vacation/blog/blog"

	"golang.org/x/crypto/acme/autocert"
)

var (
	bucketName  = flag.String("bucket_name", "nqv", "the google cloud store bucket to use.")
	projectName = flag.String("project_name", "notquitvacation", "the google cloud project.")
)

func main() {
	flag.Parse()
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	b, err := newBucket(ctx, *bucketName, *projectName)
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	m := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("www.notquitevacation.com", "notquitevacation.com"),
		Cache:      b,
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		GetCertificate: m.GetCertificate,
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(blog.FS(false)))

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

	log.Println("Starting not quite vacation...")
	err = srv.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
