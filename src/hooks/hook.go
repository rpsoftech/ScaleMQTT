package hooks

import (
	"bytes"

	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/rs/zerolog"
)

type Options struct {
	Log *zerolog.Logger // minimal no-alloc logger
}
type ClientPacket struct {
	mqtt.Client
	// id string
	// usernamae:
}

type MQTTHooks struct {
	mqtt.HookBase
}

func (h *MQTTHooks) ID() string {
	return "MQTT Auth Hook With Publish and subscribe Method"
}

func (h *MQTTHooks) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnSubscribed,
		mqtt.OnSubscribe,
		mqtt.OnUnsubscribed,
		mqtt.OnConnectAuthenticate,
		mqtt.OnPublished,
		mqtt.OnPublish,
	}, []byte{b})
}

func (h *MQTTHooks) SetOpts(l *zerolog.Logger, opts *mqtt.HookOptions) {
	println("Log Options")
	h.Log = l
	h.Log.Debug().Interface("opts", opts).Str("method", "SetOpts").Send()
}

func (h *MQTTHooks) Init(config any) error {
	h.Log.Info().Msg("initialised")
	return nil
}

func (h *MQTTHooks) OnConnect(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Msgf("client connected")
}

func (h *MQTTHooks) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	h.Log.Info().Str("client", cl.ID).Bool("expire", expire).Err(err).Msg("client disconnectedibibibib")
}

// OnUnsubscribe is called when a client unsubscribes from one or more filters.
func (h *MQTTHooks) OnUnsubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	return pk
}

func (h *MQTTHooks) OnSubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Send()
	return pk
}

func (h *MQTTHooks) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msgf("subscribed qos=%v", reasonCodes)
}

// OnConnectAuthenticate is called when a user attempts to authenticate with the server.
func (h *MQTTHooks) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	// pk.Connect.Password
	println("Here On Auth")
	h.Log.Info().Bytes("username", cl.Properties.Username).Send()
	h.Log.Info().Bytes("password", pk.Connect.Password).Send()
	return true
}

func (h *MQTTHooks) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msg("unsubscribed")
}

func (h *MQTTHooks) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	h.Log.Info().Str("client", cl.ID).Str("payload", string(pk.Payload)).Msg("received from client")

	pkx := pk
	if string(pk.Payload) == "hello" {
		pkx.Payload = []byte("hello world")
		h.Log.Info().Str("client", cl.ID).Str("payload", string(pkx.Payload)).Msg("received modified packet from client")
	}

	return pk, nil
}

func (h *MQTTHooks) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Str("payload", string(pk.Payload)).Msg("published to client")
}
