package app

import (
	"demo/cron"
	"demo/http_demo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:               "dictionarySyncService",
	Short:             "dictionarySyncService",
	Long:              "dictionarySyncService",
	DisableAutoGenTag: true, // disable displaying auto generation tag in cli docs
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func run() error {
	http_demo.StartServer()
	cron.StartCron()
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
