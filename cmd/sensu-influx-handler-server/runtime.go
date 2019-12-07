package main

import (
	"fmt"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Runtime struct {
	*viper.Viper
	log.Logger
	Config string
}

func NewRuntime() *Runtime {
	rt := &Runtime{
		Viper:  viper.New(),
		Logger: log.NewNopLogger(),
	}

	cobra.OnInitialize(rt.Init)

	return rt
}

func (rt *Runtime) Init() {
	if rt.Config != "" {
		rt.Viper.SetConfigFile(rt.Config)
	} else {
		rt.Viper.AddConfigPath("/etc")
		rt.Viper.AddConfigPath("$HOME")
		rt.Viper.SetConfigName("sensu-influx-handler-server")
	}

	rt.Viper.AutomaticEnv()
	rt.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := rt.Viper.ReadInConfig(); err == nil {
		_ = rt.Log("msg", "reading config", "file", rt.Viper.ConfigFileUsed())
	}
}

func (rt *Runtime) RegisterFlags(group string, l *pflag.FlagSet) (err error) {
	l.VisitAll(func(flag *pflag.Flag) {
		if err != nil {
			return
		}
		err = rt.Viper.BindPFlag(
			fmt.Sprintf("%s.%s", group, flag.Name),
			flag,
		)
	})

	return err
}
