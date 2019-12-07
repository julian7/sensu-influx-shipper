package main

import (
	"fmt"
	"io"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (rt *Runtime) rootCmd() (*cobra.Command, error) {
	app := &cobra.Command{
		Use:   "sensu-influx-handler-server",
		Short: "Service for forwarding sensu events to InfluxDB via TCP",
		Long: `Stand-alone service for receiving events from sensu-go TCP handlers, and
forwarding them to an InfluxDB server. This approach saves compute resources,
as it doesn't require executing a command every time a metric time has to be
forwarded to InfluxDB.`,
		PersistentPreRunE: rt.SetupLogging,
		Version:           version,
	}

	flags := app.PersistentFlags()
	flags.StringVar(&rt.Config, "config", "", "configuration file (default: /etc/sensu-influx-handler-server[.yml]")
	flags.StringP("logformat", "F", "logfmt", "log format. Possible values: logfmt, or json")
	flags.StringP("logfile", "L", "stderr", "log file. Possible values: none, stdout, stderr, or file name")

	err := viper.BindPFlags(flags)
	if err != nil {
		return nil, err
	}

	for _, cmdFunc := range []func() (*cobra.Command, error){
		rt.serveCmd,
	} {
		cmd, err := cmdFunc()
		if err != nil {
			return nil, err
		}

		app.AddCommand(cmd)
	}

	return app, nil
}

func (rt *Runtime) SetupLogging(*cobra.Command, []string) error {
	var outFD io.Writer

	var err error

	var logger log.Logger

	format := rt.Viper.GetString("logformat")
	output := rt.Viper.GetString("logfile")

	switch output {
	case "none":
		return nil
	case "stdout":
		outFD = os.Stdout
	case "stderr":
		outFD = os.Stderr
	default:
		outFD, err = os.OpenFile(output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
		if err != nil {
			return fmt.Errorf("opening logfile: %w", err)
		}
	}

	logWriter := log.NewSyncWriter(outFD)

	switch format {
	case "logfmt":
		logger = log.NewLogfmtLogger(logWriter)
	case "json":
		logger = log.NewJSONLogger(logWriter)
	default:
		return fmt.Errorf("unknown log format: %s", format)
	}

	rt.Logger = log.With(logger, "pid", os.Getpid())

	return nil
}
