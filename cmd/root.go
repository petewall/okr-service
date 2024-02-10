/*
Copyright © 2024 Pete Wall <pete@petewall.net>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/petewall/okr-service/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use: "okr-service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetBool("debug") {
			log.SetLevel(log.DebugLevel)
			log.Debug("debug mode enabled")
		}

		log.Info("Creating Datastore...")
		datastore, err := internal.InitializeDatastore()
		if err != nil {
			return err
		}

		server := &internal.Server{
			Datastore: datastore,
			Port:      viper.GetInt("port"),
		}
		return server.Start()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	log.SetLevel(log.InfoLevel)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	rootCmd.Flags().IntP("port", "p", internal.DefaultPort, "Port number to listen on (env: PORT)")
	_ = viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))

	rootCmd.Flags().String("datastore", "fs", "Datastore type (memory, fs) (env: DATASTORE_TYPE)")
	_ = viper.BindPFlag("datastore.type", rootCmd.Flags().Lookup("datastore"))

	rootCmd.Flags().String("format", "yaml", "Datastore storage format (only for fs type) (env: DATASTORE_FORMAT)")
	_ = viper.BindPFlag("datastore.format", rootCmd.Flags().Lookup("format"))

	rootCmd.Flags().String("path", "", "Datastore storage path (only for fs type) (env: DATASTORE_PATH)")
	_ = viper.BindPFlag("datastore.path", rootCmd.Flags().Lookup("path"))

	rootCmd.Flags().Bool("debug", false, "Enable debug logging (env: DEBUG)")
	_ = viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug"))
	viper.AutomaticEnv()
}
