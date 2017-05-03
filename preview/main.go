package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/not-quite-vacation/blog/blog"

	"golang.org/x/crypto/acme/autocert"
)

var (
	bucketName  = flag.String("bucket_name", "notquitevacation", "the google cloud store bucket to use.")
	projectName = flag.String("project_name", "notquitevacation", "the google cloud project.")
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
		HostPolicy: autocert.HostWhitelist("notquitevacation.com", "www.notquitevacation.com"),
		Cache:      b,
		/*
			Client: &acme.Client{
				DirectoryURL: "https://acme-staging.api.letsencrypt.org/directory",
			},
		*/
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

	lis := m.Listener()
	defer lis.Close()

	log.Println("Starting not quite vacation...")
	err = srv.Serve(lis)
	if err != nil {
		log.Fatal(err)
	}
}
