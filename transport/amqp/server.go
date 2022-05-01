package amqp

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	ramqp "github.com/rabbitmq/amqp091-go"
)

var (
	_ transport.Server = (*Server)(nil)
)

// ServerOption is an MQTT server option.
type ServerOption func(*Server)

// Url with amqp client url.
func Url(url string) ServerOption {
	return func(s *Server) {
		s.url = url
	}
}

// SASL with amqp client sasl.
func SASL(sasl []ramqp.Authentication) ServerOption {
	return func(s *Server) {
		s.config.SASL = sasl
	}
}

// Vhost with amqp client vhost.
func Vhost(vhost string) ServerOption {
	return func(s *Server) {
		s.config.Vhost = vhost
	}
}

// ChannelMax 0 max channels means 2^16 - 1
func ChannelMax(channelMax int) ServerOption {
	return func(s *Server) {
		s.config.ChannelMax = channelMax
	}
}

// FrameSize 0 max bytes means unlimited
func FrameSize(frameSize int) ServerOption {
	return func(s *Server) {
		s.config.FrameSize = frameSize
	}
}

// Heartbeat less than 1s uses the server's interval.
func Heartbeat(heartbeat time.Duration) ServerOption {
	return func(s *Server) {
		s.config.Heartbeat = heartbeat
	}
}

// TLSClientConfig
func TLSClientConfig(tlsClientConfig *tls.Config) ServerOption {
	return func(s *Server) {
		s.config.TLSClientConfig = tlsClientConfig
	}
}

// Properties
func Properties(properties ramqp.Table) ServerOption {
	return func(s *Server) {
		s.config.Properties = properties
	}
}

// Locale with amqp client
func Locale(locale string) ServerOption {
	return func(s *Server) {
		s.config.Locale = locale
	}
}
func ReconnectInterval(interval time.Duration) ServerOption {
	return func(s *Server) {
		s.reconnectInterval = interval
	}
}

// Logger with amqp client logger.
func Logger(logger log.Logger) ServerOption {
	return func(s *Server) {
		s.log = log.NewHelper(logger)
	}
}

// OnConnectHandler is called when the server is connected to the client.
func OnConnectHandler(fn OnConnect) ServerOption {
	return func(s *Server) {
		s.onConnect = fn
	}
}

// ConnectionLostHandler is called when the server is disconnected from the client.
func ConnectionLostHandler(fn ConnectionLost) ServerOption {
	return func(s *Server) {
		s.connectionLost = fn
	}
}

type ConnectionLost func(*ramqp.Connection, *ramqp.Error)
type OnConnect func(*ramqp.Connection)

func defaultOnConnect(*ramqp.Connection)                    {}
func defaultConnectionLost(*ramqp.Connection, *ramqp.Error) {}

type Server struct {
	log               *log.Helper
	url               string
	config            *ramqp.Config
	amqpConn          *ramqp.Connection
	notifyCloseChan   chan *ramqp.Error
	reconnectInterval time.Duration
	onConnect         OnConnect
	connectionLost    ConnectionLost
}

// NewServer creates an MQTT server by options.
func NewServer(opts ...ServerOption) *Server {
	srv := &Server{
		config:            &ramqp.Config{},
		log:               log.NewHelper(log.GetLogger()),
		notifyCloseChan:   make(chan *ramqp.Error),
		reconnectInterval: time.Second * 5,
		onConnect:         defaultOnConnect,
		connectionLost:    defaultConnectionLost,
	}
	for _, o := range opts {
		o(srv)
	}
	return srv
}

func (s *Server) Start(ctx context.Context) error {
	s.log.Info("[amqp] server starting")
	conn, err := ramqp.DialConfig(s.url, *s.config)
	if err != nil {
		panic("amqp connect error," + err.Error())
	}
	// set notifyClose channal
	conn.NotifyClose(s.notifyCloseChan)
	s.amqpConn = conn
	// call onConnect
	s.onConnect(s.amqpConn)
	go s.connNotifyClose(ctx)
	return nil
}

func (s *Server) connNotifyClose(ctx context.Context) {
	for {
		select {
		case err := <-s.notifyCloseChan:
			s.log.Errorf("[amqp] conn closed,error(%v)", err)
			// call connectionLost
			s.connectionLost(s.amqpConn, err)
			// reconnect
			s.reconnect(ctx)
		case <-ctx.Done():
			s.log.Info("[amqp] conn notifyClose done")
			return
		}
	}
}
func (s *Server) reconnect(ctx context.Context) {
	s.log.Info("try reconnect amqp server")
	// renew chan
	s.notifyCloseChan = make(chan *ramqp.Error)
	for {
		select {
		case <-ctx.Done():
			s.log.Info("[amqp] server reconnect done")
			return
		default:
			s.log.Info("[amqp] server reconnecting")
			conn, err := ramqp.DialConfig(s.url, *s.config)
			if err != nil {
				s.log.Errorf("[amqp] server reconnect error(%v)", err)
				time.Sleep(s.reconnectInterval)
				continue
			}
			// set notifyClose channal
			conn.NotifyClose(s.notifyCloseChan)
			s.amqpConn = conn
			// call onConnect
			s.onConnect(s.amqpConn)
			return
		}
	}
}

func (s *Server) Stop(ctx context.Context) error {
	s.log.Info("[amqp] server stopping")
	if s.amqpConn != nil {
		return s.amqpConn.Close()
	}
	return nil
}
