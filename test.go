package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	r := gin.Default()

	// Ping handler
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	HTTPPORT := os.Getenv("HTTPPORT")
	if HTTPPORT == "" {
		HTTPPORT = "1880"
	}
	m := autocert.Manager{
		Prompt: autocert.AcceptTOS,

		HostPolicy: autocert.HostWhitelist("scale.rosof.tech"),
		// Cache:      autocert.DirCache("/var/www/.cache"),
	}
	srv := &http.Server{
		Addr:    ":" + HTTPPORT,
		Handler: r,
		TLSConfig: &tls.Config{
			GetCertificate: m.GetCertificate,
		},
	}

	// service connections
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	} else {
		println("started listning: %s\n", HTTPPORT)
	}
}
