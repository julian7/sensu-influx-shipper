package tcpserver

import (
	"io"

	"github.com/go-kit/kit/log"
	"github.com/influxdata/influxdb/client/v2"
)

type Conn struct {
	log.Logger
	client.Client
	DB string
}

func (c *Conn) handle(reader io.Reader) {
}
