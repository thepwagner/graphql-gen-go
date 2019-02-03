package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "graphql-gen-go",
	Short: "Generate GraphQL queries and types",
	PersistentPreRun: func(_ *cobra.Command, _ []string) {
		if debug {
			log.SetLevel(log.DebugLevel)
			log.SetFormatter(&log.TextFormatter{
				FullTimestamp: true,
			})
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Debug logging")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
