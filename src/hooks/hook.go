package hooks

import (
	"bytes"
	"regexp"
	"rpsoftech/scaleMQTT/src/db"
	"rpsoftech/scaleMQTT/src/global"
	"strconv"
	"strings"

	"github.com/mochi-co/mqtt/v2"
	"github.com/mochi-co/mqtt/v2/packets"
	"github.com/rs/zerolog"
)

type Options struct {
	mqtt.HookOptions
	Log *zerolog.Logger // minimal no-alloc logger
	DB  *db.DbClass
}
type ClientPacket struct {
	mqtt.Client
}

type MQTTHooks struct {
	mqtt.HookBase
	config *Options
}

var NoNumaricRegEx = regexp.MustCompile(`[^0-9 ]+`)

func (h *MQTTHooks) ID() string {
	return "MQTT Auth Hook With Publish and subscribe Method"
}

func (h *MQTTHooks) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnSubscribed,
		mqtt.OnSubscribe,
		mqtt.OnUnsubscribe,
		mqtt.OnACLCheck,
		mqtt.OnUnsubscribed,
		mqtt.OnConnectAuthenticate,
		mqtt.OnPublished,
		mqtt.OnPublish,
	}, []byte{b})
}

// func (h *MQTTHooks) SetOpts(l *zerolog.Logger, opts *mqtt.HookOptions) {
// 	println("Assinging Options")
// 	// h.Db = opts.Db
// 	// h.Log = l
// 	// h.Log.Debug().Interface("opts", opts).Str("method", "SetOpts").Send()
// }

func (h *MQTTHooks) Init(config any) error {
	h.Log.Info().Msg("initialised")
	if _, ok := config.(*Options); !ok && config != nil {
		return mqtt.ErrInvalidConfigType
	}

	if config == nil {
		config = new(Options)
	}

	h.config = config.(*Options)
	return nil
}

func (h *MQTTHooks) OnConnect(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Msgf("client connected")
	h.Log.Info().Interface("Subscription", cl.State.Subscriptions.GetAll()).Send()
}

func (h *MQTTHooks) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	if val, ok := global.MQTTConnectionStatusMap[string(cl.Properties.Username)]; ok {
		val.Count -= 1
		if val.Count <= 0 {
			val.Connected = false
		}
		// global.MQTTConnectionStatusMap[string(cl.Properties.Username)] = val
	}
	h.Log.Info().Str("client", cl.ID).Bool("expire", expire).Err(err).Msg("client disconnected")
}

// OnUnsubscribe is called when a client unsubscribes from one or more filters.
func (h *MQTTHooks) OnUnsubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	return pk
}

func (h *MQTTHooks) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	allowed := strings.HasPrefix(topic, string(cl.Properties.Username)+"/")
	h.Log.Info().Str("client", string(cl.Properties.Username)).Interface("topic", topic).Interface("Allowed", allowed).Send()
	return allowed
}
func (h *MQTTHooks) OnSubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Send()
	h.Log.Info().Int("Subscription LIST", cl.State.Subscriptions.Len()).Send()
	return pk
}

func (h *MQTTHooks) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msgf("subscribed qos=%v", reasonCodes)
}

// OnConnectAuthenticate is called when a user attempts to authenticate with the server.
func (h *MQTTHooks) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	stringUserName := string(cl.Properties.Username)
	expectedPassword, readerror := db.DBClassObject.GetPasswordForScale(stringUserName)
	allowed := readerror == nil && bytes.Equal(pk.Connect.Password, expectedPassword)

	h.Log.Info().Bytes("username", cl.Properties.Username).Bytes("password", pk.Connect.Password).Bytes("expected Password", expectedPassword).Interface("Allowed", allowed).Send()
	if allowed {
		if val, ok := global.MQTTConnectionStatusMap[stringUserName]; ok {
			val.Connected = true
			val.Weight = 0.0
			val.Count += 1
		} else {
			global.MQTTConnectionStatusMap[stringUserName] = &global.MQTTConnectionMeta{
				Connected:  true,
				UserName:   stringUserName,
				LocationID: "",
				Weight:     0.0,
				Count:      1,
			}
		}
	}
	return allowed
}

func (h *MQTTHooks) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msg("unsubscribed")
}

func (h *MQTTHooks) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	h.Log.Info().Str("client", cl.ID).Str("payload", string(pk.Payload)).Msg("received from client")

	// pkx := pk
	// if string(pk.Payload) == "hello" {
	// 	pkx.Payload = []byte("hello world")
	// 	h.Log.Info().Str("client", cl.ID).Str("payload", string(pkx.Payload)).Msg("received modified packet from client")
	// }
	if strings.HasSuffix(pk.TopicName, "WeighingScale/SerialRead") {
		stringPayload := string(pk.Payload)
		negative := strings.Contains(stringPayload, "\\f")
		i, err := strconv.ParseFloat(NoNumaricRegEx.ReplaceAllString(stringPayload, ""), 32)
		if err != nil {
			h.Log.Error().Msg(err.Error())
		} else {
			if val, ok := global.MQTTConnectionStatusMap[string(cl.Properties.Username)]; ok {
				f := i / 10
				if negative {
					f = f * -1
				}
				val.Weight = f
			}
		}
	}
	return pk, nil
}

func (h *MQTTHooks) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Str("payload", string(pk.Payload)).Msg("published to client")
}
