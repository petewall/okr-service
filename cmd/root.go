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
		datastore := &internal.InMemoryDatastore{
			OKRs: []*internal.OKR{
				internal.CreateOKR("2024 Q1", "Internal Work", internal.OKR_TYPE_BOOLEAN, "Create OKR Dashboard", 1),
				internal.CreateOKR("2024 Q1", "Personal Growth", internal.OKR_TYPE_NUMBER, "Run 5 5ks", 5),
			},
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

	rootCmd.Flags().IntP("port", "p", internal.DEFAULT_PORT, "Port number to listen on (env: PORT)")
	_ = viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))

	rootCmd.Flags().Bool("debug", false, "Enable debug logging (env: DEBUG)")
	_ = viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug"))
	viper.AutomaticEnv()
}
