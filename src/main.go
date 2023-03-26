package main

import (
	"log"
	"os"
	"os/signal"
	hooks "rpsoftech/scaleMQTT/src/hooks"
	"syscall"

	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/hooks/debug"
	"github.com/mochi-co/mqtt/v2/listeners"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()

	server := mqtt.New(nil)
	// l := server.Log.Level(zerolog.DebugLevel)
	// server.Log = &l

	err := server.AddHook(new(hooks.MQTTHooks), nil)
	if err != nil {
		log.Fatal(err)
	}

	err = server.AddHook(new(debug.Hook), &debug.Options{
		// ShowPacketData: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	tcp := listeners.NewTCP("t1", ":1883", nil)
	err = server.AddListener(tcp)
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
