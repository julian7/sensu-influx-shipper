package tcpserver

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/influxdata/influxdb/client/v2"
)

const (
	errBadAddress     = "bad address"
	errDeadlineZeroed = "deadline counter zeroed"
)

var ErrDeadlineZeroed = errors.New(errDeadlineZeroed)
var ErrBadAddress = errors.New(errBadAddress)

type nopAccepter struct {
	deadlineCounter int
}

func (*nopAccepter) Accept() (net.Conn, error) {
	return &net.TCPConn{}, nil
}

func (*nopAccepter) Close() error {
	return nil
}

func (a *nopAccepter) SetDeadline(time.Time) error {
	a.deadlineCounter--

	if a.deadlineCounter <= 0 {
		return ErrDeadlineZeroed
	}

	return nil
}

type nopListener struct{}

func (*nopListener) Listen(listen string) (Accepter, error) {
	if listen == "bad" {
		return nil, ErrBadAddress
	}

	return &nopAccepter{deadlineCounter: 10}, nil
}

func TestServ_Run(t *testing.T) {
	client, _ := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://127.0.0.1:8086",
	})
	defServ := &Serv{
		Client:   client,
		DB:       "nop",
		Logger:   log.NewNopLogger(),
		Listener: &nopListener{},
	}
	tests := []struct {
		name    string
		serv    *Serv
		listen  string
		wantErr bool
	}{
		{"normal", defServ, "127.0.0.1:1234", false},
		{"bad address", defServ, "bad", true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.serv.Run(tt.listen)

			if !errors.Is(err, ErrDeadlineZeroed) && (err != nil) != tt.wantErr {
				t.Errorf("Serv.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
