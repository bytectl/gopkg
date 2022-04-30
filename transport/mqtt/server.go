package mqtt

import (
	"context"
	"time"

	pmqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

var (
	_ transport.Server = (*Server)(nil)
)

// ServerOption is an MQTT server option.
type ServerOption func(*Server)

// Broker with mqtt client broker.
func Broker(broker string) ServerOption {
	return func(s *Server) {
		s.clientOption.AddBroker(broker)
	}
}

// Username with mqtt client username.
func Username(username string) ServerOption {
	return func(s *Server) {
		s.clientOption.SetUsername(username)
	}
}

// Password with mqtt client password.
func Password(password string) ServerOption {
	return func(s *Server) {
		s.clientOption.SetPassword(password)
	}
}

// ClientId with mqtt client clientId.
func ClientId(clientId string) ServerOption {
	return func(s *Server) {
		s.clientOption.SetClientID(clientId)
	}
}

// ResumeSubs with mqtt client resumeSubs.
func ResumeSubs(resumeSubs bool) ServerOption {
	return func(s *Server) {
		s.clientOption.SetResumeSubs(resumeSubs)
	}
}

// CleanSession with mqtt client cleanSession.
func CleanSession(cleanSession bool) ServerOption {
	return func(s *Server) {
		s.clientOption.SetCleanSession(cleanSession)
	}
}

// ConnecTimeout with mqtt client connecTimeout.
func ConnectTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.clientOption.SetConnectTimeout(timeout)
	}
}

// PingTimeout with mqtt client pingTimeout.
func PingTimeout(timeout time.Duration) ServerOption {
	return func(s *Server) {
		s.clientOption.SetPingTimeout(timeout)
	}
}

// orderMatters with mqtt client orderMatters.
func OrderMatters(orderMatters bool) ServerOption {
	return func(s *Server) {
		s.clientOption.SetOrderMatters(orderMatters)
	}
}

// AutoReconnect with mqtt client autoReconnect.
func AutoReconnect(autoReconnect bool) ServerOption {
	return func(s *Server) {
		s.clientOption.SetAutoReconnect(autoReconnect)
	}
}

// Logger with server logger.
func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}

type Server struct {
	log                *log.Helper
	clientOption       *pmqtt.ClientOptions
	mqttClient         pmqtt.Client
	subscribes         []func(*Server)
	subscribeMultiples []func(*Server)
}

// NewServer creates an MQTT server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		mqttClient:   nil,
		clientOption: pmqtt.NewClientOptions(),
		log:          log.NewHelper(log.GetLogger()),
	}
	for _, o := range opts {
		o(srv)
	}
	return srv
}

// Subscribe with mqtt client subscribe. must be called before Start.
func (s *Server) Subscribe(topic string, qos byte, callback pmqtt.MessageHandler) {
	fn := func(srv *Server) {
		srv.mqttClient.Subscribe(topic, qos, callback)
	}
	s.subscribes = append(s.subscribes, fn)
}

// SubscribeMultiple with mqtt client subscribeMultiple. must be called before Start.
func (s *Server) SubscribeMultiple(filters map[string]byte, callback pmqtt.MessageHandler) {
	fn := func(srv *Server) {
		srv.mqttClient.SubscribeMultiple(filters, callback)
	}
	s.subscribeMultiples = append(s.subscribeMultiples, fn)
}

func (s *Server) Start(ctx context.Context) error {

	s.clientOption.SetConnectionLostHandler(func(c pmqtt.Client, err error) {
		//TODO: 统计错误次数
		log.Error(err)
		s.log.Debugf("mqtt connection lost: %v", c)
	})

	s.clientOption.SetOnConnectHandler(func(c pmqtt.Client) {
		// TODO: 统计重连次数
		s.log.Info("mqtt onConnect")
		s.log.Debugf("mqtt onConnect: %v", c)
		// 订阅topic
		s.subscribe()
	})
	if s.mqttClient == nil {
		s.mqttClient = pmqtt.NewClient(s.clientOption)
	}
	s.log.Info("[mqtt] server starting")
	if token := s.mqttClient.Connect(); !token.WaitTimeout(1*time.Second) || token.Error() != nil {
		panic("mqtt connect error, " + token.Error().Error())
	}

	return nil
}

func (s *Server) subscribe() {
	// subscribe mqtt topic
	for _, fn := range s.subscribes {
		fn(s)
	}
	for _, fn := range s.subscribeMultiples {
		fn(s)
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("[mqtt] server stopping")
	s.mqttClient.Disconnect(250)
	return nil
}
