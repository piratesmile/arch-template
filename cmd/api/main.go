package main

import (
	"arch-template/configs"
	_ "arch-template/ent/runtime"
	"arch-template/pkg/tlog"
	"context"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

func main() {
	command := NewCommand()
	if err := command.Execute(); err != nil {
		panic(err)
	}
}

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "server",
		Short:   "Start API server",
		Example: "app server -c config.dev.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cfgFile)
		},
	}

	cmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./config.yaml", "config file(default is ./config.yaml)")
	return cmd
}

func run(configFile string) error {
	viper.AutomaticEnv()
	// load config
	conf, err := configs.Setup(configFile)
	if err != nil {
		return err
	}
	// init log
	tlog.Init(logOptions(conf.Log))
	defer tlog.Sync()

	// init router
	gin.SetMode(ginMode(conf.APP.Env))
	server := newServer(conf)
	server.Run()

	// listen shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	tlog.Info(context.Background(), "Shutdown server", nil)

	server.Stop()
	return nil
}
