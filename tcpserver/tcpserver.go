package tcpserver

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-kit/log"
	"github.com/influxdata/influxdb/client/v2"
)

type Serv struct {
	log.Logger
	Listener
	client.Client
	DB       string
	Grouping bool
}

type Accepter interface {
	Accept() (net.Conn, error)
	Close() error
	SetDeadline(time.Time) error
}

type Listener interface {
	Listen(string) (Accepter, error)
}

type TCPListener struct{}

func (l *TCPListener) Listen(listen string) (Accepter, error) {
	addr, err := net.ResolveTCPAddr("tcp", listen)
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", addr)
}

func New(logger log.Logger) *Serv {
	return &Serv{
		Listener: &TCPListener{},
		Logger:   logger,
	}
}

func (s *Serv) InfluxConn(c client.Client, dbname string) {
	s.Client = c
	s.DB = dbname
}

func (s *Serv) Log(args ...interface{}) {
	_ = s.Logger.Log(args...)
}

func (s *Serv) Run(listen string) error {
	l, err := s.Listener.Listen(listen)
	if err != nil {
		return fmt.Errorf("retrieving listener: %w", err)
	}

	defer l.Close()

	return s.runloop(l)
}

func (s *Serv) runloop(l Accepter) error {
	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM)

	wg := sync.WaitGroup{}

	s.Log("msg", "server starting")

	for {
		select {
		case i := <-quitChan:
			l.Close()

			s.Log("msg", "server closing", "interrupt", i)

			wg.Wait()

			s.Log("msg", "server terminated")

			return nil
		default:
		}

		if err := l.SetDeadline(time.Now().Add(1 * time.Second)); err != nil {
			return fmt.Errorf("setting deadline before listen: %w", err)
		}

		conn, err := l.Accept()
		if err != nil {
			var netErr net.Error
			if !errors.As(err, &netErr) || !netErr.Timeout() {
				s.Log("error", err)
			}

			continue
		}

		wg.Add(1)

		go s.handle(conn, &wg)
	}
}

func (s *Serv) handle(conn net.Conn, wg *sync.WaitGroup) {
	data := &Conn{
		Client:   s.Client,
		DB:       s.DB,
		Grouping: s.Grouping,
		Logger:   log.With(s.Logger, "remote", conn.RemoteAddr()),
	}

	data.handle(conn)

	conn.Close()
	wg.Done()
}
