package main

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	dbPackage "rpsoftech/scaleMQTT/src/db"
	global "rpsoftech/scaleMQTT/src/global"
	"rpsoftech/scaleMQTT/src/hooks"
	"rpsoftech/scaleMQTT/src/routes"
	"syscall"

	"git.mills.io/prologic/bitcask"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/rs/zerolog"
)

func LoadEnv() {
	println(global.GetCuurentPath())
	if _, err := os.Stat(".env"); err == nil {
		// path/to/whatever exists
		godotenv.Load(".env")
	} else {
		godotenv.Load("./../.env")
	}
	defaultValue := make([]byte, 128)

	_, err := rand.Read(defaultValue)
	if err != nil {
		defaultValue = []byte("thisisjustdefaultvalue")
	}
	defaultValueString := hex.EncodeToString(defaultValue)
	envJWTKeyValue := os.Getenv("JWTKEY")
	if envJWTKeyValue == "" {
		envJWTKeyValue = defaultValueString
	}
	global.JWTKEY = []byte(envJWTKeyValue)
}

func main() {
	log.Println(global.GetCuurentPath())
	db, _ := bitcask.Open(filepath.Join(global.GetCuurentPath(), "dbcollection"))
	dbPackage.DbConnection = db
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	LoadEnv()
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

	server.Log.Debug().Bytes("JWTKEY", global.JWTKEY).Send()
	err := server.AddHook(new(hooks.MQTTHooks), &hooks.Options{
		Db: db,
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
	// m := autocert.Manager{
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist("scale.rosof.tech"),
	// 	Cache:      autocert.DirCache("/var/www/.cache"),
	// }
	srv := &http.Server{
		Addr:    ":" + HTTPPORT,
		Handler: r,
		// TLSConfig: &tls.Config{
		// 	GetCertificate: m.GetCertificate,
		// },
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
