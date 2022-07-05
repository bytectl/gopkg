package mqtt

import (
	"context"
	"strings"
	"time"

	"github.com/bytectl/gopkg/transport/mqtt/mux"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type MQTTSubscribe struct {
	srv *Server
}

func (m *MQTTSubscribe) SetServer(srv *Server) {
	m.srv = srv
}
func (m *MQTTSubscribe) Subscribe(c mqtt.Client, topic string, qos byte) {
	if m.srv == nil {
		panic("[MQTTSubscribe]: client must not be nil")
	}
	m.srv.Subscribe(c, topic, qos)
}

var (
	_ transport.Server = (*Server)(nil)
)

// ServerOption is an MQTT server option.
type ServerOption func(*Server)

// Broker with mqtt client broker.
func Broker(broker string) ServerOption {
	return func(o *Server) {
		o.clientOption.AddBroker(broker)
	}
}

// Username with mqtt client username.
func Username(username string) ServerOption {
	return func(o *Server) {
		o.clientOption.SetUsername(username)
	}
}

// Password with mqtt client password.
func Password(password string) ServerOption {
	return func(o *Server) {
		o.clientOption.SetPassword(password)
	}
}

// ClientId with mqtt client clientId.
func ClientId(clientId string) ServerOption {
	return func(o *Server) {
		o.clientOption.SetClientID(clientId)
	}
}

// ResumeSubs with mqtt client resumeSubs.
func ResumeSubs(resumeSubs bool) ServerOption {
	return func(o *Server) {
		o.clientOption.SetResumeSubs(resumeSubs)
	}
}

// CleanSession with mqtt client cleanSession.
func CleanSession(cleanSession bool) ServerOption {
	return func(o *Server) {
		o.clientOption.SetCleanSession(cleanSession)
	}
}

// ConnecTimeout with mqtt client connecTimeout.
func ConnectTimeout(timeout time.Duration) ServerOption {
	return func(o *Server) {
		o.clientOption.SetConnectTimeout(timeout)
	}
}

// PingTimeout with mqtt client pingTimeout.
func PingTimeout(timeout time.Duration) ServerOption {
	return func(o *Server) {
		o.clientOption.SetPingTimeout(timeout)
	}
}

// orderMatters with mqtt client orderMatters.
func OrderMatters(orderMatters bool) ServerOption {
	return func(o *Server) {
		o.clientOption.SetOrderMatters(orderMatters)
	}
}

// AutoReconnect with mqtt client autoReconnect.
func AutoReconnect(autoReconnect bool) ServerOption {
	return func(o *Server) {
		o.clientOption.SetAutoReconnect(autoReconnect)
	}
}

// DisconnectQuiesce with mqtt client disconnectQuiesce.
func DisconnectQuiesce(quiesce uint) ServerOption {
	return func(o *Server) {
		o.disconnectQuiesce = quiesce
	}
}

// Logger with server logger.
func Logger(logger log.Logger) ServerOption {
	return func(o *Server) {
		o.log = log.NewHelper(logger)
	}
}

// ConnectionLostHandler with mqtt client connectLostHandler.
func ConnectionLostHandler(connectLostHandler pmqtt.ConnectionLostHandler) ServerOption {
	return func(o *Server) {
		o.clientOption.SetConnectionLostHandler(connectLostHandler)
	}
}

// OnConnectHandler with mqtt client onConnectHandler.
func OnConnectHandler(onConnectHandler pmqtt.OnConnectHandler) ServerOption {
	return func(o *Server) {
		o.clientOption.SetOnConnectHandler(onConnectHandler)
	}
}

// RequestDecoder with request decoder.
func RequestDecoder(dec DecodeRequestFunc) ServerOption {
	return func(o *Server) {
		o.dec = dec
	}
}

// ResponseEncoder with response encoder.
func ResponseEncoder(en EncodeResponseFunc) ServerOption {
	return func(o *Server) {
		o.enc = en
	}
}

// ErrorEncoder with error encoder.
func ErrorEncoder(en EncodeErrorFunc) ServerOption {
	return func(o *Server) {
		o.ene = en
	}
}

// Middleware with service middleware option.
func Middleware(m ...middleware.Middleware) ServerOption {
	return func(o *Server) {
		o.ms = m
	}
}

type Server struct {
	log               *log.Helper
	clientOption      *pmqtt.ClientOptions
	mqttClient        pmqtt.Client
	disconnectQuiesce uint
	router            *mux.Router
	ms                []middleware.Middleware
	dec               DecodeRequestFunc
	enc               EncodeResponseFunc
	ene               EncodeErrorFunc
}

// NewServer creates an MQTT server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		clientOption: pmqtt.NewClientOptions(),
		log:          log.NewHelper(log.GetLogger()),
		dec:          DefaultRequestDecoder,
		enc:          DefaultResponseEncoder,
		ene:          DefaultErrorEncoder,
		router:       mux.NewRouter(),
	}
	for _, o := range opts {
		o(srv)
	}
	srv.router.NotFoundHandle = func(c pmqtt.Client, msg pmqtt.Message, ps *mux.Params) {
		srv.log.Error("not found handler topic: ", msg.Topic())
	}
	srv.mqttClient = pmqtt.NewClient(srv.clientOption)
	return srv
}

func (o *Server) Start(ctx context.Context) error {
	o.log.Info("[mqtt] server starting")
	token := o.mqttClient.Connect()
	if !token.WaitTimeout(1 * time.Second) {
		panic("mqtt connect wait timeout, address: " + o.clientOption.Servers[0].String())
	}
	if token.Error() != nil {
		panic("mqtt connect error, " + token.Error().Error())
	}
	return nil
}

func (o *Server) Stop(ctx context.Context) error {
	o.log.Info("[mqtt] server stopping")
	o.mqttClient.Disconnect(o.disconnectQuiesce)
	return nil
}

// Route registers an MQTT router.
func (s *Server) Route() *Router {
	return newRouter(s)
}

// Subscribe to topic
func (s *Server) Subscribe(c mqtt.Client, topic string, qos byte) {
	if c == nil {
		panic("server: client must not be nil")
	}
	// subscribe to topic
	subscribeTopic := s.makeSubscribeTopic(topic)
	c.Subscribe(subscribeTopic, qos, s.router.ServeMQTT)
	log.Debugf("[server] subscribe to topic: %s", subscribeTopic)
}

func (s *Server) makeSubscribeTopic(topic string) string {
	dirs := strings.Split(topic, "/")
	for i, dir := range dirs {
		if dir == "" {
			continue
		}
		if dir[0] == ':' {
			dirs[i] = "+"
		}
		if dir[0] == '*' {
			dirs[i] = "#"
		}
	}
	return strings.Join(dirs, "/")
}
