package hooks

import (
	"bytes"

	"github.com/google/uuid"
	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
)

type ClientPacket struct {
	mqtt.Client
	id string
	// usernamae:
}

type MQTTHooks struct {
	mqtt.HookBase
}

func (h *MQTTHooks) ID() string {
	return uuid.New().String()
}

func (h *MQTTHooks) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnSubscribed,
		mqtt.OnUnsubscribed,
		mqtt.OnConnectAuthenticate,
		mqtt.OnPublished,
		mqtt.OnPublish,
	}, []byte{b})
}

func (h *MQTTHooks) Init(config any) error {
	h.Log.Info().Msg("initialised")
	return nil
}

func (h *MQTTHooks) OnConnect(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Msgf("client connected")
}

func (h *MQTTHooks) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	h.Log.Info().Str("client", cl.ID).Bool("expire", expire).Err(err).Msg("client disconnected")
}

func (h *MQTTHooks) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msgf("subscribed qos=%v", reasonCodes)
}

// OnConnectAuthenticate is called when a user attempts to authenticate with the server.
func (h *MQTTHooks) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	// pk.Connect.Password
	h.Log.Log().Bytes("password", pk.Connect.Password)
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

	return pkx, nil
}

func (h *MQTTHooks) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Str("payload", string(pk.Payload)).Msg("published to client")
}
