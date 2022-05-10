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

// DisconnectQuiesce with mqtt client disconnectQuiesce.
func DisconnectQuiesce(quiesce uint) ServerOption {
	return func(s *Server) {
		s.disconnectQuiesce = quiesce
	}
}

// Logger with server logger.
func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}

// ConnectionLostHandler with mqtt client connectLostHandler.
func ConnectionLostHandler(connectLostHandler pmqtt.ConnectionLostHandler) ServerOption {
	return func(s *Server) {
		s.clientOption.SetConnectionLostHandler(connectLostHandler)
	}
}

// OnConnectHandler with mqtt client onConnectHandler.
func OnConnectHandler(onConnectHandler pmqtt.OnConnectHandler) ServerOption {
	return func(s *Server) {
		s.clientOption.SetOnConnectHandler(onConnectHandler)
	}
}
func SetRouter(router *Router) ServerOption {
	return func(s *Server) {
		s.router = router
	}
}

type Server struct {
	log               *log.Helper
	clientOption      *pmqtt.ClientOptions
	mqttClient        pmqtt.Client
	disconnectQuiesce uint
	router            *Router
}

// NewServer creates an MQTT server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		clientOption: pmqtt.NewClientOptions(),
		log:          log.NewHelper(log.GetLogger()),
	}
	for _, o := range opts {
		o(srv)
	}
	srv.mqttClient = pmqtt.NewClient(srv.clientOption)
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	s.log.Info("[mqtt] server starting")
	if token := s.mqttClient.Connect(); !token.WaitTimeout(1*time.Second) || token.Error() != nil {
		panic("mqtt connect error, " + token.Error().Error())
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("[mqtt] server stopping")
	s.mqttClient.Disconnect(s.disconnectQuiesce)
	return nil
}
