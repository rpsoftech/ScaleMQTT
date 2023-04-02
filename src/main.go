package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"rpsoftech/scaleMQTT/src/hooks"
	"syscall"

	"git.mills.io/prologic/bitcask"
	"github.com/joho/godotenv"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/listeners"
	"github.com/rs/zerolog"
)

func GetCuurentPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}

func main() {
	log.Println(GetCuurentPath())
	db, _ := bitcask.Open(filepath.Join(GetCuurentPath(), "dbcollection"))
	// defer db.Close()
	db.Put([]byte("Hello"), []byte("World11111111"))
	val, _ := db.Get([]byte("Hello"))
	db.Backup(filepath.Join(GetCuurentPath(), "dbcollection1"))
	fmt.Println(string(val))
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	godotenv.Load()
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	// r := gin.Default()

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

	<-done
	server.Log.Warn().Msg("caught signal, stopping...")
	server.Close()
	server.Log.Info().Msg("main.go finished")
}
