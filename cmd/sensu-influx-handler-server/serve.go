package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/julian7/sensu-influx-handler-server/tcpserver"
	"github.com/spf13/cobra"
)

func (rt *Runtime) serveCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:           "serve",
		Short:         "Sensu-go event consuming TCP server for InfluxDB data shipping",
		RunE:          rt.serveRun,
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	flags := cmd.Flags()
	flags.StringP("listen", "l", "127.0.0.1:3333", "TCP port to listen to")
	flags.StringP("addr", "a", "http://127.0.0.1:8086", "InfluxDB's TCP port")
	flags.StringP("database", "d", "metrics", "InfluxDB database")
	flags.StringP("user", "u", "metrics", "InfluxDB username")
	flags.StringP("pass", "p", "", "InfluxDB password")

	return cmd, rt.RegisterFlags(cmd.Use, flags)
}

func (rt *Runtime) serveRun(*cobra.Command, []string) error {
	conf := client.HTTPConfig{
		Addr:     rt.GetString("serve.addr"),
		Username: rt.GetString("serve.user"),
		Password: rt.GetString("serve.pass"),
	}

	c, err := client.NewHTTPClient(conf)
	if err != nil {
		return fmt.Errorf("initializing influxdb config: %w", err)
	}

	if _, _, err := c.Ping(10 * time.Second); err != nil {
		return fmt.Errorf("connecting influxdb server: %w", err)
	}

	serv := tcpserver.New(log.With(rt.Logger, "listen", rt.GetString("serve.listen")))
	serv.InfluxConn(c, rt.GetString("serve.name"))

	return serv.Run(rt.GetString("serve.listen"))
}
