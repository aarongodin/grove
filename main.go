package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grove",
	Short: "grove initializes Go web services",
	Long: `Grove is a modular Go project initializer that increases time-to-value for building web services.
		It includes an optional set of libraries for common usecases and removes a large amount of project boilerplate.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info().Msgf("grove v%s", Version)
		m := &Module{}
		if err := m.PromptTargetDir(); err != nil {
			log.Fatal().Err(err).Msg("failed to get target directory")
		}
		if err := m.PromptName(); err != nil {
			log.Fatal().Err(err).Msg("failed to get module name")
		}
		if err := m.Render(); err != nil {
			log.Fatal().Err(err).Msg("failed to render module")
		}
	},
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute root command")
	}
}
