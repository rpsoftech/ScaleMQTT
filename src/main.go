package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rpsoftech/scaleMQTT/src/db"
	global "rpsoftech/scaleMQTT/src/global"
	"rpsoftech/scaleMQTT/src/hooks"
	"rpsoftech/scaleMQTT/src/routes"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/rs/zerolog"
)

func main() {
	defer println("DEFER TESTINGF")
	defer db.DBClassObject.CloseConnection()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	r := gin.Default()
	routes.AddRoutes(r)
	server := mqtt.New(nil)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	l := server.Log.Level(zerolog.DebugLevel)
	server.Log = &l
	global.Logger = &l
	global.MQTTserver = server
	server.Log.Debug().Bytes("JWTKEY", global.JWTKEY).Send()
	err := server.AddHook(new(hooks.MQTTHooks), &hooks.Options{
		DB: db.DBClassObject,
	})
	if err != nil {
		log.Fatal(err)
	}

	// err = server.AddHook(new(debug.Hook), &debug.Options{
	// 	// ShowPacketData: true,
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	MQTTTCPPORT := os.Getenv("MQTTTCPPORT")
	if MQTTTCPPORT == "" {
		MQTTTCPPORT = "1883"
	}
	tcp := listeners.NewTCP("t1", ":"+MQTTTCPPORT, nil)
	err = server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	MQTTWSPORT := os.Getenv("MQTTWSPORT")
	if MQTTWSPORT == "" {
		MQTTWSPORT = "1882"
	}
	ws := listeners.NewWebsocket("ws1", ":"+MQTTWSPORT, nil)
	err = server.AddListener(ws)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()
	HTTPPORT := os.Getenv("HTTPPORT")
	if HTTPPORT == "" {
		HTTPPORT = "8891"
	}
	// cert
	// m := autocert.Manager{
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist("scale.rosof.tech"),
	// 	Cache:      autocert.DirCache("/var/www/.cache"),
	// }
	// tls.LoadX509KeyPair()
	srv := &http.Server{
		Addr:    ":" + HTTPPORT,
		Handler: r,
		// TLSConfig: &tls.Config{
		// 	GetCertificate: m.GetCertificate,
		// },
	}

	// SSLCRT string
	tlsConfigLoaded := false
	SSLKEY := os.Getenv("SSLKEY")
	SSLCRT := os.Getenv("SSLCRT")
	if SSLKEY != "" && SSLCRT != "" {
		x509, err := tls.LoadX509KeyPair(SSLCRT, SSLKEY)
		if err == nil {
			srv.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{x509},
			}
			tlsConfigLoaded = true
		} else {
			log.Fatalln(err)
		}
	}

	go func() {
		err = nil
		if tlsConfigLoaded {
			err = srv.ListenAndServeTLS("", "")
		} else {
			err = srv.ListenAndServe()
		}
		// service connections
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		} else {
			println("started listning: %s\n", HTTPPORT)
		}
	}()

	<-done

	server.Log.Warn().Msg("caught signal, stopping...")
	server.Close()
	server.Log.Info().Msg("main.go finished")
}
