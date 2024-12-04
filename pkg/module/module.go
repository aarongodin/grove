package module

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Module struct {
	TargetDir           string
	Name                string
	RelationalDatabases []string
}

func NewFromPrompt(logger *zerolog.Logger) (*Module, error) {
	if logger == nil {
		logger = &log.Logger
	}
	m := &Module{}
	if err := m.PromptTargetDir(); err != nil {
		log.Fatal().Err(err).Msg("failed to get target directory")
	}
	if err := m.PromptName(); err != nil {
		log.Fatal().Err(err).Msg("failed to get module name")
	}
	return m, nil
}
