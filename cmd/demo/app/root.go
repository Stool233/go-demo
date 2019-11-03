package app

import (
	"demo/cron"
	"demo/http_demo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

var cfg = new(http_demo.Config)

var rootCmd = &cobra.Command{
	Use:               "dictionarySyncService",
	Short:             "dictionarySyncService",
	Long:              "dictionarySyncService",
	DisableAutoGenTag: true, // disable displaying auto generation tag in cli docs
	RunE: func(cmd *cobra.Command, args []string) error {
		return run()
	},
}

func init() {
	initFlags()
}

func initFlags() {
	flagSet := rootCmd.Flags()
	// url & output
	flagSet.StringVarP(&cfg.RemoteUrl, "remoteUrl", "u", "http://dev.api.tinya.huya.com:8080/dictionary/all", "")
	flagSet.Int16VarP(&cfg.Port, "port", "p", 10080, "")

	flagSet.StringVarP(&cfg.SplitDictFilePath, "splitDictFilePath", "", "./splitDicts", "")
	flagSet.StringVarP(&cfg.StopWordFilePath, "stopWordFilePath", "", "./stopword", "")
	flagSet.StringVarP(&cfg.SynonymFilePath, "synonymFilePath", "", "./synonym", "")

	flagSet.StringVarP(&cfg.SplitDictFileName, "splitDictFileName", "", "splitDicts", "")
	flagSet.StringVarP(&cfg.StopWordFileName, "stopWordFileName", "", "stopword", "")
	flagSet.StringVarP(&cfg.SynonymFileName, "synonymFileName", "", "synonym", "")
}

func run() error {
	http_demo.StartServer(cfg)
	cron.StartCron()
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
