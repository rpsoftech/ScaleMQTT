package hooks

import (
	"bytes"
	"regexp"
	"rpsoftech/scaleMQTT/src/db"
	"rpsoftech/scaleMQTT/src/global"
	"rpsoftech/scaleMQTT/src/systypes"
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

var NoNumaricRegEx = regexp.MustCompile(`[^0-9.]+`)

func (h *MQTTHooks) ID() string {
	return "MQTT Auth Hook With Publish and subscribe Method"
}

func (h *MQTTHooks) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnDisconnect,
		mqtt.OnWillSent,
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
func (h *MQTTHooks) OnWillSent(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Msgf("client Closed")
}

func (h *MQTTHooks) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	if val, ok := global.MQTTConnectionWithUidStatusMap[string(cl.ID)]; ok {
		val.Count -= 1
		if val.Count <= 0 {
			val.Connected = false
		}
	}
	h.Log.Info().Str("client", cl.ID).Bool("expire", expire).Err(err).Msg("client disconnected")
}

// OnUnsubscribe is called when a client unsubscribes from one or more filters.
func (h *MQTTHooks) OnUnsubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	return pk
}

func (h *MQTTHooks) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {
	// allowed := strings.HasPrefix(topic, string(cl.ID)+"/")
	allowed := true
	h.Log.Info().Str("client", string(cl.ID)).Interface("topic", topic).Interface("Allowed", allowed).Send()
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
	stringId := string(cl.ID)
	DeviceConfig, readerror := db.DBClassObject.GetScaleConfigData(stringId)
	allowed := readerror == nil && stringUserName == DeviceConfig.MqttUsername && string(pk.Connect.Password) == DeviceConfig.MqttPassword

	h.Log.Info().Str("username", stringUserName).Bytes("password", pk.Connect.Password).Bytes("expected Password", []byte(DeviceConfig.MqttPassword)).Interface("Allowed", allowed).Send()
	if allowed {
		if val, ok := global.MQTTConnectionStatusMap[stringUserName]; ok {
			val.Connected = true
			val.Weight = 0.0
			val.Count += 1
			val.Cl = cl
		} else {
			global.MQTTConnectionStatusMap[stringUserName] = &systypes.MQTTConnectionMeta{
				Connected:  true,
				Config:     &DeviceConfig,
				UserName:   stringUserName,
				Cl:         cl,
				LocationID: "",
				Weight:     0.0,
				Count:      1,
			}
		}
		if val, ok := global.MQTTConnectionStatusMap[stringUserName]; ok {
			global.MQTTConnectionWithUidStatusMap[cl.ID] = val
		}
	}
	return allowed
}

func (h *MQTTHooks) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msg("unsubscribed")
}

func (h *MQTTHooks) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	h.Log.Info().Str("client", cl.ID).Str("payload", string(pk.Payload)).Msg("received from client")
	if strings.HasSuffix(pk.TopicName, "SerialRead") {
		if val, ok := global.MQTTConnectionWithUidStatusMap[string(cl.ID)]; ok {
			stringPayload := string(pk.Payload)
			val.RawWeight = stringPayload
			negative := strings.Contains(stringPayload, val.Config.NegativeChar)
			i, err := strconv.ParseFloat(NoNumaricRegEx.ReplaceAllString(stringPayload, ""), 32)
			if err != nil {
				h.Log.Error().Msg(err.Error())
			} else {
				f := 0.0
				if val.Config.DivideMultiply == systypes.Divide {
					f = i / float64(val.Config.DivideMultiplyBy)
				} else if val.Config.DivideMultiply == systypes.Multi {
					f = i * float64(val.Config.DivideMultiplyBy)
				}
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
