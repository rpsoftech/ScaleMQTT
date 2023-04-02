package main

import (
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

func main() {
	log.Println(global.GetCuurentPath())
	db, _ := bitcask.Open(filepath.Join(global.GetCuurentPath(), "dbcollection"))
	dbPackage.DbConnection = db
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	godotenv.Load()
	keyJWT := os.Getenv("JWTKEY")
	println(keyJWT)
	println(global.JWTKEY)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	r := gin.Default()
	routes.AdminRoutes(r)
	server := mqtt.New(nil)
	l := server.Log.Level(zerolog.DebugLevel)
	server.Log = &l

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

	tcp := listeners.NewTCP("t1", ":1883", nil)
	err = server.AddListener(tcp)
	if err != nil {
		log.Fatal(err)
	}

	ws := listeners.NewWebsocket("ws1", ":1882", nil)
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
	srv := &http.Server{
		Addr:    ":8891",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-done

	server.Log.Warn().Msg("caught signal, stopping...")
	server.Close()
	server.Log.Info().Msg("main.go finished")
}
