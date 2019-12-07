package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/julian7/sensu-influx-handler-server/tcpserver"
	"github.com/spf13/cobra"
)

type serve struct {
	log.Logger
	listen string
	addr   string
	name   string
	user   string
	pass   string
}

func (rt *Runtime) serveCmd() (*cobra.Command, error) {
	serv := &serve{Logger: rt.Logger}

	cmd := &cobra.Command{
		Use:           "serve",
		Short:         "runs TCP server",
		RunE:          serv.Run,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	flags := cmd.Flags()
	flags.StringVarP(&serv.listen, "listen", "l", "127.0.0.1:3333", "TCP port to listen to")
	flags.StringVarP(&serv.addr, "addr", "a", "http://127.0.0.1:8086", "InfluxDB's TCP port")
	flags.StringVarP(&serv.name, "database", "d", "metrics", "InfluxDB database")
	flags.StringVarP(&serv.user, "user", "u", "metrics", "InfluxDB username")
	flags.StringVarP(&serv.pass, "pass", "p", "", "InfluxDB password")

	return cmd, rt.RegisterFlags(cmd.Use, flags)
}

func (s *serve) Run(*cobra.Command, []string) error {
	serv := tcpserver.New(log.With(s.Logger, "listen", s.listen))
	conf := client.HTTPConfig{
		Addr:     s.addr,
		Username: s.user,
		Password: s.pass,
	}

	c, err := client.NewHTTPClient(conf)
	if err != nil {
		return fmt.Errorf("initiazilizng influxdb config: %w", err)
	}

	if _, _, err := c.Ping(10 * time.Second); err != nil {
		return fmt.Errorf("connecting influxdb server: %w", err)
	}

	serv.InfluxConn(c, s.name)

	return serv.Run(s.listen)
}
