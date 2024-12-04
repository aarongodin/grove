package main

import (
	"os"

	"github.com/aarongodin/grove/pkg/module"
	"github.com/aarongodin/grove/pkg/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var debug bool

var rootCmd = &cobra.Command{
	Use:   "grove",
	Short: "grove initializes Go web service projects",
	Long: `Grove is a modular Go project initializer that increases time-to-value for building web services.
		It includes an optional set of libraries for common usecases and removes a large amount of project boilerplate.`,
	Run: func(cmd *cobra.Command, _ []string) {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		if debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		log.Info().Msgf("grove v%s", Version)
		mod, err := module.NewFromPrompt(nil)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to prompt for module details")
		}
		if err := render.Render(mod, debug, nil); err != nil {
			log.Fatal().Err(err).Msg("failed to render module")
		}
		log.Info().Msg("module render complete!")
	},
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enabled debug logging")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute root command")
	}
}
